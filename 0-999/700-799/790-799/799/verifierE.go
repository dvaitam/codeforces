package main

import (
	"bytes"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testE struct {
	n     int
	m     int
	k     int
	cost  []int
	likeM []bool
	likeA []bool
}

func solveE(tc testE) int {
	best := math.MaxInt32
	n := tc.n
	for mask := 0; mask < (1 << n); mask++ {
		if bits.OnesCount(uint(mask)) != tc.m {
			continue
		}
		cntM, cntA, cost := 0, 0, 0
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				cost += tc.cost[i]
				if tc.likeM[i] {
					cntM++
				}
				if tc.likeA[i] {
					cntA++
				}
			}
		}
		if cntM >= tc.k && cntA >= tc.k {
			if cost < best {
				best = cost
			}
		}
	}
	if best == math.MaxInt32 {
		return -1
	}
	return best
}

func genE() (string, int) {
	n := rand.Intn(7) + 1
	m := rand.Intn(n) + 1
	k := rand.Intn(m + 1)
	cost := make([]int, n)
	likeM := make([]bool, n)
	likeA := make([]bool, n)
	for i := 0; i < n; i++ {
		cost[i] = rand.Intn(20) + 1
		likeM[i] = rand.Intn(2) == 0
		likeA[i] = rand.Intn(2) == 0
	}
	tc := testE{n, m, k, cost, likeM, likeA}
	ans := solveE(tc)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i, v := range cost {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	var idx []int
	for i, v := range likeM {
		if v {
			idx = append(idx, i+1)
		}
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(idx)))
	for i, v := range idx {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	idx = idx[:0]
	for i, v := range likeA {
		if v {
			idx = append(idx, i+1)
		}
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(idx)))
	for i, v := range idx {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
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
		fmt.Println("Usage: go run verifierE.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genE()
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
