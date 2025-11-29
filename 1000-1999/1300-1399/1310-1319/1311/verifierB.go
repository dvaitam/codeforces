package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testcase struct {
	n   int
	m   int
	arr []int
	pos []int
}

const testcasesRaw = `3 2 3 9 4 2 1
5 3 7 4 16 1 13 4 3 1
5 2 8 19 4 11 1 1 4
6 0 13 7 14 1 17 8
5 3 18 8 12 8 8 4 2 1
5 4 4 6 10 4 11 4 3 1 2
4 3 17 13 19 2 2 1 3
5 1 12 18 12 3 15 1
3 2 13 12 16 1 2
2 1 20 19 1
3 0 17 8 1
3 2 18 8 13 2 1
5 2 18 20 1 13 17 2 3
6 1 14 2 16 12 19 18 2
6 3 16 12 14 12 1 18 5 3 2
6 0 8 6 18 19 6 3
6 2 2 3 3 1 15 1 3 2
4 0 20 6 12 10
2 0 6 9
6 1 9 10 15 11 16 16 1
2 1 13 11 1
3 1 4 9 17 1
6 3 1 8 1 13 5 2 2 4 3
6 5 14 18 8 17 15 8 5 1 2 3 4
2 1 5 7 1
4 0 3 10 10 6
5 4 9 5 1 18 2 2 3 4 1
6 4 2 13 7 12 4 7 5 4 3 1
5 0 13 10 17 16 1
4 3 10 1 6 7 2 1 3
5 1 9 4 13 18 12 4
6 1 3 2 3 5 6 6 5
3 1 11 20 17 2
4 2 11 4 10 8 3 2
3 2 18 4 11 1 2
2 1 5 5 1
2 1 3 19 1
6 0 9 12 10 19 18 4
5 2 4 2 10 1 20 1 4
5 0 2 7 8 19 14
3 0 15 6 8
3 2 4 14 13 2 1
5 2 4 7 11 2 1 1 2
6 2 15 13 11 13 3 3 3 4
2 1 7 20 1
4 2 6 18 7 10 1 3
4 0 9 3 15 3
6 5 11 8 13 10 2 11 2 3 4 5 1
4 0 18 20 19 20
2 0 8 1
3 1 3 9 18 1
2 0 1 10
4 3 16 5 4 17 2 1 3
3 0 5 11 10
2 1 5 7 1
6 5 2 11 20 18 7 6 3 4 5 1 2
3 1 3 15 14 2
6 3 18 15 1 13 11 6 3 4 1
5 4 1 2 12 19 5 2 1 4 3
5 4 13 6 20 3 8 4 1 3 2
6 5 15 8 8 11 16 16 2 4 5 3 1
2 0 17 12
3 2 7 10 10 2 1
3 2 15 20 3 1 2
3 0 9 14 7
6 5 2 16 13 12 13 17 2 1 3 4 5
2 1 3 5 1
5 1 13 14 13 6 11 4
3 2 16 7 4 2 1
2 1 9 8 1
6 0 7 17 15 19 1 1
6 1 9 7 6 10 5 18 2
4 2 19 9 15 6 3 2
5 3 4 7 19 13 7 3 1 4
2 0 18 10
3 0 17 12 19
4 3 17 12 17 11 1 3 2
5 2 10 18 13 11 19 4 1
5 3 7 18 1 9 20 2 4 3
4 1 15 20 17 7 2
6 0 13 19 14 13 11 20
6 5 3 16 8 10 1 14 2 4 5 1 3
6 0 12 9 14 18 10 5
5 2 16 6 15 17 2 3 4
2 1 3 12 1
5 0 6 17 6 3 13
4 2 7 17 7 8 2 3
2 0 17 12
5 4 18 2 6 10 18 3 2 1 4
6 3 6 16 9 20 11 8 3 2 5
2 1 11 14 1
4 1 3 6 19 15 3
3 2 9 15 17 1 2
3 2 15 12 10 2 1
2 0 10 3
2 0 13 11
5 0 6 2 2 20 1
3 2 2 16 17 2 1
4 0 20 6 4 8
5 1 16 15 13 6 8 2`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []testcase {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	cases := make([]testcase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			panic(fmt.Sprintf("line %d too few fields", idx+1))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			panic(fmt.Sprintf("line %d bad n: %v", idx+1, err))
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			panic(fmt.Sprintf("line %d bad m: %v", idx+1, err))
		}
		expected := 2 + n + m
		if len(fields) != expected {
			panic(fmt.Sprintf("line %d expected %d numbers, got %d", idx+1, expected, len(fields)))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				panic(fmt.Sprintf("line %d bad arr value: %v", idx+1, err))
			}
			arr[i] = v
		}
		pos := make([]int, m)
		for i := 0; i < m; i++ {
			v, err := strconv.Atoi(fields[2+n+i])
			if err != nil {
				panic(fmt.Sprintf("line %d bad pos: %v", idx+1, err))
			}
			pos[i] = v
		}
		cases = append(cases, testcase{n: n, m: m, arr: arr, pos: pos})
	}
	if len(cases) == 0 {
		panic("no testcases parsed")
	}
	return cases
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
	return out.String(), nil
}

// Embedded solver logic from 1311B.go.
func solve(tc testcase) string {
	arr := append([]int(nil), tc.arr...)
	if tc.n > 1 {
		b := make([]bool, tc.n-1)
		for _, p := range tc.pos {
			if p >= 1 && p < tc.n {
				b[p-1] = true
			}
		}
		i := 0
		for i < tc.n-1 {
			if !b[i] {
				i++
				continue
			}
			l := i
			for i < tc.n-1 && b[i] {
				i++
			}
			sort.Ints(arr[l : i+1])
		}
	}
	for i := 0; i < tc.n-1; i++ {
		if arr[i] > arr[i+1] {
			return "NO"
		}
	}
	return "YES"
}

func checkCase(bin string, idx int, tc testcase) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, p := range tc.pos {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(p))
	}
	sb.WriteByte('\n')
	input := sb.String()

	expected := strings.ToLower(strings.TrimSpace(solve(tc)))
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got := strings.ToLower(strings.TrimSpace(out))
	if got != expected {
		return fmt.Errorf("case %d: expected %s got %s", idx+1, expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, i, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
