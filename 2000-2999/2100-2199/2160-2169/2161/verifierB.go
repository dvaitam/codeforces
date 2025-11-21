package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
	t     int
}

// reference solver from 2161B.go
type block struct {
	l      int
	length int
	cntL   int
	cntR   int
}

func (b block) cols() []int {
	if b.length == 1 {
		return []int{b.l}
	}
	return []int{b.l, b.l + 1}
}

func (b block) count(col int) (int, bool) {
	if b.length == 1 {
		if col == b.l {
			return b.cntL, true
		}
		return 0, false
	}
	if col == b.l {
		return b.cntL, true
	}
	if col == b.l+1 {
		return b.cntR, true
	}
	return 0, false
}

func makeBlock(l, length int, counts map[int]int) block {
	if length == 1 {
		return block{l: l, length: 1, cntL: counts[l]}
	}
	return block{l: l, length: 2, cntL: counts[l], cntR: counts[l+1]}
}

func genCandidates(n int, req []int, prev *block) [][2]int {
	res := make([][2]int, 0)
	seen := make(map[[2]int]struct{})
	add := func(l, length int) {
		if l < 0 || l >= n {
			return
		}
		if length == 2 && l+1 >= n {
			return
		}
		cols := [2]int{l, l}
		if length == 2 {
			cols[1] = l + 1
		}
		ok := true
		for _, r := range req {
			if length == 1 {
				if r != l {
					ok = false
					break
				}
			} else {
				if r != l && r != l+1 {
					ok = false
					break
				}
			}
		}
		if !ok {
			return
		}
		if prev != nil {
			intersects := false
			prevCols := prev.cols()
			for _, pc := range prevCols {
				if pc == cols[0] || (length == 2 && pc == cols[1]) {
					intersects = true
					break
				}
			}
			if !intersects {
				return
			}
		}
		key := [2]int{l, length}
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		res = append(res, key)
	}

	switch len(req) {
	case 2:
		c := req[0]
		if req[1] != c+1 {
			return nil
		}
		add(c, 2)
	case 1:
		c := req[0]
		add(c, 1)
		add(c, 2)
		add(c-1, 2)
	case 0:
		if prev == nil {
			return nil
		}
		prevCols := prev.cols()
		for _, c := range prevCols {
			add(c, 1)
			add(c, 2)
			add(c-1, 2)
		}
	}
	return res
}

func solveGrid(grid []string) bool {
	n := len(grid)
	rows := make([][]int, n)
	cols := make([][]int, n)
	first, last := -1, -1
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '#' {
				rows[i] = append(rows[i], j)
				cols[j] = append(cols[j], i)
			}
		}
		if len(rows[i]) > 0 {
			if first == -1 {
				first = i
			}
			last = i
		}
		if len(rows[i]) > 2 {
			return false
		}
		if len(rows[i]) == 2 && rows[i][1]-rows[i][0] != 1 {
			return false
		}
	}
	if first == -1 {
		return true
	}
	for j := 0; j < n; j++ {
		if len(cols[j]) > 2 {
			return false
		}
		if len(cols[j]) == 2 && cols[j][1]-cols[j][0] != 1 {
			return false
		}
	}

	reqFirst := rows[first]
	initial := make([]block, 0)
	for _, cand := range genCandidates(n, reqFirst, nil) {
		counts := map[int]int{}
		if cand[1] == 1 {
			counts[cand[0]] = 1
		} else {
			counts[cand[0]] = 1
			counts[cand[0]+1] = 1
		}
		initial = append(initial, makeBlock(cand[0], cand[1], counts))
	}
	if len(initial) == 0 {
		return false
	}
	dp := initial
	for i := first + 1; i <= last; i++ {
		req := rows[i]
		nextMap := make(map[block]struct{})
		for _, prev := range dp {
			prevCopy := prev
			for _, cand := range genCandidates(n, req, &prevCopy) {
				counts := make(map[int]int)
				ok := true
				if cand[1] == 1 {
					col := cand[0]
					if val, has := prevCopy.count(col); has {
						if val == 2 {
							ok = false
						} else {
							counts[col] = val + 1
						}
					} else {
						counts[col] = 1
					}
				} else {
					leftCol := cand[0]
					rightCol := cand[0] + 1
					if val, has := prevCopy.count(leftCol); has {
						if val == 2 {
							ok = false
						} else {
							counts[leftCol] = val + 1
						}
					} else {
						counts[leftCol] = 1
					}
					if val, has := prevCopy.count(rightCol); has {
						if val == 2 {
							ok = false
						} else {
							counts[rightCol] = val + 1
						}
					} else {
						counts[rightCol] = 1
					}
				}
				if !ok {
					continue
				}
				nb := makeBlock(cand[0], cand[1], counts)
				nextMap[nb] = struct{}{}
			}
		}
		if len(nextMap) == 0 {
			return false
		}
		dp = dp[:0]
		for b := range nextMap {
			dp = append(dp, b)
		}
	}
	return true
}

func solveRef(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return "", err
	}
	res := make([]string, t)
	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(reader, &n)
		grid := make([]string, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &grid[j])
		}
		if solveGrid(grid) {
			res[i] = "YES"
		} else {
			res[i] = "NO"
		}
	}
	return strings.Join(res, "\n"), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func makeCase(name string, grids []struct {
	n    int
	rows []string
}) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(grids))
	for _, g := range grids {
		fmt.Fprintf(&sb, "%d\n", g.n)
		for _, row := range g.rows {
			fmt.Fprintf(&sb, "%s\n", row)
		}
	}
	return testCase{name: name, input: sb.String(), t: len(grids)}
}

func randomGrid(n int, rng *rand.Rand) []string {
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if rng.Intn(4) == 0 {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		rows[i] = sb.String()
	}
	return rows
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for idx := 0; idx < 40; idx++ {
		tcCount := rng.Intn(3) + 1
		grids := make([]struct {
			n    int
			rows []string
		}, tcCount)
		for i := 0; i < tcCount; i++ {
			n := rng.Intn(5) + 1
			grids[i].n = n
			grids[i].rows = randomGrid(n, rng)
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", idx+1), grids))
	}
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("empty", []struct {
			n    int
			rows []string
		}{
			{n: 1, rows: []string{"."}},
		}),
		makeCase("single_row", []struct {
			n    int
			rows []string
		}{
			{n: 3, rows: []string{"#..", ".#.", "..#"}},
		}),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expLines := strings.Fields(expect)
		gotLines := strings.Fields(out)
		if len(expLines) != len(gotLines) {
			fmt.Printf("test %d (%s) mismatch in outputs count\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", idx+1, tc.name, tc.input, expect, out)
			os.Exit(1)
		}
		for i := range expLines {
			if strings.ToUpper(expLines[i]) != strings.ToUpper(gotLines[i]) {
				fmt.Printf("test %d (%s) mismatch at testcase %d\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", idx+1, tc.name, i+1, tc.input, expect, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
