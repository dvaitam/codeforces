package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func expected(nums []int) string {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return fmt.Sprintf("%d\n", sum)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		nums := make([]int, n)
		var sb strings.Builder
		w := bufio.NewWriter(&sb)
		for i := 0; i < n; i++ {
			nums[i] = rand.Intn(1001)
			fmt.Fprintln(w, nums[i])
		}
		w.Flush()
		in := sb.String()
		want := expected(nums)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != strings.TrimSpace(want) {
			fmt.Printf("test %d failed: expected %q got %q\n", t+1, strings.TrimSpace(want), out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
