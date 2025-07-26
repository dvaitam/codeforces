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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solve(n, k int, s string) string {
	L0, R0 := n, -1
	L1, R1 := n, -1
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			L0 = min(L0, i)
			R0 = max(R0, i)
		} else {
			L1 = min(L1, i)
			R1 = max(R1, i)
		}
	}
	if R0-L0+1 <= k || R1-L1+1 <= k {
		return "tokitsukaze"
	}
	prev0 := make([]int, n)
	prev1 := make([]int, n)
	last0, last1 := -1, -1
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			last0 = i
		} else {
			last1 = i
		}
		prev0[i] = last0
		prev1[i] = last1
	}
	next0 := make([]int, n)
	next1 := make([]int, n)
	nxt0, nxt1 := n, n
	for i := n - 1; i >= 0; i-- {
		if s[i] == '0' {
			nxt0 = i
		} else {
			nxt1 = i
		}
		next0[i] = nxt0
		next1[i] = nxt1
	}
	onceAgain := false
	for i := 0; i+k <= n; i++ {
		l, r := i, i+k-1
		// flip to 0
		newL0 := min(L0, l)
		newR0 := max(R0, r)
		nl1, nr1 := n, -1
		if l > 0 {
			p := prev1[l-1]
			if p >= 0 {
				nl1 = min(nl1, p)
				nr1 = max(nr1, p)
			}
		}
		if r+1 < n {
			p := next1[r+1]
			if p < n {
				nl1 = min(nl1, p)
				nr1 = max(nr1, p)
			}
		}
		if newR0-newL0+1 > k && nr1-nl1+1 > k {
			onceAgain = true
			break
		}
		// flip to 1
		newL1 := min(L1, l)
		newR1 := max(R1, r)
		nl0, nr0 := n, -1
		if l > 0 {
			p := prev0[l-1]
			if p >= 0 {
				nl0 = min(nl0, p)
				nr0 = max(nr0, p)
			}
		}
		if r+1 < n {
			p := next0[r+1]
			if p < n {
				nl0 = min(nl0, p)
				nr0 = max(nr0, p)
			}
		}
		if newR1-newL1+1 > k && nr0-nl0+1 > k {
			onceAgain = true
			break
		}
	}
	if onceAgain {
		return "once again"
	}
	return "quailty"
}

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	k := rng.Intn(n) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	s := string(b)
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	return input, solve(n, k, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
