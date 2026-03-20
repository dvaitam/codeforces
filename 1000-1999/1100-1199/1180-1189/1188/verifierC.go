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

// Embedded oracle solver from the correct CF-accepted solution.
func oracleSolve(input string) string {
	// Parse input
	words := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v := 0
		s := words[idx]
		idx++
		for _, ch := range s {
			v = v*10 + int(ch-'0')
		}
		return v
	}
	n := nextInt()
	k := nextInt()
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = nextInt()
	}
	sort.Ints(a)

	MOD := 998244353
	if k == 1 {
		// beauty is undefined for k=1 (n>1 required for beauty)
		// but problem says k >= 2, so this shouldn't happen
		return "0"
	}
	MaxV := (a[n-1] - a[0]) / (k - 1)

	ans := 0
	dp := make([]int, n)
	new_dp := make([]int, n)

	for v := 1; v <= MaxV; v++ {
		for i := 0; i < n; i++ {
			dp[i] = 1
		}
		for j := 2; j <= k; j++ {
			sum := 0
			p := 0
			allZero := true
			for i := 0; i < n; i++ {
				for p < i && a[i]-a[p] >= v {
					sum += dp[p]
					if sum >= MOD {
						sum -= MOD
					}
					p++
				}
				new_dp[i] = sum
				if sum > 0 {
					allZero = false
				}
			}
			for i := 0; i < n; i++ {
				dp[i] = new_dp[i]
			}
			if allZero {
				break
			}
		}
		for i := 0; i < n; i++ {
			ans += dp[i]
			if ans >= MOD {
				ans -= MOD
			}
		}
	}
	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 2 // at least 2
	k := rng.Intn(n-1) + 2 // k in [2, n], matching problem constraint k >= 2
	if k > n {
		k = n
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := oracleSolve(in)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
