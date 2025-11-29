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

type subset struct {
	size int
	sum  int64
	mask uint64
}

type testcase struct {
	n   int
	arr []int64
}

const testcasesRaw = `5 83 83 41 84 69
3 70 52 30
4 66 12 3 98
6 2 24 40 52 61 8
3 86 89 47
10 16 40 19 50 86 90 54 85 63 61
6 93 73 22 32 45 41
10 22 74 53 96 49 43 42 96 70 53
3 57 11 84
5 50 26 4 89 67
8 29 40 2 85 2 38 79 1
1 77
4 66 43 47 45
2 66 99
10 36 3 18 8 9 71 43 10 76 94
5 19 45 81 75 83
4 77 1 56 70
6 55 96 51 88 50 60
2 68 1
8 34 87 15 46 10 63 43 68
6 9 69 48 33 86 25
8 48 60 60 90 18 96 2 50
10 26 84 19 55 21 4 51 96 30 60
2 18 33
5 75 11 70 74 62
1 88
3 8 17 92
3 27 18 67
3 71 7 23
2 58 98
7 47 38 59 89 2 100 50
7 67 75 35 40 41 11 11
8 12 14 24 57 30 100 49 97
5 4 32 41 29 56
6 70 7 19 70 62 71
4 7 7 60 47
7 18 43 20 40 91 46 81
3 44 87 29
4 60 48 38 24
5 61 64 54 56 68
10 54 95 38 70 16 47 74 82 84 36
2 32 40
6 65 36 42 72 89 56
3 30 32 68
2 33 82
4 78 2 95 40
8 12 19 28 72 34 45 5 52
9 21 67 18 70 86 24 53 69 43
2 61 33
1 91
2 93 100
6 56 68 17 87 10 11
8 25 36 70 64 70 80 83 71
1 11
3 61 1 97
10 35 33 1 3 19 41 30 25 6 91
6 67 71 63 39 18 29
5 78 28 37 54 16
6 74 55 39 32 27 25
6 35 37 11 38 59 26
6 46 48 78 78 32 11
10 59 37 15 62 2 33 85 30 65 72
10 52 75 34 93 19 1 41 93 33 41
1 55
10 28 79 100 72 33 41 10 85 16 68
2 71 17
1 71
10 77 76 73 4 21 55 3 97 33 7
6 18 54 64 28 50 64
4 70 48 82 46
10 33 33 52 63 65 31 49 16 68 27
4 48 92 17 81
5 23 36 44 58 77
8 15 28 62 66 37 25 78 40
6 5 71 40 55 38 12
8 38 61 84 5 22 20 84 31
8 63 99 31 10 68 69 36 39
2 90 49
2 51 96
6 87 38 55 88 33 87
10 52 3 80 73 31 41 76 43 58 8
4 1 14 52 25
2 11 16
2 88 79
6 23 68 21 64 73 73
5 16 27 73 53 91
7 16 3 46 88 1 23 81
8 36 55 40 42 78 87 92 68
1 76
10 82 94 75 99 53 19 58 45 40 6
10 57 37 67 34 99 35 62 89 49 12
7 86 95 73 21 41 26 18
6 11 46 69 69 1 62
4 19 49 10 5
6 12 52 28 3 67 62
10 6 14 16 7 40 19 81 65 25 8
6 20 61 1 80 44 54
7 12 18 51 70 51 44 75
7 95 53 75 53 26 66 39
3 43 38 45
8 10 23 63 87 51 79 64 69
4 46 9 29 17
8 49 68 70 38 80 8 99 84
3 33 72 24
2 51 43`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []testcase {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	res := make([]testcase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			panic(fmt.Sprintf("line %d too short", idx+1))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			panic(fmt.Sprintf("line %d bad n: %v", idx+1, err))
		}
		if len(fields) != n+1 {
			panic(fmt.Sprintf("line %d expected %d values got %d", idx+1, n+1, len(fields)))
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[1+i], 10, 64)
			if err != nil {
				panic(fmt.Sprintf("line %d bad value: %v", idx+1, err))
			}
			arr[i] = v
		}
		res = append(res, testcase{n: n, arr: arr})
	}
	if len(res) == 0 {
		panic("no testcases parsed")
	}
	return res
}

// generateSubsets mirrors 1250I.go logic but with provided k and m.
func generateSubsets(arr []int64, k int64, m int) ([]subset, int) {
	n := len(arr)
	var subs []subset
	var dfs func(idx int, sum int64, mask uint64, size int)
	dfs = func(idx int, sum int64, mask uint64, size int) {
		if sum > k {
			return
		}
		if idx == n {
			if size > 0 {
				subs = append(subs, subset{size: size, sum: sum, mask: mask})
			}
			return
		}
		dfs(idx+1, sum, mask, size)
		dfs(idx+1, sum+arr[idx], mask|1<<uint(idx), size+1)
	}
	dfs(0, 0, 0, 0)

	sort.Slice(subs, func(i, j int) bool {
		if subs[i].size != subs[j].size {
			return subs[i].size > subs[j].size
		}
		if subs[i].sum != subs[j].sum {
			return subs[i].sum < subs[j].sum
		}
		return subs[i].mask < subs[j].mask
	})

	if m > len(subs) {
		m = len(subs)
	}
	return subs, m
}

func expectedOutput(arr []int64) string {
	k := int64(0)
	for _, v := range arr {
		k += v
	}
	totalSubs := 1<<uint(len(arr)) - 1
	subs, r := generateSubsets(arr, k, totalSubs)
	if r == 0 {
		return "0\n"
	}
	var buf bytes.Buffer
	fmt.Fprintln(&buf, r)
	for i := 0; i < r; i++ {
		fmt.Fprintf(&buf, "%d %d\n", subs[i].size, subs[i].sum)
	}
	last := subs[r-1]
	first := true
	for i := 0; i < len(arr); i++ {
		if last.mask&(1<<uint(i)) != 0 {
			if !first {
				buf.WriteByte(' ')
			}
			first = false
			fmt.Fprintf(&buf, "%d", i+1)
		}
	}
	buf.WriteByte('\n')
	return buf.String()
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

func buildInput(arr []int64) string {
	n := len(arr)
	k := int64(0)
	for _, v := range arr {
		k += v
	}
	m := (1 << uint(n)) - 1
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d %d\n", n, k, m)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func checkCase(bin string, idx int, tc testcase) error {
	input := buildInput(tc.arr)
	expected := expectedOutput(tc.arr)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != strings.TrimSpace(expected) {
		return fmt.Errorf("case %d: output mismatch", idx+1)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
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
