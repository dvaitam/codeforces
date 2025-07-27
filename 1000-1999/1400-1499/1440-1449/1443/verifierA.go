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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Test struct {
	input    string
	expected []int
}

func genTest(rng *rand.Rand) Test {
	t := rng.Intn(3) + 1 // number of cases
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	expected := make([]int, 0)
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			expected = append(expected, 4*n-2*(j+1))
		}
	}
	return Test{input: sb.String(), expected: expected}
}

func parseInts(s string) ([]int, error) {
	fields := strings.Fields(strings.TrimSpace(s))
	res := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tcase := genTest(rng)
		out, err := run(bin, tcase.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		nums, err := parseInts(out)
		if err != nil {
			fmt.Printf("test %d bad output: %v\n", i+1, err)
			os.Exit(1)
		}
		if len(nums) != len(tcase.expected) {
			fmt.Printf("test %d expected %d numbers got %d\ninput:\n%s", i+1, len(tcase.expected), len(nums), tcase.input)
			os.Exit(1)
		}
		for j, v := range tcase.expected {
			if nums[j] != v {
				fmt.Printf("test %d failed at position %d\ninput:\n%s\nexpected %d got %d\n", i+1, j+1, tcase.input, v, nums[j])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
