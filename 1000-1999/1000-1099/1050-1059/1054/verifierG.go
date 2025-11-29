package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `2 2 10 01
3 2 000 110
3 3 001 010 111
1 1 0
3 2 011 001
2 3 00 01 10
2 2 10 01
1 3 0 1 1
2 2 00 01
1 1 1
1 1 1
1 1 0
1 2 1 1
3 1 011
3 1 011
2 1 10
3 2 100 011
1 2 1 0
1 1 1
1 2 0 1
3 3 110 001 011
3 1 011
3 1 110
1 2 1 1
3 2 010 111
3 3 110 101 111
3 1 110
2 3 00 10 00
3 3 111 001 111
1 3 0 0 0
2 3 00 11 00
2 1 01
3 1 001
1 2 0 0
1 3 0 0 0
2 2 10 00
3 3 000 110 110
3 1 011
3 3 001 000 100
2 2 11 00
2 2 00 01
3 1 110
1 1 0
1 1 1
1 3 0 1 1
2 1 01
1 2 1 0
1 3 0 1 1
2 2 00 11
2 2 00 11
1 1 0
2 2 00 00
3 1 110
3 2 010 101
1 2 1 1
3 3 010 111 100
2 2 11 00
1 2 0 0
1 1 1
2 3 00 01 00
2 2 11 11
2 2 01 01
1 2 1 1
1 1 0
3 2 100 110
1 3 1 0 0
2 2 00 00
1 3 0 1 1
1 2 1 0
2 3 01 01 01
3 3 011 101 001
2 3 00 00 10
2 1 11
3 2 000 001
1 2 0 0
3 1 111
2 3 11 00 01
2 2 11 00
1 3 1 1 1
1 3 0 0 1
1 3 0 1 0
2 3 01 10 01
2 1 11
1 2 1 1
2 2 11 10
1 2 1 0
1 2 0 1
2 1 10
3 3 110 001 001
1 1 0
1 3 1 1 1
1 1 1
1 1 1
3 3 001 110 011
1 1 1
2 3 11 11 10
3 2 011 110
2 2 10 11
2 3 10 01 11
2 3 10 11 10`

type testCase struct {
	n    int
	m    int
	cols []string
}

type pair struct{ v, u int }

// embedded solver logic from 1054G.go
func solve(tc testCase) (string, bool) {
	n, m := tc.n, tc.m
	k := (m + 63) / 64
	inSet := make([][]uint64, n)
	notInSet := make([][]uint64, n)
	for i := 0; i < n; i++ {
		inSet[i] = make([]uint64, k)
		notInSet[i] = make([]uint64, k)
	}
	needCheck := make([]bool, n)
	for i := range needCheck {
		needCheck[i] = true
	}
	degForSet := make([]int, m)
	G := make([][]bool, n)
	for i := 0; i < n; i++ {
		G[i] = make([]bool, m)
	}
	for j := 0; j < m; j++ {
		s := tc.cols[j]
		for i := 0; i < n && i < len(s); i++ {
			if s[i] == '1' {
				G[i][j] = true
				degForSet[j]++
			}
		}
		if degForSet[j] <= 1 {
			continue
		}
		x := j >> 6
		z := uint64(1) << (uint(j) & 63)
		for i := 0; i < n; i++ {
			if G[i][j] {
				inSet[i][x] |= z
			} else {
				notInSet[i][x] |= z
			}
		}
	}
	alive := make([]int, n)
	for i := 0; i < n; i++ {
		alive[i] = i
	}
	ans := make([]pair, 0, n)
	for len(alive) > 2 {
		fnd := false
		for idx := 0; idx < len(alive); idx++ {
			v := alive[idx]
			if !needCheck[v] {
				continue
			}
			needCheck[v] = false
			p := -1
			for _, u := range alive {
				if u == v {
					continue
				}
				ok := true
				for b := 0; b < k; b++ {
					if inSet[v][b]&notInSet[u][b] != 0 {
						ok = false
						break
					}
				}
				if ok {
					p = u
					break
				}
			}
			if p < 0 {
				continue
			}
			fnd = true
			ans = append(ans, pair{v, p})
			for j := 0; j < m; j++ {
				if !G[v][j] {
					continue
				}
				degForSet[j]--
				if degForSet[j] == 1 {
					id := j
					x := id >> 6
					z := uint64(1) << (uint(id) & 63)
					for i := 0; i < n; i++ {
						if G[i][id] {
							inSet[i][x] ^= z
							needCheck[i] = true
						}
					}
				}
			}
			alive[idx] = alive[len(alive)-1]
			alive = alive[:len(alive)-1]
			break
		}
		if !fnd {
			return "NO", false
		}
	}
	if len(alive) == 2 {
		ans = append(ans, pair{alive[0], alive[1]})
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for _, p := range ans {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.v+1, p.u+1))
	}
	return sb.String(), true
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d", tc.n, tc.m))
	for _, s := range tc.cols {
		sb.WriteByte(' ')
		sb.WriteString(s)
	}
	sb.WriteByte('\n')
	input := sb.String()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	expected, _ := solve(tc)
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(expected)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		if len(fields) != 2+m {
			return nil, fmt.Errorf("invalid testcase, expected %d strings got %d", m, len(fields)-2)
		}
		tc := testCase{n: n, m: m, cols: make([]string, m)}
		copy(tc.cols, fields[2:])
		tests = append(tests, tc)
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
