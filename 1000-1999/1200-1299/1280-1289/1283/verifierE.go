package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n  int
	xs []int
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := range tests {
		n := rnd.Intn(20) + 1
		xs := make([]int, n)
		for j := 0; j < n; j++ {
			xs[j] = rnd.Intn(101) - 50
		}
		tests[i] = testCase{n, xs}
	}
	return tests
}

func expected(tc testCase) (int, int) {
	xs := append([]int(nil), tc.xs...)
	sort.Ints(xs)
	minAns := 0
	for i := 0; i < len(xs); {
		minAns++
		limit := xs[i] + 2
		for i < len(xs) && xs[i] <= limit {
			i++
		}
	}
	used := make(map[int]bool)
	for _, x := range xs {
		if !used[x-1] {
			used[x-1] = true
		} else if !used[x] {
			used[x] = true
		} else if !used[x+1] {
			used[x+1] = true
		}
	}
	maxAns := len(used)
	return minAns, maxAns
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), fmt.Errorf("exec error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.n)
		for i, x := range tc.xs {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", x)
		}
		input += "\n"
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		var gotMin, gotMax int
		if _, err := fmt.Sscan(out, &gotMin, &gotMax); err != nil {
			fmt.Printf("test %d: parse error\n", idx+1)
			os.Exit(1)
		}
		expMin, expMax := expected(tc)
		if gotMin != expMin || gotMax != expMax {
			fmt.Printf("test %d: expected %d %d got %d %d\n", idx+1, expMin, expMax, gotMin, gotMax)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
