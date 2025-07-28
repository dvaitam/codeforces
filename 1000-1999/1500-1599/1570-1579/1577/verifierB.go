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

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func max(nums []int) int {
	m := nums[0]
	for _, v := range nums[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		n := rng.Intn(10) + 1
		nums := make([]int, n)
		for j := 0; j < n; j++ {
			nums[j] = rng.Intn(1000)
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for j, v := range nums {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		expect := fmt.Sprintf("%d", max(nums))
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", i, expect, out, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
