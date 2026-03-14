package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseE struct {
	n   int
	m   int
	q   int
	a   []int
	b   []int
	ops []opE
}

type opE struct {
	tp  int
	pos int
	x   int
}

func generateTests() []testCaseE {
	rng := rand.New(rand.NewSource(42))
	tests := make([]testCaseE, 100)
	for i := range tests {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		a := make([]int, n)
		b := make([]int, m)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(5)
		}
		for j := 0; j < m; j++ {
			b[j] = rng.Intn(5)
		}
		ops := make([]opE, q)
		for j := 0; j < q; j++ {
			tp := rng.Intn(3) + 1
			pos := 0
			if tp != 3 {
				if tp == 1 {
					pos = rng.Intn(n) + 1
				} else {
					pos = rng.Intn(m) + 1
				}
			}
			x := rng.Intn(5)
			ops[j] = opE{tp: tp, pos: pos, x: x}
		}
		tests[i] = testCaseE{n: n, m: m, q: q, a: a, b: b, ops: ops}
	}
	return tests
}

func computeDiff(k int, a []int, b []int) int64 {
	n := len(a)
	m := len(b)
	total := n + m
	type Player struct{ p, team int }
	players := make([]Player, total)
	for i, v := range a {
		players[i] = Player{v, 1}
	}
	for i, v := range b {
		players[n+i] = Player{v, -1}
	}

	best := int64(-1 << 62)
	var dfs func(idx, prevP, device int, diff int64)
	dfs = func(idx, prevP, device int, diff int64) {
		if idx == total {
			if diff > best {
				best = diff
			}
			return
		}
		seen := make(map[[2]int]bool)
		for i := idx; i < total; i++ {
			pl := players[i]
			key := [2]int{pl.p, pl.team}
			if seen[key] {
				continue
			}
			seen[key] = true
			var nd int
			if idx == 0 {
				nd = k
			} else {
				nd = device + pl.p - prevP
				if nd < 0 {
					nd = 0
				}
			}
			players[idx], players[i] = players[i], players[idx]
			dfs(idx+1, pl.p, nd, diff+int64(pl.team)*int64(nd))
			players[idx], players[i] = players[i], players[idx]
		}
	}

	dfs(0, 0, k, 0)
	return best
}

func solveE(t testCaseE) []int64 {
	a := append([]int(nil), t.a...)
	b := append([]int(nil), t.b...)
	res := make([]int64, 0, t.q)
	for _, op := range t.ops {
		if op.tp == 1 {
			if op.pos >= 1 && op.pos <= len(a) {
				a[op.pos-1] = op.x
			}
		} else if op.tp == 2 {
			if op.pos >= 1 && op.pos <= len(b) {
				b[op.pos-1] = op.x
			}
		} else {
			res = append(res, computeDiff(op.x, a, b))
		}
	}
	return res
}

func buildInput(t testCaseE) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", t.n, t.m, t.q)
	for i := 0; i < t.n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", t.a[i])
	}
	b.WriteByte('\n')
	for i := 0; i < t.m; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", t.b[i])
	}
	b.WriteByte('\n')
	for _, op := range t.ops {
		if op.tp == 1 || op.tp == 2 {
			fmt.Fprintf(&b, "%d %d %d\n", op.tp, op.pos, op.x)
		} else {
			fmt.Fprintf(&b, "%d %d\n", op.tp, op.x)
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for idx, t := range tests {
		input := buildInput(t)
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "execution failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(&out)
		outputs := []string{}
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				outputs = append(outputs, line)
			}
		}
		expected := solveE(t)
		if len(outputs) != len(expected) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d lines, got %d\n", idx+1, len(expected), len(outputs))
			os.Exit(1)
		}
		for i, v := range expected {
			if outputs[i] != fmt.Sprintf("%d", v) {
				fmt.Fprintf(os.Stderr, "test %d failed at line %d: expected %d got %s\n", idx+1, i+1, v, outputs[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}
