package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type instance struct {
	n, m int
	grid [][]int
}

type test struct {
	input    string
	expected string
}

func solveInstance(inst instance) string {
	n, m := inst.n, inst.m
	N := n * m
	K := N / 2
	pos := make([][2]int, K)
	cnt := make([]int, K)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := inst.grid[i][j] - 1
			idx := i*m + j
			if cnt[val] < 2 {
				pos[val][cnt[val]] = idx
			}
			cnt[val]++
		}
	}
	edges := make([][2]int, 0, n*(m-1)+m*(n-1))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			u := i*m + j
			if j+1 < m {
				edges = append(edges, [2]int{u, u + 1})
			}
			if i+1 < n {
				edges = append(edges, [2]int{u, u + m})
			}
		}
	}
	E := len(edges)
	costMat := make([][]int, K)
	cand := make([][]int, K)
	for k := 0; k < K; k++ {
		costMat[k] = make([]int, E)
		cand[k] = make([]int, E)
		p0 := pos[k][0]
		p1 := pos[k][1]
		for e := 0; e < E; e++ {
			u := edges[e][0]
			v := edges[e][1]
			c1 := 0
			if p0 != u {
				c1++
			}
			if p1 != v {
				c1++
			}
			c2 := 0
			if p0 != v {
				c2++
			}
			if p1 != u {
				c2++
			}
			if c2 < c1 {
				c1 = c2
			}
			costMat[k][e] = c1
		}
		for e := 0; e < E; e++ {
			cand[k][e] = e
		}
		sort.Slice(cand[k], func(i, j int) bool {
			return costMat[k][cand[k][i]] < costMat[k][cand[k][j]]
		})
	}
	order := make([]int, K)
	for i := 0; i < K; i++ {
		order[i] = i
	}
	type od struct{ k, c0, c1 int }
	ods := make([]od, K)
	for _, k := range order {
		c0, c1 := 0, 0
		for _, e := range cand[k] {
			c := costMat[k][e]
			if c == 0 {
				c0++
			} else if c == 1 {
				c1++
			} else {
				break
			}
		}
		ods[k] = od{k: k, c0: c0, c1: c1}
	}
	for i := 1; i < K; i++ {
		tmp := ods[i]
		j := i
		for j > 0 {
			o := ods[j-1]
			if tmp.c0 < o.c0 || (tmp.c0 == o.c0 && tmp.c1 < o.c1) {
				ods[j] = o
				j--
			} else {
				break
			}
		}
		ods[j] = tmp
	}
	for i := 0; i < K; i++ {
		order[i] = ods[i].k
	}
	minCost := make([]int, K)
	for idx, k := range order {
		minCost[idx] = costMat[k][cand[k][0]]
	}
	lbSum := make([]int, K+1)
	for i := K - 1; i >= 0; i-- {
		lbSum[i] = lbSum[i+1] + minCost[i]
	}
	best := 2 * K
	var dfs func(idx int, used uint64, cur int)
	dfs = func(idx int, used uint64, cur int) {
		if cur+lbSum[idx] >= best {
			return
		}
		if idx == K {
			if cur < best {
				best = cur
			}
			return
		}
		k := order[idx]
		for _, e := range cand[k] {
			c := costMat[k][e]
			if cur+c >= best {
				continue
			}
			u := edges[e][0]
			v := edges[e][1]
			mask := (uint64(1) << uint(u)) | (uint64(1) << uint(v))
			if used&mask != 0 {
				continue
			}
			dfs(idx+1, used|mask, cur+c)
		}
	}
	dfs(0, 0, 0)
	return fmt.Sprintf("%d", best)
}

func formatInput(inst instance) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", inst.n, inst.m))
	for i := 0; i < inst.n; i++ {
		for j := 0; j < inst.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", inst.grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func fixedTests() []test {
	cases := []instance{
		{n: 2, m: 2, grid: [][]int{{1, 1}, {2, 2}}},
		{n: 2, m: 2, grid: [][]int{{1, 2}, {1, 2}}},
		{n: 2, m: 4, grid: [][]int{{1, 1, 2, 2}, {3, 3, 4, 4}}},
		{n: 3, m: 4, grid: [][]int{{1, 2, 1, 2}, {3, 4, 3, 4}, {5, 6, 5, 6}}},
	}
	var tests []test
	for _, inst := range cases {
		tests = append(tests, test{
			input:    formatInput(inst),
			expected: solveInstance(inst),
		})
	}
	return tests
}

func randomInstance(rng *rand.Rand, nMin, nMax int) instance {
	var n, m int
	for {
		n = rng.Intn(nMax-nMin+1) + nMin
		m = rng.Intn(nMax-nMin+1) + nMin
		if (n*m)%2 == 0 {
			break
		}
	}
	N := n * m
	K := N / 2
	vals := make([]int, 0, N)
	for v := 1; v <= K; v++ {
		vals = append(vals, v, v)
	}
	rng.Shuffle(len(vals), func(i, j int) {
		vals[i], vals[j] = vals[j], vals[i]
	})
	grid := make([][]int, n)
	idx := 0
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = vals[idx]
			idx++
		}
	}
	return instance{n: n, m: m, grid: grid}
}

func randomTests(rng *rand.Rand, count, minSize, maxSize int) []test {
	tests := make([]test, 0, count)
	for len(tests) < count {
		inst := randomInstance(rng, minSize, maxSize)
		tests = append(tests, test{
			input:    formatInput(inst),
			expected: solveInstance(inst),
		})
	}
	return tests
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(316))
	tests := fixedTests()
	tests = append(tests, randomTests(rng, 40, 2, 4)...)
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\nInput:%s\n", i+1, err, t.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
