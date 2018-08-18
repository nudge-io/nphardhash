# NP Hard Hash

ASIC resistant, 256Bit, cryptographic hashing algorithm.
Works by requiring the client to solve an np-hard problem (https://en.wikipedia.org/wiki/NP-hardness)
before they are able to determine the hash value. In this case we use the traveling salesman problem.

The algorithm generates a list of 2D points derived from the input data bytes.
The client is required to find all unique permutations of these points and
sort them in ascending order by travel distance.
This ordering and the point bytes are then used generate a fingerprint of the original data
which is then hashed to generate the resulting digest.

This algorithm is both memory and cpu intensive compared to other algorithms and the difficulty can be configured to
scale exponentially by increasing the pointCount that needs to be calculated. This is where the problem for ASICs lies.

In the context of a blockchain Proof-of-Work, we increase the calculate point count relative to the network hash rate.
So as ASICs come online eclipsing the hash rate of cpu devices, the memory requirements to solve a hash will jump.
This translates to only those devices with significant memory capacity are able to keep up
and ASICs are quickly dropped off the network.


Usage:
```
import "github.com/nudgeplatform/nphardhash"

func main() {
    //number of points to generate from source hash bytes
    pointCount:=6

    //reusable hash generator
    nph := nphardhash.New(pointCount)

    //resulting hash digest - 32 bytes (256 bits)
    digest := nph.HashBytes( []byte("hello world") )
}
```
