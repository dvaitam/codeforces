package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
)

func ceilSqrt(x int64) int {
	if x <= 0 {
		return 0
	}
	y := int(math.Sqrt(float64(x)))
	if int64(y)*int64(y) < x {
		y++
	}
	return y
}

func solve(n int, m int64) int64 {
	var irr int64
	for b := 1; b <= n; b++ {
		b2 := int64(b) * int64(b)
		var X int64
		if m < b2-1 {
			X = m
		} else {
			X = b2 - 1
		}
		if X <= 0 {
			continue
		}
		t := b2 - X
		sMin := ceilSqrt(t)
		sq := b - sMin
		if sq < 0 {
			sq = 0
		}
		irrC := X - int64(sq)
		irr += 2 * irrC
	}
	maxS := 2 * n
	mark := make([]byte, maxS+1)
	for u := 1; u <= n; u++ {
		vMax1 := int(m / int64(u))
		vMax2 := 2*n - u
		vMax := vMax1
		if vMax2 < vMax {
			vMax = vMax2
		}
		if vMax < u {
			continue
		}
		mark[u] = 1
		for v := u; v <= vMax; v += 2 {
			mark[v] = 1
		}
	}
	var intCnt int64
	for s := 1; s <= maxS; s++ {
		if mark[s] != 0 {
			intCnt++
		}
	}
	return irr + intCnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(45))
	const cases = 120
	for t := 0; t < cases; t++ {
		n := r.Intn(200) + 1
		m := int64(r.Intn(200) + 1)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, m)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Fscan(&out, &got); err != nil {
			fmt.Printf("Test %d parse error: %v\n", t+1, err)
			os.Exit(1)
		}
		want := solve(n, m)
		if got != want {
			fmt.Printf("Test %d failed: expected %d got %d\n", t+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
