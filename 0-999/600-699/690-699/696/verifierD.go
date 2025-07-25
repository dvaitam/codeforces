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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
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
