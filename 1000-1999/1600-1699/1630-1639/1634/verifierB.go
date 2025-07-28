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

type testCase struct {
	n   int
	x   int64
	y   int64
	arr []int64
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(43))
	cases := make([]testCase, 100)
	for i := range cases {
		n := rng.Intn(8) + 1
		x := rng.Int63n(20)
		y := rng.Int63n(40)
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = rng.Int63n(20)
		}
		cases[i] = testCase{n: n, x: x, y: y, arr: arr}
	}
	return cases
}

func expected(tc testCase) string {
	parity := int64(0)
	for _, v := range tc.arr {
		parity ^= v & 1
	}
	alice := (tc.x & 1) ^ parity
	if alice == (tc.y & 1) {
		return "Alice"
	}
	return "Bob"
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		var in strings.Builder
		fmt.Fprintf(&in, "1\n%d %d %d\n", tc.n, tc.x, tc.y)
		for j, v := range tc.arr {
			if j > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", v)
		}
		in.WriteByte('\n')
		exp := expected(tc)
		out, err := run(bin, in.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, exp, out, in.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
