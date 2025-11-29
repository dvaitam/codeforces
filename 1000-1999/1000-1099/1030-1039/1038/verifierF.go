package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases previously stored in testcasesF.txt.
const testcasesFData = `100
10
10100
3
1
8
1000
4
1010
2
1
3
0
1
0
3
1
6
00
4
1011
3
1
2
10
10
011111
3
10
1
0
6
0111
10
1
1
0
2
1
6
11101
5
1
2
1
9
011011
5
00
5
0001
1
1
5
1000
8
000100
7
001
5
10
9
11111110
2
01
6
01
5
10101
6
11111
3
011
6
111111
6
110111
9
001
8
01011
10
100010000
8
0000
6
10
9
1
8
100010
6
1101
1
0
5
11101
7
1100110
3
010
7
1
10
0111101111
7
0111110
6
110
8
0
10
01001011
4
0011
3
111
10
00011110
8
100011
2
0
7
1111101
8
10
7
11100
8
10100010
2
1
7
1
8
0
2
0
3
10
3
010
2
0
9
10
5
0
7
00101
10
11100101
9
010011001
8
0001000
2
01
8
10010010
2
1
4
11
4
111
7
11
7
1000100
1
0
4
01
8
011111
8
00000111
3
101
3
00
8
1
6
001
8
011
8
0000
10
0111
3
11
2
00
3
001
5
1`

type testCase struct {
	n int
	s string
}

// solve mirrors the 1038F solution logic to produce the expected answer.
func solve(n int, s string) string {
	m := len(s)
	a := make([]int, m+2)
	for i := 1; i <= m; i++ {
		a[i] = int(s[i-1] - '0')
	}
	Next := make([]int, m+2)
	j := 0
	for i := 2; i <= m; i++ {
		for j > 0 && a[j+1] != a[i] {
			j = Next[j]
		}
		if a[j+1] == a[i] {
			j++
		}
		Next[i] = j
	}
	used := make([]int, n+2)
	F := make([][][]int64, n+2)
	for i := 0; i <= n; i++ {
		F[i] = make([][]int64, m)
		for k := 0; k < m; k++ {
			F[i][k] = make([]int64, 2)
		}
	}
	tmp := make([][]int64, m)
	for k := 0; k < m; k++ {
		tmp[k] = make([]int64, 2)
	}
	var ans int64
	clear := func() {
		for i := 0; i <= n; i++ {
			for k := 0; k < m; k++ {
				F[i][k][0] = 0
				F[i][k][1] = 0
			}
		}
	}
	var dp func(st, ed int)
	dp = func(st, ed int) {
		for i := st; i <= ed; i++ {
			for k := 0; k < m; k++ {
				for ok := 0; ok < 2; ok++ {
					cnt := F[i-1][k][ok]
					if cnt == 0 {
						continue
					}
					for cur := 0; cur < 2; cur++ {
						if used[i] != -1 && used[i] != cur {
							continue
						}
						j := k
						for j > 0 && a[j+1] != cur {
							j = Next[j]
						}
						if a[j+1] == cur {
							j++
						}
						if j == m {
							F[i][0][1] += cnt
						} else {
							F[i][j][ok] += cnt
						}
					}
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		used[i] = -1
	}
	clear()
	F[0][0][0] = 1
	dp(1, n)
	for k := 0; k < m; k++ {
		ans += F[n][k][1]
	}
	for length := 1; length < m; length++ {
		for i := 1; i <= n; i++ {
			used[i] = -1
		}
		for i := 1; i <= length; i++ {
			used[n-length+i] = a[i]
		}
		for i := length + 1; i <= m; i++ {
			used[i-length] = a[i]
		}
		clear()
		start := n - length + 1
		F[start][0][0] = 1
		dp(start+1, n)
		for k := 0; k < m; k++ {
			for ok := 0; ok < 2; ok++ {
				tmp[k][ok] = F[n][k][ok]
			}
		}
		clear()
		for k := 0; k < m; k++ {
			for ok := 0; ok < 2; ok++ {
				F[0][k][ok] = tmp[k][ok]
			}
		}
		dp(1, n)
		for k := 0; k < m; k++ {
			ans += F[n][k][0]
		}
	}
	return fmt.Sprintf("%d", ans)
}

func parseTestCases(data string) ([]testCase, error) {
	tokens := strings.Fields(data)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no embedded testcases found")
	}
	idx := 0
	t, err := strconv.Atoi(tokens[idx])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if idx+1 >= len(tokens) {
			return nil, fmt.Errorf("test %d missing data", i+1)
		}
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("test %d has invalid n: %w", i+1, err)
		}
		s := tokens[idx+1]
		idx += 2
		cases = append(cases, testCase{n: n, s: s})
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("embedded data has %d extra tokens", len(tokens)-idx)
	}
	return cases, nil
}

func runCase(exe string, tc testCase) error {
	input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solve(tc.n, tc.s)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	cases, err := parseTestCases(testcasesFData)
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for caseIdx, tc := range cases {
		if err := runCase(exe, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
