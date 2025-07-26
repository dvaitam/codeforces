package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func maxDiff(nums []int) int {
	m := 0
	for i := 0; i < len(nums)-1; i++ {
		d := nums[i] - nums[i+1]
		if d < 0 {
			if -d > m {
				m = -d
			}
		} else if d > m {
			m = d
		}
	}
	return m
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		n := rand.Intn(9) + 2
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i] = rand.Intn(100) + 1
		}
		fields := make([]string, n)
		for i, v := range nums {
			fields[i] = fmt.Sprintf("%d", v)
		}
		in := strings.Join(fields, " ") + "\n"
		want := fmt.Sprintf("%d", maxDiff(nums))
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != want {
			fmt.Printf("test %d failed: expected %q got %q\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
