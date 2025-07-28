package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const numTestsE = 100

type Grid struct {
	n, m int
	val  []int
}

func (g *Grid) good(idx int) bool {
	v := g.val[idx]
	if v == 1 {
		return true
	}
	r := idx / g.m
	c := idx % g.m
	if r > 0 && g.val[idx-g.m] < v {
		return true
	}
	if r+1 < g.n && g.val[idx+g.m] < v {
		return true
	}
	if c > 0 && g.val[idx-1] < v {
		return true
	}
	if c+1 < g.m && g.val[idx+1] < v {
		return true
	}
	return false
}

func (g *Grid) goodAfter(idx, a, b, valA, valB int) bool {
	v := g.val[idx]
	if idx == a {
		v = valB
	} else if idx == b {
		v = valA
	}
	if v == 1 {
		return true
	}
	r := idx / g.m
	c := idx % g.m
	if r > 0 {
		nv := g.val[idx-g.m]
		if idx-g.m == a {
			nv = valB
		} else if idx-g.m == b {
			nv = valA
		}
		if nv < v {
			return true
		}
	}
	if r+1 < g.n {
		nv := g.val[idx+g.m]
		if idx+g.m == a {
			nv = valB
		} else if idx+g.m == b {
			nv = valA
		}
		if nv < v {
			return true
		}
	}
	if c > 0 {
		nv := g.val[idx-1]
		if idx-1 == a {
			nv = valB
		} else if idx-1 == b {
			nv = valA
		}
		if nv < v {
			return true
		}
	}
	if c+1 < g.m {
		nv := g.val[idx+1]
		if idx+1 == a {
			nv = valB
		} else if idx+1 == b {
			nv = valA
		}
		if nv < v {
			return true
		}
	}
	return false
}

func solveE(input string) string {
	reader := strings.NewReader(input)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	g := Grid{n: n, m: m, val: make([]int, n*m)}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			idx := i*m + j
			fmt.Fscan(reader, &g.val[idx])
		}
	}

	good := make([]bool, n*m)
	badList := make([]int, 0)
	for i := 0; i < n*m; i++ {
		if g.good(i) {
			good[i] = true
		} else {
			badList = append(badList, i)
		}
	}
	if len(badList) == 0 {
		return "0\n"
	}
	if len(badList) > 10 {
		return "2\n"
	}

	candidate := make([]bool, n*m)
	addCand := func(x int) {
		if x >= 0 && x < n*m {
			candidate[x] = true
		}
	}
	for _, idx := range badList {
		addCand(idx)
		r := idx / m
		c := idx % m
		if r > 0 {
			addCand(idx - m)
		}
		if r+1 < n {
			addCand(idx + m)
		}
		if c > 0 {
			addCand(idx - 1)
		}
		if c+1 < m {
			addCand(idx + 1)
		}
	}

	candList := make([]int, 0)
	for i := 0; i < n*m; i++ {
		if candidate[i] {
			candList = append(candList, i)
		}
	}

	visited := make([]int, n*m)
	cur := 0
	count := 0

	buildUnion := func(a, b int) []int {
		cur++
		res := make([]int, 0, 10)
		var add func(int)
		add = func(x int) {
			if x >= 0 && x < n*m {
				if visited[x] != cur {
					visited[x] = cur
					res = append(res, x)
				}
			}
		}
		add(a)
		add(b)
		ra := a / m
		ca := a % m
		if ra > 0 {
			add(a - m)
		}
		if ra+1 < n {
			add(a + m)
		}
		if ca > 0 {
			add(a - 1)
		}
		if ca+1 < m {
			add(a + 1)
		}

		rb := b / m
		cb := b % m
		if rb > 0 {
			add(b - m)
		}
		if rb+1 < n {
			add(b + m)
		}
		if cb > 0 {
			add(b - 1)
		}
		if cb+1 < m {
			add(b + 1)
		}
		return res
	}

	badCount := len(badList)

	checkSwap := func(a, b int) bool {
		if a == b {
			return false
		}
		va := g.val[a]
		vb := g.val[b]
		cells := buildUnion(a, b)
		nb := badCount
		for _, p := range cells {
			old := good[p]
			nw := g.goodAfter(p, a, b, va, vb)
			if old != nw {
				if old {
					nb++
				} else {
					nb--
				}
			}
		}
		return nb == 0
	}

	totalCells := n * m
	for _, a := range candList {
		for b := 0; b < totalCells; b++ {
			if a >= b {
				continue
			}
			if checkSwap(a, b) {
				count++
			}
		}
	}

	if count > 0 {
		return fmt.Sprintf("1 %d\n", count)
	}
	return "2\n"
}

func generateTestsE() []string {
	rng := rand.New(rand.NewSource(5))
	tests := make([]string, numTestsE)
	for i := 0; i < numTestsE; i++ {
		n := rng.Intn(3) + 2
		m := rng.Intn(3) + 2
		vals := rand.Perm(n * m)
		for j := range vals {
			vals[j]++
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		idx := 0
		for r := 0; r < n; r++ {
			for c := 0; c < m; c++ {
				sb.WriteString(fmt.Sprintf("%d", vals[idx]))
				idx++
				if c+1 < m {
					sb.WriteByte(' ')
				}
			}
			sb.WriteByte('\n')
		}
		tests[i] = sb.String()
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsE()
	for i, tc := range tests {
		expected := solveE(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("Test %d: error running binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, tc, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
