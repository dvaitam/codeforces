package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int = 1000000007

type State struct {
	x, y int
	grid []string
}

type testcase struct {
	n    int
	m    int
	grid []string
}

const testcasesRaw = `4 2 RR RR .. .R
3 3 .RR RR. ...
2 3 RRR RRR
4 1 R R . .
3 2 .. .R RR
3 2 R. .R R.
3 1 . . .
4 4 .R.. .R.R .... R...
3 1 R . R
4 1 . R . R
3 1 . . .
2 3 .R. ..R
2 3 ..R ...
3 4 R..R ..RR ..R.
3 3 RRR RR. ...
3 2 .. R. RR
1 1 R
3 3 ..R .RR R..
4 4 .... .... .RR. RRR.
4 1 . . . .
1 2 RR
4 3 ..R .RR R.R ..R
3 3 ... RRR RR.
1 2 R.
4 3 ..R .R. RR. R.R
4 2 .. .. .R .R
4 4 RRRR R..R .RRR ...R
1 4 ....
4 2 RR R. RR .R
1 3 RR.
2 4 .R.. ..R.
4 4 R..R ..RR .... ....
4 1 R R R R
2 3 RR. RR.
1 1 .
4 1 . R R .
3 4 .R.. .RR. RRR.
3 2 .. .R .R
1 1 .
4 1 . . . .
1 1 .
4 1 R R R .
2 2 .. .R
1 1 .
2 1 R R
3 2 .. .R .R
2 1 . R
1 4 R..R
3 3 .RR RRR ...
1 3 RRR
3 1 . R R
2 3 RRR ...
1 4 .RRR
1 3 R..
4 1 R R . R
1 1 .
4 2 RR .. .. RR
3 2 R. .. ..
3 4 R... R... R.RR
2 3 R.. .R.
2 2 .R .R
4 2 R. .R .. RR
3 1 R . R
1 4 R...
1 1 .
2 3 .R. ..R
2 2 .. R.
2 1 R .
4 4 .R.. RRRR R... RR..
4 1 R . . R
3 3 .RR ... RR.
2 1 R .
3 1 . . R
3 3 ..R .R. .R.
1 4 ..RR
2 1 R R
2 1 . R
1 2 RR
2 4 RRRR .R..
2 1 R .
1 3 .R.
2 3 .R. R.R
4 1 . R R R
4 4 RRR. RR.. .RR. .RR.
1 1 R
2 4 RR.. RRRR
4 1 . . R .
3 2 R. R. .R
1 3 ..R
4 2 RR RR RR ..
1 2 R.
2 1 . R
1 2 .R
4 2 .R .R .. R.
1 4 .R..
4 3 RR. ... R.R ...
1 1 .
4 4 .RRR RRR. .RRR ....
3 2 .R R. .R
3 2 .. .R ..`

var testcases = mustParseTestcases(testcasesRaw)

func serialize(g []string) string {
	s := make([]byte, 0, len(g)*(len(g[0])+1))
	for _, row := range g {
		s = append(s, row...)
		s = append(s, '\n')
	}
	return string(s)
}

func countPaths(g []string) int {
	n := len(g)
	m := len(g[0])
	start := State{0, 0, append([]string(nil), g...)}
	q := []State{start}
	ways := map[string]int{serialize(start.grid) + fmt.Sprintf("%d,%d", 0, 0): 1}
	res := 0

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		key := serialize(cur.grid) + fmt.Sprintf("%d,%d", cur.x, cur.y)
		w := ways[key]
		if cur.x == n-1 && cur.y == m-1 {
			res = (res + w) % MOD
			continue
		}
		for _, d := range [][2]int{{1, 0}, {0, 1}} {
			nx := cur.x + d[0]
			ny := cur.y + d[1]
			if nx >= n || ny >= m {
				continue
			}
			ng := make([]string, n)
			copy(ng, cur.grid)
			if d[0] == 1 {
				if ng[nx][ny] == 'R' {
					k := nx
					for k < n && ng[k][ny] == 'R' {
						k++
					}
					if k == n {
						continue
					}
					b := []byte(ng[k])
					b[ny] = 'R'
					ng[k] = string(b)
					for t := k - 1; t >= nx; t-- {
						row := []byte(ng[t])
						row[ny] = '.'
						ng[t] = string(row)
					}
				}
			} else {
				if ng[nx][ny] == 'R' {
					k := ny
					for k < m && ng[nx][k] == 'R' {
						k++
					}
					if k == m {
						continue
					}
					row := []byte(ng[nx])
					row[k] = 'R'
					for t := k - 1; t >= ny; t-- {
						row[t] = '.'
					}
					ng[nx] = string(row)
				}
			}
			nk := serialize(ng) + fmt.Sprintf("%d,%d", nx, ny)
			if _, ok := ways[nk]; !ok {
				q = append(q, State{nx, ny, ng})
			}
			ways[nk] = (ways[nk] + w) % MOD
		}
	}
	return res
}

func mustParseTestcases(raw string) []testcase {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	res := make([]testcase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			panic(fmt.Sprintf("line %d too short", idx+1))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			panic(fmt.Sprintf("line %d bad n: %v", idx+1, err))
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			panic(fmt.Sprintf("line %d bad m: %v", idx+1, err))
		}
		if len(fields) != 2+n {
			panic(fmt.Sprintf("line %d expected %d grid rows got %d", idx+1, n, len(fields)-2))
		}
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			row := fields[2+i]
			if len(row) != m {
				panic(fmt.Sprintf("line %d row %d length mismatch", idx+1, i+1))
			}
			grid[i] = row
		}
		res = append(res, testcase{n: n, m: m, grid: grid})
	}
	return res
}

func solve(tc testcase) int {
	if tc.grid[0][0] == 'R' {
		return 0
	}
	return countPaths(tc.grid)
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseCandidateOutput(out string) (int, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return 0, fmt.Errorf("no output")
	}
	v, err := strconv.Atoi(sc.Text())
	if err != nil {
		return 0, fmt.Errorf("failed to parse output: %v", err)
	}
	if err := sc.Err(); err != nil {
		return 0, fmt.Errorf("scanner error: %v", err)
	}
	return v, nil
}

func checkCase(bin string, idx int, tc testcase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	input := sb.String()
	expected := solve(tc)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got, err := parseCandidateOutput(out)
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("case %d: expected %d got %d", idx+1, expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, i, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
