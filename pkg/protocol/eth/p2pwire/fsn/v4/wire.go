package v4wire

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover/v4wire"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zhiqiangxu/devp2p"
)

const (
	pingPacket = iota + 40 // zero is 'reserved'
	pongPacket
	findnodePacket
	neighborsPacket
)

const (
	macSize  = 32
	sigSize  = crypto.SignatureLength
	headSize = macSize + sigSize // space of packet frame data
)

type Wire struct {
}

func (w *Wire) Decode(input []byte) (v4wire.Packet, v4wire.Pubkey, []byte, error) {
	if len(input) < headSize+1 {
		return nil, v4wire.Pubkey{}, nil, v4wire.ErrPacketTooSmall
	}
	hash, sig, sigdata := input[:macSize], input[macSize:headSize], input[headSize:]
	shouldhash := crypto.Keccak256(input[macSize:])
	if !bytes.Equal(hash, shouldhash) {
		return nil, v4wire.Pubkey{}, nil, v4wire.ErrBadHash
	}
	fromKey, err := recoverNodeKey(crypto.Keccak256(input[headSize:]), sig)
	if err != nil {
		return nil, fromKey, hash, err
	}

	var req v4wire.Packet
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
		return nil, fromKey, hash, fmt.Errorf("unknown type: %d", ptype)
	}
	s := rlp.NewStream(bytes.NewReader(sigdata[1:]), 0)
	err = s.Decode(req)
	return req, fromKey, hash, err
}

const (
	expiration     = 20 * time.Second
	bondExpiration = 24 * time.Hour
)

func (w *Wire) WrapPacket(t *devp2p.UDPv4, p v4wire.Packet) *devp2p.PacketHandlerV4 {

	var h devp2p.PacketHandlerV4
	h.Packet = p

	switch p.(type) {
	case *ping:
		h.Preverify = verifyPing
		h.Handle = func(h *devp2p.PacketHandlerV4, from *net.UDPAddr, fromID enode.ID, mac []byte) {
			req := h.Packet.(*ping)

			// Reply.
			t.Send(from, fromID, &pong{
				To:         makeEndpoint(from, req.From.TCP),
				ReplyTok:   mac,
				Expiration: uint64(time.Now().Add(expiration).Unix()),
			})

			// Ping back if our last pong on file is too far in the past.
			n := devp2p.WrapNode(enode.NewV4(h.SenderKey, from.IP, int(req.From.TCP), from.Port))
			if time.Since(t.DB.LastPongReceived(n.ID(), from.IP)) > bondExpiration {
				t.SendPing(fromID, from, func() {
					t.Tab.AddVerifiedNode(n)
				})
			} else {
				t.Tab.AddVerifiedNode(n)
			}

			// Update node database and endpoint predictor.
			t.DB.UpdateLastPingReceived(n.ID(), from.IP, time.Now())
			t.LocalNode.UDPEndpointStatement(from, &net.UDPAddr{IP: req.To.IP, Port: int(req.To.UDP)})
		}
	case *pong:
		h.Preverify = func(p *devp2p.PacketHandlerV4, from *net.UDPAddr, fromID enode.ID, fromKey v4wire.Pubkey) error {
			req := h.Packet.(*pong)

			if v4wire.Expired(req.Expiration) {
				return errExpired
			}
			if !t.HandleReply(fromID, from.IP, req) {
				return errUnsolicitedReply
			}
			t.LocalNode.UDPEndpointStatement(from, &net.UDPAddr{IP: req.To.IP, Port: int(req.To.UDP)})
			t.DB.UpdateLastPongReceived(fromID, from.IP, time.Now())
			return nil
		}
	case *findnode:
		h.Preverify = func(p *devp2p.PacketHandlerV4, from *net.UDPAddr, fromID enode.ID, fromKey v4wire.Pubkey) error {
			req := h.Packet.(*findnode)

			if v4wire.Expired(req.Expiration) {
				return errExpired
			}
			if !t.CheckBond(fromID, from.IP) {
				// No endpoint proof pong exists, we don't process the packet. This prevents an
				// attack vector where the discovery protocol could be used to amplify traffic in a
				// DDOS attack. A malicious actor would send a findnode request with the IP address
				// and UDP port of the target as the source address. The recipient of the findnode
				// packet would then send a neighbors packet (which is a much bigger packet than
				// findnode) to the victim.
				return errUnknownNode
			}
			return nil
		}
		h.Handle = func(h *devp2p.PacketHandlerV4, from *net.UDPAddr, fromID enode.ID, mac []byte) {
			req := h.Packet.(*findnode)

			// Determine closest nodes.
			target := enode.ID(crypto.Keccak256Hash(req.Target[:]))
			closest := t.Tab.FindnodeByID(target, devp2p.BucketSize, true)

			// Send neighbors in chunks with at most maxNeighbors per packet
			// to stay below the packet size limit.
			p := neighbors{Expiration: uint64(time.Now().Add(expiration).Unix())}
			var sent bool
			for _, n := range closest {
				if netutil.CheckRelayIP(from.IP, n.IP()) == nil {
					p.Nodes = append(p.Nodes, nodeToRPC(n))
				}
				if len(p.Nodes) == v4wire.MaxNeighbors {
					t.Send(from, fromID, &p)
					p.Nodes = p.Nodes[:0]
					sent = true
				}
			}
			if len(p.Nodes) > 0 || !sent {
				t.Send(from, fromID, &p)
			}
		}
	case *neighbors:
		h.Preverify = func(p *devp2p.PacketHandlerV4, from *net.UDPAddr, fromID enode.ID, fromKey v4wire.Pubkey) error {
			req := h.Packet.(*neighbors)

			if v4wire.Expired(req.Expiration) {
				return devp2p.ErrExpired
			}
			if !t.HandleReply(fromID, from.IP, h.Packet) {
				return devp2p.ErrUnsolicitedReply
			}
			return nil
		}
	default:
		panic("unknown packet")
	}

	return &h
}

func (w *Wire) HandlePendingNeighborsPacket(t *devp2p.UDPv4, neighborsPacket v4wire.Packet, toaddr *net.UDPAddr) (nreceived int, nodes []*enode.Node) {
	reply := neighborsPacket.(*neighbors)
	for _, rn := range reply.Nodes {
		nreceived++
		n, err := enodeFromRPC(t, toaddr, rn)
		if err != nil {
			t.Log.Trace("Invalid neighbor node received", "ip", rn.IP, "addr", toaddr, "err", err)
			continue
		}
		nodes = append(nodes, n)
	}
	return
}

func (w *Wire) HandlePendingPongPacket(pongPacket v4wire.Packet, hash []byte, callback func()) (bool, bool) {
	if p, ok := pongPacket.(*pong); ok {
		matched := bytes.Equal(p.ReplyTok, hash)
		if matched && callback != nil {
			callback()
		}
		return matched, matched
	}
	return false, false
}

func (w *Wire) HandlePendingENRResponsePacket(pongPacket v4wire.Packet, hash []byte) (bool, bool) {
	panic("this should not be called")
}

func (w *Wire) PingPacket(self *enode.Node, toaddr *net.UDPAddr) v4wire.Packet {
	return &ping{
		Version:    4,
		From:       enodeToRPC(self),
		To:         makeEndpoint(toaddr, 0),
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	}
}

func (w *Wire) ENRRequestPacket() v4wire.Packet {
	panic("this should not be called")
}

func (w *Wire) FindnodePacket(target v4wire.Pubkey) v4wire.Packet {
	return &findnode{Target: NodeID(target), Expiration: uint64(time.Now().Add(expiration).Unix())}
}

func (w *Wire) PongPacketType() byte {
	return pongPacket
}

func (w *Wire) NeighborsPacketType() byte {
	return neighborsPacket
}

func (w *Wire) ENRResponsePacketType() byte {
	panic("this should not be called")
}

func (w *Wire) Resolve(n *enode.Node) *enode.Node {
	// disable Resolve by return n directly
	return n
}

var (
	errExpired          = errors.New("expired")
	errUnsolicitedReply = errors.New("unsolicited reply")
	errUnknownNode      = errors.New("unknown node")
)

func enodeFromRPC(t *devp2p.UDPv4, sender *net.UDPAddr, rn rpcNode) (*enode.Node, error) {
	if rn.UDP <= 1024 {
		return nil, devp2p.ErrLowPort
	}
	if err := netutil.CheckRelayIP(sender.IP, rn.IP); err != nil {
		return nil, err
	}
	if t.NetRestrict != nil && !t.NetRestrict.Contains(rn.IP) {
		return nil, errors.New("not contained in netrestrict list")
	}
	key, err := v4wire.DecodePubkey(crypto.S256(), v4wire.Pubkey(rn.ID))
	if err != nil {
		return nil, err
	}
	n := enode.NewV4(key, rn.IP, int(rn.TCP), int(rn.UDP))
	err = n.ValidateComplete()
	return n, err
}

func nodeToRPC(n *devp2p.Node) rpcNode {
	var key ecdsa.PublicKey
	var ekey v4wire.Pubkey
	if err := n.Load((*enode.Secp256k1)(&key)); err == nil {
		ekey = v4wire.EncodePubkey(&key)
	}
	return rpcNode{ID: NodeID(ekey), IP: n.IP(), UDP: uint16(n.UDP()), TCP: uint16(n.TCP())}
}

func verifyPing(h *devp2p.PacketHandlerV4, from *net.UDPAddr, fromID enode.ID, fromKey v4wire.Pubkey) error {
	req := h.Packet.(*ping)

	senderKey, err := v4wire.DecodePubkey(crypto.S256(), fromKey)
	if err != nil {
		return err
	}
	if v4wire.Expired(req.Expiration) {
		return errExpired
	}
	h.SenderKey = senderKey
	return nil
}

// recoverNodeKey computes the public key used to sign the given hash from the signature.
func recoverNodeKey(hash, sig []byte) (key v4wire.Pubkey, err error) {
	pubkey, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		return key, err
	}
	copy(key[:], pubkey[1:])
	return key, nil
}

func makeEndpoint(addr *net.UDPAddr, tcpPort uint16) rpcEndpoint {
	ip := addr.IP.To4()
	if ip == nil {
		ip = addr.IP.To16()
	}
	return rpcEndpoint{IP: ip, UDP: uint16(addr.Port), TCP: tcpPort}
}

func enodeToRPC(n *enode.Node) rpcEndpoint {
	return rpcEndpoint{IP: n.IP(), UDP: uint16(n.UDP()), TCP: uint16(n.TCP())}
}
