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

const MOD int = 998244353

func checkGood(p, q string) bool {
	n := len(p)
	for i := 0; i < n; i++ {
		ch := p[i]
		found := false
		for l := 0; l <= i && !found; l++ {
			for r := i; r < n && !found; r++ {
				if l <= i && i <= r {
					sub := q[l : r+1]
					m := r - l + 1
					count := 0
					for j := 0; j < m; j++ {
						if sub[j] == ch {
							count++
						}
					}
					need := m / 2
					if m%2 != 0 {
						need = (m + 1) / 2
					}
					if count >= need {
						found = true
					}
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func countGoodStrings(p string) int {
	n := len(p)
	total := 0
	cur := make([]byte, n)
	var dfs func(idx int)
	dfs = func(idx int) {
		if idx == n {
			if checkGood(p, string(cur)) {
				total++
				if total >= MOD {
					total -= MOD
				}
			}
			return
		}
		cur[idx] = '0'
		dfs(idx + 1)
		cur[idx] = '1'
		dfs(idx + 1)
	}
	dfs(0)
	return total % MOD
}

func generateCaseI(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	bytes := make([]byte, n)
	for i := range bytes {
		if rng.Intn(2) == 0 {
			bytes[i] = '0'
		} else {
			bytes[i] = '1'
		}
	}
	p := string(bytes)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, p))
	input := sb.String()
	ans := countGoodStrings(p)
	exp := fmt.Sprintf("%d\n", ans)
	return input, exp
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseI(rng)
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
