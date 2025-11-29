package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `1 94
2 370
3 754
5 258
10 218
10 37
10 698
3 442
7 823
9 973
6 558
8 515
5 923
1 892
1 373
8 955
6 930
7 434
9 169
9 182
4 237
1 181
6 178
3 523
9 369
9 691
9 187
8 816
7 753
9 929
6 809
10 363
6 880
8 166
7 733
8 671
9 256
8 286
8 513
9 852
6 678
8 922
8 360
10 744
9 742
8 499
4 964
6 835
3 898
10 275
8 317
5 981
9 576
9 520
10 603
7 320
4 501
9 376
10 904
2 804
6 744
1 930
4 763
2 61
10 669
1 280
10 233
2 773
9 140
5 251
4 967
1 434
1 59
6 369
3 256
1 85
2 978
2 26
1 747
1 383
5 131
3 753
3 536
1 395
10 45
4 156
1 5
6 961
10 643
2 293
6 501
1 316
8 565
10 758
1 924
5 774
7 884
10 723
3 485
4 96`

// Embedded reference logic from 121C.go.
func solve(n, k int64) int {
	const capVal = int64(1e18)
	fact := make([]int64, 21)
	fact[0] = 1
	for i := 1; i <= 20; i++ {
		fact[i] = fact[i-1] * int64(i)
		if fact[i] > capVal {
			fact[i] = capVal
		}
	}

	remMin := -1
	for i := 0; i <= 20; i++ {
		if fact[i] >= k {
			remMin = i
			break
		}
	}
	if remMin < 0 || int64(remMin) > n {
		return -1
	}
	m := remMin
	prefixLen := n - int64(m)

	ans := 0
	for i := int64(1); i <= prefixLen; i++ {
		if isLucky(i) {
			ans++
		}
	}

	k--
	rem := make([]int64, m)
	for i := 0; i < m; i++ {
		rem[i] = prefixLen + int64(i) + 1
	}
	for i := 0; i < m; i++ {
		f := fact[m-1-i]
		idx := int(k / f)
		val := rem[idx]
		pos := prefixLen + int64(i) + 1
		if isLucky(pos) && isLucky(val) {
			ans++
		}
		rem = append(rem[:idx], rem[idx+1:]...)
		k %= f
	}
	return ans
}

func isLucky(x int64) bool {
	if x <= 0 {
		return false
	}
	for x > 0 {
		d := x % 10
		if d != 4 && d != 7 {
			return false
		}
		x /= 10
	}
	return true
}

type testCase struct {
	n int64
	k int64
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields", idx+1)
		}
		n, err1 := strconv.ParseInt(parts[0], 10, 64)
		k, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: invalid numbers", idx+1)
		}
		res = append(res, testCase{n: n, k: k})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.k)

		expected := fmt.Sprintf("%d", solve(tc.n, tc.k))

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", i+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed. Expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
