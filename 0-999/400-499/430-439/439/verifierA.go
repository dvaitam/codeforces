package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testA struct {
	n int
	d int
	t []int
}

func solveA(tc testA) int {
	sum := 0
	for _, v := range tc.t {
		sum += v
	}
	if sum+10*(tc.n-1) > tc.d {
		return -1
	}
	return (tc.d - sum) / 5
}

func generateA(rng *rand.Rand) testA {
	n := rng.Intn(100) + 1
	d := rng.Intn(10000) + 1
	t := make([]int, n)
	for i := 0; i < n; i++ {
		t[i] = rng.Intn(100) + 1
	}
	return testA{n: n, d: d, t: t}
}

func runCase(bin string, tc testA) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.d)
	for i, v := range tc.t {
		if i > 0 {
			fmt.Fprint(&sb, " ")
		}
		fmt.Fprint(&sb, v)
	}
	fmt.Fprint(&sb, "\n")
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testA, 0, 102)
	// add some edge cases
	cases = append(cases, testA{n: 1, d: 5, t: []int{5}})
	cases = append(cases, testA{n: 2, d: 30, t: []int{2, 2}})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateA(rng))
	}
	for i, tc := range cases {
		expect := solveA(tc)
		got, err := runCase(exe, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if fmt.Sprint(expect) != got {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
