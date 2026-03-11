package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// bruteHasSolution returns true if there exists any pair of disjoint intervals with equal XOR
func bruteHasSolution(k int, g []int) bool {
	n := 1 << (k + 1)
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] ^ g[i-1]
	}
	for l1 := 1; l1 <= n; l1++ {
		for r1 := l1; r1 <= n; r1++ {
			x1 := pref[r1] ^ pref[l1-1]
			for l2 := r1 + 1; l2 <= n; l2++ {
				for r2 := l2; r2 <= n; r2++ {
					x2 := pref[r2] ^ pref[l2-1]
					if x1 == x2 {
						return true
					}
				}
			}
		}
	}
	return false
}

// validateAnswer checks if a candidate answer (a b c d) is valid for the given test case
func validateAnswer(k int, g []int, output string) error {
	n := 1 << (k + 1)
	trimmed := strings.TrimSpace(output)

	// Check if -1
	if trimmed == "-1" {
		// Verify that indeed no solution exists
		if bruteHasSolution(k, g) {
			return fmt.Errorf("candidate says -1 but a solution exists")
		}
		return nil
	}

	tokens := strings.Fields(trimmed)
	if len(tokens) != 4 {
		return fmt.Errorf("expected 4 integers or -1, got: %s", trimmed)
	}

	vals := make([]int, 4)
	for i := 0; i < 4; i++ {
		v, err := strconv.Atoi(tokens[i])
		if err != nil {
			return fmt.Errorf("failed to parse token %d: %v", i, err)
		}
		vals[i] = v
	}
	a, b, c, d := vals[0], vals[1], vals[2], vals[3]

	// Validate ranges
	if a < 1 || a > n || b < 1 || b > n || c < 1 || c > n || d < 1 || d > n {
		return fmt.Errorf("indices out of range [1,%d]: %d %d %d %d", n, a, b, c, d)
	}
	if a > b {
		return fmt.Errorf("invalid interval [%d,%d]: a > b", a, b)
	}
	if c > d {
		return fmt.Errorf("invalid interval [%d,%d]: c > d", c, d)
	}

	// Check disjoint: intervals [a,b] and [c,d] must be disjoint
	if !(b < c || d < a) {
		return fmt.Errorf("intervals [%d,%d] and [%d,%d] are not disjoint", a, b, c, d)
	}

	// Compute XOR of each interval
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] ^ g[i-1]
	}
	xor1 := pref[b] ^ pref[a-1]
	xor2 := pref[d] ^ pref[c-1]

	if xor1 != xor2 {
		return fmt.Errorf("XOR mismatch: [%d,%d] has XOR %d, [%d,%d] has XOR %d", a, b, xor1, c, d, xor2)
	}

	return nil
}

func generateCase(rng *rand.Rand) (string, int, []int) {
	k := rng.Intn(3) // 0..2
	n := 1 << (k + 1)
	g := make([]int, n)
	maxVal := 1 << (2 * k) // 4^k
	if maxVal < 1 {
		maxVal = 1
	}
	for i := 0; i < n; i++ {
		if maxVal == 1 {
			g[i] = 0
		} else {
			g[i] = rng.Intn(maxVal)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", g[i])
	}
	sb.WriteByte('\n')
	return sb.String(), k, g
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, k, g := generateCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validateAnswer(k, g, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
