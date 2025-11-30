package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	n int
	s string
}

// Embedded testcases from testcasesD.txt.
const testcasesRaw = `100
7
KBVAACB
9
KAKVBVVCK
2
BV
9
ACKAVBAKV
2
VC
9
BVCBCVKAA
7
BAKVVBC
2
BA
8
CBKCKKKB
8
KKBBKABV
4
KBBC
8
VBVVACCC
5
BBVCV
7
ACCVCAV
9
KVKABAAKA
8
BAKVKBVV
6
CAKVBA
9
BAVBVBBAA
6
BVKVVV
2
VK
2
BC
4
VACK
9
KBVKBCKBC
9
KBBCCKBBV
2
BB
5
VCBCC
7
VKKCVKA
4
CCCC
7
AVAABKA
3
BKC
2
BK
2
KA
3
KKB
2
CB
6
VAAAKA
7
CKKKBCC
2
KV
2
VK
5
VAKAK
9
KAAACAVVC
3
AVA
4
AKCC
4
KBKK
6
CVVBAC
6
VVVVVC
8
KACBBAKA
3
CKA
5
VCCAA
3
BVV
8
AVBKAVVB
4
CCAB
9
AKCCKBCBK
6
AKAAKK
7
KAAKCCB
5
BKKBA
5
CKCBV
3
BBK
6
CCBKAA
9
KAVKBBKCK
3
KVA
6
BAACBV
4
CCKB
5
VCKAA
3
AVB
2
CA
4
VVAV
4
BKCA
8
BCAVVKBC
7
ABKBAAV
4
KCAV
6
CKCVVK
9
CVABKCAKV
6
BCKBCK
6
KKKCAC
2
CV
9
KVAAKBCVV
2
KV
8
CKCVBBVB
8
KBABKCAB
7
VBBVBAK
2
AA
2
AV
3
CBK
9
BVACBBBAB
4
KVCC
7
KCBKVCA
4
AAKB
6
CKBKCK
8
BVAVAKVV
5
KAVVV
9
AKKBAVKCC
7
BVCCCAK
4
VKAB
3
BVC
2
BC
6
CBVVKB
6
CCVVVC
2
BV
4
AKKA
4
CKAA`

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Embedded solver logic from 771D.go.
func solveCase(n int, s string) int {
	pos := make([][]int, 3)
	prefix := make([][]int, 3)
	for i := 0; i < 3; i++ {
		prefix[i] = make([]int, n+1)
	}
	for i := 0; i < n; i++ {
		for t := 0; t < 3; t++ {
			prefix[t][i+1] = prefix[t][i]
		}
		var cat int
		if s[i] == 'V' {
			cat = 0
		} else if s[i] == 'K' {
			cat = 1
		} else {
			cat = 2
		}
		pos[cat] = append(pos[cat], i)
		prefix[cat][i+1]++
	}
	cntV, cntK, cntO := len(pos[0]), len(pos[1]), len(pos[2])
	size0, size1, size2 := cntV+1, cntK+1, cntO+1
	const INF = int(1e9)
	dp := make([]int, size0*size1*size2*4)
	for i := range dp {
		dp[i] = INF
	}
	idx := func(a, b, c, last int) int { return ((a*size1+b)*size2+c)*4 + last }
	dp[idx(0, 0, 0, 3)] = 0
	for a := 0; a <= cntV; a++ {
		for b := 0; b <= cntK; b++ {
			for c := 0; c <= cntO; c++ {
				for last := 0; last < 4; last++ {
					cur := dp[idx(a, b, c, last)]
					if cur == INF {
						continue
					}
					for t := 0; t < 3; t++ {
						if last == 0 && t == 1 {
							continue
						}
						var posIdx int
						if t == 0 {
							if a >= cntV {
								continue
							}
							posIdx = pos[0][a]
						} else if t == 1 {
							if b >= cntK {
								continue
							}
							posIdx = pos[1][b]
						} else {
							if c >= cntO {
								continue
							}
							posIdx = pos[2][c]
						}
						used := min(a, prefix[0][posIdx]) + min(b, prefix[1][posIdx]) + min(c, prefix[2][posIdx])
						cost := posIdx - used
						na, nb, nc := a, b, c
						if t == 0 {
							na++
						} else if t == 1 {
							nb++
						} else {
							nc++
						}
						id2 := idx(na, nb, nc, t)
						if cur+cost < dp[id2] {
							dp[id2] = cur + cost
						}
					}
				}
			}
		}
	}
	ans := INF
	for last := 0; last < 3; last++ {
		if dp[idx(cntV, cntK, cntO, last)] < ans {
			ans = dp[idx(cntV, cntK, cntO, last)]
		}
	}
	return ans
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func parseTestcases() ([]testcase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	var cases []testcase
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d missing string", i+1)
		}
		s := scan.Text()
		cases = append(cases, testcase{n: n, s: s})
	}
	return cases, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for caseIdx, tc := range cases {
		expected := solveCase(tc.n, tc.s)
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
		out, err := runCandidate(exe, []byte(input))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != expected {
			fmt.Printf("case %d failed: expected %d got %s\n", caseIdx+1, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
