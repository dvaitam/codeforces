package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

type Point struct {
	X int64
	Y int64
}

type Pair struct {
	i, j  int
	dist2 int64
}

func sampleIndices(n, k int) []int {
	if k > n {
		k = n
	}
	perm := rand.Perm(n)
	return perm[:k]
}

func buildPairs(pts []Point, idx []int) []Pair {
	var pairs []Pair
	for a := 0; a < len(idx); a++ {
		for b := a + 1; b < len(idx); b++ {
			i := idx[a]
			j := idx[b]
			dx := pts[j].X - pts[i].X
			dy := pts[j].Y - pts[i].Y
			d2 := dx*dx + dy*dy
			pairs = append(pairs, Pair{i: i, j: j, dist2: d2})
		}
	}
	return pairs
}

func attemptTransform(p1 Pair, p2 Pair, arr1 []Point, arr2 []Point, lookup map[[2]int64]int, target int) ([][2]int, bool) {
	A := arr1[p1.i]
	B := arr1[p1.j]
	C := arr2[p2.i]
	D := arr2[p2.j]

	dx1 := float64(B.X - A.X)
	dy1 := float64(B.Y - A.Y)
	dx2 := float64(D.X - C.X)
	dy2 := float64(D.Y - C.Y)

	len1 := math.Hypot(dx1, dy1)
	len2 := math.Hypot(dx2, dy2)
	if len1 == 0 || len2 == 0 {
		return nil, false
	}
	if math.Abs(len1-len2) > 1e-6 {
		return nil, false
	}

	cosv := (dx1*dx2 + dy1*dy2) / (len1 * len2)
	sinv := (dx1*dy2 - dy1*dx2) / (len1 * len2)
	tx := float64(C.X) - (cosv*float64(A.X) - sinv*float64(A.Y))
	ty := float64(C.Y) - (sinv*float64(A.X) + cosv*float64(A.Y))

	used := make([]bool, len(arr2))
	mapping := make([][2]int, 0, target)
	for i, p := range arr1 {
		x := cosv*float64(p.X) - sinv*float64(p.Y) + tx
		y := sinv*float64(p.X) + cosv*float64(p.Y) + ty
		xi := int64(math.Round(x))
		yi := int64(math.Round(y))
		key := [2]int64{xi, yi}
		if j, ok := lookup[key]; ok && !used[j] {
			used[j] = true
			mapping = append(mapping, [2]int{i + 1, j + 1})
		}
	}
	return mapping, len(mapping) >= target
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	rand.Seed(time.Now().UnixNano())

	var N, N1 int
	if _, err := fmt.Fscan(reader, &N); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &N1); err != nil {
		return
	}
	arr1 := make([]Point, N1)
	for i := 0; i < N1; i++ {
		var x, y float64
		fmt.Fscan(reader, &x, &y)
		arr1[i] = Point{X: int64(math.Round(x * 100)), Y: int64(math.Round(y * 100))}
	}
	var N2 int
	fmt.Fscan(reader, &N2)
	arr2 := make([]Point, N2)
	for i := 0; i < N2; i++ {
		var x, y float64
		fmt.Fscan(reader, &x, &y)
		arr2[i] = Point{X: int64(math.Round(x * 100)), Y: int64(math.Round(y * 100))}
	}

	lookup := make(map[[2]int64]int, N2)
	for i, p := range arr2 {
		key := [2]int64{p.X, p.Y}
		if _, ok := lookup[key]; !ok {
			lookup[key] = i
		}
	}

	sampleSize := 25
	idx1 := sampleIndices(N1, sampleSize)
	idx2 := sampleIndices(N2, sampleSize)
	pairs1 := buildPairs(arr1, idx1)
	pairs2 := buildPairs(arr2, idx2)

	byDist1 := make(map[int64][]Pair)
	byDist2 := make(map[int64][]Pair)
	for _, pr := range pairs1 {
		byDist1[pr.dist2] = append(byDist1[pr.dist2], pr)
	}
	for _, pr := range pairs2 {
		byDist2[pr.dist2] = append(byDist2[pr.dist2], pr)
	}

	var dists []int64
	for d := range byDist1 {
		if _, ok := byDist2[d]; ok {
			dists = append(dists, d)
		}
	}

	target := int(math.Ceil(0.9 * float64(N)))
	var result [][2]int
	found := false

	attempts := 300
	for a := 0; a < attempts && !found && len(dists) > 0; a++ {
		d := dists[rand.Intn(len(dists))]
		p1list := byDist1[d]
		p2list := byDist2[d]
		p1 := p1list[rand.Intn(len(p1list))]
		p2 := p2list[rand.Intn(len(p2list))]
		if mapping, ok := attemptTransform(p1, p2, arr1, arr2, lookup, target); ok {
			result = mapping
			found = true
		}
	}

	if !found {
		result = make([][2]int, N)
		for i := 0; i < N; i++ {
			result[i] = [2]int{1, 1}
		}
	}

	for len(result) < N {
		result = append(result, result[len(result)%len(result)])
	}

	for i := 0; i < N; i++ {
		fmt.Fprintln(writer, result[i][0], result[i][1])
	}
}
