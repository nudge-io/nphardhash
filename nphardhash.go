package nphardhash

//nudge.io: hashing algorithm used to generate asic resistant hashes.
//uses bruteforce approach to generate path permutations. This could easily be optimized.

//nphardhash
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math"
	"sort"
)

//helper funcs
//func Factorial(n int) int {
//	if n == 0 {
//		return 1
//	}
//	return n * Factorial(n-1)
//}

func internalHash(data []byte) [32]byte {
	var buf [32]byte
	sha3.ShakeSum256(buf[:], data)
	return buf
	//return sha256.Sum256(data)
}

func ArgSortConcatByteArrs(src_arr [][]uint8, arg_indexes []int) []uint8 {

	//preArr
	preArr := make([][]uint8, len(arg_indexes))
	for i, idx := range arg_indexes {
		preArr[i] = src_arr[idx]
	}

	//postArr
	postArr := bytes.Join(preArr, nil)

	//
	return postArr
}

func scaledfloat64FromBytes(bytes []uint8) float64 {
	//
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)

	//normalize
	val := math.Log(math.Abs(float))

	//return
	return val
}

func calcDistances(points *[]point) {
	for i, v := range *points {
		(*points)[i].distances = make([]float64, len(*points))

		for k, v2 := range *points {
			dx, dy := v.x-v2.x, v.y-v2.y
			(*points)[i].distances[k] = math.Sqrt(dx*dx + dy*dy)
		}
	}
}

func genPermutations(src *[]int, c chan []int) {
	p, err := NewPerm(*src, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, err := p.Next(); err == nil; i, err = p.Next() {
		perm := i.([]int)
		//fmt.Printf("%3d permutation: %v left %d\n",p.Index()-1,perm,p.Left())
		c <- perm
	}

	close(c)
}

func calculateScore(path *[]int, points *[]point) float64 {
	score := 0.0
	for idx, k := range *path {
		var j int
		if idx+1 < len(*path) {
			j = (*path)[idx+1]
		} else {
			j = (*path)[0]
		}
		score += (*points)[k].distances[j]
	}
	return score
}

//
type point struct {
	x         float64
	y         float64
	bytes     []uint8
	distances []float64
}

func NewPointFromBytes(bytes []uint8) *point {
	if len(bytes) != 16 {
		return nil
	}

	//
	x := scaledfloat64FromBytes(bytes[:8])
	y := scaledfloat64FromBytes(bytes[8:])

	//
	p := &point{x, y, bytes, nil}
	return p
}

//
type npHash struct {
	pointCount int
}

func New(pointCount int) *npHash {
	//return result
	return &npHash{pointCount: pointCount}
}

func (n *npHash) generatePoints(bytes []uint8) []point {

	//pointCount
	pointCount := n.pointCount
	remainder := pointCount % 2
	if remainder > 0 {
		pointCount += 1
	}

	//hash once
	byteArr := internalHash(bytes)

	//iterate for points
	points := make([]point, pointCount)
	for x := 0; x < pointCount; x += 2 {
		byteArr = internalHash(byteArr[:])

		//new points
		pointA := NewPointFromBytes(byteArr[:16])
		pointB := NewPointFromBytes(byteArr[16:])
		points[x] = *pointA
		points[x+1] = *pointB
	}

	//trim to actual pointCount if needed
	realPointCount := n.pointCount
	if pointCount > realPointCount {
		points = points[:realPointCount]
	}

	//
	return points
}

//
type pathScore struct {
	index int
	score float64
	bytes []uint8
}

type byPathScore []*pathScore

func (s byPathScore) Len() int {
	return len(s)
}

func (s byPathScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPathScore) Less(i, j int) bool {

	pI := s[i]
	pJ := s[j]

	//compare scores
	if pI.score < pJ.score {
		return true
	} else if pI.score > pJ.score {
		return false
	}

	//compare equal items by bytes
	pR := bytes.Compare(pI.bytes, pJ.bytes)
	if pR < 0 {
		return true
	} else {
		return false
	}

}

func (n *npHash) HashBytes(hashBytes []uint8) [32]byte {
	//fmt.Println("hashBytes")
	//fmt.Println(hashBytes)

	//generate points
	points := n.generatePoints(hashBytes)
	pointCount := len(points)
	//fmt.Println(points)

	//calculate distances to other points
	calcDistances(&points)

	//path permuations/////////
	paths := make([]int, pointCount)
	for x := 0; x < pointCount; x += 1 {
		paths[x] = x
	}
	ch := make(chan []int)
	go genPermutations(&paths, ch)

	//pointBytes
	pointBytes := make([][]uint8, pointCount)
	for i, pnt := range points {
		pointBytes[i] = pnt.bytes
	}

	//pathScores
	pathScores := make([]*pathScore, 0)

	idx := 0
	score := -1.0
	for i := range ch {
		//fmt.Println(i)
		//get score
		score = calculateScore(&i, &points)

		//get bytes
		pathByteArr := ArgSortConcatByteArrs(pointBytes, i)

		//pathScore
		pScore := pathScore{idx, score, pathByteArr}
		pathScores = append(pathScores, &pScore)

		//
		idx++
	}
	sort.Sort(byPathScore(pathScores))

	//combine []uint8
	pathByteArrs := make([][]uint8, len(pathScores))
	for i, pathScore := range pathScores {
		pathByteArrs[i] = pathScore.bytes
	}

	//
	combineBytes := bytes.Join(pathByteArrs, nil)
	resultHash := internalHash(combineBytes)

	return resultHash
}
