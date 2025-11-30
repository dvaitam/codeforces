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

type testCase struct {
	n, c, k int
	s       string
}

// Embedded testcases (previously stored in testcasesD.txt) so the verifier stays self contained.
const rawTestcasesD = `
6 3 6 BAACAB
2 1 2 AA
8 1 7 AAAAAAAA
7 3 1 AACBACB
6 3 1 CACCAB
1 1 1 A
4 3 3 CBBC
2 2 2 AA
4 4 4 ABAC
2 1 1 AA
6 4 1 BDCCDD
8 1 3 AAAAAAAA
5 1 2 AAAAA
3 3 1 CBC
1 1 1 A
7 2 4 BABBAAA
1 1 1 A
7 4 6 CDCDCAB
4 2 1 AABA
8 2 5 AABBBBBA
1 1 1 A
6 4 3 DBDCCC
7 4 5 DCBBDBC
3 3 1 CCC
8 3 5 BBBCCACA
5 4 1 ABACA
4 2 2 BBAA
1 1 1 A
7 4 6 CCADAAB
1 1 1 A
7 2 3 AABABAA
1 1 1 A
2 2 2 BA
1 1 1 A
5 4 4 CCADB
2 2 2 BB
3 3 2 ACC
2 1 1 AA
2 2 2 BA
6 1 1 AAAAAA
8 1 6 AAAAAAAA
1 1 1 A
7 4 7 ADBBBBA
8 3 5 BABCACAC
2 1 2 AA
3 1 2 AAA
6 1 1 AAAAAA
8 1 2 AAAAAAAA
3 3 2 CBA
2 1 1 AA
2 1 2 AA
8 4 6 DCBADCCA
4 1 4 AAAA
1 1 1 A
1 1 1 A
5 4 4 DABDC
3 1 3 AAA
5 4 2 ADBDC
6 4 5 BDDBCD
5 2 2 BABBA
6 2 3 AABAAB
8 2 3 BAABABBB
4 3 4 BCAA
4 2 2 BABB
8 4 6 DAAACBCD
5 4 4 CBAAA
1 1 1 A
1 1 1 A
2 2 1 AB
5 2 2 AAAAB
8 3 3 CCBAAACB
3 1 3 AAA
3 3 2 ACC
5 2 1 AAABA
3 2 3 ABA
3 2 2 BAB
4 1 3 AAAA
4 3 4 AABA
2 1 2 AA
5 2 2 ABBBB
1 1 1 A
3 3 3 CCA
1 1 1 A
4 4 2 DAAD
2 2 1 AA
4 2 1 BBAA
4 2 4 AAAA
4 4 4 CDAC
6 2 6 AABAAA
3 1 2 AAA
5 4 3 CAADC
2 1 1 AA
1 1 1 A
6 2 1 BBABAA
5 4 4 DABDA
5 2 5 ABBAA
7 3 6 AAAABCB
3 1 2 AAA
2 1 2 AA
1 1 1 A
`

func loadTestcases() ([]testCase, error) {
	fields := strings.Fields(rawTestcasesD)
	if len(fields)%4 != 0 {
		return nil, fmt.Errorf("unexpected testcase token count %d (want multiple of 4)", len(fields))
	}
	testcases := make([]testCase, 0, len(fields)/4)
	for i := 0; i < len(fields); i += 4 {
		n, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("parse n at token %d (%q): %w", i+1, fields[i], err)
		}
		c, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("parse c at token %d (%q): %w", i+2, fields[i+1], err)
		}
		k, err := strconv.Atoi(fields[i+2])
		if err != nil {
			return nil, fmt.Errorf("parse k at token %d (%q): %w", i+3, fields[i+2], err)
		}
		testcases = append(testcases, testCase{n: n, c: c, k: k, s: fields[i+3]})
	}
	return testcases, nil
}

// min and solve1995DCase are lifted directly from 1995D.go so the verifier does
// not depend on an external oracle.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve1995DCase(n, c, k int, s string) int {
	tag := make([]bool, 1<<uint(c))
	sum := make([][]int, c)
	for j := 0; j < c; j++ {
		sum[j] = make([]int, n+1)
		for i := 1; i <= n; i++ {
			val := 0
			if int(s[i-1]-'A') == j {
				val = 1
			}
			sum[j][i] = sum[j][i-1] + val
		}
	}
	for i := 1; i <= n-k+1; i++ {
		t := 0
		for j := 0; j < c; j++ {
			if sum[j][i+k-1]-sum[j][i-1] == 0 {
				t |= 1 << uint(j)
			}
		}
		tag[t] = true
	}
	for i := (1 << uint(c)) - 1; i > 0; i-- {
		if tag[i] {
			for j := 0; j < c; j++ {
				if (i>>uint(j))&1 != 0 {
					tag[i^(1<<uint(j))] = true
				}
			}
		}
	}
	ans := int(1e9)
	last := int(s[len(s)-1] - 'A')
	for i := 0; i < (1 << uint(c)); i++ {
		if !tag[i] && ((i>>uint(last))&1) != 0 {
			cnt := bits.OnesCount(uint(i))
			ans = min(ans, cnt)
		}
	}
	return ans
}

func solve1995D(tc testCase) int {
	return solve1995DCase(tc.n, tc.c, tc.k, tc.s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expect := solve1995D(tc)
		input := fmt.Sprintf("1\n%d %d %d %s\n", tc.n, tc.c, tc.k, tc.s)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx+1, res)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
