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

const mod = 998244353

func add(a, b int) int {
	a += b
	if a >= mod {
		a -= mod
	}
	return a
}
func mul(a, b int) int { return int((int64(a) * int64(b)) % mod) }

func solve(n int, parents []int) int {
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		children[p] = append(children[p], i)
	}
	leafCount := make([]int, n+1)
	dp := make([]int, n+1)
	type fr struct{ u, state int }
	stack := []fr{{1, 0}}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		u, state := v.u, v.state
		if state == 0 {
			stack = append(stack, fr{u, 1})
			for _, ch := range children[u] {
				stack = append(stack, fr{ch, 0})
			}
		} else {
			if len(children[u]) == 0 {
				leafCount[u] = 1
				dp[u] = 1
			} else {
				tot := 0
				prod := 1
				for _, ch := range children[u] {
					tot += leafCount[ch]
					prod = mul(prod, dp[ch])
				}
				leafCount[u] = tot
				if tot >= 2 {
					dp[u] = add(prod, 1)
				} else {
					dp[u] = prod
				}
			}
		}
	}
	return dp[1]
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 2
	parents := make([]int, n-1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		parents[i-2] = p
		fmt.Fprintf(&sb, "%d ", p)
	}
	sb.WriteByte('\n')
	exp := solve(n, parents)
	return sb.String(), fmt.Sprint(exp)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
