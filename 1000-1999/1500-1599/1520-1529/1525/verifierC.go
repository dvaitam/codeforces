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

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `100
1 4 1 R
2 12 5 10 LL
5 16 3 7 11 14 9 RRRLL
3 18 11 13 7 LLL
2 3 1 2 LL
5 14 6 9 11 12 3 RRRRR
4 10 7 8 6 5 LRRR
5 14 13 6 11 8 10 RRRLR
6 20 6 9 16 10 5 13 RRLRRL
3 4 1 3 2 LRL
6 8 7 5 2 3 1 4 LRLLRR
2 10 1 2 LL
1 3 1 R
3 8 7 2 6 LRL
2 7 1 6 RL
3 14 8 1 5 RLR
4 9 8 2 1 3 LLRL
5 15 13 7 8 9 6 LRRRR
6 7 6 5 2 3 1 4 LLLLLR
6 10 9 1 2 7 4 8 RLLRRR
3 20 1 5 2 RRL
1 18 3 L
1 5 1 L
2 6 2 1 RR
3 16 4 11 13 LRR
5 6 5 1 2 4 3 RRLLR
3 15 5 1 11 RLL
3 10 1 8 9 RRR
2 5 1 2 LR
3 6 2 4 1 LRR
6 8 1 8 2 3 6 5 RLRLLR
6 8 2 6 3 4 1 5 RRLRLR
3 4 4 3 1 RRR
3 9 2 9 3 LRR
5 7 4 6 1 3 5 RLRLR
4 13 6 8 11 2 RLRL
1 1 1 R
4 8 4 5 7 6 LLLL
5 12 2 6 8 7 5 LRLLL
3 8 2 3 8 LLR
4 15 8 2 1 9 RLRR
6 19 9 8 5 16 10 12 LRRLLR
3 18 1 2 16 RLL
1 20 15 R
3 6 5 6 3 RLR
5 14 8 11 5 6 1 RLRRR
4 20 17 11 18 8 LLRR
4 16 12 15 13 5 LLLR
5 20 4 9 18 5 6 RLRLR
5 20 7 13 8 11 3 LLRLR
3 13 1 9 10 RRL
1 1 1 R
3 6 6 4 1 RLL
3 17 12 3 7 LRR
3 12 9 4 7 RRR
6 9 1 8 2 4 9 6 RRRLLR
5 6 3 1 4 6 2 RLRLR
1 2 2 R
1 13 13 R
6 20 5 7 3 6 19 17 LLLLRL
6 3 3 2 3 2 1 RRRLLL
3 9 1 4 8 LRL
6 17 7 9 16 13 3 1 LRRRLL
6 4 4 3 3 4 1 RLLLRR
2 18 9 11 RL
3 11 9 10 7 RLR
6 3 2 1 2 1 3 LRLLRR
2 7 2 4 LR
3 6 5 2 3 LLR
3 6 4 5 6 LRL
1 6 6 R
3 13 9 4 8 RRL
1 3 3 L
6 4 1 3 2 4 1 RLLLRL
6 11 2 7 10 9 6 11 LRRRRL
6 5 2 3 4 5 1 RRLRLR
4 16 2 8 10 16 RLLL
2 4 3 4 LR
6 1 1 1 1 1 1 LLLLLL
4 16 5 7 6 2 RRRLL
3 5 5 2 4 LRR
3 7 4 5 3 RRR
4 10 7 4 9 4 RLLR
6 18 6 2 7 3 11 1 RRRRRR
5 19 17 2 16 11 6 LRLLR
4 14 10 7 8 11 RRLR
6 7 2 1 5 7 6 1 RLLRLR
3 9 6 1 9 LRR
2 13 5 6 LR
1 1 1 L
2 16 4 12 LR
1 3 3 L
5 12 10 6 4 1 3 LRLRR
4 12 10 11 6 7 RRRR
6 9 5 2 4 3 9 8 LRLLLL
1 4 4 R
5 18 17 8 11 1 4 RLRRR
4 16 12 9 7 14 LRLL
4 13 7 10 8 6 LRRR
5 15 1 10 12 6 9 LLRRL
4 8 7 8 7 6 RLRL
3 3 1 2 3 RRR
5 8 3 4 6 1 2 LLRRL
5 6 3 2 1 5 4 RLRRR
2 8 5 7 RL
2 1 1 2 LR
3 3 1 2 3 LLL`

type robot struct {
	x   int
	dir byte
	idx int
}

func solveCase(n, m int, xs []int, dirs string) []int {
	groups := [2][]robot{}
	for i := 0; i < n; i++ {
		p := xs[i] % 2
		groups[p] = append(groups[p], robot{xs[i], dirs[i], i})
	}

	ans := make([]int, n)
	for i := range ans {
		ans[i] = -1
	}

	for p := 0; p < 2; p++ {
		arr := groups[p]
		sort.Slice(arr, func(i, j int) bool { return arr[i].x < arr[j].x })
		stack := make([]robot, 0)
		for _, r := range arr {
			if r.dir == 'R' {
				stack = append(stack, r)
			} else {
				if len(stack) > 0 && stack[len(stack)-1].dir == 'R' {
					prev := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					time := (r.x - prev.x) / 2
					ans[r.idx] = time
					ans[prev.idx] = time
				} else {
					r.x = -r.x
					r.dir = 'R'
					stack = append(stack, r)
				}
			}
		}
		for len(stack) >= 2 {
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			time := (2*m - a.x - b.x) / 2
			ans[a.idx] = time
			ans[b.idx] = time
		}
	}
	return ans
}

func parseTestcases() ([][2]int, [][]int, []string, error) {
	in := bufio.NewReader(strings.NewReader(testcasesRaw))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, nil, nil, fmt.Errorf("parse t: %v", err)
	}
	ns := make([][2]int, 0, t)
	xsList := make([][]int, 0, t)
	dirList := make([]string, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			return nil, nil, nil, fmt.Errorf("case %d: parse n,m: %v", caseIdx+1, err)
		}
		xs := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &xs[i])
		}
		var dirs string
		if _, err := fmt.Fscan(in, &dirs); err != nil {
			return nil, nil, nil, fmt.Errorf("case %d: parse dirs: %v", caseIdx+1, err)
		}
		ns = append(ns, [2]int{n, m})
		xsList = append(xsList, xs)
		dirList = append(dirList, dirs)
	}
	return ns, xsList, dirList, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	casesNM, xsList, dirList, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i := range casesNM {
		n, m := casesNM[i][0], casesNM[i][1]
		xs := xsList[i]
		dirs := dirList[i]
		expected := solveCase(n, m, xs, dirs)

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j, v := range xs {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		sb.WriteString(dirs)
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(strings.TrimSpace(got))
		if len(gotFields) != n {
			fmt.Printf("case %d failed\nexpected %d numbers\ngot %d numbers\n", i+1, n, len(gotFields))
			os.Exit(1)
		}
		for idx, f := range gotFields {
			val, err := strconv.Atoi(f)
			if err != nil || val != expected[idx] {
				fmt.Printf("case %d failed\nexpected: %v\ngot: %s\n", i+1, expected, strings.TrimSpace(got))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(casesNM))
}
