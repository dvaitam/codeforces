package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const modC = 51123987

const testcaseData = `100
11 aacbabbccaa
13 bbbbaaccbcaac
6 baaccb
4 abca
13 bbcacaaaacabc
4 bacc
4 acbb
6 cbcbbc
6 cbbaac
3 acc
12 ccbcabcbbbaa
12 abaaaabbabaa
11 aaaccbcbbcc
16 ccbbacccabacccab
15 aaabccccccccbbb
20 babbacccacbaaacaaaab
10 caacaacabb
17 ccbbabcaccacabcab
5 aacca
11 bbbbcbbbbcb
15 cbbcbbbcaacbbbb
5 abaac
7 abccaca
16 ababaaacaccaccac
5 acacb
12 aabbabccacac
19 bcbaabacacaacbbcbab
11 caabaaabbab
14 cababcbbabcabb
6 acbbac
6 acccbc
15 cabbbbababaabbb
20 abbacbababacaabcbcca
14 baccbcaaaaaabc
7 abbbbcb
15 ccabbbaccbacccb
11 acbbaaccbaa
12 cbabbccabaca
6 bbaacc
8 caabcbcc
20 aaacccabbacabaacacaa
16 cacbcccaabcacbcb
18 acbacbbabbbabcabac
1 c
20 bababcbbaaccbcacacaa
5 bbcac
8 ccaaabbb
6 bacbcb
17 acbcacabbccaabbca
19 aaaaabbbcbaaccbbaba
19 bababbaccbccbacccaa
17 bbaccacaaababaacc
17 cacbbabaacacbccbb
2 ac
9 abcacbbbb
2 bc
4 bbab
8 aaacbacb
6 bbacca
18 bcabcacbcaacbccbaa
4 bccb
6 bbaccc
8 ababbbba
5 bcbca
3 bbb
11 cbbababcabc
6 aababc
4 bccb
18 acabcbcbaaaccaacca
10 bbcabaccab
6 cabaac
11 caacccaccbc
20 cbbbcaaacbbccbacbbcc
3 baa
4 bcbb
9 babcbcabc
12 bcbaabcccbbc
1 c
4 baba
3 aac
18 bacbacaacbbccccbaa
15 bcbaacbaaaaccaa
12 acaabbbcbccc
11 ccbccaaacaa
4 aaac
18 acccbbcbbacccccbcc
7 acbcbcb
20 cabccabbccbcabccbccb
15 bacbbcacbbcaacb
19 acacccaaabcbaccabca
9 cbbacbcbb
10 bacccbabca
5 bbbbb
20 ccacacbbbbaacacabbbb
2 bc
11 abbaccbbbcc
11 ababcabcbcb
16 abaabaacabbbccab
1 c
1 b`

type testCase struct {
	input    string
	expected string
}

func solveCaseC(s string) int {
	n := len(s)
	nextOcc := make([][3]int, n+2)
	for x := 0; x < 3; x++ {
		nextOcc[n][x] = n + 1
		nextOcc[n+1][x] = n + 1
	}
	for i := n - 1; i >= 0; i-- {
		for x := 0; x < 3; x++ {
			nextOcc[i][x] = nextOcc[i+1][x]
		}
		var xi int
		switch s[i] {
		case 'a':
			xi = 0
		case 'b':
			xi = 1
		default:
			xi = 2
		}
		nextOcc[i][xi] = i + 1
	}
	type triple struct{ a, b, c int }
	var dists []triple
	base := n / 3
	rem := n % 3
	if rem == 0 {
		dists = append(dists, triple{base, base, base})
	} else if rem == 1 {
		dists = append(dists, triple{base + 1, base, base})
		dists = append(dists, triple{base, base + 1, base})
		dists = append(dists, triple{base, base, base + 1})
	} else {
		dists = append(dists, triple{base + 1, base + 1, base})
		dists = append(dists, triple{base + 1, base, base + 1})
		dists = append(dists, triple{base, base + 1, base + 1})
	}
	ans := 0
	for _, d := range dists {
		va, vb, vc := d.a, d.b, d.c
		dp := make([][][]int, n+1)
		for i := 0; i <= n; i++ {
			dp[i] = make([][]int, va+1)
			for ca := 0; ca <= va; ca++ {
				dp[i][ca] = make([]int, vb+1)
			}
		}
		dp[0][0][0] = 1
		for pos := 0; pos < n; pos++ {
			for ca := 0; ca <= va; ca++ {
				for cb := 0; cb <= vb; cb++ {
					cur := dp[pos][ca][cb]
					if cur == 0 {
						continue
					}
					ccDone := pos - ca - cb
					for x := 0; x < 3; x++ {
						var remCap int
						switch x {
						case 0:
							remCap = va - ca
						case 1:
							remCap = vb - cb
						default:
							remCap = vc - ccDone
						}
						if remCap <= 0 {
							continue
						}
						nxt := nextOcc[pos][x]
						if nxt > n {
							continue
						}
						Lmin := nxt - pos
						maxLen := n - pos
						if remCap < maxLen {
							maxLen = remCap
						}
						for L := Lmin; L <= maxLen; L++ {
							pos2 := pos + L
							ca2 := ca
							cb2 := cb
							if x == 0 {
								ca2 += L
							} else if x == 1 {
								cb2 += L
							}
							dp[pos2][ca2][cb2] += cur
							if dp[pos2][ca2][cb2] >= modC {
								dp[pos2][ca2][cb2] -= modC
							}
						}
					}
				}
			}
		}
		ans = (ans + dp[n][va][vb]) % modC
	}
	return ans
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing string", caseNum+1)
		}
		s := fields[pos]
		pos++
		if len(s) != n {
			return nil, fmt.Errorf("case %d: length mismatch", caseNum+1)
		}
		input := fmt.Sprintf("%d\n%s\n", n, s)
		cases = append(cases, testCase{
			input:    input,
			expected: strconv.Itoa(solveCaseC(s)),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
