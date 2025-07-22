package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n int
	p []int
	m int
}

func (tc testCase) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, x := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", x))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.m))
	return sb.String()
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(7) + 1
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	m := rng.Intn(n)
	return testCase{n: n, p: perm, m: m}
}

func computeOps(n int, p []int, m int) [][2]int {
	comp := make([]int, n+1)
	compList := make(map[int][]int)
	visited := make([]bool, n+1)
	compID := 0
	for i := 1; i <= n; i++ {
		if !visited[i] {
			compID++
			u := i
			for !visited[u] {
				visited[u] = true
				comp[u] = compID
				compList[compID] = append(compList[compID], u)
				u = p[u-1]
			}
		}
	}
	for id := 1; id <= compID; id++ {
		sort.Ints(compList[id])
	}
	cycles := compID
	targetCycles := n - m
	var ops [][2]int
	var rebuild func([]int)
	rebuild = func(nodes []int) {
		for _, u := range nodes {
			visited[u] = false
			comp[u] = 0
		}
		for _, u := range nodes {
			if !visited[u] {
				compID++
				var group []int
				v := u
				for !visited[v] {
					visited[v] = true
					comp[v] = compID
					group = append(group, v)
					v = p[v-1]
				}
				sort.Ints(group)
				compList[compID] = group
			}
		}
	}
	for cycles < targetCycles {
		var ai, bi int
		for i := 1; i <= n; i++ {
			cl := compList[comp[i]]
			if len(cl) >= 2 {
				ai = i
				for _, x := range cl {
					if x > i {
						bi = x
						break
					}
				}
				break
			}
		}
		ops = append(ops, [2]int{ai, bi})
		p[ai-1], p[bi-1] = p[bi-1], p[ai-1]
		old := compList[comp[ai]]
		delete(compList, comp[ai])
		rebuild(old)
		cycles++
	}
	for cycles > targetCycles {
		var ai, bi int
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				if comp[i] != comp[j] {
					ai, bi = i, j
					break
				}
			}
			if ai != 0 {
				break
			}
		}
		ops = append(ops, [2]int{ai, bi})
		p[ai-1], p[bi-1] = p[bi-1], p[ai-1]
		id1, id2 := comp[ai], comp[bi]
		nodes := append(compList[id1], compList[id2]...)
		delete(compList, id1)
		delete(compList, id2)
		rebuild(nodes)
		cycles--
	}
	return ops
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	input := tc.Input()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	r := strings.NewReader(out.String())
	var k int
	if _, err := fmt.Fscan(r, &k); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	ops := make([][2]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(r, &ops[i][0], &ops[i][1]); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
	}
	expect := computeOps(tc.n, append([]int(nil), tc.p...), tc.m)
	if len(ops) != len(expect) {
		return fmt.Errorf("expected %d ops got %d", len(expect), len(ops))
	}
	for i := range expect {
		if ops[i][0] != expect[i][0] || ops[i][1] != expect[i][1] {
			return fmt.Errorf("op %d expected %v got %v", i+1, expect[i], ops[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		{n: 1, p: []int{1}, m: 0},
		{n: 2, p: []int{2, 1}, m: 1},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.Input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
