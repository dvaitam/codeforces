package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
)

func solve(n, k, e, ix, iy int, xs, ys []int) float64 {
	prob := make([]float64, n)
	prev := make([]float64, n+2)
	cur := make([]float64, n+2)
	chk := func(x float64) bool {
		for i := 0; i < n; i++ {
			dx := float64(xs[i] - ix)
			dy := float64(ys[i] - iy)
			d := math.Hypot(dx, dy)
			switch {
			case d < x:
				prob[i] = 1
			case d > x*1000:
				prob[i] = 0
			default:
				prob[i] = math.Exp(1 - (d*d)/(x*x))
			}
		}
		for i := range prev {
			prev[i] = 0
			cur[i] = 0
		}
		prev[0] = 1
		for i := 0; i < n; i++ {
			p := prob[i]
			for j := 0; j <= i+1; j++ {
				cur[j] = 0
			}
			for j := 0; j <= i; j++ {
				cur[j] += prev[j] * (1 - p)
				cur[j+1] += prev[j] * p
			}
			prev, cur = cur, prev
		}
		limit := float64(e) / 1000.0
		sum := 0.0
		for j := 0; j < k && j < len(prev); j++ {
			sum += prev[j]
		}
		return sum <= limit
	}
	lo, hi := 0.0, 10000.0
	eps := 1e-7
	for hi-lo > eps {
		mid := (lo + hi) * 0.5
		if chk(mid) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return lo
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(44))
	const cases = 110
	for t := 0; t < cases; t++ {
		n := r.Intn(8) + 1
		k := r.Intn(n) + 1
		e := r.Intn(999) + 1
		ix := r.Intn(21) - 10
		iy := r.Intn(21) - 10
		xs := make([]int, n)
		ys := make([]int, n)
		for i := 0; i < n; i++ {
			xs[i] = r.Intn(21) - 10
			ys[i] = r.Intn(21) - 10
		}
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d %d %d\n", n, k, e, ix, iy)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&input, "%d %d\n", xs[i], ys[i])
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		var got float64
		if _, err := fmt.Fscan(&out, &got); err != nil {
			fmt.Printf("Test %d parse error: %v\n", t+1, err)
			os.Exit(1)
		}
		want := solve(n, k, e, ix, iy, xs, ys)
		if math.Abs(got-want) > 1e-4 {
			fmt.Printf("Test %d failed: expected %.5f got %.5f\n", t+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
