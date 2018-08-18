package nphash

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {

	pointCount := 8
	nph := New(pointCount)

	for i := 0; i < 100000; i += 1 {
		h := fmt.Sprintf("%x", i)

		hash := nph.HashBytes([]byte("hello worldX" + h))
		if hash[0] == 0 {
			fmt.Println(hash)
			fmt.Println(i)
			//break
		}

		//fmt.Println(hash)
	}

	//IntToHex(int64(nonce))

}
