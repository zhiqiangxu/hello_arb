package v4wire

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/rlp"
)

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
		Target     NodeID // doesn't need to be an actual public key
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
		ID  NodeID
	}

	rpcEndpoint struct {
		IP  net.IP // len 4 for IPv4 or 16 for IPv6
		UDP uint16 // for discovery protocol
		TCP uint16 // for RLPx protocol
	}
)

const nodeIDBits = 512

// NodeID is a unique identifier for each node.
// The node identifier is a marshaled elliptic curve public key.
type NodeID [nodeIDBits / 8]byte

// NodeID prints as a long hexadecimal number.
func (n NodeID) String() string {
	return fmt.Sprintf("%x", n[:])
}

// The Go syntax representation of a NodeID is a call to HexID.
func (n NodeID) GoString() string {
	return fmt.Sprintf("discover.HexID(\"%x\")", n[:])
}

// TerminalString returns a shortened hex string for terminal logging.
func (n NodeID) TerminalString() string {
	return hex.EncodeToString(n[:8])
}

func (req *ping) Name() string { return "PING/v4" }
func (req *ping) Kind() byte   { return pingPacket }

func (req *pong) Name() string { return "PONG/v4" }
func (req *pong) Kind() byte   { return pongPacket }

func (req *findnode) Name() string { return "FINDNODE/v4" }
func (req *findnode) Kind() byte   { return findnodePacket }

func (req *neighbors) Name() string { return "NEIGHBORS/v4" }
func (req *neighbors) Kind() byte   { return neighborsPacket }
