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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveCase(a []int) string {
	n := len(a)
	if n <= 3 {
		return "YES"
	}
	g := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		g[i] = gcd(a[i], a[i+1])
	}
	pre := make([]bool, n-1)
	pre[0] = true
	for i := 1; i < n-1; i++ {
		pre[i] = pre[i-1] && g[i-1] <= g[i]
	}
	suf := make([]bool, n-1)
	suf[n-2] = true
	for i := n - 3; i >= 0; i-- {
		suf[i] = suf[i+1] && g[i] <= g[i+1]
	}
	for i := 0; i < n; i++ {
		var ok bool
		if i == 0 {
			if n-2 <= 0 {
				ok = true
			} else {
				ok = suf[1]
			}
		} else if i == n-1 {
			if n-3 < 0 {
				ok = true
			} else {
				ok = pre[n-3]
			}
		} else {
			newG := gcd(a[i-1], a[i+1])
			ok = true
			if i-2 >= 0 {
				ok = ok && pre[i-2] && g[i-2] <= newG
			}
			if i+1 <= n-2 {
				ok = ok && suf[i+1] && newG <= g[i+1]
			}
		}
		if ok {
			return "YES"
		}
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	a := make([]int, n)
	var sb strings.Builder
	for i := range a {
		a[i] = rng.Intn(20) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	input := fmt.Sprintf("1\n%d\n%s\n", n, sb.String())
	expect := solveCase(a)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
