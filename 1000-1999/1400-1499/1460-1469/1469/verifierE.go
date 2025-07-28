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

func expected(n, k int, s string) string {
	m := k
	if m > 20 {
		m = 20
	}
	size := 1 << m
	seen := make([]bool, size)
	zeros := make([]int, n+1)
	for i := 0; i < n; i++ {
		zeros[i+1] = zeros[i]
		if s[i] == '0' {
			zeros[i+1]++
		}
	}
	mask := 0
	for i := 0; i < n; i++ {
		mask = ((mask << 1) & (size - 1))
		if s[i] == '1' {
			mask |= 1
		}
		if i >= k-1 {
			start := i - k + 1
			if k > m {
				if zeros[start+k-m]-zeros[start] == 0 {
					seen[mask] = true
				}
			} else {
				seen[mask] = true
			}
		}
	}
	ans := -1
	for y := 0; y < size; y++ {
		var idx int
		if k > m {
			idx = (^y) & (size - 1)
		} else {
			idx = (size - 1) ^ y
		}
		if !seen[idx] {
			ans = y
			break
		}
	}
	if ans == -1 {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	if k > m {
		for i := 0; i < k-m; i++ {
			sb.WriteByte('0')
		}
	}
	sb.WriteString(fmt.Sprintf("%0*b", m, ans))
	return sb.String()
}

func generateCase(rng *rand.Rand) (int, int, string) {
	n := rng.Intn(50) + 1
	k := rng.Intn(n) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return n, k, sb.String()
}

func runCase(bin string, n, k int, s string) error {
	input := fmt.Sprintf("1\n%d %d\n%s\n", n, k, s)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	expect := expected(n, k, s)
	if strings.TrimSpace(out) != expect {
		return fmt.Errorf("expected %q got %q", expect, strings.TrimSpace(out))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, s := generateCase(rng)
		if err := runCase(bin, n, k, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
