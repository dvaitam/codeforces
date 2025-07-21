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

type Pair struct{ first, second int }

func (a Pair) add(b Pair) Pair { return Pair{a.first + b.first, a.second + b.second} }
func (a Pair) sub(b Pair) Pair { return Pair{a.first - b.first, a.second - b.second} }
func (a Pair) less(b Pair) bool {
	if a.first != b.first {
		return a.first < b.first
	}
	return a.second < b.second
}

var (
	n        int
	f, s     []int
	tag, pre []int
	c        [][]int
	dp, cdp  [][]Pair
	cpre     []int
	ans      Pair
	todo     []Pair
)

func gao(v int) {
	if tag[v] == 0 {
		tag[v] = 1
	}
	dp[v][0] = Pair{0, 0}
	dp[v][1] = Pair{0, 0}
	pre[v] = -1
	for _, w := range c[v] {
		if tag[w] == -1 {
			continue
		}
		gao(w)
		dp[v][0] = dp[v][0].add(dp[w][1])
		diff := dp[w][1].sub(dp[w][0])
		base := Pair{1, 0}
		if s[v] != s[w] {
			base.second = 1
		}
		tmp := base.sub(diff)
		if dp[v][1].less(tmp) {
			dp[v][1] = tmp
			pre[v] = w
		}
	}
	dp[v][1] = dp[v][1].add(dp[v][0])
}

func dump(v int, flag bool, ret *[]Pair) {
	for _, w := range c[v] {
		if tag[w] == -1 {
			continue
		}
		if !flag || w != pre[v] {
			dump(w, true, ret)
		} else {
			*ret = append(*ret, Pair{v, w})
			dump(w, false, ret)
		}
	}
}

func solve(v0 int) {
	v := v0
	for tag[v] == 0 {
		tag[v] = 1
		v = f[v]
	}
	circ := []int{v}
	for u := f[v]; u != v; u = f[u] {
		circ = append(circ, u)
	}
	for _, u := range circ {
		tag[u] = -1
	}
	for _, u := range circ {
		gao(u)
	}
	m := len(circ)
	cdp[circ[0]][0] = dp[circ[0]][0]
	cdp[circ[0]][1] = dp[circ[0]][1]
	cpre[circ[0]] = -1
	for i := 1; i < m; i++ {
		u := circ[i]
		p := circ[i-1]
		cdp[u][0] = cdp[p][1].add(dp[u][0])
		cdp[u][1] = cdp[p][1].add(dp[u][1])
		cpre[u] = -1
		base := Pair{1, 0}
		if s[p] != s[u] {
			base.second = 1
		}
		tmp := cdp[p][0].add(base).add(dp[u][0])
		if cdp[u][1].less(tmp) {
			cdp[u][1] = tmp
			cpre[u] = 1
		}
	}
	best := cdp[circ[m-1]][1]
	how := []Pair{}
	for i := m - 1; i >= 0; {
		u := circ[i]
		if cpre[u] == -1 {
			dump(u, true, &how)
			i--
		} else {
			prev := circ[i-1]
			dump(u, false, &how)
			dump(prev, false, &how)
			how = append(how, Pair{u, prev})
			i -= 2
		}
	}
	circ = append(circ[1:], circ[0])
	cdp[circ[0]][0] = dp[circ[0]][0]
	cdp[circ[0]][1] = dp[circ[0]][1]
	cpre[circ[0]] = -1
	for i := 1; i < m; i++ {
		u := circ[i]
		p := circ[i-1]
		cdp[u][0] = cdp[p][1].add(dp[u][0])
		cdp[u][1] = cdp[p][1].add(dp[u][1])
		cpre[u] = -1
		base := Pair{1, 0}
		if s[p] != s[u] {
			base.second = 1
		}
		tmp := cdp[p][0].add(base).add(dp[u][0])
		if cdp[u][1].less(tmp) {
			cdp[u][1] = tmp
			cpre[u] = 1
		}
	}
	cand := cdp[circ[m-1]][1]
	if best.less(cand) {
		best = cand
		how = how[:0]
		for i := m - 1; i >= 0; {
			u := circ[i]
			if cpre[u] == -1 {
				dump(u, true, &how)
				i--
			} else {
				prev := circ[i-1]
				dump(u, false, &how)
				dump(prev, false, &how)
				how = append(how, Pair{u, prev})
				i -= 2
			}
		}
	}
	ans = ans.add(best)
	todo = append(todo, how...)
}

func expected(lines []string) string {
	n, _ = strconv.Atoi(strings.TrimSpace(lines[0]))
	f = make([]int, n+1)
	s = make([]int, n+1)
	tag = make([]int, n+1)
	pre = make([]int, n+1)
	c = make([][]int, n+1)
	dp = make([][]Pair, n+1)
	cdp = make([][]Pair, n+1)
	cpre = make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = make([]Pair, 2)
		cdp[i] = make([]Pair, 2)
	}
	for i := 0; i < n; i++ {
		parts := strings.Fields(lines[1+i])
		fi, _ := strconv.Atoi(parts[0])
		si, _ := strconv.Atoi(parts[1])
		f[i+1] = fi
		s[i+1] = si
		c[fi] = append(c[fi], i+1)
	}
	ans = Pair{0, 0}
	todo = []Pair{}
	for i := 1; i <= n; i++ {
		if tag[i] == 0 {
			solve(i)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", ans.first, ans.second))
	for _, p := range todo {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.first, p.second))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if len(lines) == 0 {
				continue
			}
			idx++
			expect := expected(lines)
			var input bytes.Buffer
			for _, l := range lines {
				input.WriteString(strings.TrimSpace(l))
				input.WriteByte('\n')
			}
			cmd := exec.Command(bin)
			cmd.Stdin = bytes.NewReader(input.Bytes())
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
				os.Exit(1)
			}
			got := strings.TrimSpace(out.String())
			if got != expect {
				fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expect, got)
				os.Exit(1)
			}
			lines = lines[:0]
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) > 0 {
		idx++
		expect := expected(lines)
		var input bytes.Buffer
		for _, l := range lines {
			input.WriteString(strings.TrimSpace(l))
			input.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
