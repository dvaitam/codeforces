package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseB struct {
	p []int
}

func solveB(p []int) (int, int) {
	n := len(p)
	bestVal := int(^uint(0) >> 1)
	bestIdx := 0
	for k := 0; k < n; k++ {
		cur := 0
		for i := 0; i < n; i++ {
			v := p[(i+n-k)%n]
			d := v - (i + 1)
			if d < 0 {
				d = -d
			}
			cur += d
		}
		if cur < bestVal {
			bestVal = cur
			bestIdx = k
		}
	}
	return bestVal, bestIdx
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCaseB {
	rand.Seed(42)
	tests := make([]testCaseB, 100)
	for i := range tests {
		n := rand.Intn(7) + 2
		perm := rand.Perm(n)
		for j := range perm {
			perm[j]++
		}
		tests[i] = testCaseB{perm}
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.p)))
		for j, v := range tc.p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		val, idx := solveB(tc.p)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		var gotVal, gotIdx int
		fmt.Sscan(output, &gotVal, &gotIdx)
		if gotVal != val || gotIdx != idx {
			fmt.Printf("test %d failed:\ninput:%sexpected %d %d got %s\n", i+1, input, val, idx, output)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
