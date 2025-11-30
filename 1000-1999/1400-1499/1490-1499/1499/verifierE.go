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
ba
baaa
b
aab
aab
ab
aab
bab
aaab
aba
bab
bbba
bb
a
b
bba
ab
bb
bb
aba
a
aab
bb
a
abb
b
bab
b
bb
bb
ba
aba
bba
abaa
a
a
ab
ba
b
abaa
bbaa
baba
aaab
aaa
baba
aba
baba
bb
bbab
bba
aab
ba
bb
baaa
ab
b
aaaa
bba
a
b
aab
a
baa
babb
bb
a
b
bbaa
bba
ba
ab
bbaa
b
a
baab
baa
a
bbb
babb
aba
ab
aba
b
bb
aba
a
aaa
aaab
bb
ba
bab
bbbb
baba
bbba
ab
bab
abaa
ab
aa
baaa
a
b
a
aba
ba
b
a
aabb
b
ab
bbba
abab
b
aa
a
aa
abbb
aa
bbba
bb
baab
a
bba
ab
aa
b
ba
a
ba
a
aaba
bbb
a
babb
bb
aa
bbbb
babb
a
bbab
a
a
aba
aba
aa
a
baa
aa
ba
bba
aa
a
b
aaab
bbbb
aaaa
aba
abb
aaa
aaa
a
bb
abba
bb
bab
ba
aabb
a
aaaa
aab
aaaa
ba
bbab
bba
ba
baaa
aa
b
a
ab
b
babb
a
baa
baa
abab
aa
aa
a
baa
aaa
a
ba
b
baa
a
ab
bbab
abab
a
`

const mod int = 998244353

type testCase struct {
	x string
	y string
}

func solveCase(x, y string) int {
	n := len(x)
	m := len(y)

	dp0Only := make([][]int, n+1)
	dp1Only := make([][]int, n+1)
	dp0Both := make([][]int, n+1)
	dp1Both := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp0Only[i] = make([]int, m+1)
		dp1Only[i] = make([]int, m+1)
		dp0Both[i] = make([]int, m+1)
		dp1Both[i] = make([]int, m+1)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			dp0Only[i+1][j] = (dp0Only[i+1][j] + 1) % mod
			dp1Only[i][j+1] = (dp1Only[i][j+1] + 1) % mod
		}
	}

	ans := 0
	for i := 0; i <= n; i++ {
		for j := 0; j <= m; j++ {
			if i > 0 {
				last := x[i-1]
				if val := dp0Only[i][j]; val != 0 {
					if i < n && x[i] != last {
						dp0Only[i+1][j] = (dp0Only[i+1][j] + val) % mod
					}
					if j < m && y[j] != last {
						dp1Both[i][j+1] = (dp1Both[i][j+1] + val) % mod
					}
				}
				if val := dp0Both[i][j]; val != 0 {
					if i < n && x[i] != last {
						dp0Both[i+1][j] = (dp0Both[i+1][j] + val) % mod
					}
					if j < m && y[j] != last {
						dp1Both[i][j+1] = (dp1Both[i][j+1] + val) % mod
					}
				}
			}
			if j > 0 {
				last := y[j-1]
				if val := dp1Only[i][j]; val != 0 {
					if j < m && y[j] != last {
						dp1Only[i][j+1] = (dp1Only[i][j+1] + val) % mod
					}
					if i < n && x[i] != last {
						dp0Both[i+1][j] = (dp0Both[i+1][j] + val) % mod
					}
				}
				if val := dp1Both[i][j]; val != 0 {
					if j < m && y[j] != last {
						dp1Both[i][j+1] = (dp1Both[i][j+1] + val) % mod
					}
					if i < n && x[i] != last {
						dp0Both[i+1][j] = (dp0Both[i+1][j] + val) % mod
					}
				}
			}
			ans += dp0Both[i][j] + dp1Both[i][j]
			ans %= mod
		}
	}
	return ans % mod
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	if len(lines) < 1+2*t {
		return nil, fmt.Errorf("not enough lines for %d cases", t)
	}
	res := make([]testCase, 0, t)
	idx := 1
	for i := 0; i < t; i++ {
		x := strings.TrimSpace(lines[idx])
		y := strings.TrimSpace(lines[idx+1])
		idx += 2
		res = append(res, testCase{x: x, y: y})
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
		input := tc.x + "\n" + tc.y + "\n"
		expected := strconv.Itoa(solveCase(tc.x, tc.y))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
