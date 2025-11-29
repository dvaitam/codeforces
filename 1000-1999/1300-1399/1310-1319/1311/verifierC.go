package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1311CSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       var s string
       fmt.Fscan(reader, &s)
       // cnt[i]: number of wrong tries where p_j == i
       cnt := make([]int, n+2)
       for i := 0; i < m; i++ {
           var p int
           fmt.Fscan(reader, &p)
           if p <= n {
               cnt[p]++
           }
       }
       // build suffix sums: cnt[i] = number of wrong tries with p_j >= i
       for i := n; i >= 1; i-- {
           cnt[i] += cnt[i+1]
       }
       // compute answers for each letter
       ans := make([]int64, 26)
       for i := 1; i <= n; i++ {
           times := int64(cnt[i] + 1)
           ans[s[i-1]-'a'] += times
       }
       // output
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
`

// Keep the embedded reference solution reachable.
var _ = solution1311CSource

type testCase struct {
	n   int
	m   int
	s   string
	pos []int
}

const testcasesRaw = `8 7 aaabaccb 3 5 2 1 7 4 6
7 4 ccbcbcb 1 6 3 2
4 2 bcac 1 3
3 1 aba 1
6 5 bcccab 4 3 5 2 1
8 4 abccbcca 4 3 7 6
7 4 bbccccb 4 2 3 6
8 2 cbbbbccc 5 7
7 5 cbbcabc 3 5 1 2 4
8 2 caaccabc 2 6
2 1 ba 1
2 1 ca 1
4 2 aaca 1 3
2 1 ac 1
4 2 aaca 3 1
5 1 aaaab 1
4 2 babb 3 1
4 2 ccab 1 3
7 6 baabacc 4 6 3 1 2 5
4 3 bcac 3 1 2
4 1 aaaa 2
7 2 ccaaacb 1 3
2 1 cc 1
4 3 bbca 1 3 2
5 2 accaa 1 4
2 1 aa 1
2 1 bb 1
3 1 cbb 1
6 5 abccaa 4 3 1 5 2
4 3 bbac 2 1 3
4 1 caba 2
7 4 baccaab 1 3 5 4
2 1 ba 1
6 3 bcbabb 1 3 5
2 1 cb 1
3 1 acc 2
8 1 ccbbbbab 3
5 4 ccbbc 3 2 1 4
5 3 cbccc 1 3 4
2 1 ac 1
8 4 aaccaabb 2 6 3 4
3 2 aac 2 1
4 3 accb 1 3 2
7 2 bbccabb 6 1
8 5 bcabcabb 4 3 7 5 2
6 3 cccaca 5 2 4
8 6 bcabbccb 6 2 1 5 7 4
7 5 cbaabca 3 5 2 4 1
2 1 bb 1
5 2 ababc 1 3
5 3 bbcba 1 3 2
3 2 bcc 1 2
3 2 bbc 2 1
2 1 cb 1
4 1 cbcb 1
8 6 bbcabccc 3 1 7 5 2 4
7 3 ccbcacc 2 4 6
8 5 cabaabca 5 3 1 7 4
7 4 babbabb 3 2 4 5
8 3 bacbcabb 2 1 7
5 2 cbaba 2 1
4 3 cbbc 2 3 1
7 4 ccccbcc 6 1 3 2
5 4 abbbb 1 4 2 3
3 1 aba 2
2 1 ac 1
3 2 cbb 1 2
4 1 baaa 2
4 1 bbba 2
3 1 cac 1
4 2 ccaa 1 3
5 4 babcc 1 2 3 4
2 1 bb 1
6 5 cacccc 4 5 3 2 1
5 4 cbaab 2 1 3 4
5 3 bcaba 3 1 4
2 1 bc 1
4 3 baba 1 2 3
6 1 acccaa 2
6 1 cbccaa 3
5 1 acaaa 4
3 1 cac 2
4 3 cbcc 3 1 2
8 7 acaababb 4 2 7 6 5 3 1
8 1 cabbbcbb 2
5 2 cccaa 2 1
3 2 cac 1 2
7 3 ccbcaca 2 5 3
2 1 aa 1
8 1 cbcbaabc 7
4 3 cabb 3 1 2
4 3 accc 2 1 3
5 3 bccbc 2 4 1
7 2 bbbbbcb 5 2
6 4 babcbb 1 2 4 3
3 1 cba 2
2 1 aa 1
2 1 cb 1
8 4 abbcabcc 2 7 6 5
6 4 caccbb 1 4 2 5
`

func parseTestcases() []testCase {
	fields := strings.Fields(testcasesRaw)
	res := []testCase{}
	for i := 0; i < len(fields); {
		if i+2 >= len(fields) {
			break
		}
		n, _ := strconv.Atoi(fields[i])
		m, _ := strconv.Atoi(fields[i+1])
		s := fields[i+2]
		i += 3
		if i+m > len(fields) {
			break
		}
		pos := make([]int, m)
		for j := 0; j < m; j++ {
			pos[j], _ = strconv.Atoi(fields[i+j])
		}
		i += m
		res = append(res, testCase{n: n, m: m, s: s, pos: pos})
	}
	return res
}

func solveExpected(tc testCase) string {
	cnt := make([]int, tc.n+2)
	for _, p := range tc.pos {
		if p <= tc.n {
			cnt[p]++
		}
	}
	for i := tc.n; i >= 1; i-- {
		cnt[i] += cnt[i+1]
	}
	ans := make([]int64, 26)
	for i := 1; i <= tc.n; i++ {
		times := int64(cnt[i] + 1)
		ans[tc.s[i-1]-'a'] += times
	}
	parts := make([]string, 26)
	for i, v := range ans {
		parts[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(parts, " ")
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n%s\n", tc.n, tc.m, tc.s)
	for i, v := range tc.pos {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, idx int, tc testCase) error {
	input := buildInput(tc)
	expect := solveExpected(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s\ninput:\n%s", idx, err, string(out), input)
	}
	got := strings.TrimSpace(string(out))
	if got != expect {
		return fmt.Errorf("case %d failed:\nexpected: %s\ngot: %s\ninput:\n%s", idx, expect, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for i, tc := range testcases {
		if err := runCase(bin, i+1, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
