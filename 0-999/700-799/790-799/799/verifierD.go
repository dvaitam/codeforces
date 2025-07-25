package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type pair struct{ w, h int64 }

func solve(targetA, targetB, h, w int64, factors []int64) int {
	if h >= targetA && w >= targetB {
		return 0
	}
	sort.Slice(factors, func(i, j int) bool { return factors[i] > factors[j] })
	limit := 10
	if len(factors) < limit {
		limit = len(factors)
	}
	factors = factors[:limit]
	state := map[int64]int64{h: w}
	for i, f := range factors {
		next := make(map[int64]int64, len(state)*2)
		for width, height := range state {
			if height > next[width] {
				next[width] = height
			}
			nw := width * f
			if nw > targetA {
				nw = targetA
			}
			if height > next[nw] {
				next[nw] = height
			}
			nh := height * f
			if nh > targetB {
				nh = targetB
			}
			if nh > next[width] {
				next[width] = nh
			}
		}
		arr := make([]pair, 0, len(next))
		for wv, hv := range next {
			arr = append(arr, pair{wv, hv})
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].w < arr[j].w })
		state = make(map[int64]int64)
		maxH := int64(0)
		for _, p := range arr {
			if p.h > maxH {
				maxH = p.h
				state[p.w] = p.h
			}
		}
		for wv, hv := range state {
			if wv >= targetA && hv >= targetB {
				return i + 1
			}
		}
	}
	return -1
}

func solveD(a, b, h, w int64, factors []int64) int {
	r1 := solve(a, b, h, w, factors)
	r2 := solve(b, a, h, w, factors)
	if r1 == -1 {
		return r2
	}
	if r2 == -1 {
		return r1
	}
	if r1 < r2 {
		return r1
	}
	return r2
}

func genD() (string, int) {
	a := int64(rand.Intn(10) + 1)
	b := int64(rand.Intn(10) + 1)
	h := int64(rand.Intn(10) + 1)
	w := int64(rand.Intn(10) + 1)
	n := rand.Intn(5) + 1
	factors := make([]int64, n)
	for i := 0; i < n; i++ {
		factors[i] = int64(rand.Intn(5) + 2)
	}
	ans := solveD(a, b, h, w, factors)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", a, b, h, w, n))
	for i, f := range factors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", f))
	}
	sb.WriteByte('\n')
	return sb.String(), ans
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genD()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s", i+1, err, in, got)
			return
		}
		got = strings.TrimSpace(got)
		var val int
		if _, err := fmt.Sscan(got, &val); err != nil || val != exp {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%d\nGot:\n%s", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
