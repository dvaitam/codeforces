package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `
100
0 1 0
3.88 2.92 1.36
1 7
2 1 0
1.89 3.15 2.11
3 2
1 0 0
4.32 4.23 4.2
4 3
1 1 2
4.47 4.84 4.39
6 2
2 1 2
2.55 2.0 1.99
5 2
2 1 0
4.64 3.29 4.53
9 4
1 1 2
2.15 2.81 1.93
5 1
0 0 1
3.51 2.12 3.14
8 6
0 2 0
1.27 4.65 3.54
8 5
0 1 1
3.99 2.28 3.23
6 2
0 2 0
2.11 3.33 4.45
2 6
0 1 1
1.1 2.43 1.33
5 6
0 1 1
2.29 1.61 3.61
10 2
1 2 0
4.58 2.17 2.0
10 3
1 2 0
2.45 2.82 2.46
6 5
2 0 1
4.93 2.7 1.83
1 1
0 2 0
3.38 4.88 3.43
9 8
2 0 1
4.99 1.49 3.12
7 4
1 0 0
2.76 2.97 1.88
8 4
2 1 0
2.99 1.13 2.02
4 9
0 0 1
4.48 1.57 1.21
6 10
0 2 1
4.85 3.61 4.48
1 8
1 0 1
1.84 4.5 4.6
3 6
1 2 1
4.2 3.57 4.26
9 4
2 2 1
2.35 2.57 2.99
5 4
0 1 2
1.51 4.87 3.67
1 3
2 2 1
3.28 3.99 4.71
4 1
0 0 0
3.44 2.03 2.58
7 4
2 0 0
1.65 3.43 4.27
9 8
2 1 0
1.32 3.78 1.45
9 5
2 0 0
2.45 4.07 3.09
1 5
1 0 0
4.84 2.82 1.82
5 7
0 1 1
1.38 1.46 3.48
6 9
1 1 2
4.11 1.27 4.53
4 5
1 1 0
4.19 1.67 4.5
3 3
1 1 1
2.03 4.9 1.02
3 1
2 1 0
3.18 2.95 3.86
8 9
0 2 0
2.64 2.17 1.92
3 1
2 0 2
4.63 3.18 4.64
10 5
2 1 2
3.47 2.57 2.01
10 6
2 1 0
2.73 3.42 3.7
10 3
1 1 0
3.3 2.4 3.48
2 2
1 2 0
2.32 2.49 1.7
1 10
0 2 0
4.2 4.23 4.81
3 3
2 1 2
3.3 4.75 4.04
2 3
2 1 2
4.55 1.89 4.15
5 7
2 0 1
3.81 2.24 1.92
6 9
2 2 1
4.6 2.6 2.6
6 5
1 1 2
1.05 1.74 3.16
9 10
1 1 1
3.49 1.62 1.27
6 10
1 0 1
2.93 3.56 2.91
10 2
0 0 0
2.2 4.79 1.65
10 10
2 2 2
2.05 3.1 1.63
2 4
0 1 0
1.49 1.25 4.97
5 10
2 2 0
3.93 3.62 4.81
2 9
1 2 0
3.89 3.21 3.01
3 4
1 1 2
1.27 1.67 4.5
5 9
1 2 2
2.21 2.6 1.68
1 7
0 1 2
1.08 4.84 1.61
3 2
2 0 0
1.93 1.03 4.64
9 3
1 2 2
4.48 3.73 3.79
6 3
2 2 0
1.43 3.48 4.65
5 7
2 2 0
4.28 1.85 3.32
10 5
1 2 1
4.23 1.33 1.94
6 2
2 2 0
3.42 3.11 4.95
6 2
0 2 0
2.99 3.18 3.39
4 1
0 2 0
2.35 3.41 1.55
5 2
2 2 0
1.29 1.9 2.1
1 1
0 1 2
1.58 2.44 4.56
6 5
1 2 2
2.61 1.69 2.97
2 9
2 1 2
3.76 4.16 3.82
7 6
0 0 0
2.91 2.75 4.89
6 7
2 1 1
4.0 1.84 2.12
2 1
1 1 2
2.15 4.66 1.53
6 2
0 2 0
3.8 4.74 2.01
4 8
0 2 1
1.03 3.28 1.81
2 4
2 0 0
4.19 2.26 4.48
10 3
2 0 2
4.39 1.72 3.87
1 3
2 0 1
1.14 2.3 2.31
10 6
0 0 0
1.2 1.21 1.64
1 5
2 2 2
4.37 1.42 1.8
3 7
1 2 2
2.98 2.29 2.72
7 1
0 0 0
2.43 2.19 2.45
4 10
2 1 1
2.45 3.83 4.78
4 6
0 2 0
1.08 1.47 2.19
6 9
1 0 1
4.2 4.21 1.65
7 10
1 1 1
1.02 3.79 2.79
7 4
0 1 1
3.39 1.22 1.31
3 9
0 0 2
1.96 2.33 1.08
1 7
1 0 2
1.55 1.17 1.73
5 4
1 2 1
3.16 2.49 3.92
9 10
1 1 2
2.34 3.15 3.18
6 10
0 1 1
4.31 1.3 3.51
4 3
0 2 2
1.56 4.35 3.15
6 2
2 0 2
1.6 1.7 3.5
5 8
1 0 2
2.89 2.13 3.76
1 3
1 1 0
1.01 3.12 4.41
6 2
0 1 2
1.32 2.7 4.65
5 4
2 0 0
1.81 4.67 4.69
2 1
0 1 2
3.07 2.7 1.37
6 3
`

type testCase struct {
	nf, ne, ns    int
	rfi, rei, rsi float64
	df, de        int
}

func Len(a1, b1, c1, d1 float64) float64 {
	lo := math.Max(a1, c1)
	hi := math.Min(b1, d1)
	if hi > lo {
		return hi - lo
	}
	return 0
}

func parseTestcases(raw string) ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Split(bufio.ScanWords)
	scanInt := func() (int, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(sc.Text())
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	scanFloat := func() (float64, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.ParseFloat(sc.Text(), 64)
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	t, err := scanInt()
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		nf, err := scanInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: nf: %w", i+1, err)
		}
		ne, err := scanInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: ne: %w", i+1, err)
		}
		ns, err := scanInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: ns: %w", i+1, err)
		}
		rfi, err := scanFloat()
		if err != nil {
			return nil, fmt.Errorf("case %d: rfi: %w", i+1, err)
		}
		rei, err := scanFloat()
		if err != nil {
			return nil, fmt.Errorf("case %d: rei: %w", i+1, err)
		}
		rsi, err := scanFloat()
		if err != nil {
			return nil, fmt.Errorf("case %d: rsi: %w", i+1, err)
		}
		df, err := scanInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: df: %w", i+1, err)
		}
		de, err := scanInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: de: %w", i+1, err)
		}
		tests = append(tests, testCase{nf: nf, ne: ne, ns: ns, rfi: rfi, rei: rei, rsi: rsi, df: df, de: de})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}
	return tests, nil
}

func solve(tc testCase) float64 {
	nf, ne, ns := tc.nf, tc.ne, tc.ns
	rf := math.Sqrt(tc.rfi*tc.rfi - 1)
	re := math.Sqrt(tc.rei*tc.rei - 1)
	rs := math.Sqrt(tc.rsi*tc.rsi - 1)
	total := nf + ne + ns
	U := make([]bool, total)
	a := make([]int, ns)
	b := make([]float64, total)
	ans := 0.0

	calc := func() float64 {
		Fc := 2*float64(nf)*rf*float64(tc.df) + 2*float64(ne)*re*float64(tc.de)
		m := 0
		for i := 0; i < total; i++ {
			if !U[i] {
				xi := float64(i) / 2.0
				var Df, DeF float64
				for j := 0; j < ns; j++ {
					Df += float64(tc.df) * Len(xi-rf, xi+rf, float64(a[j])-rs, float64(a[j])+rs)
					DeF += float64(tc.de) * Len(xi-re, xi+re, float64(a[j])-rs, float64(a[j])+rs)
				}
				Fc += Df
				b[m] = DeF - Df
				m++
			}
		}
		if m > 0 {
			sort.Slice(b[:m], func(i, j int) bool { return b[i] > b[j] })
			limit := ne
			if limit > m {
				limit = m
			}
			for i := 0; i < limit; i++ {
				Fc += b[i]
			}
		}
		return Fc
	}

	var dfs func(x, y int)
	dfs = func(x, y int) {
		if nf+ne+y < x {
			return
		}
		if x == total {
			if val := calc(); val > ans {
				ans = val
			}
			return
		}
		U[x] = false
		dfs(x+1, y)
		if y < ns && (x%2 == 0 || U[x-1]) {
			U[x] = true
			a[y] = x / 2
			dfs(x+1, y+1)
		}
	}

	dfs(0, 0)
	return ans
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.nf, tc.ne, tc.ns))
	sb.WriteString(fmt.Sprintf("%.15g %.15g %.15g\n", tc.rfi, tc.rei, tc.rsi))
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.df, tc.de))
	return sb.String()
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		expected := solve(tc)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseFloat(strings.TrimSpace(gotStr), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if math.Abs(got-expected) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %.10f got %.10f\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
