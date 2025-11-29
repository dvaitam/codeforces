package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	k int
}

const testcasesRaw = `100
138 73
868 98
65 33
121 64
780 58
484 84
389 27
97 63
30 50
444 78
781 99
3 90
457 35
739 30
606 14
924 41
32 3
27 84
555 2
962 49
703 28
993 55
744 4
541 29
783 57
962 64
567 30
354 30
694 29
780 59
976 38
949 3
427 72
945 83
103 24
645 93
881 38
124 96
341 93
997 92
513 55
520 86
195 39
291 76
997 64
867 65
403 76
874 5
492 32
762 52
425 86
178 47
562 90
795 87
756 48
89 57
680 66
111 100
168 67
861 51
380 63
751 4
481 6
316 91
869 79
608 75
404 83
175 22
515 30
13 99
205 70
943 71
238 52
527 45
976 74
362 59
932 35
676 71
624 94
6 50
803 95
525 17
532 100
575 27
437 8
493 47
584 71
205 65
424 63
833 46
425 45
2 69
554 80
806 79
340 59
615 4
824 30
651 23
564 75
186 12`

func parseTestcases(raw string) []testCase {
	fields := strings.Fields(raw)
	if len(fields) < 1 {
		panic("no testcase data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(fmt.Sprintf("bad test count: %v", err))
	}
	fields = fields[1:]
	if len(fields)%2 != 0 || len(fields)/2 != t {
		panic("testcase count mismatch")
	}
	res := make([]testCase, 0, t)
	for i := 0; i < len(fields); i += 2 {
		n, err := strconv.Atoi(fields[i])
		if err != nil {
			panic(fmt.Sprintf("bad n at pair %d: %v", i/2+1, err))
		}
		k, err := strconv.Atoi(fields[i+1])
		if err != nil {
			panic(fmt.Sprintf("bad k at pair %d: %v", i/2+1, err))
		}
		res = append(res, testCase{n: n, k: k})
	}
	return res
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func validRepresentation(n, k int, nums []int) bool {
	if len(nums) != k {
		return false
	}
	sum := 0
	parity := nums[0] % 2
	for _, v := range nums {
		if v <= 0 || v%2 != parity {
			return false
		}
		sum += v
	}
	return sum == n
}

func possible(n, k int) bool {
	if n%2 == k%2 && n >= k {
		return true
	}
	if n%2 == 0 && n >= 2*k {
		return true
	}
	return false
}

func runCase(bin string, n, k int) error {
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	out, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	if strings.ToLower(fields[0]) == "no" {
		if possible(n, k) {
			return fmt.Errorf("expected YES but got NO")
		}
		return nil
	}
	if strings.ToLower(fields[0]) != "yes" {
		return fmt.Errorf("first word should be YES or NO")
	}
	if !possible(n, k) {
		return fmt.Errorf("expected NO but got YES")
	}
	nums := make([]int, 0, k)
	for _, f := range fields[1:] {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number %q", f)
		}
		nums = append(nums, v)
	}
	if !validRepresentation(n, k, nums) {
		return fmt.Errorf("invalid representation")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := parseTestcases(testcasesRaw)
	for i, tc := range tests {
		if err := runCase(bin, tc.n, tc.k); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
