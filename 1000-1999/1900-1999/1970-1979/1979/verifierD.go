package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesD = `100
5 1
11001
6 2
011100
4 4
1000
2 2
01
9 9
111101000
9 1
111111101
2 2
01
10 1
0110011010
8 2
01110001
7 7
0010000
5 5
11001
7 7
0111101
8 2
11111011
2 2
00
7 7
1111000
7 7
1110101
7 7
1100010
6 2
101001
7 1
1001011
5 1
00101
7 1
1011111
8 4
11010001
10 10
0011100100
2 1
01
2 1
10
4 4
1001
2 1
10
6 6
010000
4 1
1001
8 1
10110110
2 1
00
2 1
10
3 3
101
3 3
111
7 7
0110001
3 1
011
10 10
0101111110
5 1
10101
4 1
1110
9 9
010100100
8 1
00100100
3 3
110
2 1
00
10 1
1000001111
8 4
10101011
4 4
1000
9 9
110010010
10 5
1101011100
7 1
1101101
10 1
1010100001
5 1
10001
4 1
1010
3 1
011
2 2
11
2 1
11
8 8
10011010
6 1
101101
7 7
1101011
2 2
00
4 2
1000
7 7
1011100
3 1
110
6 3
100011
5 5
11101
9 1
010000011
4 4
0100
4 1
1110
9 3
111111000
3 1
101
5 1
00101
4 1
1000
7 1
0101111
4 1
0010
3 1
101
7 1
0100011
9 9
010110111
9 1
101111110
2 1
11
4 2
0000
4 4
1000
3 1
110
2 1
11
7 7
0001011
2 1
00
8 8
00101001
6 6
110101
7 7
0010000
3 1
001
2 1
10
9 3
110111010
8 1
11101101
8 8
01111011
2 1
00
9 1
010001111
2 1
01
3 3
101
3 1
000
4 4
0001
5 1
11010
3 1
000`

type testCase struct {
	n, k int
	s    string
}

func solveOne(n, k int, s string) int {
	pre := make([]int, n+1)
	for i := 0; i < n; i++ {
		pre[i+1] = pre[i] + int(s[i]-'0')
	}
	f := make([]bool, n+1)
	f[n] = true
	for i := n - 1; i >= 0; i-- {
		j := i + k
		if j > n {
			j = n
		}
		ones := pre[j] - pre[i]
		allSame := ones == 0 || ones == j-i
		okNext := f[j]
		okChange := (j == n) || (s[i] != s[j])
		f[i] = okNext && allSame && okChange
	}

	b0 := int(s[0] - '0')
	lastParity := ((n - 1) / k) % 2
	expectedLast := byte('0' + byte(b0^lastParity))
	ans := -1
	for p := 1; p <= n; p++ {
		parity := ((p - 1) / k) % 2
		expected := byte('0' + byte(b0^parity))
		if s[p-1] != expected {
			break
		}
		if f[p] && (p == n || s[p] == expectedLast) {
			ans = p
			break
		}
	}
	return ans
}

func parseTests() ([]testCase, error) {
	reader := strings.NewReader(testcasesD)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n, k int
		if _, err := fmt.Fscan(reader, &n, &k); err != nil {
			return nil, err
		}
		var s string
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return nil, err
		}
		tests[i] = testCase{n: n, k: k, s: s}
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n%s\n", tc.n, tc.k, tc.s)
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}
	allInput := buildAllInput(tests)
	allOutput, err := runCandidate(bin, allInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime error:", err)
		os.Exit(1)
	}
	outLines := strings.Split(strings.TrimSpace(allOutput), "\n")
	if len(outLines) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outLines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solveOne(tc.n, tc.k, tc.s)
		if strings.TrimSpace(outLines[i]) != strconv.Itoa(want) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput: %d %d %s\nexpected: %d\ngot: %s\n", i+1, tc.n, tc.k, tc.s, want, outLines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
