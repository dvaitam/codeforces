package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type item struct {
	x, y int
	s    string
}

func solve(n, m, k int) string {
	h := make([]item, 0, n*m)
	h = append(h, item{1, 1, "(1,1)"})
	done := 1

outer:
	for st := 0; st < len(h); st++ {
		if done >= k {
			break
		}
		cur := h[st]
		if cur.x == 1 && cur.y < m {
			ns := fmt.Sprintf("%s (%d,%d)", cur.s, 1, cur.y+1)
			h = append(h, item{1, cur.y + 1, ns})
			done++
			if done >= k {
				break outer
			}
		}
		if cur.x < n {
			ns := fmt.Sprintf("%s (%d,%d)", cur.s, cur.x+1, cur.y)
			h = append(h, item{cur.x + 1, cur.y, ns})
			done++
			if done >= k {
				break outer
			}
		}
	}

	if len(h) > k {
		h = h[:k]
	}

	sort.Slice(h, func(i, j int) bool {
		return h[i].x+h[i].y > h[j].x+h[j].y
	})

	ans := 0
	for _, it := range h {
		ans += it.x + it.y - 1
	}

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(ans))
	for _, it := range h {
		sb.WriteByte('\n')
		sb.WriteString(it.s)
	}
	return sb.String()
}

var testcasesRaw = `5 5 20
4 3 3
1 3 1
3 3 2
1 4 3
3 3 6
1 2 2
2 2 2
4 2 2
5 3 15
1 1 1
1 4 4
4 2 3
1 1 1
1 4 1
3 1 3
5 1 2
3 4 3
5 4 9
3 2 6
4 3 3
1 1 2
1 1 1
4 4 9
1 3 1
2 1 2
1 4 4
1 1 1
2 1 2
1 1 3
1 3 2
2 2 2
1 3 2
4 4 8
4 3 3
1 2 2
5 3 9
3 3 4
1 2 2
2 1 3
3 3 3
4 3 3
3 1 2
4 2 2
4 3 6
2 2 3
4 2 2
2 2 2
4 2 3
2 1 1
3 1 3
3 2 3
3 3 3
2 1 2
1 3 3
1 1 1
2 2 3
2 3 3
3 2 2
3 3 2
3 2 3
5 4 8
2 3 3
5 3 4
2 3 3
2 1 1
2 2 2
3 2 3
2 3 4
2 1 2
2 2 2
2 2 1
3 2 3
1 1 1
4 3 6
3 2 3
4 3 4
1 1 3
2 1 2
1 1 2
1 3 3
3 3 3
4 1 2
2 2 1
5 2 5
4 3 3
2 1 2
5 3 3
2 1 3
1 2 2
5 3 4
3 2 2
2 2 2
4 3 3
2 3 3
5 2 3
1 1 1
4 3 5
4 2 1
5 3 2
1 1 1
3 3 3
1 2 1
3 2 2
4 3 2
2 3 1
5 1 2
3 1 1
4 1 3
4 2 2
3 2 2
4 1 2
1 1 2
5 2 3
3 3 4
5 3 3
5 2 3
5 3 3
1 3 1
5 4 4
3 3 3
3 2 2
5 3 4`

type testcase struct {
	n int
	m int
	k int
}

func parseTestcases() ([]testcase, error) {
	lines := strings.Fields(testcasesRaw)
	if len(lines)%3 != 0 {
		return nil, fmt.Errorf("invalid embedded testcases")
	}
	res := make([]testcase, 0, len(lines)/3)
	for i := 0; i < len(lines); i += 3 {
		n, err1 := strconv.Atoi(lines[i])
		m, err2 := strconv.Atoi(lines[i+1])
		k, err3 := strconv.Atoi(lines[i+2])
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, fmt.Errorf("invalid number in testcase")
		}
		res = append(res, testcase{n: n, m: m, k: k})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected := solve(tc.n, tc.m, tc.k)

		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected:\n%s\n got:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
