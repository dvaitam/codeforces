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
	n, k int
	arr  []int
}

func solveB(tc testB) string {
	sorted := make([]int, tc.n)
	copy(sorted, tc.arr)
	sort.Ints(sorted)
	pos := make(map[int]int, tc.n)
	for i, v := range tc.arr {
		pos[v] = i
	}
	segments := 1
	for i := 0; i < tc.n-1; i++ {
		if pos[sorted[i]]+1 != pos[sorted[i+1]] {
			segments++
		}
	}
	if segments <= tc.k {
		return "YES"
	}
	return "NO"
}

func generateB(rng *rand.Rand) testB {
	n := rng.Intn(8) + 2
	k := rng.Intn(n) + 1
	perm := rng.Perm(n * 3)
	arr := make([]int, n)
	for i := range arr {
		arr[i] = perm[i]
	}
	return testB{n, k, arr}
}

func runCase(bin string, tc testB) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.k)
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testB, 0, 102)
	// add edge cases
	cases = append(cases, testB{5, 4, []int{6, 3, 4, 2, 1}})
	cases = append(cases, testB{3, 1, []int{3, 2, 1}})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateB(rng))
	}
	for i, tc := range cases {
		expect := solveB(tc)
		out, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
