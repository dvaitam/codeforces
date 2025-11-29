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

const mod = 1000000007
const testcases = `100
2 6
4 1
2 3
1 2
2 1
5 4
4 3
1 2
4 2
4 2
5 1
4 2
4 1
1 3
3 3
4 2
4 12
1 3
3 1
1 2
5 2
1 11
1 2
2 12
1 2
6 3
2 10
2 3
1 3
1 8
6 1
2 10
4 3
5 2
3 9
3 2
2 1
5 1
5 12
3 3
3 3
4 3
5 1
3 1
3 9
3 3
6 1
5 3
5 4
3 1
1 2
3 1
3 1
4 2
4 5
6 2
2 3
4 1
2 3
1 4
2 2
1 1
4 2
2 10
4 1
2 2
5 8
4 2
1 3
1 2
6 1
1 1
6 2
5 3
3 1
5 2
6 1
5 1
2 1
2 6
4 3
5 1
6 5
3 2
3 1
6 2
2 1
5 1
5 1
2 6
4 2
1 1
2 10
2 2
4 1
4 2
1 3
6 2
6 1
6 2
5 12
1 3
4 3
4 2
4 1
1 1
1 8
1 1
2 12
1 3
5 3
1 11
3 2
5 2
6 3
5 3
5 1
5 1
3 3
4 12
5 3
2 1
3 3
3 1
6 6
4 1
6 2
1 3
4 3
3 3
3 3
3 11
3 2
6 1
2 2
3 7
5 2
2 2
4 3
4 11
5 3
6 2
6 3
5 2
6 7
2 2
5 3
1 2
2 1
4 3
3 2
1 10
3 2
3 9
3 2
6 2
2 2
2 4
5 3
4 1
2 10
5 1
3 2
1 2
4 1
6 3
5 2
1 1
1 3
2 3
3 1
1 2
4 11
3 3
1 1
5 2
6 1
3 1
1 1
6 2
1 1
6 1
5 2
2 2
6 2
3 1
5 1
6 3
3 7
2 1
3 2
2 2
6 2
1 3
2 1
3 1
3 3
2 1
1 1
5 12
2 1
5 3
3 2
2 3
5 2
1 10
6 1
2 12
2 1
3 1
4 3
2 1
4 1
6 2
5 3
5 5
4 2
5 3
6 2
6 3
6 2
3 12
5 1
3 2
6 2
2 5
4 3
5 3
5 5
2 3
5 2
2 2
6 2
1 2
4 1
1 2
5 2
3 2
6 1
5 3
6 1
1 1
4 3
2 2
5 3
5 2
3 3
1 3
6 3
3 2
1 2
6 2
4 1
6 1
4 2
5 2
3 2
3 3
6 11
1 3
2 3
6 3
3 1
1 1
5 2
4 6
4 3
2 2
2 3
3 1
4 12
3 3
6 2
4 2
3 2
4 8
6 3
1 2
6 2
3 1
5 12
1 2
6 3
6 2
3 2
5 2
2 9
5 3
5 1
2 9
1 1
5 3
3 4
2 3
4 2
2 1
1 1
4 3
2 6
3 1
1 3
1 12
2 3
3 11
3 1
1 1
3 2
1 7
2 1
6 2
5 3
2 2
6 2
1 2
5 1
1 1
6 12
6 3
3 3
4 3
1 2
4 3
6 1
1 5
1 2
3 12
2 1
1 2
6 2
3 5
5 2
5 1
4 1
6 12
6 3
2 1
1 3
4 1
6 2
2 3
1 7
3 3
5 11
4 2
5 2
5 1
5 1
4 1
2 9
3 3
2 3
2 4
4 3
6 3
2 6
6 1
5 2
1 8
3 3
2 11
2 3
5 2
1 6
4 2
3 2
1 2
6 3
4 2
5 1
4 2
5 3
2 1
1 3
5 3
5 6
3 3
4 3
1 3
2 1
1 2
3 2
5 1
2 1
6 1
1 6
3 2
6 6
4 3
2 2
6 2
6 1
4 2
5 1
5 11
6 1
5 3
6 2
4 3
4 3
2 2
5 3
5 3
2 6
6 2
5 3
4 10
1 2
5 1
3 1
6 2
3 9
1 3
4 3
3 1
2 9
6 3
4 2
3 12
2 2
5 1
6 1
6 5
1 1
4 1
5 1
2 1
6 2
3 3
`

type song struct {
	t int
	g int
}

type testCase struct {
	n     int
	T     int
	songs []song
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcases))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("empty testcases")
	}
	total, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("parse T: %w", err)
	}
	cases := make([]testCase, total)
	for i := 0; i < total; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("missing n for case %d", i+1)
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n case %d: %w", i+1, err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing T for case %d", i+1)
		}
		T, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse T case %d: %w", i+1, err)
		}
		songs := make([]song, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("missing t for case %d song %d", i+1, j+1)
			}
			t, err := strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("parse t case %d song %d: %w", i+1, j+1, err)
			}
			if !scan.Scan() {
				return nil, fmt.Errorf("missing g for case %d song %d", i+1, j+1)
			}
			g, err := strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("parse g case %d song %d: %w", i+1, j+1, err)
			}
			songs[j] = song{t: t, g: g}
		}
		cases[i] = testCase{n: n, T: T, songs: songs}
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return cases, nil
}

func referenceSolve(tc testCase) string {
	n, T := tc.n, tc.T
	_ = n
	d1 := []int{}
	d2 := []int{}
	d3 := []int{}
	for _, s := range tc.songs {
		switch s.g {
		case 1:
			d1 = append(d1, s.t)
		case 2:
			d2 = append(d2, s.t)
		case 3:
			d3 = append(d3, s.t)
		}
	}
	m1, m2, m3 := len(d1), len(d2), len(d3)
	dp1 := make([][]int, m1+1)
	for i := range dp1 {
		dp1[i] = make([]int, T+1)
	}
	dp1[0][0] = 1
	for _, t := range d1 {
		for k := m1; k >= 1; k-- {
			rowPrev := dp1[k-1]
			row := dp1[k]
			for d := T; d >= t; d-- {
				row[d] = (row[d] + rowPrev[d-t]) % mod
			}
		}
	}
	dp12 := make([][][]int, m1+1)
	for i := 0; i <= m1; i++ {
		dp12[i] = make([][]int, m2+1)
		for j := 0; j <= m2; j++ {
			dp12[i][j] = make([]int, T+1)
		}
		copy(dp12[i][0], dp1[i])
	}
	for _, t := range d2 {
		for j := m2; j >= 1; j-- {
			for i := 0; i <= m1; i++ {
				prev := dp12[i][j-1]
				cur := dp12[i][j]
				for d := T; d >= t; d-- {
					cur[d] = (cur[d] + prev[d-t]) % mod
				}
			}
		}
	}
	dp3 := make([][]int, m3+1)
	for i := range dp3 {
		dp3[i] = make([]int, T+1)
	}
	dp3[0][0] = 1
	for _, t := range d3 {
		for k := m3; k >= 1; k-- {
			prev := dp3[k-1]
			cur := dp3[k]
			for d := T; d >= t; d-- {
				cur[d] = (cur[d] + prev[d-t]) % mod
			}
		}
	}
	maxm := m1
	if m2 > maxm {
		maxm = m2
	}
	if m3 > maxm {
		maxm = m3
	}
	fact := make([]int, maxm+1)
	fact[0] = 1
	for i := 1; i <= maxm; i++ {
		fact[i] = fact[i-1] * i % mod
	}
	G := make([][][][]int, m1+1)
	for a := 0; a <= m1; a++ {
		G[a] = make([][][]int, m2+1)
		for b := 0; b <= m2; b++ {
			G[a][b] = make([][]int, m3+1)
			for c := 0; c <= m3; c++ {
				G[a][b][c] = make([]int, 4)
			}
		}
	}
	if m1 > 0 {
		G[1][0][0][1] = 1
	}
	if m2 > 0 {
		G[0][1][0][2] = 1
	}
	if m3 > 0 {
		G[0][0][1][3] = 1
	}
	for a := 0; a <= m1; a++ {
		for b := 0; b <= m2; b++ {
			for c := 0; c <= m3; c++ {
				if a+b+c <= 1 {
					continue
				}
				if a > 0 {
					v := 0
					if b > 0 {
						v = (v + G[a-1][b][c][2]) % mod
					}
					if c > 0 {
						v = (v + G[a-1][b][c][3]) % mod
					}
					G[a][b][c][1] = v
				}
				if b > 0 {
					v := 0
					if a > 0 {
						v = (v + G[a][b-1][c][1]) % mod
					}
					if c > 0 {
						v = (v + G[a][b-1][c][3]) % mod
					}
					G[a][b][c][2] = v
				}
				if c > 0 {
					v := 0
					if a > 0 {
						v = (v + G[a][b][c-1][1]) % mod
					}
					if b > 0 {
						v = (v + G[a][b][c-1][2]) % mod
					}
					G[a][b][c][3] = v
				}
			}
		}
	}
	H := make([][][]int, m1+1)
	for a := 0; a <= m1; a++ {
		H[a] = make([][]int, m2+1)
		for b := 0; b <= m2; b++ {
			H[a][b] = make([]int, m3+1)
			for c := 0; c <= m3; c++ {
				sum := (G[a][b][c][1] + G[a][b][c][2] + G[a][b][c][3]) % mod
				if sum == 0 {
					H[a][b][c] = 0
				} else {
					H[a][b][c] = int(int64(sum) * int64(fact[a]) % mod * int64(fact[b]) % mod * int64(fact[c]) % mod)
				}
			}
		}
	}
	var ans int64
	for a := 0; a <= m1; a++ {
		for b := 0; b <= m2; b++ {
			for c := 0; c <= m3; c++ {
				h := H[a][b][c]
				if h == 0 {
					continue
				}
				for d := 0; d <= T; d++ {
					cnt12 := dp12[a][b][d]
					if cnt12 == 0 {
						continue
					}
					d3 := T - d
					if d3 < 0 {
						continue
					}
					cnt3 := dp3[c][d3]
					if cnt3 == 0 {
						continue
					}
					ans = (ans + int64(cnt12)*int64(cnt3)%mod*int64(h)) % mod
				}
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.T))
		for _, s := range tc.songs {
			sb.WriteString(fmt.Sprintf("%d %d\n", s.t, s.g))
		}
		input := sb.String()

		want := referenceSolve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
