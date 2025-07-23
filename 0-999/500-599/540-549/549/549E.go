package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	type point struct{ x, y int }
	aPts := make([]point, n)
	bPts := make([]point, m)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &aPts[i].x, &aPts[i].y)
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &bPts[i].x, &bPts[i].y)
	}
	// Prepare data for perceptron: features [x,y,z], label +1 for bPts, -1 for aPts
	total := n + m
	type sample struct {
		f     [3]float64
		label float64
	}
	data := make([]sample, 0, total)
	// find normalization factors
	maxC := 1.0
	maxZ := 1.0
	for _, p := range aPts {
		maxC = math.Max(maxC, math.Abs(float64(p.x)))
		maxC = math.Max(maxC, math.Abs(float64(p.y)))
		z := float64(p.x*p.x + p.y*p.y)
		maxZ = math.Max(maxZ, z)
	}
	for _, p := range bPts {
		maxC = math.Max(maxC, math.Abs(float64(p.x)))
		maxC = math.Max(maxC, math.Abs(float64(p.y)))
		z := float64(p.x*p.x + p.y*p.y)
		maxZ = math.Max(maxZ, z)
	}
	// build samples
	for _, p := range aPts {
		x := float64(p.x) / maxC
		y := float64(p.y) / maxC
		z := (float64(p.x*p.x + p.y*p.y)) / maxZ
		data = append(data, sample{f: [3]float64{x, y, z}, label: -1})
	}
	for _, p := range bPts {
		x := float64(p.x) / maxC
		y := float64(p.y) / maxC
		z := (float64(p.x*p.x + p.y*p.y)) / maxZ
		data = append(data, sample{f: [3]float64{x, y, z}, label: +1})
	}
	// shuffle data for perceptron
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })
	// perceptron
	var w [3]float64
	var d float64
	// maximum number of epochs for perceptron
	const maxEpochs = 200
	for epoch := 0; epoch < maxEpochs; epoch++ {
		changed := false
		for _, s := range data {
			// dot product
			dp := w[0]*s.f[0] + w[1]*s.f[1] + w[2]*s.f[2] + d
			if s.label*dp <= 0 {
				// update
				w[0] += s.label * s.f[0]
				w[1] += s.label * s.f[1]
				w[2] += s.label * s.f[2]
				d += s.label
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	// check separation
	ok := true
	const eps = 1e-9
	for _, s := range data {
		dp := w[0]*s.f[0] + w[1]*s.f[1] + w[2]*s.f[2] + d
		if s.label*dp <= eps {
			ok = false
			break
		}
	}
	if ok {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
