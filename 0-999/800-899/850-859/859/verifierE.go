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

const modE int = 1000000007

type CaseE struct {
	input    string
	expected int
}

func solveE(n int, pairs [][2]int) int {
	seatIndex := make(map[int]int)
	idx := 0
	for _, p := range pairs {
		if _, ok := seatIndex[p[0]]; !ok {
			seatIndex[p[0]] = idx
			idx++
		}
		if _, ok := seatIndex[p[1]]; !ok {
			seatIndex[p[1]] = idx
			idx++
		}
	}
	m := idx
	parent := make([]int, n+m)
	size := make([]int, n+m)
	for i := range parent {
		parent[i] = i
		size[i] = 1
	}
	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(a, b int) {
		a = find(a)
		b = find(b)
		if a == b {
			return
		}
		if size[a] < size[b] {
			a, b = b, a
		}
		parent[b] = a
		size[a] += size[b]
	}
	for i, p := range pairs {
		union(i, n+seatIndex[p[0]])
		union(i, n+seatIndex[p[1]])
	}
	type comp struct {
		t, s int
		self bool
	}
	comps := make(map[int]*comp)
	for i, p := range pairs {
		r := find(i)
		if comps[r] == nil {
			comps[r] = &comp{}
		}
		c := comps[r]
		c.t++
		if p[0] == p[1] {
			c.self = true
		}
	}
	for _, idx := range seatIndex {
		r := find(n + idx)
		if comps[r] == nil {
			comps[r] = &comp{}
		}
		comps[r].s++
	}
	res := 1
	for _, c := range comps {
		if c.s == c.t+1 {
			res = res * c.s % modE
		} else if c.s == c.t {
			if c.self {
				res = res * 1 % modE
			} else {
				res = res * 2 % modE
			}
		}
	}
	return res % modE
}

func generateCaseE(rng *rand.Rand) CaseE {
	n := rng.Intn(20) + 1
	pairs := make([][2]int, n)
	for i := range pairs {
		x := rng.Intn(50) + 1
		y := rng.Intn(50) + 1
		pairs[i] = [2]int{x, y}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, p := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return CaseE{sb.String(), solveE(n, pairs)}
}

func runCase(exe, input string, expected int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(strings.TrimSpace(out.String())), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	pairs := [][2]int{{1, 2}}
	cases := []CaseE{{"1\n1 2\n", solveE(1, pairs)}}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseE(rng))
	}
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
