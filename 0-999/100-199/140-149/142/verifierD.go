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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// referenceSolve mirrors 142D.go, including its I/O parsing behaviour.
func referenceSolve(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return "", err
	}
	draw := false
	rows := make([][]byte, n)
	for i := 0; i < n; i++ {
		line := make([]byte, 0, m)
		for len(line) < m {
			buf := make([]byte, m)
			cnt, err := reader.Read(buf)
			if err != nil {
				break
			}
			for j := 0; j < cnt && len(line) < m; j++ {
				if buf[j] != '\n' && buf[j] != '\r' {
					line = append(line, buf[j])
				}
			}
		}
		rows[i] = line
		cnt := 0
		for _, c := range line {
			if c == 'G' || c == 'R' {
				cnt++
			}
		}
		if cnt <= 1 {
			draw = true
		}
	}
	if draw || k >= 2 {
		return "Draw", nil
	}
	xor := 0
	for i := 0; i < n; i++ {
		gpos, rpos := -1, -1
		for j, c := range rows[i] {
			if c == 'G' {
				gpos = j
			} else if c == 'R' {
				rpos = j
			}
		}
		if gpos >= 0 && rpos >= 0 {
			d := abs(gpos-rpos) - 1
			xor ^= d
		}
	}
	if xor != 0 {
		return "First", nil
	}
	return "Second", nil
}

type testCase struct {
	n, m, k int
	rows    []string
}

// Embedded testcases from testcasesD.txt (n m k board).
const testcaseData = `
7 2 1 ----RG---RR---
6 4 2 -G-R--------R-G--------R
4 1 1 ---R
3 2 3 R-RRG-
1 6 3 --G--G
5 3 2 -GGG-G---R--G--
8 6 3 ----G-G---R-----------R----G--------RR----R-----
5 5 2 ---G-------R--R----------
8 7 1 ------G-R---R-----------R--G-------------G--------------
3 6 1 --R-G-G--------G-R
5 3 2 --GR-RR----G---
1 8 2 --------
4 4 1 --GR-------G----
7 2 1 G--R---R--G-R-
6 3 1 ---G-----------R-G
6 1 3 ---RG-
4 6 2 --R-----G--G--------R-G-
3 3 1 --G------
5 6 2 GR----------RR-----R-R--R-G---
6 1 2 RGR---
5 2 3 RG-G---G-G
8 4 2 ------GR----R----R-GR--RR--G----
5 7 3 ---G-R---R---GG-----R----R---------
4 3 1 -R-G-R-RR-G-
3 2 3 ---R--
5 6 3 -----R---------R-------G--G---
8 1 3 --R---R-
2 1 2 GG
7 8 3 ------G-------R-----------R-----------------------------
8 6 2 ---R---R---GG----R------RG-----R---GR---R----G--
7 8 2 -------G----G---------------------------R-------GR------
4 7 3 --------R-----G---R------R--
6 5 1 RG---------RG----R-R--G-------
5 2 1 -R---RG---
4 2 1 R-RG----
8 4 1 --R-R----RR---RR------------R--G
4 5 3 G--G---G----R-G-----
4 2 3 -G-RGRG-
7 3 1 G-RG-G----G-RG--R----
4 3 3 --RGR-RR-G--
1 6 3 -RG---
5 6 3 ------R-----R---G------------G
1 4 1 ----
4 6 1 -----------R------------
8 7 1 --G------R--R----------R-----------R--G-----------G--G--
6 4 1 -------------GR-R-G-RR--
1 1 2 -
6 5 3 ------R----------RR-R------GR-
1 8 2 -R-R----
1 1 2 -
5 5 3 -------R-R-----------GR--
3 6 3 ------G---R---R---
6 8 2 -G----G-G--R-----------------G-R----------G---G-
2 3 3 ----RR
5 5 1 --R----G-R-G-G----RR-----
7 4 1 -G----RG-GR-----R---------GG
3 8 1 ------R-------R-G----R--
8 1 2 -RRRRG-G
2 4 3 ----G--R
6 7 1 --GG---------RGG----------G-----R-G-------
6 2 1 GR-RGR--RG--
8 4 2 -RG-G---G-G---R----RR--G----GR--
1 7 3 -------
5 7 1 -----RG-------R----------R----R-R--
3 1 2 R-G
1 6 1 ------
7 1 2 R--G--G
8 1 1 GG-----R
7 1 1 -G-GR-G
1 7 3 -------
3 1 1 G--
3 3 3 R----RR-G
7 2 1 -G-G-GR-R--GGR
6 7 3 -------R-----G------------R-R-------R-G---
8 7 1 --R---G-R----------G---------R--------G-------R---------
7 5 2 -GG--------RG-----GG----G----------
4 3 1 ---RR-RR----
6 6 3 -R-G----R-------R----G---G---G-R---G
3 7 1 ---------------------
6 8 2 ----------R----R-----------------G----G---R-G---
8 7 2 --------G-R---G--R----R----------------G--------R---R-R-
3 4 3 -R---G-G---R
2 5 2 ----------
5 3 2 ------G-G-----R
2 8 3 -----G--GR------
5 5 1 -------R----------------R
4 6 2 -------------R------R--G
2 7 2 ---R--G--RR---
7 5 1 --G---R---------RG----------RR-R--R
8 4 2 R--R-R------GR------G-G-G-G-----
7 1 1 ---GG--
7 7 2 --------R------------R-R---------R--G--R---------
2 5 1 -R--G-GR--
2 8 1 ---R---G-R---R--
7 2 3 GR-R--R-R----R
3 4 2 --R--GG-G---
3 6 1 ----G----G-----G--
3 7 1 ------G-GR--------GR-
8 1 1 G-----GR
8 4 1 --R----RG------------------GG---
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			return nil, fmt.Errorf("line %d: invalid testcase", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m: %v", i+1, err)
		}
		k, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k: %v", i+1, err)
		}
		board := fields[3]
		if len(board) != n*m {
			return nil, fmt.Errorf("line %d: board length %d != %d", i+1, len(board), n*m)
		}
		rows := make([]string, n)
		for r := 0; r < n; r++ {
			rows[r] = board[r*m : (r+1)*m]
		}
		cases = append(cases, testCase{n: n, m: m, k: k, rows: rows})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
		for _, row := range tc.rows {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		expected, err := referenceSolve(sb.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, sb.String(), expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
