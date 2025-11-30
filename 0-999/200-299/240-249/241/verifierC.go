package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const (
	MOD = 1000000007
	H   = 100.0
	Lx  = 105.0
)

// Embedded testcases from testcasesC.txt.
const embeddedTestcasesC = `100
19 1 1 19 T 23 53
5 72 3 20 F 1 25 17 F 43 56 12 T 1 48
32 74 1 9 F 48 75
10 74 3 8 T 32 78 4 F 10 39 3 T 40 66
35 33 3 12 T 5 25 1 T 0 49 9 F 48 74
20 24 2 7 F 16 19 9 T 74 89
67 49 3 14 F 20 30 11 F 14 61 12 F 20 25
74 49 2 17 T 12 74 15 F 64 96
84 33 1 2 F 26 32
24 22 3 13 T 65 85 19 T 11 72 1 T 72 83
15 4 1 9 F 2 100
37 69 3 5 T 2 100 10 T 67 99 19 F 49 78
73 12 3 5 F 81 92 9 T 7 23 13 F 89 100
64 52 1 3 T 0 100
32 29 1 4 T 6 69
67 25 1 5 F 39 95
81 27 2 6 F 64 70 2 F 33 82
26 70 2 3 F 53 84 5 F 26 33
17 31 1 18 F 11 63
79 90 3 10 F 26 75 5 T 5 8 3 F 33 96
21 81 2 11 F 44 85 9 F 20 48
18 75 3 13 T 10 35 13 F 1 69 9 T 48 90
99 13 3 7 T 30 90 13 F 32 96 20 F 0 13
19 7 3 15 T 36 80 6 F 65 75 1 F 40 92
49 100 1 13 F 1 79
9 15 2 3 T 83 100 15 T 25 92
90 71 2 2 F 11 95 7 F 59 88
80 11 3 17 F 47 73 2 T 32 55 11 F 18 50
11 7 1 11 T 75 85
80 81 2 12 T 50 97 20 F 57 81
64 51 3 11 T 2 81 6 F 60 69 3 F 12 25
81 84 2 11 T 1 80 19 T 6 26
9 12 3 14 F 50 81 2 F 18 32 9 T 80 86
75 58 3 9 F 34 76 18 T 10 80 7 F 30 38
93 11 3 18 F 18 81 17 T 41 100 8 T 5 64
37 2 2 20 F 4 80 7 F 12 48
55 10 1 15 T 7 73
34 79 3 12 T 24 60 7 T 38 81 18 T 21 100
77 67 1 11 T 3 43
96 2 3 5 F 12 17 18 F 57 83 4 F 53 73
9 12 2 5 T 70 99 16 T 16 73
28 64 3 10 T 1 84 15 T 70 92 19 T 38 56
88 23 2 3 F 47 66 6 F 38 61
74 41 1 10 T 22 52
7 86 2 5 T 10 42 18 T 6 29
50 81 2 11 T 11 31 20 F 3 90
64 62 1 19 T 0 2
5 31 2 19 T 2 27 8 T 9 76
8 83 3 8 T 1 25 10 F 12 29 9 F 35 69
88 23 2 20 T 0 7 20 T 6 99
42 34 2 19 T 62 86 18 T 11 62
25 89 1 9 F 6 83
14 10 3 9 T 7 52 8 F 13 77 17 F 70 81
22 78 3 16 T 22 59 6 F 1 84 7 F 15 21
8 76 1 6 F 2 60
37 20 2 17 T 16 87 18 F 32 50
83 51 2 1 T 16 38 1 T 38 84
17 53 2 6 F 7 61 13 F 1 14
61 85 3 14 T 1 25 8 T 15 39 4 T 66 72
22 6 1 3 T 24 99
99 70 3 16 F 35 55 13 T 22 72 1 T 32 99
43 37 2 11 T 0 75 17 F 32 93
52 53 1 13 T 61 87
53 36 2 10 T 9 72 17 F 31 42
52 63 2 7 T 9 37 20 T 52 66
38 35 3 7 T 17 85 16 T 27 74 8 T 2 54
13 38 1 4 T 26 78
61 29 3 8 F 66 85 8 F 58 86 15 T 5 87
83 26 1 6 T 23 79
52 67 1 6 F 29 35
84 16 3 3 F 55 70 6 F 32 97 14 F 23 100
49 79 3 15 F 20 76 5 T 19 29 1 T 71 100
72 14 3 13 F 0 19 8 F 29 92 10 T 1 68
10 26 2 18 T 38 43 7 T 3 59
64 56 2 5 T 3 21 5 T 28 77
91 9 1 10 F 20 97
95 39 1 1 F 33 66
82 10 2 3 T 46 81 16 F 43 93
90 99 2 14 T 0 29 8 T 3 48
73 9 2 8 T 1 43 8 F 11 40
12 78 2 2 T 7 57 14 T 36 87
28 45 2 20 T 27 41 9 T 11 38
15 55 2 9 F 17 66 7 T 6 9
75 63 1 18 F 10 16
21 72 1 3 F 18 30
64 24 3 1 F 9 14 6 F 51 77 2 F 12 17
95 13 3 10 F 13 88 16 T 13 60 18 F 49 71
63 69 1 18 F 60 100
24 36 2 10 T 10 98 20 T 31 55
22 80 2 15 F 12 99 11 F 20 77
80 86 1 15 T 54 73
59 26 2 3 T 8 12 6 T 22 91
50 6 2 14 F 29 71 7 F 20 90
25 87 2 13 F 5 8 9 F 11 57
84 10 2 14 F 47 83 10 T 9 11
34 81 1 2 F 29 83
59 31 2 3 T 47 62 12 F 10 57
18 99 2 18 F 5 98 8 F 24 89
8 49 1 14 T 30 93
1 6 2 17 F 1 99 15 T 34 72
45 26 3 5 F 29 48 20 F 28 44 20 F 29 73
45 95 3 15 F 16 83 4 F 5 90 11 T 14 99
52 46 2 3 F 33 76 9 F 25 90
39 43 2 20 F 27 96 16 T 8 35
96 32 3 5 T 16 50 6 T 62 98 6 T 25 35`

type mirror struct {
	v int
	c byte
	a float64
	b float64
	rawV string
	rawC string
	rawA string
	rawB string
}

type testCase struct {
	hl, hr float64
	hlRaw  string
	hrRaw  string
	n      int
	ms     []mirror
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(embeddedTestcasesC))
	scan.Split(bufio.ScanWords)
	next := func() (string, error) {
		if !scan.Scan() {
			return "", fmt.Errorf("unexpected EOF")
		}
		return scan.Text(), nil
	}
	nextInt := func() (int, string, error) {
		tok, err := next()
		if err != nil {
			return 0, "", err
		}
		v, err := strconv.Atoi(tok)
		if err != nil {
			return 0, tok, err
		}
		return v, tok, nil
	}
	nextFloat := func() (float64, string, error) {
		tok, err := next()
		if err != nil {
			return 0, "", err
		}
		v, err := strconv.ParseFloat(tok, 64)
		if err != nil {
			return 0, tok, err
		}
		return v, tok, nil
	}

	t, _, err := nextInt()
	if err != nil {
		return nil, err
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		hlVal, hlTok, err := nextFloat()
		if err != nil {
			return nil, err
		}
		hrVal, hrTok, err := nextFloat()
		if err != nil {
			return nil, err
		}
		nVal, nTok, err := nextInt()
		if err != nil {
			return nil, err
		}
		ms := make([]mirror, nVal)
		for j := 0; j < nVal; j++ {
			vInt, vTok, err := nextInt()
			if err != nil {
				return nil, err
			}
			cTok, err := next()
			if err != nil {
				return nil, err
			}
			aVal, aTok, err := nextFloat()
			if err != nil {
				return nil, err
			}
			bVal, bTok, err := nextFloat()
			if err != nil {
				return nil, err
			}
			ms[j] = mirror{
				v:    vInt,
				c:    cTok[0],
				a:    aVal,
				b:    bVal,
				rawV: vTok,
				rawC: cTok,
				rawA: aTok,
				rawB: bTok,
			}
		}
		cases = append(cases, testCase{
			hl:    hlVal,
			hr:    hrVal,
			hlRaw: hlTok,
			hrRaw: hrTok,
			n:     nVal,
			ms:    ms,
		})
		_ = nTok // silence unused if not needed
	}
	return cases, nil
}

type Mirror struct {
	a, b float64
	v    int
}

func solveCase(tc testCase) int {
	var floorMs, ceilMs []Mirror
	for _, m := range tc.ms {
		mir := Mirror{a: m.a, b: m.b, v: m.v}
		if m.c == 'F' {
			floorMs = append(floorMs, mir)
		} else {
			ceilMs = append(ceilMs, mir)
		}
	}
	sort.Slice(floorMs, func(i, j int) bool { return floorMs[i].a < floorMs[j].a })
	sort.Slice(ceilMs, func(i, j int) bool { return ceilMs[i].a < ceilMs[j].a })

	findMirror := func(ms []Mirror, x float64) *Mirror {
		i := sort.Search(len(ms), func(i int) bool { return ms[i].a > x })
		if i == 0 {
			return nil
		}
		m := &ms[i-1]
		if x+1e-9 >= m.a && x-1e-9 <= m.b {
			return m
		}
		return nil
	}

	maxScore := 0
	for s := 0; s < 2; s++ {
		for mcnt := 1; mcnt <= tc.n; mcnt++ {
			y := tc.hr
			for j := mcnt; j >= 1; j-- {
				var btype int
				if j%2 == 1 {
					btype = s
				} else {
					btype = 1 - s
				}
				if btype == 0 {
					y = -y
				} else {
					y = 2*H - y
				}
			}
			dy := y - tc.hl
			if dy == 0 {
				continue
			}
			score := 0
			ok := true
			for k := 1; k <= mcnt; k++ {
				var btype int
				if k%2 == 1 {
					btype = s
				} else {
					btype = 1 - s
				}
				var planeY float64
				if dy > 0 {
					planeY = float64(k) * H
				} else {
					planeY = float64(k-1) * H
				}
				xk := Lx * (planeY - tc.hl) / dy
				var m *Mirror
				if btype == 0 {
					m = findMirror(floorMs, xk)
				} else {
					m = findMirror(ceilMs, xk)
				}
				if m == nil {
					ok = false
					break
				}
				score += m.v
			}
			if ok && score > maxScore {
				maxScore = score
			}
		}
	}
	return maxScore
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		want := solveCase(tc)
		var input strings.Builder
		fmt.Fprintf(&input, "%s %s %d\n", tc.hlRaw, tc.hrRaw, tc.n)
		for _, m := range tc.ms {
			fmt.Fprintf(&input, "%s %s %s %s\n", m.rawV, m.rawC, m.rawA, m.rawB)
		}
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strconv.Itoa(want) != got {
			fmt.Printf("case %d failed: expected %d got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
