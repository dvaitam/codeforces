package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesE = `2 2 0 0
6 2 4 0 2 6 0 1
2 2 0 0
6 0 5 1 3 4 6 2
6 0 3 4 6 1 5 2
3 1 1 2 0
4 1 0 2 1 4
1 0 1
4 0 1 2 4 3
6 1 1 2 3 6 0 4
5 1 3 1 0 2 5
5 2 5 0 0 3 2
3 0 3 2 1
6 1 0 4 2 6 3 1
3 1 0 1 2
5 2 5 2 0 4 0
5 0 2 4 5 3 1
1 0 1
2 1 1 0
2 2 0 0
3 0 1 3 2
4 2 2 0 0 4
6 2 1 4 0 3 0 6
5 1 0 2 3 4 1
4 2 3 0 1 0
6 1 2 6 1 4 3 0
1 1 0
6 1 4 0 1 6 5 3
3 1 1 0 3
4 0 2 4 3 1
3 2 0 0 1
2 2 0 0
5 1 0 1 4 2 3
3 2 2 0 0
3 0 2 3 1
5 0 1 5 4 2 3
6 0 1 6 4 2 3 5
5 1 1 2 0 3 4
4 2 0 0 1 4
1 1 0
5 2 4 5 0 1 0
5 1 0 5 3 1 2
2 0 2 1
6 2 0 6 5 1 4 0
2 2 0 0
5 2 2 0 0 1 4
3 2 0 0 2
2 0 2 1
2 0 1 2
2 0 1 2
6 2 6 2 5 3 0 0
5 1 0 5 2 4 1
2 0 2 1
1 0 1
4 0 4 2 3 1
3 0 3 1 2
3 1 1 2 0
3 1 3 0 2
4 0 4 1 3 2
6 2 3 6 5 0 0 4
6 2 3 2 5 0 1 0
6 1 2 0 1 5 6 4
5 1 3 4 2 0 5
2 0 1 2
3 2 0 0 2
6 2 6 2 3 1 0 0
1 0 1
4 0 3 1 4 2
5 0 4 2 3 1 5
2 0 1 2
1 0 1
1 1 0
5 1 0 4 3 5 1
5 2 4 0 3 5 0
1 0 1
2 2 0 0
3 1 1 0 3
5 1 2 5 3 0 1
6 0 3 2 1 6 5 4
2 0 2 1
4 2 0 0 3 4
3 0 3 2 1
2 2 0 0
2 1 2 0
4 1 3 0 1 2
3 1 3 2 0
4 2 3 1 0 0
1 1 0
4 0 4 1 3 2
2 1 0 1
3 0 2 3 1
3 0 3 2 1
2 1 0 1
1 0 1
3 1 3 0 1
5 2 5 2 0 0 3
4 1 3 2 4 0
2 2 0 0
1 1 0
4 2 4 0 0 3`

const modE = 1000000007

type perm []int

func clone(p perm) perm {
	q := make(perm, len(p))
	copy(q, p)
	return q
}

func permKey(p perm) string {
	b := make([]byte, len(p)*2)
	idx := 0
	for _, v := range p {
		b[idx] = byte(v)
		idx++
		b[idx] = ','
		idx++
	}
	return string(b)
}

func nextStates(p perm) []perm {
	n := len(p)
	used := make([]bool, n)
	resMap := make(map[string]perm)
	var dfs func(int, perm)
	dfs = func(i int, cur perm) {
		for i < n && used[i] {
			i++
		}
		if i == n {
			key := permKey(cur)
			if _, ok := resMap[key]; !ok {
				resMap[key] = clone(cur)
			}
			return
		}
		used[i] = true
		dfs(i+1, cur)
		used[i] = false
		for j := i + 1; j < n; j++ {
			if !used[j] {
				used[i] = true
				used[j] = true
				cur[i], cur[j] = cur[j], cur[i]
				dfs(i+1, cur)
				cur[i], cur[j] = cur[j], cur[i]
				used[i] = false
				used[j] = false
			}
		}
	}
	dfs(0, clone(p))
	res := make([]perm, 0, len(resMap))
	for _, v := range resMap {
		res = append(res, v)
	}
	return res
}

func countWays(start perm) int {
	n := len(start)
	goal := make(perm, n)
	for i := 0; i < n; i++ {
		goal[i] = i + 1
	}
	startKey := permKey(start)
	goalKey := permKey(goal)
	dist := map[string]int{startKey: 0}
	ways := map[string]int{startKey: 1}
	queue := []perm{start}
	best := -1
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		d := dist[permKey(cur)]
		if best != -1 && d >= best {
			continue
		}
		if permKey(cur) == goalKey {
			best = d
			continue
		}
		for _, nxt := range nextStates(cur) {
			k := permKey(nxt)
			if prev, ok := dist[k]; !ok {
				dist[k] = d + 1
				ways[k] = ways[permKey(cur)]
				queue = append(queue, nxt)
			} else if prev == d+1 {
				ways[k] = (ways[k] + ways[permKey(cur)]) % modE
			}
		}
	}
	if best == -1 {
		if startKey == goalKey {
			return 1
		}
		return 0
	}
	return ways[goalKey]
}

type testCase struct {
	n   int
	k   int
	arr []int
}

func solve(tc testCase) int {
	n := tc.n
	arr := append([]int(nil), tc.arr...)
	used := make([]bool, n+1)
	for _, v := range arr {
		if v != 0 {
			used[v] = true
		}
	}
	missing := make([]int, 0)
	for i := 1; i <= n; i++ {
		if !used[i] {
			missing = append(missing, i)
		}
	}
	pos := make([]int, 0)
	for i, v := range arr {
		if v == 0 {
			pos = append(pos, i)
		}
	}
	total := 0
	var dfs func(int)
	dfs = func(idx int) {
		if idx == len(pos) {
			p := make(perm, n)
			copy(p, arr)
			val := countWays(p)
			total += val
			if total >= modE {
				total %= modE
			}
			return
		}
		for i, num := range missing {
			if num == -1 {
				continue
			}
			arr[pos[idx]] = num
			missing[i] = -1
			dfs(idx + 1)
			missing[i] = num
			arr[pos[idx]] = 0
		}
	}
	if len(pos) == 0 {
		total = countWays(perm(arr)) % modE
	} else {
		dfs(0)
	}
	return total % modE
}

func parseCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesE), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err1 := strconv.Atoi(fields[0])
		k, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: bad n or k", idx+1)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n, len(fields)-2)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value %d", idx+1, i+1)
			}
			arr[i] = val
		}
		cases = append(cases, testCase{n: n, k: k, arr: arr})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Println("failed to load testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc)
		input := buildInput(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		g, err := strconv.Atoi(strings.TrimSpace(got))
		if err != nil {
			fmt.Printf("case %d: bad output\n", idx+1)
			os.Exit(1)
		}
		if g%modE != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, g)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d cases passed\n", len(cases))
}
