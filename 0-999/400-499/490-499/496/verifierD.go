package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	n   int
	seq []int
}

func parseTestcases(path string) ([]testCase, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cases []testCase
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		if len(fields)-1 != n {
			return nil, fmt.Errorf("bad line: %s", line)
		}
		seq := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[i+1])
			seq[i] = v
		}
		cases = append(cases, testCase{n: n, seq: seq})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

type pair struct{ s, t int }

func solveCase(tc testCase) []pair {
	a := tc.seq
	n := len(a)
	p1 := make([]int, n+1)
	p2 := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p1[i] = p1[i-1]
		p2[i] = p2[i-1]
		if a[i-1] == 1 {
			p1[i]++
		} else {
			p2[i]++
		}
	}
	max1, max2 := p1[n], p2[n]
	pos1 := make([][]int, max1+1)
	pos2 := make([][]int, max2+1)
	for i := 1; i <= n; i++ {
		pos1[p1[i]] = append(pos1[p1[i]], i)
		pos2[p2[i]] = append(pos2[p2[i]], i)
	}
	var res []pair
	INF := n + 5
	for t := 1; t <= n; t++ {
		if t > max1 && t > max2 {
			break
		}
		i := 0
		w1, w2 := 0, 0
		lastWin := 0
		valid := true
		for i < n {
			c1, c2 := p1[i], p2[i]
			tgt1, tgt2 := c1+t, c2+t
			next1, next2 := INF, INF
			if tgt1 <= max1 {
				arr := pos1[tgt1]
				j := sort.Search(len(arr), func(j int) bool { return arr[j] > i })
				if j < len(arr) {
					next1 = arr[j]
				}
			}
			if tgt2 <= max2 {
				arr := pos2[tgt2]
				j := sort.Search(len(arr), func(j int) bool { return arr[j] > i })
				if j < len(arr) {
					next2 = arr[j]
				}
			}
			j := next1
			winner := 1
			if next2 < next1 {
				j = next2
				winner = 2
			}
			if j == INF {
				valid = false
				break
			}
			if winner == 1 {
				w1++
			} else {
				w2++
			}
			lastWin = winner
			i = j
		}
		if !valid || i != n || w1 == w2 {
			continue
		}
		s := w1
		if w2 > s {
			s = w2
		}
		if (lastWin == 1 && w1 != s) || (lastWin == 2 && w2 != s) {
			continue
		}
		res = append(res, pair{s: s, t: t})
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].s != res[j].s {
			return res[i].s < res[j].s
		}
		return res[i].t < res[j].t
	})
	return res
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.seq {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		ans := solveCase(tc)
		var exp strings.Builder
		exp.WriteString(fmt.Sprintf("%d\n", len(ans)))
		for _, p := range ans {
			exp.WriteString(fmt.Sprintf("%d %d\n", p.s, p.t))
		}
		expected := strings.TrimSpace(exp.String())
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
