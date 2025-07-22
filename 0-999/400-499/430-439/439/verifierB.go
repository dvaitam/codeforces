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

type testB struct {
	n int
	x int
	c []int
}

func solveB(tc testB) int64 {
	vals := make([]int, len(tc.c))
	copy(vals, tc.c)
	sort.Ints(vals)
	var total int64
	for i, v := range vals {
		t := tc.x - i
		if t < 1 {
			t = 1
		}
		total += int64(t) * int64(v)
	}
	return total
}

func genB(rng *rand.Rand) testB {
	n := rng.Intn(50) + 1
	x := rng.Intn(100) + 1
	c := make([]int, n)
	for i := range c {
		c[i] = rng.Intn(100) + 1
	}
	return testB{n: n, x: x, c: c}
}

func runCase(bin string, tc testB) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.x)
	for i, v := range tc.c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testB, 0, 102)
	cases = append(cases, testB{n: 1, x: 3, c: []int{5}})
	cases = append(cases, testB{n: 2, x: 1, c: []int{3, 4}})
	for i := 0; i < 100; i++ {
		cases = append(cases, genB(rng))
	}
	for i, tc := range cases {
		expect := solveB(tc)
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
