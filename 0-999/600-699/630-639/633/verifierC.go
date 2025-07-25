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

type Node struct {
	next [26]int32
	vis  int
}

func solveC(n int, s string, m int, words []string) string {
	sb := []byte(s)
	for i, j := 0, len(sb)-1; i < j; i, j = i+1, j-1 {
		sb[i], sb[j] = sb[j], sb[i]
	}
	ws := make([]string, m+1)
	copy(ws[1:], words)
	nodes := make([]Node, 1)
	for i := 1; i <= m; i++ {
		w := ws[i]
		x := int32(0)
		for j := 0; j < len(w); j++ {
			c := w[j]
			if c >= 'A' && c <= 'Z' {
				c |= 32
			}
			idx := c - 'a'
			if nodes[x].next[idx] == 0 {
				nodes = append(nodes, Node{})
				nodes[x].next[idx] = int32(len(nodes) - 1)
			}
			x = nodes[x].next[idx]
		}
		nodes[x].vis = i
	}
	dp := make([]int8, n+1)
	for i := range dp {
		dp[i] = -1
	}
	ans := []int{}
	var calc func(int) bool
	calc = func(i int) bool {
		if i == n {
			return true
		}
		if dp[i] != -1 {
			return dp[i] == 1
		}
		x := int32(0)
		for j := i; j < n; j++ {
			c := sb[j]
			idx := c - 'a'
			if idx < 0 || idx >= 26 || nodes[x].next[idx] == 0 {
				dp[i] = 0
				return false
			}
			x = nodes[x].next[idx]
			if nodes[x].vis != 0 {
				if calc(j + 1) {
					ans = append(ans, nodes[x].vis)
					dp[i] = 1
					return true
				}
			}
		}
		dp[i] = 0
		return false
	}
	calc(0)
	var out strings.Builder
	for i := 0; i < len(ans); i++ {
		out.WriteString(ws[ans[i]])
		if i+1 < len(ans) {
			out.WriteByte(' ')
		}
	}
	out.WriteByte('\n')
	return out.String()
}

func generateCaseC(rng *rand.Rand) (string, string) {
	m := rng.Intn(3) + 1
	words := make([]string, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(3) + 1
		var sb strings.Builder
		for j := 0; j < l; j++ {
			ch := byte('a' + rng.Intn(26))
			if rng.Intn(2) == 1 {
				ch -= 'a' - 'A'
			}
			sb.WriteByte(ch)
		}
		words[i] = sb.String()
	}
	k := rng.Intn(3) + 1
	var cipher strings.Builder
	for i := 0; i < k; i++ {
		w := strings.ToLower(words[rng.Intn(m)])
		for j := len(w) - 1; j >= 0; j-- {
			cipher.WriteByte(w[j])
		}
	}
	s := cipher.String()
	n := len(s)
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n%s\n%d\n", n, s, m))
	for i := 0; i < m; i++ {
		input.WriteString(words[i])
		if i+1 < m {
			input.WriteByte('\n')
		}
	}
	input.WriteByte('\n')
	expect := solveC(n, s, m, words)
	return input.String(), expect
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
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseC(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
