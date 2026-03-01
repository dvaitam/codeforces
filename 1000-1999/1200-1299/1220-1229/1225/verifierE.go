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

func countPaths(g []string) int {
	n := len(g)
	m := len(g[0])
	if n == 1 && m == 1 {
		return 1
	}
	D := make([][]int, n+2)
	R := make([][]int, n+2)
	f := make([][]int, n+2)
	gdp := make([][]int, n+2)
	for i := 0; i < n+2; i++ {
		D[i] = make([]int, m+2)
		R[i] = make([]int, m+2)
		f[i] = make([]int, m+2)
		gdp[i] = make([]int, m+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if g[i-1][j-1] == 'R' {
				D[i][j] = 1
				R[i][j] = 1
			}
		}
	}
	for i := n; i >= 1; i-- {
		for j := m; j >= 1; j-- {
			D[i][j] += D[i+1][j]
			R[i][j] += R[i][j+1]
		}
	}
	f[1][1], gdp[1][1] = 1, 1
	if n >= 2 {
		f[2][1] = MOD - 1
	}
	if m >= 2 {
		gdp[1][2] = MOD - 1
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			f[i][j] = (f[i][j] + f[i-1][j] + gdp[i-1][j]) % MOD
			tgtCol := m - R[i][j+1] + 1
			gdp[i][tgtCol] = (gdp[i][tgtCol] - f[i][j] + MOD) % MOD
			gdp[i][j] = (gdp[i][j] + gdp[i][j-1] + f[i][j-1]) % MOD
			tgtRow := n - D[i+1][j] + 1
			f[tgtRow][j] = (f[tgtRow][j] - gdp[i][j] + MOD) % MOD
		}
	}
	return (f[n][m] + gdp[n][m]) % MOD
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
