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

type caseC struct {
	n       int
	lengths []int
}

func generateCase(rng *rand.Rand) caseC {
	n := rng.Intn(5) + 1
	lengths := make([]int, n)
	for i := range lengths {
		lengths[i] = rng.Intn(4) + 1
	}
	return caseC{n, lengths}
}

func solveCase(tc caseC) (bool, []string) {
	n := tc.n
	a := tc.lengths
	maxD := 0
	for _, v := range a {
		if v > maxD {
			maxD = v
		}
	}
	b := make([]int, maxD+2)
	q := make([][]int, maxD+2)
	for i, d := range a {
		b[d]++
		q[d] = append(q[d], i)
	}
	cur := int64(2)
	const inf = 100000000
	for depth := 1; depth <= maxD; depth++ {
		if int64(b[depth]) > cur {
			return false, nil
		}
		cur = (cur - int64(b[depth])) * 2
		if cur > inf {
			break
		}
	}
	ans := make([][]byte, n)
	for i := range ans {
		ans[i] = make([]byte, a[i])
	}
	c := make([]byte, maxD+2)
	q1 := make([]int, maxD+2)
	all := n
	var dfs func(int)
	END := false
	dfs = func(depth int) {
		if END {
			return
		}
		if depth <= maxD && b[depth] > 0 {
			b[depth]--
			all--
			idx := q[depth][q1[depth]]
			for i := 0; i < depth; i++ {
				ans[idx][i] = c[i]
			}
			q1[depth]++
			if all == 0 {
				END = true
			}
			return
		}
		if depth > maxD {
			return
		}
		c[depth] = '0'
		dfs(depth + 1)
		if END {
			return
		}
		c[depth] = '1'
		dfs(depth + 1)
	}
	dfs(0)
	words := make([]string, n)
	for i := range ans {
		words[i] = string(ans[i])
	}
	return true, words
}

func runCase(bin string, tc caseC) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", tc.n)
	for i, v := range tc.lengths {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", v)
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	ok, _ := solveCase(tc)
	if !ok {
		if len(lines) == 0 || strings.TrimSpace(lines[0]) != "NO" {
			return fmt.Errorf("expected NO, got %s", out.String())
		}
		return nil
	}
	if len(lines) != tc.n+1 {
		return fmt.Errorf("expected %d lines, got %d", tc.n+1, len(lines))
	}
	if strings.TrimSpace(lines[0]) != "YES" {
		return fmt.Errorf("expected YES")
	}
	for i := 0; i < tc.n; i++ {
		w := strings.TrimSpace(lines[i+1])
		if len(w) != tc.lengths[i] {
			return fmt.Errorf("word %d length mismatch", i)
		}
		for j := 0; j < len(w); j++ {
			if w[j] != '0' && w[j] != '1' {
				return fmt.Errorf("bad char in word")
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
