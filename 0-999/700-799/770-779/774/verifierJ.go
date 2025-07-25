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

func expected(n, k int, s string) string {
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true
	for idx := 0; idx < n; idx++ {
		ndp := make([][]bool, n+1)
		for i := range ndp {
			ndp[i] = make([]bool, n+1)
		}
		for j := 0; j <= n; j++ {
			for m := 0; m <= n; m++ {
				if !dp[j][m] {
					continue
				}
				c := s[idx]
				if c != 'N' {
					newJ := 0
					newM := m
					if j > newM {
						newM = j
					}
					ndp[newJ][newM] = true
				}
				if c != 'Y' {
					newJ := j + 1
					newM := m
					if newJ > newM {
						newM = newJ
					}
					if newJ <= n && newM <= n {
						ndp[newJ][newM] = true
					}
				}
			}
		}
		dp = ndp
	}
	for j := 0; j <= n; j++ {
		if dp[j][k] {
			return "YES"
		}
	}
	return "NO"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	if rng.Float64() < 0.1 {
		n = rng.Intn(20) + 1
	}
	k := rng.Intn(n + 1)
	letters := []byte{'Y', 'N', '?'}
	sb := make([]byte, n)
	for i := 0; i < n; i++ {
		sb[i] = letters[rng.Intn(len(letters))]
	}
	s := string(sb)
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	return input, expected(n, k, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
