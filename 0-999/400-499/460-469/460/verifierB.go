package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func powInt(x, a int) int64 {
	res := int64(1)
	for i := 0; i < a; i++ {
		res *= int64(x)
	}
	return res
}

func sumDigits(x int64) int {
	s := 0
	for x > 0 {
		s += int(x % 10)
		x /= 10
	}
	return s
}

func expectedB(a, b, c int) []int {
	var res []int
	for s := 1; s <= 81; s++ {
		x := int64(b)*powInt(s, a) + int64(c)
		if x > 0 && x < 1000000000 && sumDigits(x) == s {
			res = append(res, int(x))
		}
	}
	sort.Ints(res)
	return res
}

func runCase(bin string, a, b, c int) error {
	input := fmt.Sprintf("%d %d %d\n", a, b, c)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	var nums []int
	for scanner.Scan() {
		var x int
		fmt.Sscan(scanner.Text(), &x)
		nums = append(nums, x)
	}
	if len(nums) < 1 {
		return fmt.Errorf("no output")
	}
	n := nums[0]
	if n != len(nums)-1 {
		return fmt.Errorf("expected %d numbers but got %d", n, len(nums)-1)
	}
	nums = nums[1:]
	expect := expectedB(a, b, c)
	if len(nums) != len(expect) {
		return fmt.Errorf("expected %v got %v", expect, nums)
	}
	for i := range nums {
		if nums[i] != expect[i] {
			return fmt.Errorf("expected %v got %v", expect, nums)
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (int, int, int) {
	a := rng.Intn(5) + 1
	b := rng.Intn(10000) + 1
	c := rng.Intn(20001) - 10000
	return a, b, c
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	edges := []struct{ a, b, c int }{
		{1, 1, 0}, {5, 10000, -10000}, {5, 1, 10000}, {3, 100, 0},
	}
	for i, e := range edges {
		if err := runCase(bin, e.a, e.b, e.c); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		a, b, c := generateCase(rng)
		if err := runCase(bin, a, b, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d %d %d\n", i+1, err, a, b, c)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
