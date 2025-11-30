package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `2 2 2 1
1 2 1
5 0 1 2 3 5 4
3 1 2 3 1
4 1 3 1 4 2
4 2 2 4 3 1
5 2 3 1 2 5 4
2 2 1 2
4 1 3 1 4 2
2 2 2 1
3 2 2 3 1
2 1 2 1
4 1 1 3 2 4
4 1 4 3 2 1
3 1 1 2 3
4 2 3 4 2 1
4 1 2 1 3 4
1 1 1
5 0 5 2 4 1 3
4 2 1 2 4 3
5 1 2 1 3 4 5
5 0 4 2 1 3 5
4 0 4 3 1 2
5 1 1 3 5 4 2
2 1 1 2
4 1 3 4 1 2
5 1 3 5 2 1 4
5 0 3 4 5 1 2
1 0 1
4 1 2 4 3 1
2 2 1 2
5 0 4 2 1 3 5
3 0 1 3 2
1 1 1
1 0 1
5 2 5 2 4 3 1
5 1 4 5 3 1 2
4 1 3 4 1 2
3 1 1 3 2
4 0 3 2 1 4
2 1 2 1
2 1 1 2
3 1 1 3 2
2 0 1 2
2 2 1 2
2 1 2 1
1 2 1
5 2 2 5 1 3 4
3 0 1 3 2
1 0 1
4 1 1 4 2 3
1 0 1
3 2 2 3 1
5 0 3 5 1 2 4
1 0 1
4 2 4 3 2 1
4 2 4 1 3 2
1 0 1
2 2 2 1
4 2 3 4 2 1
5 2 5 4 3 2 1
4 0 1 2 3 4
3 2 3 2 1
3 0 2 3 1
2 2 1 2
4 2 4 3 2 1
3 1 2 3 1
3 0 3 1 2
2 0 2 1
4 1 3 4 2 1
1 1 1
2 0 2 1
5 2 1 5 2 4 3
5 2 2 4 1 5 3
3 0 2 1 3
3 1 2 3 1
1 1 1
5 2 1 3 2 4 5
5 1 4 1 5 3 2
5 2 2 3 5 4 1
1 2 1
1 0 1
1 0 1
3 1 2 1 3
1 2 1
1 0 1
2 1 1 2
5 2 5 2 1 4 3
1 2 1
2 2 2 1
4 2 3 4 2 1
2 0 2 1
3 1 1 3 2
3 2 1 3 2
3 1 1 2 3
4 2 1 4 2 3
4 1 3 2 4 1
3 1 3 1 2
2 1 2 1
2 2 2 1`

type testCase struct {
	n int
	k int
	p []int
}

func runCandidate(bin, input string) (string, error) {
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

func minInversions(p []int, k int) int {
	n := len(p)
	arr := make([]int, n+1)
	for i, v := range p {
		arr[v] = i + 1
	}
	w := k + 1
	future := make([][]int, n+1)
	for i := range future {
		future[i] = make([]int, w)
	}
	bit := make([]int, n+2)
	update := func(i, val int) {
		for i <= n {
			bit[i] += val
			i += i & -i
		}
	}
	query := func(i int) int {
		s := 0
		for i > 0 {
			s += bit[i]
			i -= i & -i
		}
		return s
	}
	for i := n; i >= 1; i-- {
		limit := i - k
		if limit < 1 {
			limit = 1
		}
		for j := i; j >= limit; j-- {
			future[i][i-j] = query(arr[j] - 1)
		}
		update(arr[i], 1)
	}

	const INF = int(1e18)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, 1<<uint(w))
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}

	var dfs func(int, int) int
	dfs = func(i, mask int) int {
		if i == n && mask == 0 {
			return 0
		}
		if dp[i][mask] != -1 {
			return dp[i][mask]
		}
		res := INF
		if mask != 0 {
			for j := 0; j < w; j++ {
				if mask>>uint(j)&1 == 1 {
					idx := i - j
					if idx <= 0 {
						continue
					}
					val := arr[idx]
					newMask := mask & ^(1 << uint(j))
					small := future[i][i-idx]
					for q := 0; q < w; q++ {
						if newMask>>uint(q)&1 == 1 {
							idx2 := i - q
							if idx2 > 0 && arr[idx2] < val {
								small++
							}
						}
					}
					cand := dfs(i, newMask) + small
					if cand < res {
						res = cand
					}
				}
			}
		}
		if i < n && bits.OnesCount(uint(mask)) < w && (mask>>uint(w-1)&1) == 0 {
			newMask := (mask << 1) | 1
			cand := dfs(i+1, newMask)
			if cand < res {
				res = cand
			}
		}
		dp[i][mask] = res
		return res
	}

	return dfs(0, 0)
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d parse k: %v", idx+1, err)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d expected %d numbers got %d", idx+1, 2+n, len(fields))
		}
		p := make([]int, n)
		for i := 0; i < n; i++ {
			p[i], err = strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d parse p[%d]: %v", idx+1, i, err)
			}
		}
		cases = append(cases, testCase{n: n, k: k, p: p})
	}
	return cases, nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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

	for idx, tc := range cases {
		expect := minInversions(tc.p, tc.k)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.p {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := strconv.Atoi(strings.Fields(got)[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse output %q\n", idx+1, got)
			os.Exit(1)
		}
		if gotVal != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, gotVal)
			fmt.Printf("input:\n%s", input.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
