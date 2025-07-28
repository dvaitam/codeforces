package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	arr []int
}

func generateTestCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 100)
	// deterministic cases
	cases = append(cases, testCase{arr: []int{1, 1}})
	cases = append(cases, testCase{arr: []int{1, 2}})
	for len(cases) < 100 {
		n := rng.Intn(10) + 2 // 2..11
		arr := make([]int, n)
		if rng.Intn(3) == 0 {
			v := rng.Intn(20) + 1
			for i := range arr {
				arr[i] = v
			}
		} else {
			for i := range arr {
				arr[i] = rng.Intn(20) + 1
			}
		}
		cases = append(cases, testCase{arr: arr})
	}
	return cases
}

func validPartition(a, b, c []int) bool {
	if len(b) == 0 || len(c) == 0 {
		return false
	}
	if len(b)+len(c) != len(a) {
		return false
	}
	cntInput := make(map[int]int)
	for _, v := range a {
		cntInput[v]++
	}
	for _, v := range b {
		cntInput[v]--
		if cntInput[v] < 0 {
			return false
		}
	}
	for _, v := range c {
		cntInput[v]--
		if cntInput[v] < 0 {
			return false
		}
	}
	for _, v := range cntInput {
		if v != 0 {
			return false
		}
	}
	for _, bi := range b {
		for _, cj := range c {
			if bi%cj == 0 {
				return false
			}
		}
	}
	return true
}

func allEqual(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i] != a[0] {
			return false
		}
	}
	return true
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", len(tc.arr))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	output := strings.TrimSpace(out.String())
	if output == "-1" {
		if allEqual(tc.arr) {
			return nil
		}
		return fmt.Errorf("expected valid partition, got -1")
	}
	nums := strings.Fields(output)
	if len(nums) < 3 {
		return fmt.Errorf("output too short: %q", output)
	}
	lb, err1 := strconv.Atoi(nums[0])
	lc, err2 := strconv.Atoi(nums[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("invalid sizes")
	}
	if lb <= 0 || lc <= 0 {
		return fmt.Errorf("empty arrays")
	}
	if len(nums) != 2+lb+lc {
		return fmt.Errorf("expected %d numbers got %d", 2+lb+lc, len(nums))
	}
	b := make([]int, lb)
	c := make([]int, lc)
	for i := 0; i < lb; i++ {
		v, err := strconv.Atoi(nums[2+i])
		if err != nil {
			return fmt.Errorf("bad number %q", nums[2+i])
		}
		b[i] = v
	}
	for i := 0; i < lc; i++ {
		v, err := strconv.Atoi(nums[2+lb+i])
		if err != nil {
			return fmt.Errorf("bad number %q", nums[2+lb+i])
		}
		c[i] = v
	}
	if !validPartition(tc.arr, b, c) {
		return fmt.Errorf("invalid partition for input %v", tc.arr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTestCases()
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
