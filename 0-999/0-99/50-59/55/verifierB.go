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

func dfs(nums []int64, ops []string, idx int) int64 {
	if len(nums) == 1 {
		return nums[0]
	}
	minRes := int64(1<<63 - 1)
	op := ops[idx]
	n := len(nums)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			var t int64
			if op == "*" {
				t = nums[i] * nums[j]
			} else {
				t = nums[i] + nums[j]
			}
			next := make([]int64, 0, n-1)
			for k := 0; k < n; k++ {
				if k == i || k == j {
					continue
				}
				next = append(next, nums[k])
			}
			next = append(next, t)
			res := dfs(next, ops, idx+1)
			if res < minRes {
				minRes = res
			}
		}
	}
	return minRes
}

func expected(a, b, c, d int64, ops []string) string {
	nums := []int64{a, b, c, d}
	res := dfs(nums, ops, 0)
	return fmt.Sprintf("%d", res)
}

func generateCase() (string, string) {
	a := rand.Intn(1001)
	b := rand.Intn(1001)
	c := rand.Intn(1001)
	d := rand.Intn(1001)
	ops := make([]string, 3)
	for i := 0; i < 3; i++ {
		if rand.Intn(2) == 0 {
			ops[i] = "+"
		} else {
			ops[i] = "*"
		}
	}
	input := fmt.Sprintf("%d %d %d %d\n%s %s %s\n", a, b, c, d, ops[0], ops[1], ops[2])
	return input, expected(int64(a), int64(b), int64(c), int64(d), ops)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input, exp := generateCase()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
