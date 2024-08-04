package cmd

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/FusionFoundation/efsn/v4/p2p/discover"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	"github.com/ethereum/go-ethereum/rlp"
	zlog "github.com/rs/zerolog/log"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd/flag"
)

var FsnCmd = cli.Command{
	Name:  "fsn",
	Usage: "fsn actions",
	Subcommands: []cli.Command{
		fsnDiscoverCmd,
	},
}

var fsnDiscoverCmd = cli.Command{
	Name:   "discover",
	Usage:  "manual discover",
	Action: fsnDiscover,
	Flags: []cli.Flag{
		flag.OptionalListenFlag,
	},
}

// RPC request structures
type (
	ping struct {
		Version    uint
		From, To   rpcEndpoint
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}

	// pong is the reply to ping.
	pong struct {
		// This field should mirror the UDP envelope address
		// of the ping packet, which provides a way to discover the
		// the external address (after NAT).
		To rpcEndpoint

		ReplyTok   []byte // This contains the hash of the ping packet.
		Expiration uint64 // Absolute timestamp at which the packet becomes invalid.
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}

	// findnode is a query for nodes close to the given target.
	findnode struct {
		Target     discover.NodeID // doesn't need to be an actual public key
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}

	// reply to findnode
	neighbors struct {
		Nodes      []rpcNode
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}

	rpcNode struct {
		IP  net.IP // len 4 for IPv4 or 16 for IPv6
		UDP uint16 // for discovery protocol
		TCP uint16 // for RLPx protocol
		ID  discover.NodeID
	}

	rpcEndpoint struct {
		IP  net.IP // len 4 for IPv4 or 16 for IPv6
		UDP uint16 // for discovery protocol
		TCP uint16 // for RLPx protocol
	}
)

func (req *ping) Name() string { return "PING/v4" }
func (req *ping) Kind() byte   { return pingPacket }

func (req *pong) Name() string { return "PONG/v4" }
func (req *pong) Kind() byte   { return pongPacket }

func (req *findnode) Name() string { return "FINDNODE/v4" }
func (req *findnode) Kind() byte   { return findnodePacket }

func (req *neighbors) Name() string { return "NEIGHBORS/v4" }
func (req *neighbors) Kind() byte   { return neighborsPacket }

func fsnDiscover(ctx *cli.Context) (err error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	db, err := enode.OpenDB("")
	if err != nil {
		return err
	}
	localnode := enode.NewLocalNode(db, key)
	localnode.SetFallbackIP(net.IP{127, 0, 0, 1})

	addr := ctx.String(flag.OptionalListenFlag.Name)
	if addr == "" {
		addr = ":40408"
	}
	listenAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	a := &net.UDPAddr{IP: localnode.Node().IP(), Port: localnode.Node().UDP()}
	ourEndpoint := makeEndpoint(a, uint16(localnode.Node().TCP()))

	enodeStr := ctx.Args()[0]
	n := discover.MustParseNode(enodeStr)
	toid := n.ID
	toaddr := &net.UDPAddr{IP: n.IP, Port: int(n.UDP)}

	expiration := 20 * time.Second
	req := &ping{
		Version:    4,
		From:       ourEndpoint,
		To:         makeEndpoint(toaddr, 0), // TODO: maybe use known TCP port from DB
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	}
	packetBytes, hash, err := encodePacket(key, req.Kind(), req)
	if err != nil {
		return err
	}

	_, err = conn.WriteToUDP(packetBytes, toaddr)
	if err != nil {
		return err
	}
	zlog.Info().Str("toaddr", toaddr.String()).Msg("udp packet sent")

	cb := func(p packet, fromID discover.NodeID) error {
		if fromID != toid {
			return fmt.Errorf("id mismatch")
		}
		if p.Kind() != pongPacket {
			return fmt.Errorf("type mismatch")
		}
		if !bytes.Equal(p.(*pong).ReplyTok, hash) {
			return fmt.Errorf("hash mismatch")
		}
		return nil
	}
	// Discovery packets are defined to be no larger than 1280 bytes.
	// Packets larger than this size will be cut at the end and treated
	// as invalid because their hash won't match.
	buf := make([]byte, 1280)
	var (
		nbytes int
		from   *net.UDPAddr
	)
	for {
		nbytes, from, err = conn.ReadFromUDP(buf)
		if netutil.IsTemporaryError(err) {
			// Ignore temporary read errors.
			zlog.Info().Err(err).Msg("Temporary UDP read error")
			continue
		} else if err != nil {
			// Shut down the loop for permament errors.
			zlog.Info().Err(err).Msg("UDP read erro")
			return
		}
		if err = handlePacket(from, buf[:nbytes], cb); err != nil {
			zlog.Info().Err(err).Msg("UDP unhandled error")
			return
		}
		break
	}
	return
}

func makeEndpoint(addr *net.UDPAddr, tcpPort uint16) rpcEndpoint {
	ip := addr.IP.To4()
	if ip == nil {
		ip = addr.IP.To16()
	}
	return rpcEndpoint{IP: ip, UDP: uint16(addr.Port), TCP: tcpPort}
}

func handlePacket(from *net.UDPAddr, buf []byte, cb func(p packet, fromID discover.NodeID) error) error {
	packet, fromID, _, err := decodePacket(buf)
	if err != nil {
		zlog.Info().Err(err).Str("addr", from.String()).Msg("Bad discv4 packet")
		return err
	}

	zlog.Info().Str("type", packet.Name()).Msg("got packet")

	err = cb(packet, fromID)
	if err != nil {
		return err
	}

	zlog.Info().Msg("<< " + packet.Name())
	return nil
}

const (
	macSize  = 256 / 8
	sigSize  = 520 / 8
	headSize = macSize + sigSize // space of packet frame data
)

var (
	headSpace = make([]byte, headSize)
)

var (
	errPacketTooSmall = errors.New("too small")
	errBadHash        = errors.New("bad hash")
)

// RPC packet types
const (
	pingPacket = iota + 40 // zero is 'reserved'
	pongPacket
	findnodePacket
	neighborsPacket
)

type packet interface {
	Kind() byte
	Name() string
}

func encodePacket(priv *ecdsa.PrivateKey, ptype byte, req interface{}) (packet, hash []byte, err error) {
	b := new(bytes.Buffer)
	b.Write(headSpace)
	b.WriteByte(ptype)
	if err := rlp.Encode(b, req); err != nil {
		zlog.Error().Err(err).Msg("Can't encode discv4 packet")
		return nil, nil, err
	}
	packet = b.Bytes()
	sig, err := crypto.Sign(crypto.Keccak256(packet[headSize:]), priv)
	if err != nil {
		zlog.Error().Err(err).Msg("Can't sign discv4 packe")
		return nil, nil, err
	}
	copy(packet[macSize:], sig)
	// add the hash to the front. Note: this doesn't protect the
	// packet in any way. Our public key will be part of this hash in
	// The future.
	hash = crypto.Keccak256(packet[macSize:])
	copy(packet, hash)
	return packet, hash, nil
}

func decodePacket(buf []byte) (packet, discover.NodeID, []byte, error) {
	if len(buf) < headSize+1 {
		return nil, discover.NodeID{}, nil, errPacketTooSmall
	}
	hash, sig, sigdata := buf[:macSize], buf[macSize:headSize], buf[headSize:]
	shouldhash := crypto.Keccak256(buf[macSize:])
	if !bytes.Equal(hash, shouldhash) {
		return nil, discover.NodeID{}, nil, errBadHash
	}
	fromID, err := recoverNodeID(crypto.Keccak256(buf[headSize:]), sig)
	if err != nil {
		return nil, discover.NodeID{}, hash, err
	}
	var req packet
	switch ptype := sigdata[0]; ptype {
	case pingPacket:
		req = new(ping)
	case pongPacket:
		req = new(pong)
	case findnodePacket:
		req = new(findnode)
	case neighborsPacket:
		req = new(neighbors)
	default:
		return nil, fromID, hash, fmt.Errorf("unknown type: %d", ptype)
	}
	s := rlp.NewStream(bytes.NewReader(sigdata[1:]), 0)
	err = s.Decode(req)
	return req, fromID, hash, err
}

// recoverNodeID computes the public key used to sign the
// given hash from the signature.
func recoverNodeID(hash, sig []byte) (id discover.NodeID, err error) {
	pubkey, err := secp256k1.RecoverPubkey(hash, sig)
	if err != nil {
		return id, err
	}
	if len(pubkey)-1 != len(id) {
		return id, fmt.Errorf("recovered pubkey has %d bits, want %d bits", len(pubkey)*8, (len(id)+1)*8)
	}
	for i := range id {
		id[i] = pubkey[i+1]
	}
	return id, nil
}
