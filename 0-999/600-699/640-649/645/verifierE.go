package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt.
const embeddedTestcasesE = `7 4 cad 1
8 4 c 1
20 1 aaaaaa 1
7 5 cdbbcdaceb 1
3 2 aabb 1
9 1 aaaaaaaa 1
19 5 b 1
1 4 cbc 1
20 4 aacbad 1
7 4 badacb 1
13 2 aaaababab 1
10 3 cacaac 1
12 1 a 1
5 2 baaabaaa 1
0 1 aaaa 1
5 1 aaaaaaaaa 1
10 2 bbaab 1
4 1 aaaaaaaaaa 1
7 2 a 1
12 5 cedeea 1
12 2 b 1
1 1 aaaaaaa 1
5 4 bd 1
12 5 ecccbcabec 1
5 4 ca 1
6 3 caacbccacc 1
9 4 dbacbdbc 1
1 5 eabdbc 1
19 5 bb 1
7 5 caebecb 1
3 1 aaaaaa 1
20 2 aabb 1
18 4 ddcbd 1
16 1 aaaaaa 1
1 5 cdebe 1
14 3 b 1
4 5 aad 1
6 1 aaaaaaaaaa 1
8 3 cbcbca 1
13 3 a 1
14 2 aaabba 1
13 1 aaaaaa 1
14 1 aaaaaaa 1
2 1 aaa 1
19 1 aaaaaa 1
8 1 aaaaa 1
7 3 accaaabbb 1
15 5 adcad 1
13 5 deb 1
1 2 aabababaaa 1
13 3 bbccbbabc 1
7 5 adbaeebaa 1
1 4 dacbadbbad 1
17 4 bcadcabcad 1
0 1 aaaa 1
16 3 babcaabaaa 1
4 1 aaaa 1
9 2 abbbbaaaab 1
9 1 a 1
12 3 ccbccabca 1
7 5 ddcaebdc 1
7 4 bbaa 1
10 1 aaa 1
6 3 cabbacaac 1
11 2 ba 1
16 3 ccabcbabcc 1
3 1 aaaa 1
11 1 aaaa 1
14 5 d 1
2 4 abad 1
10 2 ab 1
6 4 d 1
13 1 aa 1
14 4 d 1
2 5 a 1
9 3 bcbcca 1
9 2 bbba 1
12 1 aaaaaaa 1
13 3 bbabacaaca 1
8 4 bdcbcbacdb 1
2 1 aa 1
2 5 ce 1
6 4 bbd 1
9 4 dd 1
2 3 aaba 1
17 2 abbabbb 1
15 3 ba 1
10 3 abaccba 1
1 1 aaaaaaa 1
17 3 baabbabcac 1
8 4 bbcb 1
20 1 aaaaaa 1
5 5 dea 1
20 1 a 1
11 2 aabaa 1
4 2 ab 1
14 4 cbc 1
3 4 bad 1
3 1 aaa 1
6 2 abbbba 1`

const mod int64 = 1000000007

func solve645E(n, k int, t string) int64 {
	m := len(t)
	size := m + n + 1
	dp := make([]int64, size)
	dp[0] = 1
	last := make([]int, k)
	for i := 1; i <= m; i++ {
		c := int(t[i-1] - 'a')
		dp[i] = (2*dp[i-1] - dp[last[c]] + mod) % mod
		last[c] = i
	}
	for i := m + 1; i <= m+n; i++ {
		best := 0
		minPos := last[0]
		for j := 1; j < k; j++ {
			if last[j] < minPos {
				best = j
				minPos = last[j]
			}
		}
		dp[i] = (2*dp[i-1] - dp[last[best]] + mod) % mod
		last[best] = i
	}
	return dp[m+n] % mod
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesE), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 4 {
			fmt.Fprintf(os.Stderr, "case %d malformed\n", idx+1)
			os.Exit(1)
		}
		n, err1 := strconv.Atoi(fields[0])
		k, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid n or k\n", idx+1)
			os.Exit(1)
		}
		t := fields[2]
		want := strconv.FormatInt(solve645E(n, k, t), 10)
		input := fmt.Sprintf("%d %d\n%s\n", n, k, t)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
