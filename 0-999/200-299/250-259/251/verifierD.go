package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseD struct {
	nums []uint64
	ans  []int
}

func bruteD(nums []uint64) []int {
	n := len(nums)
	bestVal := uint64(0)
	bestX1 := uint64(0)
	bestMask := 0
	for mask := 0; mask < (1 << n); mask++ {
		var x1, x2 uint64
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				x1 ^= nums[i]
			} else {
				x2 ^= nums[i]
			}
		}
		val := x1 + x2
		if val > bestVal || (val == bestVal && x1 < bestX1) {
			bestVal = val
			bestX1 = x1
			bestMask = mask
		}
	}
	res := make([]int, n)
	for i := 0; i < n; i++ {
		if bestMask&(1<<i) != 0 {
			res[i] = 1
		} else {
			res[i] = 2
		}
	}
	return res
}

func genCaseD() testCaseD {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10) + 1
	nums := make([]uint64, n)
	for i := range nums {
		nums[i] = uint64(rand.Intn(30))
	}
	return testCaseD{nums, bruteD(nums)}
}

func buildInputD(cs []testCaseD) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cs))
	for _, c := range cs {
		fmt.Fprintln(&sb, len(c.nums))
		for i, v := range c.nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutputD(out string, n int) ([]int, bool) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, false
	}
	res := make([]int, n)
	for i, f := range fields {
		var v int
		fmt.Sscan(f, &v)
		if v != 1 && v != 2 {
			return nil, false
		}
		res[i] = v
	}
	return res, true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := make([]testCaseD, 100)
	for i := range cases {
		cases[i] = genCaseD()
	}
	input := buildInputD(cases)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outputs := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(outputs) != len(cases) {
		fmt.Printf("expected %d lines got %d\n", len(cases), len(outputs))
		os.Exit(1)
	}
	for i, line := range outputs {
		res, ok := parseOutputD(line, len(cases[i].nums))
		if !ok {
			fmt.Printf("bad output format on case %d\n", i+1)
			os.Exit(1)
		}
		expect := cases[i].ans
		for j := range expect {
			if res[j] != expect[j] {
				fmt.Printf("mismatch on case %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
