package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solve(n, m int, puzzles []int) int {
	sort.Ints(puzzles)
	ans := puzzles[n-1] - puzzles[0]
	for i := 1; i+n-1 < m; i++ {
		diff := puzzles[i+n-1] - puzzles[i]
		if diff < ans {
			ans = diff
		}
	}
	return ans
}

func generateTest(rng *rand.Rand) (string, int) {
	m := rng.Intn(49) + 2  // 2..50
	n := rng.Intn(m-1) + 2 // 2..m
	puzzles := make([]int, m)
	for i := range puzzles {
		puzzles[i] = rng.Intn(997) + 4 // 4..1000
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, v := range puzzles {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	expected := solve(n, m, puzzles)
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	const tests = 100
	for t := 1; t <= tests; t++ {
		inp, want := generateTest(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(inp)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n%s", t, err, out.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: failed to parse output: %v\nOutput: %s\n", t, err, gotStr)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %d\nGot: %d\n", t, inp, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
