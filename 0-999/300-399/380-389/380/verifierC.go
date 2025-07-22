package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Node struct {
	pairs int
	open  int
	close int
}

func merge(a, b Node) Node {
	match := a.open
	if b.close < match {
		match = b.close
	}
	return Node{pairs: a.pairs + b.pairs + match, open: a.open + b.open - match, close: a.close + b.close - match}
}

func solveC(r io.Reader) string {
	reader := bufio.NewReader(r)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return ""
	}
	n := len(s)
	N := 1
	for N < n {
		N <<= 1
	}
	tree := make([]Node, 2*N)
	for i := 0; i < n; i++ {
		idx := N + i
		if s[i] == '(' {
			tree[idx] = Node{open: 1}
		} else {
			tree[idx] = Node{close: 1}
		}
	}
	for i := N - 1; i > 0; i-- {
		tree[i] = merge(tree[2*i], tree[2*i+1])
	}
	var m int
	fmt.Fscan(reader, &m)
	var sb strings.Builder
	for q := 0; q < m; q++ {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		l--
		r--
		l += N
		r += N
		resL := Node{}
		resR := Node{}
		for l <= r {
			if l&1 == 1 {
				resL = merge(resL, tree[l])
				l++
			}
			if r&1 == 0 {
				resR = merge(tree[r], resR)
				r--
			}
			l >>= 1
			r >>= 1
		}
		res := merge(resL, resR)
		fmt.Fprintln(&sb, res.pairs*2)
	}
	return sb.String()
}

func runCaseC(bin string, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(out.String()))
	}
	return nil
}

func genCaseC(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('(')
		} else {
			sb.WriteByte(')')
		}
	}
	sb.WriteByte('\n')
	m := rng.Intn(20) + 1
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCaseC(rng)
		expect := solveC(strings.NewReader(in))
		if err := runCaseC(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
