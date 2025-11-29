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

const testcasesData = `
7 98 54 6 34 66 63 52
5 62 46 75 28 65
3 37 18 97
2 80 33
9 91 78 19 40 13 94 10 88 43
8 72 13 46 56 41 79 82 27
9 62 57 67 34 8 71 2 12 93
7 91 86 81 1 79 64 43
4 94 42 91 9
4 73 29 31 19
9 58 12 11 41 66 63 14 39 71
5 91 16 71 43 70
4 78 71 76 37
8 12 77 50 41 74 31 38 24
4 24 5 79 85
5 61 9 12 87 97
3 20 5 11
9 88 51 91 68 36 67 31 28 87
10 54 75 36 58 64 85 83 90 46 11
6 79 15 63 76 81 43
4 32 3 94 35
2 91 29
6 22 43 55 8 13 19
4 6 74 82 69
10 88 10 4 16 82 25 78 74 16 51
2 48 15
1 78
1 25
3 92 16 62
4 94 8 87 3
9 55 80 13 34 9 29 10 83 39
6 56 24 8 65 60 6
10 13 90 51 26 34 46 94 61 73 22
4 99 8 87 21
3 44 68 33
2 77 57
3 2 61 88
7 73 66 40 84 46 50 85
5 20 72 89 2 59
2 43 95
1 70
5 18 31 98 62 46
10 37 87 46 76 82 80 17 92 40 50
7 84 11 1 77 25 90 43
3 31 29 82
8 49 91 87 73 54 5 52 90
10 54 99 85 91 6 22 58 9 34 90
3 58 68 63
9 78 97 1 5 64 42 40 60 7
7 25 71 82 11 93 17 2
7 87 54 41 1 28 2 92
1 87
9 79 13 25 16 78 84 26 39 36
3 13 61 51
2 3 36
8 15 33 18 84 67 84 83 45
2 20 36
1 6
1 27
5 72 41 47 73 6
10 84 64 92 83 59 82 56 48 69 23
4 49 76 38 2
3 20 35 43
6 48 92 12 44 100 80
1 6
5 21 20 75 38 47
7 71 17 38 15 62 94 31
1 40
3 67 94 10
5 52 43 39 54 14
2 72 62
8 44 44 16 62 15 90 64 55
1 39
6 95 88 20 22 81 73
7 82 12 9 11 26 96 29
1 50
1 13
7 72 67 38 58 63 75 92
4 55 11 48 29
5 75 100 22 56 25
6 15 9 90 4 68 58
4 16 64 51 33
4 83 6 28 80
3 14 26 59
7 47 70 20 14 77 63 19
10 52 82 88 55 67 64 87 42 64 64
4 70 79 29 2
6 91 96 41 42 5 68
3 33 78 20
7 75 38 92 91 61 9 11
9 6 9 29 17 6 39 2 98 58
6 21 20 84 59 48 65
7 68 65 5 74 12 87 67
10 10 96 55 97 27 38 69 77 54 62
7 78 76 30 3 85 1 95
3 39 65 73
5 43 9 64 34 39
7 50 50 8 21 83 17 31
5 94 43 8 5 62
7 19 63 78 92 11 87 90
`

func run(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseLine(line string) ([]int, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty testcase")
	}
	nums := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid int %q: %w", f, err)
		}
		nums[i] = v
	}
	return nums, nil
}

// solve mirrors the logic from 1007A.go so we do not depend on an external binary.
func solve(nums []int) (int, error) {
	if len(nums) == 0 {
		return 0, fmt.Errorf("missing n")
	}
	n := nums[0]
	if n < 0 {
		return 0, fmt.Errorf("negative n")
	}
	if len(nums)-1 < n {
		return 0, fmt.Errorf("need %d numbers, got %d", n, len(nums)-1)
	}
	a := make([]int, n)
	copy(a, nums[1:n+1])
	sort.Ints(a)

	cnt := 0
	j := 0
	for i := 0; i < n; i++ {
		for j < n && a[j] <= a[i] {
			j++
		}
		if j < n {
			cnt++
			j++
		} else {
			break
		}
	}
	return cnt, nil
}

func expectedOutput(line string) (string, error) {
	nums, err := parseLine(line)
	if err != nil {
		return "", err
	}
	ans, err := solve(nums)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(ans), nil
}

func loadTestcases() []string {
	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	var cases []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		cases = append(cases, line)
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	testcases := loadTestcases()
	for idx, line := range testcases {
		want, err := expectedOutput(line)
		if err != nil {
			fmt.Printf("failed to compute expected output on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput: %s\nexpected: %s\ngot: %s\n", idx+1, line, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
