package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"100",
	"5 8 9 2 20 11",
	"5 4 18 2 5 13",
	"2 14 13",
	"5 4 15 20 15 6",
	"3 11 16 14",
	"3 19 10 17",
	"2 12 12",
	"3 12 16 17",
	"2 7 9",
	"3 19 11 10",
	"5 2 10 18 14 2",
	"5 9 13 7 12 5",
	"3 4 20 12",
	"3 1 14 19",
	"5 15 3 3 14 18",
	"6 5 6 5 7 6 8",
	"2 17 5",
	"5 12 20 10 11 4",
	"5 9 6 11 17 11",
	"6 5 13 18 10 8 13",
	"4 13 16 17 10",
	"5 14 4 5 5 1",
	"6 20 17 4 7 20 20",
	"6 4 9 20 6 13 3",
	"2 1 4",
	"4 16 11 4 15",
	"4 19 9 16 8",
	"3 18 18 3",
	"6 6 1 6 17 14 20",
	"3 15 13 9",
	"2 19 5",
	"5 6 15 19 2 13",
	"2 19 13",
	"4 8 17 15 2",
	"5 20 4 9 17 16",
	"6 13 16 9 6 8 18",
	"4 6 10 20 5",
	"5 3 3 16 13 19",
	"5 18 3 9 16 8",
	"2 10 5",
	"4 4 5 2 5",
	"6 18 7 1 2 13 18",
	"6 16 4 16 18 12 11",
	"2 1 8",
	"3 16 10 9",
	"3 1 16 12",
	"6 11 3 3 10 19 14",
	"3 12 13 5",
	"3 10 7 16",
	"4 10 13 20 5",
	"2 13 12",
	"6 16 8 12 12 14 9",
	"4 13 10 4 16",
	"4 4 15 5 12",
	"3 6 11 16",
	"3 4 13 13",
	"5 17 15 19 20 8",
	"5 17 10 16 8 11",
	"6 1 3 16 11 13 8",
	"5 2 19 2 14 3",
	"4 7 11 6 4",
	"3 12 1 8",
	"2 1 13",
	"6 1 5 4 20 20 7",
	"2 15 7",
	"2 17 20",
	"5 3 18 6 8 8",
	"5 13 16 1 14 7",
	"5 2 20 9 1 19",
	"4 12 11 15 5",
	"6 17 3 9 4 4 9",
	"2 5 20",
	"3 13 7 19",
	"4 7 14 17 17",
	"2 18 4",
	"5 4 17 15 16 6",
	"5 18 11 5 14 9",
	"5 3 19 17 11 8",
	"5 8 12 16 14 1",
	"5 1 18 13 15 8",
	"5 8 9 16 16 5",
	"3 15 10 12",
	"5 20 5 17 3 7",
	"4 17 4 2 5",
	"4 2 11 20 6",
	"5 13 7 14 16 7",
	"3 13 2 11",
	"6 7 11 6 17 18 20",
	"6 14 5 20 16 19 7",
	"6 15 2 8 16 20 11",
	"6 7 1 2 2 5 8",
	"5 20 8 4 17 14",
	"4 11 16 7 6",
	"4 18 12 15 13",
	"5 11 15 3 5 8",
	"2 5 13",
	"5 18 7 12 2 1",
	"5 7 3 14 19 19",
	"5 18 7 1 5 17",
	"5 10 13 18 18 3",
}

type testCase struct {
	n, m, k, q int
	pairs      [][2]int
	queries    [][2]int
	input      string
	want       string
}

func parseCases() []testCase {
	words := rawTestcases
	pos := 0
	nextLine := func() string {
		if pos >= len(words) {
			return ""
		}
		s := strings.TrimSpace(words[pos])
		pos++
		return s
	}
	nextInts := func(line string) []int {
		fields := strings.Fields(line)
		res := make([]int, len(fields))
		for i, f := range fields {
			val, _ := strconv.Atoi(f)
			res[i] = val
		}
		return res
	}

	t, _ := strconv.Atoi(nextLine())
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		line := nextLine()
		for line == "" && pos < len(words) {
			line = nextLine()
		}
		if line == "" {
			break
		}
		ints := nextInts(line)
		if len(ints) < 4 {
			continue
		}
		n, m, k, q := ints[0], ints[1], ints[2], ints[3]
		expectedLen := 4 + 2*k + 2*q
		if len(ints) != expectedLen {
			continue
		}
		pairs := make([][2]int, k)
		ptr := 4
		for j := 0; j < k; j++ {
			pairs[j] = [2]int{ints[ptr], ints[ptr+1]}
			ptr += 2
		}
		queries := make([][2]int, q)
		for j := 0; j < q; j++ {
			queries[j] = [2]int{ints[ptr], ints[ptr+1]}
			ptr += 2
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, k, q)
		for _, p := range pairs {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
		for _, qu := range queries {
			fmt.Fprintf(&sb, "%d %d\n", qu[0], qu[1])
		}
		cases = append(cases, testCase{
			n:       n,
			m:       m,
			k:       k,
			q:       q,
			pairs:   pairs,
			queries: queries,
			input:   sb.String(),
		})
	}
	return cases
}

func solve(tc testCase) string {
	n, m := tc.n, tc.m
	owner := make([]int, m+1)
	parent := make([]int, n+1)
	for _, p := range tc.pairs {
		a, b := p[0], p[1]
		if owner[b] == 0 {
			owner[b] = a
		} else {
			parent[a] = owner[b]
		}
	}
	children := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if parent[i] != 0 {
			children[parent[i]] = append(children[parent[i]], i)
		}
	}
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	size := make([]int, n+1)
	timer := 1
	type pair struct{ v, stage int }
	for i := 1; i <= n; i++ {
		if parent[i] == 0 {
			stack := []pair{{i, 0}}
			for len(stack) > 0 {
				cur := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				v := cur.v
				if cur.stage == 0 {
					tin[v] = timer
					timer++
					stack = append(stack, pair{v, 1})
					for j := len(children[v]) - 1; j >= 0; j-- {
						stack = append(stack, pair{children[v][j], 0})
					}
				} else {
					sz := 1
					for _, u := range children[v] {
						sz += size[u]
					}
					size[v] = sz
					tout[v] = timer
					timer++
				}
			}
		}
	}

	isAncestor := func(a, b int) bool {
		return tin[a] <= tin[b] && tout[b] <= tout[a]
	}

	results := make([]string, len(tc.queries))
	for i, qu := range tc.queries {
		x, y := qu[0], qu[1]
		w := owner[y]
		if w != 0 && isAncestor(x, w) {
			results[i] = strconv.Itoa(size[x])
		} else {
			results[i] = "0"
		}
	}
	return strings.Join(results, "\n")
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC <solution-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseCases()
	for idx, tc := range cases {
		exp := solve(tc)
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("case %d failed:\nexpected: %s\ngot: %s\ninput:\n%s", idx+1, exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
