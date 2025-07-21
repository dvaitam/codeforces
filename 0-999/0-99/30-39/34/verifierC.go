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

func generateCase() string {
	length := rand.Intn(100) + 1
	var nums []string
	for i := 0; i < length; i++ {
		nums = append(nums, fmt.Sprintf("%d", rand.Intn(1000)+1))
		if rand.Intn(3) == 0 { // duplicate sometimes
			nums = append(nums, nums[len(nums)-1])
		}
	}
	return strings.Join(nums, ",")
}

func expected(in string) string {
	parts := strings.Split(in, ",")
	seen := make(map[int]bool)
	var nums []int
	for _, p := range parts {
		v, _ := strconv.Atoi(p)
		if !seen[v] {
			seen[v] = true
			nums = append(nums, v)
		}
	}
	if len(nums) == 0 {
		return ""
	}
	var ranges []string
	start := nums[0]
	prev := nums[0]
	for i := 1; i < len(nums); i++ {
		v := nums[i]
		if v == prev+1 {
			prev = v
			continue
		}
		if start == prev {
			ranges = append(ranges, fmt.Sprintf("%d", start))
		} else {
			ranges = append(ranges, fmt.Sprintf("%d-%d", start, prev))
		}
		start = v
		prev = v
	}
	if start == prev {
		ranges = append(ranges, fmt.Sprintf("%d", start))
	} else {
		ranges = append(ranges, fmt.Sprintf("%d-%d", start, prev))
	}
	return strings.Join(ranges, ",")
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := generateCase()
		exp := expected(tc)
		got, err := run(bin, tc)
		if err != nil {
			fmt.Printf("case %d: error executing binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed\ninput:%s\nexpected:%s\ngot:%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
