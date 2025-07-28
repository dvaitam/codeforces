package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	rand.Seed(42)
	tests := make([]testCaseE, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		q := rand.Intn(5) + 1
		a := make([]int, n)
		b := make([]int, m)
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(5)
		}
		for j := 0; j < m; j++ {
			b[j] = rand.Intn(5)
		}
		ops := make([]opE, q)
		for j := 0; j < q; j++ {
			tp := rand.Intn(3) + 1
			pos := 0
			if tp != 3 {
				if tp == 1 {
					pos = rand.Intn(n) + 1
				} else {
					pos = rand.Intn(m) + 1
				}
			}
			x := rand.Intn(5)
			ops[j] = opE{tp: tp, pos: pos, x: x}
		}
		tests[i] = testCaseE{n: n, m: m, q: q, a: a, b: b, ops: ops}
	}
	return tests
}

func computeDiff(k int, a []int, b []int) int64 {
	n := len(a)
	m := len(b)
	type Player struct {
		p    int
		team int
	}
	players := make([]Player, 0, n+m)
	for _, v := range a {
		players = append(players, Player{p: v, team: 1})
	}
	for _, v := range b {
		players = append(players, Player{p: v, team: -1})
	}
	sort.Slice(players, func(i, j int) bool {
		if players[i].p == players[j].p {
			return players[i].team > players[j].team
		}
		return players[i].p > players[j].p
	})
	if len(players) == 0 {
		return 0
	}
	device := k
	prev := players[0].p
	var diff int64
	diff += int64(players[0].team) * int64(device)
	for i := 1; i < len(players); i++ {
		d := players[i].p - prev
		device += d
		if device < 0 {
			device = 0
		}
		prev = players[i].p
		diff += int64(players[i].team) * int64(device)
	}
	return diff
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
