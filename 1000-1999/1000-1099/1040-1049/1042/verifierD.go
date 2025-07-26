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

type testCaseD struct {
	n   int
	t   int64
	arr []int64
}

func solveD(tc testCaseD) int64 {
	var cnt int64
	for l := 0; l < tc.n; l++ {
		sum := int64(0)
		for r := l; r < tc.n; r++ {
			sum += tc.arr[r]
			if sum < tc.t {
				cnt++
			}
		}
	}
	return cnt
}

func run(bin, input string) (string, error) {
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

func generateTests() []testCaseD {
	rng := rand.New(rand.NewSource(45))
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rng.Intn(20) + 1
		t := int64(rng.Intn(201) - 100)
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = int64(rng.Intn(101) - 50)
		}
		tests[i] = testCaseD{n, t, arr}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for idx, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.t)
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		expected := solveD(tc)
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", idx+1, err)
			return
		}
		var got int64
		fmt.Sscan(out, &got)
		if got != expected {
			fmt.Printf("test %d failed:\ninput:%sexpected %d got %s\n", idx+1, sb.String(), expected, out)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
