package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `110
14
11101000010010
2
11
8
10001101
10
1001110000
16
1000100010001111
16
1110011000001001
2
01
6
001100
10
0011111001
8
11000011
6
101000
16
0110110011001111
16
1000101000100110
8
00000001
12
100011010010
2
11
10
0111111001
2
11
2
01
16
0111111010110011
16
1010110110011011
20
10001001011011010000
12
001101110101
4
0101
2
01
12
011000011110
2
01
6
110111
6
001011
18
001100000110101101
16
0011100110011010
2
01
12
111101111111
20
10010001000110101111
20
11011011000110011001
4
1100
12
100100101100
6
100010
20
11110011100111101010
4
1101
4
0100
16
0101100011010000
20
11010001111110100011
4
1110
2
01
12
110010110100
16
0110111110101010
16
0110110100011100
18
110000011010000101
2
10
20
01101001110101010100
8
10101001
18
101000000000001001
6
111111
20
01011111111101111100
10
1010101001
18
000110001101011101
20
10111010100001100001
20
10110110110010000100
2
10
12
110010010110
10
1101011001
18
010001000010000011
18
111010101011101111
14
01110101110111
20
11110111010010111110
18
011001010000111101
10
1100110100
10
1010111110
6
010010
8
10110111
8
11110100
10
1011000100
14
10010110111111
8
01010111
10
0101111000
16
1001001001111101
8
11011100
4
1111
16
1000111000100001
16
1010001100110011
18
010110000100011110
14
11010100010000
10
0000011111
20
10011101001010001101
8
10100011
18
001011011000110111
18
001000010110100010
10
1010111011
8
01110010
16
0110101100001010
4
1110
12
110111110101
8
11111001
18
110000110111100010
8
11001010
14
10110010010110
8
01001110
12
010110111000
8
11110101
14
10100111111000
20
00010101111000110100
20
11010011010011110010
10
1010110110
18
101101000000010100
2
00
6
101000
14
10001110000110
8
01101110
8
00010101`

type testCase struct {
	n int
	s string
}

func solveCase(tc testCase) (bool, string, string) {
	n := tc.n
	s := tc.s
	ones := 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			ones++
		}
	}
	a := make([]byte, n)
	b := make([]byte, n)
	half := ones / 2
	sw := true
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			if half > 0 {
				a[i], b[i] = '(', '('
			} else {
				a[i], b[i] = ')', ')'
			}
			half--
		} else {
			if sw {
				a[i], b[i] = '(', ')'
			} else {
				a[i], b[i] = ')', '('
			}
			sw = !sw
		}
	}
	balA, balB := 0, 0
	for i := 0; i < n; i++ {
		if a[i] == '(' {
			balA++
		} else {
			balA--
		}
		if b[i] == '(' {
			balB++
		} else {
			balB--
		}
		if balA < 0 || balB < 0 {
			return false, "", ""
		}
	}
	if balA != 0 || balB != 0 {
		return false, "", ""
	}
	return true, string(a), string(b)
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	res := make([]testCase, 0, t)
	idx := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: missing data", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", caseIdx+1, err)
		}
		s := fields[idx+1]
		if len(s) != n {
			return nil, fmt.Errorf("case %d: expected string of length %d got %d", caseIdx+1, n, len(s))
		}
		res = append(res, testCase{n: n, s: s})
		idx += 2
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra data after parsing")
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		expectOK, expA, expB := solveCase(tc)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
		sb.WriteByte('\n')

		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) == 0 {
			fmt.Printf("case %d: empty output\n", i+1)
			os.Exit(1)
		}
		first := strings.ToUpper(strings.TrimSpace(lines[0]))
		if !expectOK {
			if first != "NO" {
				fmt.Printf("case %d failed\nexpected: NO\ngot: %s\n", i+1, lines[0])
				os.Exit(1)
			}
			continue
		}
		if first != "YES" {
			fmt.Printf("case %d failed\nexpected: YES\ngot: %s\n", i+1, lines[0])
			os.Exit(1)
		}
		if len(lines) < 3 {
			fmt.Printf("case %d failed: insufficient output lines\n", i+1)
			os.Exit(1)
		}
		gotA := strings.TrimSpace(lines[1])
		gotB := strings.TrimSpace(lines[2])
		if gotA != expA || gotB != expB {
			fmt.Printf("case %d failed\nexpected: %s / %s\ngot: %s / %s\n", i+1, expA, expB, gotA, gotB)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
