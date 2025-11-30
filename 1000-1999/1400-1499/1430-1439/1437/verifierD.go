package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `2 5 5 6 12 7 19 10 6 2 4 5 9 19 10 10 6 2
3 10 12 4 14 9 11 20 18 17 9 4 9 18 2 1 11 14 20 16 11 19 4 4 14 14 3
5 5 18 12 5 2 1 3 13 12 1 6 11 17 19 19 13 18 2 4 20 9 7 7 11 18 13 5 2 3 2
3 4 14 19 4 14 4 11 3 4 5 5 8 6 7 13 14
2 4 1 17 9 14 4 17 4 20 12
3 5 17 8 16 19 7 10 2 19 4 9 14 4 2 4 11 11 2 20 10
1 4 10 19 3 4
4 8 5 11 8 5 8 5 2 19 1 15 8 18 18 8 14 8 7 20 4 6 6 10 4 15 15 7
4 5 16 7 15 6 20 4 19 13 12 14 5 13 6 5 20 18 8 4 8 20 12 19 9 1 10
4 4 13 8 10 19 6 1 12 2 14 1 19 9 1 2 2 14 20 10 15 11 1 9 13 4 14 9 2 11 2 2 12
1 10 1 17 4 19 18 2 7 19 19 14
3 2 4 6 2 5 18 4 8 16 20 9
4 4 2 17 17 16 2 8 5 7 15 3 8 8 4 5 14 8 14 9 20 7 18 2 18 13
2 8 19 5 2 10 13 19 16 15 10 14 11 18 8 3 10 18 17 8 11
3 4 17 14 9 10 7 11 11 7 4 10 10 7 2 8 14
5 3 11 7 14 8 8 9 17 12 13 3 19 16 8 4 6 4 17 9 9 1 13 7 7 1 20 6 2 14 14 6 14 13 2 9 16 19
5 1 3 1 5 10 18 8 6 11 14 5 17 10 16 8 2 9 5 10 7 8 5 1 1 15 9 10 17 1
4 8 20 6 2 6 1 17 13 12 2 11 1 8 17 8 9 13 2 14 15 16 3 7 4 18
5 2 17 7 9 13 11 9 8 20 2 11 9 3 10 19 6 4 15 1 16 12 9 14 6 1 17 3 6 14 7
3 10 18 19 2 2 8 1 19 12 9 15 10 10 9 17 11 12 5 3 15 18 4 2 5 6
3 7 7 18 9 12 12 9 17 1 17 2 14 4
5 7 15 3 4 15 5 13 6 7 18 2 12 14 9 10 5 2 19 15 9 11 20 7 18 2 8 15 17 14 3 18 18 2
2 5 18 7 9 2 17 6 10 13 5 2 7 20
3 4 15 10 1 2 7 14 18 1 6 4 12 14 4 9 13 17 12
5 3 3 11 7 1 13 7 9 11 7 17 8 19 19 10 2 16 8 7 9 6 6 15 8 2 2 8 13
4 7 1 5 20 10 10 2 19 4 10 16 13 4 10 12 3 15 12 8 20 19 2 4 12 8 7 10 10 4 1 3 20 14
2 8 7 16 18 13 17 16 6 5 1 17
4 2 14 11 3 3 19 20 6 15 20 7 16 14 2 3 9 10 18
3 7 19 1 16 3 16 10 15 9 12 16 20 3 10 7 14 15 12 8 5 16 12 3 9 19 7 19
2 9 10 3 6 1 18 19 18 15 18 8 9 6 1 14 12 10 14 2
1 1 8
4 2 2 13 8 20 7 4 9 1 6 16 20 1 3 1 9
2 8 6 7 8 8 18 1 2 15 8 11 5 3 2 20 2 2 8
4 5 1 16 9 4 18 3 3 12 11 9 7 18 20 20 7 18 17 10 17 4 6 6 12 3
1 5 14 14 17 9 17
2 2 12 13 2 6 2
4 9 1 19 9 19 11 8 1 2 2 1 8 6 4 3 10 9 15 2 2 14 8
3 5 19 20 4 5 18 2 12 12 5 12 2 3 14 10
3 5 1 5 16 11 8 4 3 14 12 18 2 6 6
5 4 18 3 10 17 7 2 12 7 15 19 19 7 8 17 4 15 7 12 20 11 11 10 20 1 14 4 11 13 14 12 8 5 6 13 13 18 10 7 7
1 7 7 6 7 2 11 13 3
1 10 18 20 2 12 5 6 12 15 14 11
3 5 11 14 6 15 5 8 5 14 18 12 11 4 1 5 6 11 7 4 10 12 20
1 2 15 18
5 4 1 15 16 3 4 3 7 7 20 4 12 12 1 10 5 1 11 20 2 12 3 14 3 14
5 3 15 16 18 6 1 12 1 3 15 16 1 8 9 2 4 13 4 14 5 16 19 19 5 19 13 1 18 2
5 4 6 20 6 15 3 15 7 1 9 2 17 16 16 13 9 15 11 20 10 18 4 6 19 2 15 16 16 9 14 4 2 2 11 13
4 4 20 19 1 18 4 19 11 7 5 5 18 18 6 14 20 2 9 13
2 2 20 12 7 5 14 18 16 4 9 14
3 6 13 12 14 9 1 7 8 12 7 3 11 1 13 14 7 10 14 8 14 11 4 1 7 10 3 7
3 2 7 6 4 19 14 11 10 3 8 17 16
1 3 16 20 3
1 5 15 9 7 17 18
2 6 10 5 6 11 13 10 6 15 16 14 6 20 15
2 5 14 14 16 2 11 8 11 2 11 16 12 15 3 3
5 5 10 20 14 11 15 2 12 15 7 18 2 13 18 12 1 12 1 19 8 17 20 14 8 2 10 17 18
2 10 18 7 4 17 13 4 18 2 9 11 6 2 15 5 17 19 1
5 10 20 14 9 16 9 20 18 8 10 6 4 7 16 8 8 5 12 7 19 11 5 2 13 7 4 20 20 7 13
2 8 10 20 3 4 10 7 1 2 6 3 15 10 10 16 13
3 4 18 12 17 20 2 11 13 6 4 18 5 20 16 20
5 7 17 12 9 18 5 5 3 3 3 9 2 4 17 12 13 18 6 4 17 19 16 1 7 9 14 1 14 14 19 18 18 8 4
1 6 14 18 10 7 16 15
2 10 19 10 17 18 16 12 4 1 11 12 8 5 18 9 2 3 15 6 16
3 9 5 6 19 3 8 17 6 15 2 2 17 12 6 14 12 8 4 14 16
5 10 17 3 19 20 19 15 10 1 11 18 8 4 18 10 4 20 12 2 19 3 10 6 11 8 8 13 11 18 1 12 11 19 8 7 15 14 19 19 11 14 14
4 4 6 16 2 7 1 16 1 6 2 6 17
2 2 4 14 4 17 1 6 13
2 2 14 11 9 8 15 2 4 4 7 17 4 4
5 2 20 3 1 6 5 15 19 17 5 10 5 19 19 20 8 11 5 5 16 6 11 7
2 3 20 13 16 7 2 19 6 20 1 10 4
1 7 2 5 7 10 3 7 1
5 8 13 14 16 7 9 3 12 16 7 11 20 1 2 4 19 20 8 15 7 15 17 1 9 8 2 4 11 17 9 2 7 12 7 8 7 11 1 20
1 1 3
2 6 20 4 18 13 9 7 3 15 15 15
2 10 3 17 18 6 4 11 19 6 3 9 7 10 7 5 4 8 7 10
3 1 4 3 18 16 9 3 3 4 7
5 9 8 8 5 15 7 16 20 16 4 6 14 18 6 1 14 19 10 5 19 5 13 12 17 18 7 7 12 4 20 16 2 14 10 15 2 18 11 7 11 6 1 11 16
1 9 20 13 10 19 18 6 8 1 16
1 4 15 10 19 10
3 1 20 6 9 3 1 20 1 15 1 5
1 1 9
3 2 3 3 4 19 20 2 11 6 20 20 15 1 18 20
4 7 10 20 1 8 12 12 11 3 14 7 13 1 1 5 20 18 7 10 5
5 1 13 5 8 1 2 20 6 8 20 5 8 20 11 10 13 17 2 17 16 4 9 9 2 12
5 5 19 9 8 10 11 9 10 9 13 16 4 3 14 7 6 4 9 12 17 15 8 20 14 7 18 15 11 4 7 10 16 10 19 4 1 2 3 8 18 14
3 9 4 3 16 10 16 15 2 10 18 9 16 12 20 16 5 18 8 2 15 10 8 8 19 8 2 10 4 13 11 14
3 3 19 13 15 3 14 6 8 1 12
4 8 4 11 9 20 5 10 2 7 9 10 9 8 3 15 4 18 16 11 3 10 17 4 1 7
2 6 18 6 12 11 8 10 4 17 2 13 9
3 7 9 3 2 1 19 5 8 9 1 18 5 19 20 10 6 6 19 8 2 16 20 10 18 8 12 4
2 4 5 5 9 17 8 10 15 10 7 11 7 5 15
4 3 17 8 14 7 9 7 12 10 15 5 14 7 15 14 7 20 2 9 19 9 2 16 19 2 14 20 20 6 1
1 10 9 6 3 2 9 11 1 10 13 4
4 7 11 7 8 6 20 15 13 6 15 18 5 19 4 10 8 20 4 5 10 18 8 17 11 7 19 11 3 12 16 1 4
2 3 9 9 5 4 6 12 11 17
2 4 19 2 17 7 9 16 12 6 13 16 16 7 15 18
4 5 8 4 17 15 14 10 9 11 3 4 5 8 16 18 6 2 8 6 17 16 5 1 8 6 7 10 2 11 12 16 10 11 19 11 11 3
3 7 3 16 20 11 19 12 14 2 4 1 3 18 19 12
2 1 3 9 4 10 14 13 10 11 16 2 8
3 6 14 11 2 3 3 2 6 19 5 2 10 19 9 10 8 17 11 18 7 7 5 20 15 13
`

type testCase struct {
	n   int
	arr []int
}

type bundle struct {
	cases []testCase
}

func solveCase(arr []int) int {
	n := len(arr)
	head, tail := 0, 0
	idx := 1
	levels := 0
	for idx < n {
		levels++
		newTail := idx
		for i := head; i <= tail && idx < n; i++ {
			j := idx
			for j+1 < n && arr[j+1] > arr[j] {
				j++
			}
			idx = j + 1
			if j > newTail {
				newTail = j
			}
		}
		head = tail + 1
		tail = newTail
	}
	return levels
}

func parseTestcases() ([]bundle, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []bundle
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		t, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid t: %v", idx+1, err)
		}
		nums := make([]int, len(fields)-1)
		for i := 1; i < len(fields); i++ {
			v, err := strconv.Atoi(fields[i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value %d: %v", idx+1, i, err)
			}
			nums[i-1] = v
		}
		cur := 0
		b := bundle{}
		for caseIdx := 0; caseIdx < t; caseIdx++ {
			if cur >= len(nums) {
				return nil, fmt.Errorf("line %d: missing data for case %d", idx+1, caseIdx+1)
			}
			n := nums[cur]
			cur++
			if cur+n > len(nums) {
				return nil, fmt.Errorf("line %d: not enough numbers for case %d", idx+1, caseIdx+1)
			}
			arr := make([]int, n)
			copy(arr, nums[cur:cur+n])
			cur += n
			b.cases = append(b.cases, testCase{n: n, arr: arr})
		}
		if cur != len(nums) {
			return nil, fmt.Errorf("line %d: extra numbers leftover", idx+1)
		}
		res = append(res, b)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	bundles, err := parseTestcases()
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

	for i, b := range bundles {
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", len(b.cases))
		for _, tc := range b.cases {
			fmt.Fprintf(&input, "%d", tc.n)
			for _, v := range tc.arr {
				fmt.Fprintf(&input, " %d", v)
			}
			input.WriteByte('\n')
		}

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		var expected strings.Builder
		for idx, tc := range b.cases {
			if idx > 0 {
				expected.WriteByte('\n')
			}
			expected.WriteString(strconv.Itoa(solveCase(tc.arr)))
		}

		if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expected.String(), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(bundles))
}
