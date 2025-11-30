package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const N = 80

func solve(input string) (string, error) {
	reader := strings.NewReader(strings.TrimSpace(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return "", err
	}
	L := make([]int, n+1)
	R := make([]int, n+1)
	v := make([]int, 0, 2*n)
	for i := 1; i <= n; i++ {
		if _, err := fmt.Fscan(reader, &L[i], &R[i]); err != nil {
			return "", err
		}
		v = append(v, L[i], R[i])
	}

	var (
		p        [N + 9]int
		pre, pn  int
		inArr    [N + 9]float64
		smallArr [N + 9]float64
		f        [N + 9][N + 9][N + 9]float64
		g        [N + 9][N + 9]float64
		b        [8][N + 9][N + 9]float64
		ans      [N + 9][N + 9]float64
		cnt      int
	)

	backup := func(d int) {
		for i := 0; i < pn; i++ {
			for j := 0; i+j < pn; j++ {
				b[d][i][j] = g[i][j]
			}
		}
	}

	recoverG := func(d int) {
		for i := 0; i < pn; i++ {
			for j := 0; i+j < pn; j++ {
				g[i][j] = b[d][i][j]
			}
		}
	}

	var ins func(int, int)
	ins = func(l, r int) {
		for i := l; i <= r; i++ {
			x := inArr[i]
			y := smallArr[i]
			for a := cnt; a >= 0; a-- {
				for bb := cnt - a; bb >= 0; bb-- {
					g[a+1][bb] += g[a][bb] * y
					g[a][bb+1] += g[a][bb] * x
					g[a][bb] *= 1 - x - y
				}
			}
			cnt++
		}
	}

	var dac func(int, int, int)
	dac = func(l, r, d int) {
		if l > r {
			return
		}
		if l == r {
			for i := 0; i < pn; i++ {
				for j := 0; i+j < pn; j++ {
					f[p[l]][i+pre][j+1] += g[i][j] * inArr[l]
				}
			}
			return
		}
		mid := (l + r) >> 1
		backup(d)
		ins(l, mid)
		dac(mid+1, r, d+1)
		recoverG(d)
		cnt -= mid - l + 1

		backup(d)
		ins(mid+1, r)
		dac(l, mid, d+1)
		recoverG(d)
		cnt -= r - mid
	}

	work := func(lCoord, rCoord int) {
		pn = 0
		pre = 0
		for i := 1; i <= n; i++ {
			if L[i] <= lCoord && rCoord <= R[i] {
				pn++
				p[pn] = i
				length := float64(R[i] - L[i])
				inArr[pn] = float64(rCoord-lCoord) / length
				smallArr[pn] = float64(lCoord-L[i]) / length
			} else if R[i] <= lCoord {
				pre++
			}
		}
		for i := 0; i <= pn; i++ {
			for j := 0; j <= pn; j++ {
				g[i][j] = 0
			}
		}
		cnt = 0
		g[0][0] = 1.0
		dac(1, pn, 0)
	}

	sort.Ints(v)
	uniq := v[:0]
	for i, x := range v {
		if i == 0 || x != v[i-1] {
			uniq = append(uniq, x)
		}
	}
	v = uniq

	for idx := 0; idx+1 < len(v); idx++ {
		work(v[idx], v[idx+1])
	}

	for i := 1; i <= n; i++ {
		for j := 0; j < n; j++ {
			for k := 1; k+j <= n; k++ {
				t := f[i][j][k] / float64(k)
				for r := 1; r <= k; r++ {
					ans[i][j+r] += t
				}
			}
		}
	}

	var sb strings.Builder
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			sb.WriteString(fmt.Sprintf("%.10f", ans[i][j]))
			if j < n {
				sb.WriteByte(' ')
			}
		}
		if i < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String(), nil
}

// Embedded testcases (originally from testcasesE.txt).
const rawTestcasesData = `100
1
26 32
3
12 59
20 32
7 19
3
1 16
2 3
24 59
1
8 16
5
38 51
5 53
29 42
0 34
39 66
1
34 46
2
14 41
24 55
1
27 41
4
48 51
39 57
1 39
22 68
3
21 65
29 71
8 47
5
5 22
6 52
6 24
1 46
9 49
2
24 38
36 79
3
12 39
32 65
7 43
1
44 75
1
32 61
4
11 41
35 57
8 35
16 41
1
36 69
3
14 44
15 38
30 80
4
1 30
46 47
35 61
48 98
4
14 42
15 32
30 61
9 24
4
18 42
41 73
38 48
33 77
1
12 32
5
40 84
7 11
9 31
2 24
39 50
4
24 68
12 39
31 44
10 36
1
50 100
3
48 92
33 46
46 67
2
33 68
45 84
5
27 37
40 85
39 71
42 79
13 52
4
2 18
31 70
20 56
46 59
1
2 6
2
15 44
43 83
2
42 86
6 55
5
27 69
40 59
43 65
31 75
12 23
3
35 58
45 74
25 54
3
29 34
9 24
48 56
2
25 55
45 80
2
23 69
2 4
4
13 19
26 64
36 63
34 48
1
42 51
5
41 73
19 44
34 69
4 10
38 85
4
48 49
13 40
39 88
23 43
1
47 49
2
39 41
37 70
3
9 39
16 24
37 66
3
14 39
24 72
32 54
2
37 54
5 8
2
20 61
37 74
4
4 38
30 58
20 45
46 48
1
27 66
2
48 58
1 9
4
15 18
40 58
16 41
7 30
3
23 63
40 53
47 82
5
3 25
32 68
0 11
31 33
40 57
3
8 31
6 46
35 62
1
16 53
1
33 38
5
8 40
2 36
50 65
22 29
35 81
2
11 56
42 57
4
12 51
44 92
28 51
10 28
3
35 48
13 43
20 29
3
35 78
49 99
0 6
3
19 65
50 71
2 8
5
0 19
7 18
45 63
12 45
35 58
4
5 32
32 67
26 34
5 43
4
20 32
18 34
40 68
43 88
3
2 4
0 41
17 22
3
20 56
41 60
3 49
2
44 78
8 25
5
3 28
26 44
37 52
16 37
34 66
4
45 53
27 51
10 44
3 18
4
24 26
15 21
19 26
11 20
4
0 49
23 57
26 34
33 49
4
18 36
25 53
30 75
3 14
5
22 35
15 25
14 34
50 60
11 47
3
34 36
11 26
46 73
2
15 64
35 44
5
43 54
47 63
33 83
19 65
47 49
4
50 63
3 30
10 35
32 57
2
24 59
33 65
1
8 14
5
0 6
0 48
15 63
25 54
43 61
4
22 65
25 42
46 51
48 59
5
15 63
13 16
48 52
32 37
44 73
2
28 74
49 98
2
38 79
32 51
1
38 48
3
42 65
23 43
32 74
3
8 10
30 77
12 51
4
21 52
19 48
8 39
18 63
3
22 63
38 63
25 63
4
44 77
19 63
48 85
15 28
4
13 30
32 43
32 41
38 72
5
49 82
41 72
21 38
27 47
29 71
5
3 50
12 29
29 74
4 6
12 51
1
16 17
2
11 51
11 54
3
39 57
0 44
42 44
1
38 49
3
15 41
31 79
23 46
1
37 65`

func parseTestcases() ([]string, error) {
	r := strings.NewReader(strings.TrimSpace(rawTestcasesData))
	var t int
	if _, err := fmt.Fscan(r, &t); err != nil {
		return nil, err
	}
	cases := make([]string, 0, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(r, &n); err != nil {
			return nil, err
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			var l, rr int
			if _, err := fmt.Fscan(r, &l, &rr); err != nil {
				return nil, err
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", l, rr))
		}
		cases = append(cases, strings.TrimSpace(sb.String()))
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected, err := solve(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc + "\n")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
