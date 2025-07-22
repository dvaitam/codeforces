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

type testCaseA struct {
	m   int
	q   []int
	n   int
	arr []int
}

func expectedA(tc testCaseA) int64 {
	minQ := tc.q[0]
	for _, v := range tc.q {
		if v < minQ {
			minQ = v
		}
	}
	sort.Slice(tc.arr, func(i, j int) bool { return tc.arr[i] > tc.arr[j] })
	k := minQ + 2
	var res int64
	for i, v := range tc.arr {
		if i%k < minQ {
			res += int64(v)
		}
	}
	return res
}

func genCaseA(rng *rand.Rand) (string, string) {
	m := rng.Intn(5) + 1
	q := make([]int, m)
	for i := range q {
		q[i] = rng.Intn(5) + 1
	}
	n := rng.Intn(6) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(20) + 1
	}
	tc := testCaseA{m, q, n, arr}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", m)
	for i, v := range q {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d\n", expectedA(tc))
	return sb.String(), exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", strings.TrimSpace(expected), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseA(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
