package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type nodeD struct {
	next [26]int
	fail int
	val  int64
}

func buildAC(patterns []string, weights []int64) []nodeD {
	nodes := make([]nodeD, 1)
	for i, s := range patterns {
		v := 0
		for j := 0; j < len(s); j++ {
			c := int(s[j] - 'a')
			if nodes[v].next[c] == 0 {
				nodes = append(nodes, nodeD{})
				nodes[v].next[c] = len(nodes) - 1
			}
			v = nodes[v].next[c]
		}
		nodes[v].val += weights[i]
	}
	queue := make([]int, 0)
	for c := 0; c < 26; c++ {
		v := nodes[0].next[c]
		if v != 0 {
			queue = append(queue, v)
		}
	}
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		f := nodes[v].fail
		nodes[v].val += nodes[f].val
		for c := 0; c < 26; c++ {
			u := nodes[v].next[c]
			if u != 0 {
				nodes[u].fail = nodes[f].next[c]
				queue = append(queue, u)
			} else {
				nodes[v].next[c] = nodes[f].next[c]
			}
		}
	}
	return nodes
}

func solveD(n int, L int, weights []int64, patterns []string) string {
	ac := buildAC(patterns, weights)
	m := len(ac)
	const NEG = int64(-1 << 60)
	dp := make([]int64, m)
	for i := 1; i < m; i++ {
		dp[i] = NEG
	}
	for step := 0; step < L; step++ {
		ndp := make([]int64, m)
		for i := 0; i < m; i++ {
			ndp[i] = NEG
		}
		for s := 0; s < m; s++ {
			if dp[s] == NEG {
				continue
			}
			for c := 0; c < 26; c++ {
				ns := ac[s].next[c]
				val := dp[s] + ac[ns].val
				if val > ndp[ns] {
					ndp[ns] = val
				}
			}
		}
		dp = ndp
	}
	ans := dp[0]
	for i := 1; i < m; i++ {
		if dp[i] > ans {
			ans = dp[i]
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

const testcasesDRaw = `100
1 5
5
b
3 4
5 1 5
b
ca
c
2 5
5 4
ca
c
1 5
4
aca
1 5
1
ab
2 5
4 4
cc
ab
1 1
2
ab
3 4
3 4 5
cb
cbc
b
3 1
3 5 2
bcc
acc
c
3 3
3 1 1
cb
b
b
1 1
3
ba
1 5
5
b
3 5
3 5 3
aab
a
c
3 1
2 4 3
bac
b
ba
2 4
4 5
cc
cac
3 3
4 2 3
bc
cb
b
3 3
1 4 5
caa
cbb
cb
3 3
4 1 5
c
b
cb
2 5
5 3
b
b
2 5
3 3
aa
cca
2 5
2 3
b
c
2 1
1 5
bc
b
1 1
3
cac
2 3
2 1
c
b
3 2
3 3 1
bca
bc
bb
3 4
3 4 5
ab
a
b
3 5
4 5 2
c
cc
bcb
1 1
5
aa
1 1
5
b
3 1
1 4 1
c
ac
c
3 4
1 5 1
ab
bab
a
1 5
1
a
2 2
1 4
cba
ab
3 5
5 4 1
ba
a
a
1 1
4
c
1 5
5
ba
2 1
3 4
bcb
ba
2 4
1 2
acc
ac
1 1
3
cc
3 4
1 5 4
b
bcb
cb
2 1
2 2
bcc
b
1 4
2
b
2 5
3 1
ca
ccb
3 1
3 5 5
c
aba
b
1 5
1
acb
1 1
3
a
1 4
3
baa
3 5
5 2 1
cac
c
ca
1 2
2
abc
3 4
3 3 5
bc
ab
abc
1 4
5
ccc
2 2
4 2
a
cb
3 5
4 5 2
b
a
cba
3 5
3 4 5
cba
ab
ba
1 5
3
b
2 2
4 4
cab
cca
2 1
4 1
aa
cca
1 3
2
b
1 2
5
cab
1 1
3
b
1 1
1
b
2 1
3 4
ccb
a
2 3
4 4
aa
ccb
2 2
5 3
b
c
2 1
4 5
ac
bcb
2 3
3 1
cc
cac
2 5
3 1
bcc
cac
3 2
2 3 4
a
abc
ccb
3 3
2 4 4
ac
a
acc
3 4
3 1 3
ba
a
cb
1 5
3
bb
3 1
3 4 1
b
aca
caa
3 2
5 4 4
acc
b
c
3 2
5 2 2
bb
b
bb
3 5
3 4 5
cbc
ac
c
3 1
1 2 4
c
baa
aac
3 4
1 5 3
aab
c
a
2 5
2 5
b
abb
2 3
2 1
aca
baa
1 4
2
ba
1 4
1
c
1 4
4
ca
3 4
2 4 3
b
cb
ab
1 3
3
bc
2 5
4 2
a
ba
2 5
2 2
a
c
2 5
2 1
a
abb
1 1
1
bb
2 1
1 2
ab
aa
1 1
4
aa
2 5
1 3
c
c
2 3
4 3
ac
ab
3 3
3 1 1
bac
aaa
aaa
3 1
4 1 2
cba
bb
ca
2 3
5 4
b
b
1 5
4
bcc
2 2
4 5
ab
ba
2 4
3 3
c
ccc
2 4
1 2
c
b
1 1
3
c`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data := []byte(testcasesDRaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		L, _ := strconv.Atoi(scan.Text())
		weights := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			w, _ := strconv.ParseInt(scan.Text(), 10, 64)
			weights[i] = w
		}
		patterns := make([]string, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			patterns[i] = scan.Text()
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, L))
		for i, w := range weights {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(w, 10))
		}
		sb.WriteByte('\n')
		for _, s := range patterns {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
		input := sb.String()
		expected := solveD(n, L, weights, patterns)
		if err := runCase(exe, input, expected); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
