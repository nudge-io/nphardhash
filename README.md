# nphardhash

ASIC resistant, 256Bit, cryptographic hashing algorithm.
Works by requiring the client to solve an np-hard problem before they are able to determine the hash value
(https://en.wikipedia.org/wiki/NP-hardness). In this case we use the traveling salesman problem. This allows for


Usage:
```
import "github.com/nudgeplatform/nphardhash"

func main() {
    //number of points to generate from source hash bytes
    pointCount:=6

    //reusable hash generator
    nph := nphardhash.New(pointCount)

    //resulting hash bytes - 32 bytes (256 bits)
    hashBytes := nph.HashBytes( []byte("hello world") )

}
```
