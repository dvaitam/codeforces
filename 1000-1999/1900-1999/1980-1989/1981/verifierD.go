package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesD = `100
5
12
6
5
3
8
4
13
13
20
7
18
17
21
20
16
20
2
13
9
10
2
11
4
17
3
15
2
5
19
13
11
11
20
13
9
13
21
5
5
7
8
10
7
17
5
14
2
3
8
10
17
9
18
18
7
17
20
10
11
9
14
7
18
6
3
14
21
6
4
8
5
5
16
9
20
19
12
5
12
16
14
18
21
6
3
18
11
18
20
12
10
15
8
21
14
19
15
19
13
17`

type edge struct {
	to, id int
}

var pr []int

func init() {
	const maxP = 200005
	mark := make([]bool, maxP)
	for i := 2; i < maxP; i++ {
		if mark[i] {
			continue
		}
		pr = append(pr, i)
		for j := i * 2; j < maxP; j += i {
			mark[j] = true
		}
	}
}

func solveCase(n int) []int {
	if n == 2 {
		return []int{2, 2}
	}
	var vec []int
	for i := 0; i < len(pr); i++ {
		vec = append(vec, pr[i])
		k := len(vec)
		if k&1 == 1 {
			if k*(k+1)/2 >= n-1 {
				break
			}
		} else {
			if k*(k+1)/2-k/2+1 >= n-1 {
				break
			}
		}
	}
	k := len(vec)
	adj := make([][]edge, k)
	nume := 0
	if k&1 == 1 {
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				nume++
				adj[i] = append(adj[i], edge{j, nume})
				adj[j] = append(adj[j], edge{i, nume})
			}
		}
	} else {
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				if i&1 == 0 && j == i+1 {
					continue
				}
				nume++
				adj[i] = append(adj[i], edge{j, nume})
				adj[j] = append(adj[j], edge{i, nume})
			}
		}
		nume++
		adj[0] = append(adj[0], edge{1, nume})
		adj[1] = append(adj[1], edge{0, nume})
	}
	ptr := make([]int, k)
	used := make([]bool, k)
	evis := make([]bool, nume+2)
	a := make([]int, n)
	ptr2 := 0
	var tour func(int)
	tour = func(v int) {
		if ptr2 == n {
			return
		}
		for ptr[v] < len(adj[v]) {
			e := adj[v][ptr[v]]
			if evis[e.id] {
				ptr[v]++
				continue
			}
			evis[e.id] = true
			tour(e.to)
			if ptr2 == n {
				return
			}
			ptr[v]++
		}
		if !used[v] {
			used[v] = true
			if ptr2 < n {
				a[ptr2] = vec[v]
				ptr2++
			}
		}
		if ptr2 == n {
			return
		}
		if ptr2 < n {
			a[ptr2] = vec[v]
			ptr2++
		}
	}
	tour(0)
	return a
}

type testCase struct {
	n int
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesD)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t := len(fields) - 1
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		val, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, err
		}
		tests[i] = testCase{n: val}
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}
	allInput := buildAllInput(tests)
	rawOut, err := runCandidate(bin, allInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime error:", err)
		os.Exit(1)
	}
	outFields := strings.Fields(rawOut)
	idx := 0
	for i, tc := range tests {
		if idx+tc.n > len(outFields) {
			fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
			os.Exit(1)
		}
		expected := solveCase(tc.n)
		for j := 0; j < tc.n; j++ {
			got, err := strconv.Atoi(outFields[idx+j])
			if err != nil {
				fmt.Fprintf(os.Stderr, "bad output for test %d\n", i+1)
				os.Exit(1)
			}
			if got != expected[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at position %d\nexpected: %v\ngot: %v\n", i+1, j+1, expected, outFields[idx:idx+tc.n])
				os.Exit(1)
			}
		}
		idx += tc.n
	}
	if idx != len(outFields) {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
