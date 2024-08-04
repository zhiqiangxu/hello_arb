package v4wire

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/p2p/discover/v4wire"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestCompatibility(t *testing.T) {
	p1 := &v4wire.Findnode{
		Target:     hexPubkey("ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31387574077f301b421bc84df7266c44e9e6d569fc56be00812904767bf5ccd1fc7f"),
		Expiration: 1136239445,
		Rest:       []rlp.RawValue{{0x82, 0x99, 0x99}, {0x83, 0x99, 0x99, 0x99}},
	}
	p2 := &findnode{
		Target:     NodeID(hexPubkey("ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31387574077f301b421bc84df7266c44e9e6d569fc56be00812904767bf5ccd1fc7f")),
		Expiration: 1136239445,
		Rest:       []rlp.RawValue{{0x82, 0x99, 0x99}, {0x83, 0x99, 0x99, 0x99}},
	}
	p1Bytes, err := rlp.EncodeToBytes(p1)
	if err != nil {
		t.Fatal(err)
	}
	p2Bytes, err := rlp.EncodeToBytes(p2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p1Bytes, p2Bytes) {
		t.Fatal()
	}
}

func hexPubkey(h string) (ret v4wire.Pubkey) {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	if len(b) != len(ret) {
		panic("invalid length")
	}
	copy(ret[:], b)
	return ret
}
