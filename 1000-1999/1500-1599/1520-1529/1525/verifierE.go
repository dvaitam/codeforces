package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `100
2 3 1 3 2 2 1 1
1 1 2
3 1 2 3 3
2 1 2 1
1 3 2 1 1
3 3 3 1 3 4 2 2 2 4 3
1 5 2 1 2 2 1
4 4 5 3 4 4 2 2 3 3 1 1 1 4 3 5 5 4
3 2 2 1 4 2 4 3
2 3 2 3 3 2 3 3
2 3 1 1 3 1 2 3
2 1 2 1
3 4 1 1 3 1 3 3 1 3 3 3 2 4
1 3 1 2 2
2 3 2 3 1 2 3 1
3 1 4 2 3
3 3 1 4 2 4 2 1 1 1 1
2 5 3 1 3 1 3 2 3 1 2 1
1 5 2 2 1 2 1
2 4 2 2 1 1 2 2 1 3
4 2 4 2 1 1 3 3 2 5
2 2 2 2 1 2
1 3 1 2 1
4 4 1 4 2 5 2 3 3 4 3 4 5 2 3 3 4 4
1 3 1 1 2
2 3 3 1 1 3 3 2
4 4 4 2 1 2 2 1 5 3 1 4 4 2 5 1 2 2
3 5 4 4 1 1 1 1 4 3 2 1 3 1 1 3 3
1 1 2
4 2 3 4 2 4 4 1 1 1
3 5 4 4 4 1 2 3 4 4 1 2 3 2 2 2 3
4 3 3 5 1 2 1 3 1 5 1 4 5 4
1 5 1 2 2 2 1
2 1 3 1
3 5 4 3 4 4 4 2 3 3 3 2 4 1 2 2 3
3 2 3 1 1 4 2 3
3 3 2 3 1 1 1 3 1 2 2
4 5 1 1 2 4 2 5 3 4 5 2 4 2 3 3 2 3 5 5 4 4
4 3 3 4 4 5 1 3 2 5 4 5 5 3
4 4 5 1 2 5 1 4 3 5 3 1 3 4 2 4 5 1
4 4 5 1 2 5 1 4 3 5 3 1 3 4 2 4 5 1
2 2 1 2 1 1
2 4 3 3 3 3 3 3 2 1
2 2 1 1 1 2
1 1 2
1 3 1 1 2
2 1 3 2
1 5 1 2 1 1 2
4 5 1 4 2 2 3 5 4 5 3 4 3 2 4 1 4 1 3 1 3 2
1 2 1 1
2 2 3 1 2 3
2 4 3 3 3 3 2 2 1 3
2 1 2 3
3 3 4 1 2 3 3 4 3 1 2
3 1 1 3 4
3 1 2 1 4
2 1 1 3
1 3 1 2 1
1 1 1
3 5 1 1 1 4 2 2 3 2 3 3 4 4 1 2 4
4 1 5 5 3 5
3 2 4 3 2 2 1 4
4 3 3 4 4 3 5 2 3 3 5 1 1 4
3 5 3 2 3 1 2 1 2 3 2 4 1 3 1 4 2
2 1 1 3
1 1 2
3 5 2 1 2 1 2 2 3 1 3 3 3 2 2 2 1
1 2 1 2
4 1 2 2 4 4
4 5 3 3 4 4 1 1 2 2 3 2 3 2 3 5 2 5 3 4 3 5
1 2 1 2
2 5 1 1 1 1 1 2 3 3 2 3
4 1 3 4 2 2
4 5 4 3 3 1 4 2 4 1 4 2 1 3 3 5 5 1 1 4 2 5
1 2 1 2
1 1 2
4 2 5 2 5 1 2 2 5 3
2 4 3 2 3 3 2 1 3 2
3 3 3 1 2 3 2 4 4 1 1
2 2 1 3 3 1
1 5 2 1 1 1 1
1 5 2 1 1 1 1
3 5 1 3 4 3 2 4 3 1 2 3 4 1 1 3 3
1 2 2 1
4 1 3 2 2 2
2 1 1 1
1 3 2 1 1
3 2 2 1 1 1 2 4
2 1 3 2
3 3 2 1 1 1 2 1 3 4 3
3 2 3 4 2 4 2 4
1 5 1 2 2 2 1
3 2 3 3 3 2 4 3
3 3 2 3 1 2 3 2 2 4 1
3 1 1 1 1
2 2 1 2 3 2
4 2 4 4 3 5 1 3 1 1
2 1 2 1
4 2 4 1 5 2 3 3 5 3
2 1 2 2
1 4 1 2 2 2
1 4 1 1 2 2`

const mod int64 = 998244353

type testCase struct {
	n    int
	m    int
	grid [][]int
}

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func solveCase(tc testCase) int64 {
	n := tc.n
	m := tc.m

	fact := int64(1)
	for i := 1; i <= n; i++ {
		fact = fact * int64(i) % mod
	}
	invFact := modPow(fact, mod-2)

	ans := int64(0)

	for j := 0; j < m; j++ {
		// Count occurrences of each distance value in this column
		cnt := make([]int, n+2)
		for i := 0; i < n; i++ {
			d := tc.grid[i][j]
			if d <= n+1 {
				cnt[d]++
			}
		}

		// Compute number of permutations where all ranks avoid this column
		// (i.e., for every row i, rank p_i > d_{i,j})
		// Using the same approach as the accepted solution:
		// iterate r from n down to 1, accumulating suffix of rows with d >= r+1
		fail := int64(1)
		suffix := 0
		for r := n; r >= 1; r-- {
			if r+1 <= n+1 {
				suffix += cnt[r+1]
			}
			factor := suffix - (n - r)
			if factor <= 0 {
				fail = 0
				break
			}
			fail *= int64(factor)
		}

		val := int64(1)
		if fail != 0 {
			val = (val - (fail%mod)*invFact%mod + mod) % mod
		}

		ans += val
		if ans >= mod {
			ans -= mod
		}
	}

	return ans
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	res := make([]testCase, 0, 100)
	idx := 0
	caseIdx := 0
	for idx < len(fields) {
		if idx+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: not enough data", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", caseIdx+1, err)
		}
		m, err := strconv.Atoi(fields[idx+1])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse m: %v", caseIdx+1, err)
		}
		idx += 2
		if idx+n*m > len(fields) {
			return nil, fmt.Errorf("case %d: insufficient grid numbers", caseIdx+1)
		}
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				v, err := strconv.Atoi(fields[idx+i*m+j])
				if err != nil {
					return nil, fmt.Errorf("case %d: parse grid %d,%d: %v", caseIdx+1, i+1, j+1, err)
				}
				grid[i][j] = v
			}
		}
		idx += n * m
		res = append(res, testCase{n: n, m: m, grid: grid})
		caseIdx++
	}
	return res, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
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
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for r := 0; r < tc.n; r++ {
			for c := 0; c < tc.m; c++ {
				if c > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(tc.grid[r][c]))
			}
			sb.WriteByte('\n')
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil || gotVal != expected {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
