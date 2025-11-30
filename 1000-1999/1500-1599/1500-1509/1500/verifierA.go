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

// solve mirrors 1500A.go for one test case.
func solve(a []int) (string, []int) {
	cnt := make(map[int]int)
	num := make(map[int]int)
	fresh := make(map[int][2]int)
	for i, v := range a {
		if cnt[v] == 1 {
			fresh[2*v] = [2]int{num[v], i}
		}
		cnt[v]++
		num[v] = i
	}
	keys := make([]int, 0, len(cnt))
	for k := range cnt {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	c := 0
	mr := make([]int, 0, 4)
	for _, v := range keys {
		c += cnt[v] / 2
		repeats := (cnt[v] / 2) * 2
		for j := 0; j < repeats && len(mr) < 4; j++ {
			mr = append(mr, v)
		}
	}
	if c >= 2 {
		need := make(map[int]int)
		for _, v := range mr {
			need[v]++
		}
		res := make([]int, 0, 4)
		for i, v := range a {
			if need[v] > 0 {
				res = append(res, i)
				need[v]--
				if len(res) == 4 {
					break
				}
			}
		}
		if a[res[0]] == a[res[1]] {
			res[0], res[3] = res[3], res[0]
		}
		return "Yes", res
	}
	for i1 := 0; i1 < len(keys); i1++ {
		v1 := keys[i1]
		for _, v2 := range keys[i1+1:] {
			sum := v1 + v2
			if old, exist := fresh[sum]; exist {
				return "Yes", []int{old[0], old[1], num[v1], num[v2]}
			}
			fresh[sum] = [2]int{num[v1], num[v2]}
		}
	}
	return "No", nil
}

type testCase struct {
	n   int
	arr []int
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
7 14 2 9 17 16 13 10
7 12 19 7 17 5 10 5
4 20 9 18 20
5 10 4 3 11 16
8 4 12 14 11 20 7 18 16
7 17 9 2 18 1 3 13
4 20 16 11 8
6 3 7 19 8 8 5
8 15 3 3 11 17 16 4 10
8 10 4 18 11 18 7 20 18
8 10 15 3 20 13 11 19 8
6 6 7 6 2 20 9
7 3 3 5 5 2 3 18
7 17 9 17 8 7 19 14
8 9 15 16 12 3 11 20 4
7 19 11 7 8 1 9 4
5 12 6 11 14 2
4 5 8 2 19
8 20 3 1 4 7 20 19 4
7 3 12 4 2 20 1 7
5 4 16 7 2 1
8 14 20 4 9 3 8 3 10
6 14 6 2 17 15 2
8 4 13 7 9 12 16 19 6
5 2 6 6 11 17
6 4 20 15 6 1 16
7 19 17 10 12 13 9 5
8 1 15 3 11 2 18 9 5
5 16 12 20 10 12
8 20 5 10 13 14 3 1 20
5 11 6 8 8 15
7 19 14 2 13 19 14 2
5 15 3 9 6 15
8 16 18 20 1 2 16 11 10
7 2 14 7 18 3 5 1
7 14 11 1 7 1 1 17
8 4 7 4 20 7 10 9 6
4 16 13 3 1
6 15 4 9 5 17 12
4 5 9 1 2
4 7 9 18 11
6 19 2 20 16 15 14
6 18 6 7 13 19 10
4 5 5 9 11
6 12 3 11 20 2 2
6 6 5 19 10 12 13
8 5 10 4 16 8 2 10 6
8 3 10 13 11 10 14 4 4
8 16 16 11 11 4 16 4 16
7 2 10 11 5 6 19 13
4 3 3 7 8
4 13 1 4 13
8 17 10 15 16 19 7 14 3
6 8 9 19 6 14 7
6 4 3 1 17 15 7
4 16 13 9 7
4 7 20 5 4
5 15 13 12 18 5
4 20 16 5 19
7 14 17 16 11 16 16 7
8 20 8 1 11 11 11 2 17
5 9 20 5 13 19
6 16 3 3 17 2 3
5 5 2 10 1 15
6 6 5 15 12 17 13
8 17 2 19 3 17 20 3 14
5 10 18 20 14 16
7 20 19 8 1 1 6 10
8 19 9 11 3 16 9 10 14
7 13 2 6 5 8 10 11
4 2 16 14 5
7 20 3 5 12 14 2 20
7 13 15 2 4 16 5 1
4 20 20 5 11
4 18 12 7 13
7 4 2 20 15 20 11 4
8 10 5 13 10 4 17 7 2
7 15 12 7 15 12 3 2
4 16 9 1 17
8 19 7 8 3 17 17 14 17
6 4 5 14 19 14 3
4 14 3 4 14
5 1 15 14 14 1
7 11 9 3 12 3 4 12
4 12 12 6 1
5 12 3 20 5 7
4 7 4 1 10
6 1 20 8 5 6 15
4 16 12 9 5
4 7 12 11 16
6 10 18 11 6 19 3
4 18 19 10 6
7 5 5 8 11 17 8 8
5 10 12 14 2 5
8 1 13 3 3 5 14 10 18
7 5 19 14 10 12 3 8
7 12 17 2 13 14 1 14
6 15 7 12 10 16 3
5 4 9 4 18 20
5 15 13 6 14 14
`

func parseTestcases() ([]testCase, error) {
	data := strings.TrimSpace(testcaseData)
	if data == "" {
		return nil, fmt.Errorf("no test data")
	}
	lines := strings.Split(data, "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d expected %d values, got %d", i+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[j+1])
			if err != nil {
				return nil, fmt.Errorf("case %d bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = v
		}
		res = append(res, testCase{n: n, arr: arr})
	}
	return res, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expectAns, _ := solve(tc.arr)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		words := strings.Fields(got)
		if len(words) == 0 || strings.ToLower(words[0]) != strings.ToLower(expectAns) {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expectAns, got)
			os.Exit(1)
		}
		if expectAns == "No" {
			continue
		}
		if len(words) != 5 {
			fmt.Printf("test %d failed: expected 4 indices after Yes\n", idx+1)
			os.Exit(1)
		}
		idxs := make([]int, 4)
		for i := 0; i < 4; i++ {
			v, err := strconv.Atoi(words[i+1])
			if err != nil || v < 1 || v > tc.n {
				fmt.Printf("test %d failed: bad index\n", idx+1)
				os.Exit(1)
			}
			idxs[i] = v - 1
		}
		if idxs[0] == idxs[1] || idxs[0] == idxs[2] || idxs[0] == idxs[3] ||
			idxs[1] == idxs[2] || idxs[1] == idxs[3] || idxs[2] == idxs[3] {
			fmt.Printf("test %d failed: indices not distinct\n", idx+1)
			os.Exit(1)
		}
		sum1 := tc.arr[idxs[0]] + tc.arr[idxs[1]]
		sum2 := tc.arr[idxs[2]] + tc.arr[idxs[3]]
		if sum1 != sum2 {
			fmt.Printf("test %d failed: sums mismatch\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
