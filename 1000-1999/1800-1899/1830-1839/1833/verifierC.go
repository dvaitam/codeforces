package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expected(nums []int) string {
	even, odd := 0, 0
	minVal := int(^uint(0) >> 1)
	for _, x := range nums {
		if x%2 == 0 {
			even++
		} else {
			odd++
		}
		if x < minVal {
			minVal = x
		}
	}
	res := "NO"
	if even == len(nums) || odd == len(nums) || minVal%2 == 1 {
		res = "YES"
	}
	return res + "\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1000)
	}
	input := fmt.Sprintf("1\n%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", arr[i])
	}
	input += "\n"
	return input, expected(arr)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
