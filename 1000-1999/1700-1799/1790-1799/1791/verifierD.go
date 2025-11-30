package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
4 srel
10 pusctapirh
4 wprr
8 muehueqm
1 v
2 fy
10 bjyaiptxmw
7 mxzsoel
2 be
8 givnyujn
9 mslrsnshk
1 i
10 vwfwkrssdw
4 usij
2 cp
8 clzcneaj
7 yndbtty
1 m
10 kriqhbjacd
7 igbqtob
9 qklkdneuw
7 tbgcjzw
8 mbdbsygp
6 mrsqcf
8 lgaeesjw
10 kyzfwoayml
9 qtdtlnwmv
5 enfxd
8 iznjshhn
9 cpmlwsxaq
9 tlcgarfaj
4 vjjo
4 irwf
5 nkdaq
5 llwul
7 nksjykt
2 jk
1 k
2 yx
7 gvobpky
4 ozga
9 shxtdyltv
5 udfkg
8 hsixklvl
10 xeabtqyftl
2 ef
9 yoypxncoj
8 vkcphpuy
9 uxcpkjcgt
3 ufx
2 oj
7 uzgcjjb
9 exqwouina
9 lupwagvno
1 y
4 olbm
6 kgostp
5 dlqhy
6 fmbsra
2 mu
2 my
3 qwe
9 wjwqrncmj
2 nd
1 c
10 aewgbfjntd
5 aecne
6 vuoghc
5 gmuop
8 uahkbvoj
1 z
7 yhmnsua
5 gbpgk
5 gywrs
2 xa
7 fptvhky
1 z
6 ovwzsu
9 etozmhtxr
2 nq
6 xubgrf
8 lvvzfwog
10 jibencnjai
3 fml
3 myb
3 ssb
6 blkxku
3 bnp
9 iiucfjjvs
7 ohzbrsj
5 chhvn
2 vp
7 uwrzflb
8 fxskppza
9 ctsvtmrqb
8 hllvhyiq
1 y
2 ch
1 h
7 pklnrik
8 hutjmsbl
7 ikmuxzp
5 alyyn
10 wmybcxtuny
5 wkccj
4 ncto
10 caacoarayn`

type testCase struct {
	input    string
	expected string
}

func solve(n int, s string) int {
	prefix := make([]int, n)
	seen := [26]bool{}
	cnt := 0
	for i := 0; i < n; i++ {
		c := s[i] - 'a'
		if !seen[c] {
			seen[c] = true
			cnt++
		}
		prefix[i] = cnt
	}
	seen = [26]bool{}
	suffix := make([]int, n)
	cnt = 0
	for i := n - 1; i >= 0; i-- {
		c := s[i] - 'a'
		if !seen[c] {
			seen[c] = true
			cnt++
		}
		suffix[i] = cnt
	}
	ans := 0
	for i := 0; i < n-1; i++ {
		val := prefix[i] + suffix[i+1]
		if val > ans {
			ans = val
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	cases := make([]testCase, 0, t)
	for idx, line := range lines[1:] {
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) != 2 {
			return nil, fmt.Errorf("case %d: expected n and string", idx+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", idx+1, err)
		}
		s := parts[1]
		if len(s) != n {
			return nil, fmt.Errorf("case %d: string length mismatch", idx+1)
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		sb.WriteString(s)
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.Itoa(solve(n, s)),
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
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(bin, tc.input)
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
