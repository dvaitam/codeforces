package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) int {
	return rng.Intn(20) + 1
}

func checkSequence(n int, nums []int) error {
	if len(nums) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(nums))
	}
	prev := 0
	for i, x := range nums {
		if x < 1 || x > 10000000 {
			return fmt.Errorf("number out of range: %d", x)
		}
		if i > 0 && x <= prev {
			return fmt.Errorf("sequence not strictly increasing")
		}
		prev = x
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if nums[j]%nums[i] == 0 {
				return fmt.Errorf("%d divides %d", nums[i], nums[j])
			}
		}
	}
	return nil
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	var nums []int
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid integer %q", scanner.Text())
		}
		nums = append(nums, v)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := checkSequence(n, nums); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := generateCase(rng)
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
