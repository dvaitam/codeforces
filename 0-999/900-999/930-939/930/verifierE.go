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

const MOD int64 = 1000000007

// bruteForce enumerates all 2^k bit strings and counts valid ones.
// Only suitable for small k (≤ 20).
func bruteForce(k int, A, B [][2]int) int64 {
	count := int64(0)
	total := 1 << k
	for mask := 0; mask < total; mask++ {
		valid := true
		for _, ab := range A {
			hasOne := false
			for j := ab[0]; j <= ab[1]; j++ {
				if (mask>>(j-1))&1 == 1 {
					hasOne = true
					break
				}
			}
			if !hasOne {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		for _, kb := range B {
			hasZero := false
			for j := kb[0]; j <= kb[1]; j++ {
				if (mask>>(j-1))&1 == 0 {
					hasZero = true
					break
				}
			}
			if !hasZero {
				valid = false
				break
			}
		}
		if valid {
			count++
		}
	}
	return count % MOD
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// genCase generates a random test with k ≤ 15 so brute-force stays fast.
func genCase(rng *rand.Rand) (string, int, [][2]int, [][2]int) {
	k := rng.Intn(15) + 1
	n := rng.Intn(6)
	m := rng.Intn(6)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", k, n, m)
	A := make([][2]int, n)
	B := make([][2]int, m)
	for i := 0; i < n; i++ {
		l := rng.Intn(k) + 1
		r := l + rng.Intn(k-l+1)
		A[i] = [2]int{l, r}
		fmt.Fprintf(&sb, "%d %d\n", l, r)
	}
	for i := 0; i < m; i++ {
		l := rng.Intn(k) + 1
		r := l + rng.Intn(k-l+1)
		B[i] = [2]int{l, r}
		fmt.Fprintf(&sb, "%d %d\n", l, r)
	}
	return sb.String(), k, A, B
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const total = 500
	for i := 1; i <= total; i++ {
		input, k, A, B := genCase(rng)
		exp := fmt.Sprintf("%d", bruteForce(k, A, B))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", total)
}
