package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `
2 6 7 0 9 3 0 5 7 6
1 5 1 1 6 0 5 0
1 5 0 4 4 3 2 9
3 2 1 6 7 5 6
2 3 6 6 2 7 2
5 3 2 3 2 7 5 6 6 7
4 6 3 3 7 3 9 0 6 0 3 1
2 3 0 2 3 9 4
5 1 8 4 5 6 7 0
6 6 8 8 6 9 7 7 4 7 3 5 4 0
1 1 2 5
1 3 0 2 1 6
6 2 9 6 8 3 7 3 5 9
1 5 1 5 5 8 7 5
3 1 8 0 3 5
1 2 8 5 3
2 3 4 4 8 6 4
4 3 3 0 4 8 1 0 7
4 6 7 0 6 7 7 7 1 1 1 3
1 2 6 3 7
5 1 6 8 6 0 2 3
4 2 2 4 5 5 6 1
5 3 9 8 3 4 7 8 9 7
5 6 4 4 3 0 1 9 1 2 6 3 3
3 6 0 8 8 6 0 1 6 4 1
6 5 5 3 8 4 3 3 1 8 4 5 3
3 6 7 4 9 2 2 0 8 8 5
3 5 0 2 6 2 2 8 1 2
2 4 9 3 3 2 3 6
3 5 9 2 7 1 9 0 8 9
3 4 7 4 0 3 8 2 7
6 4 8 5 1 4 2 9 6 3 5 4
4 1 3 0 5 3 5
4 6 3 4 5 2 4 0 5 9 8 0
6 6 2 5 0 7 0 0 3 0 0 3 5 1
1 3 6 2 3 7
4 2 5 4 2 5 6 6
1 4 4 8 8 7 0
5 1 6 6 2 0 8 2
5 6 8 2 1 5 3 2 3 0 2 8 2
6 1 6 9 1 9 7 2 9
5 1 4 5 6 0 0 7
1 3 4 2 7 3
5 3 2 6 5 4 7 6 0 4
5 3 8 7 0 8 9 8 4 0
4 4 1 6 5 7 0 0 4 0
3 6 9 4 3 8 8 5 6 4 3
1 5 5 3 9 8 5 2
2 3 0 9 0 9 2
3 3 4 4 5 7 6 9
4 2 0 2 9 0 7 2
3 1 7 4 9 3
1 5 6 4 2 8 2 1
6 6 2 9 1 8 8 9 6 6 4 4 4 0
4 6 4 4 8 8 8 5 5 3 6 2
1 5 2 9 6 5 7 0
5 4 9 3 0 5 8 2 3 5 7
1 6 3 9 3 4 2 6 1
5 4 3 7 8 1 3 2 7 1 6
6 4 4 4 6 5 9 5 1 4 0 7
1 3 3 6 6 6
6 6 6 0 9 7 5 9 2 9 4 5 0 6
4 5 2 0 1 9 5 5 0 1 3
6 1 8 7 0 5 0 5 6
2 6 4 6 2 9 2 6 4 8
1 2 2 2 7
6 6 0 8 0 8 6 2 5 9 1 1 8 2
3 2 4 5 4 4 8
4 2 7 8 2 0 9 2
6 5 0 5 1 3 7 9 3 7 8 2 5
6 2 7 8 0 8 1 8 5 0
1 1 6 9
3 5 7 5 6 8 5 1 2 5
1 2 2 0 5
5 2 0 6 0 4 6 0 9
6 2 5 1 6 0 7 5 9 9
3 6 4 9 7 6 2 0 7 4 3
4 1 5 1 1 0 5
1 2 6 9 0
3 4 8 7 7 1 0 8 6
3 1 8 1 1 5
3 1 7 0 2 8
6 3 0 0 6 5 2 8 2 2 2
2 6 3 9 5 0 7 6 0 3
2 6 4 5 2 3 5 3 2 6
4 3 9 2 6 9 0 2 9
1 6 6 2 2 0 0 5 8
1 1 0 1
5 5 2 2 6 0 6 6 9 5 3 2
3 5 3 8 6 1 2 6 9 5
1 4 6 3 7 6 3
4 2 7 6 9 1 4 4
5 3 8 0 9 9 7 3 4 0
5 3 6 1 8 0 2 6 0 6
6 4 6 1 7 9 7 2 2 5 7 6
2 5 4 8 1 5 5 2 5
4 6 8 0 3 4 2 9 5 4 6 0
3 5 6 0 6 4 6 3 5 2
2 1 9 5 2
1 4 9 6 7 1 1
`

type testCase struct {
	la int
	lb int
	a  []int
	b  []int
}

func solveCase(tc testCase) int {
	// map each value in b to its position
	posB := make(map[int]int, tc.lb)
	for i, v := range tc.b {
		posB[v] = i
	}

	// map values of a to positions in b (or -1 if absent)
	pos := make([]int, tc.la)
	for i, v := range tc.a {
		if p, ok := posB[v]; ok {
			pos[i] = p
		} else {
			pos[i] = -1
		}
	}

	// handle rotation of a by doubling the array
	pos2 := make([]int, 2*tc.la)
	copy(pos2, pos)
	copy(pos2[tc.la:], pos)

	const INF int64 = -1
	arr := make([]int64, len(pos2)) // unwrapped positions in b
	var prev int64 = INF
	lb64 := int64(tc.lb)
	for i, p := range pos2 {
		if p == -1 {
			arr[i] = INF
			prev = INF
			continue
		}
		cur := int64(p)
		if prev != INF && cur < prev {
			cur += ((prev-cur)/lb64 + 1) * lb64
		}
		arr[i] = cur
		prev = cur
	}

	best := 0
	l := 0
	for r := 0; r < len(arr); r++ {
		if arr[r] == INF {
			l = r + 1
			continue
		}
		for l <= r && (arr[r]-arr[l] >= lb64 || r-l+1 > tc.la) {
			l++
		}
		if r-l+1 > best {
			best = r - l + 1
		}
	}

	if best > tc.la {
		best = tc.la
	}
	if best > tc.lb {
		best = tc.lb
	}
	return best
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		la, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse la: %v", idx+1, err)
		}
		lb, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse lb: %v", idx+1, err)
		}
		expected := 2 + la + lb
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, expected, len(fields))
		}
		tc := testCase{la: la, lb: lb, a: make([]int, la), b: make([]int, lb)}
		for i := 0; i < la; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", idx+1, i, err)
			}
			tc.a[i] = v
		}
		for i := 0; i < lb; i++ {
			v, err := strconv.Atoi(fields[2+la+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse b[%d]: %v", idx+1, i, err)
			}
			tc.b[i] = v
		}
		cases = append(cases, tc)
	}
	return cases, nil
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
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

	for i, tc := range cases {
		expected := solveCase(tc)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.la, tc.lb))
		for idx, v := range tc.a {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for idx, v := range tc.b {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(expected) {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
