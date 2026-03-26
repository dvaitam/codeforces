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

const MOD int64 = 998244353

// Brute-force solver: enumerate all permutations, simulate the tournament,
// check if all results are predictable.
// Key insight: when |x-y| > k, the higher-rated player always wins.
// So the current champion is always the max-rated player among those seen so far.
// For game i (i=1..n-1): champion = max of first i players, challenger = (i+1)-th.
// The game is predictable iff |champion_rating - challenger_rating| > k.
func solveBrute(n int, k int64, a []int64) int64 {
	if n == 1 {
		return 1
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i
	}
	count := int64(0)
	var generate func(int)
	generate = func(start int) {
		if start == n {
			// Simulate tournament: champion is max of first i, challenger is perm[i]
			predictable := true
			curMax := a[perm[0]]
			for i := 1; i < n; i++ {
				challenger := a[perm[i]]
				diff := curMax - challenger
				if diff < 0 {
					diff = -diff
				}
				if diff <= k {
					predictable = false
					break
				}
				if challenger > curMax {
					curMax = challenger
				}
			}
			if predictable {
				count++
			}
			return
		}
		for i := start; i < n; i++ {
			perm[start], perm[i] = perm[i], perm[start]
			generate(start + 1)
			perm[start], perm[i] = perm[i], perm[start]
		}
	}
	generate(0)
	return count % MOD
}

type testCaseG struct {
	n int
	k int64
	a []int64
}

func generateCaseG(rng *rand.Rand) (string, testCaseG) {
	n := rng.Intn(5) + 2 // 2..6
	k := int64(rng.Intn(10))
	a := make([]int64, n)
	cur := int64(0)
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(5))
		a[i] = cur
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), testCaseG{n, k, a}
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCaseG(rng)
		expect := fmt.Sprintf("%d", solveBrute(tc.n, tc.k, tc.a))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
