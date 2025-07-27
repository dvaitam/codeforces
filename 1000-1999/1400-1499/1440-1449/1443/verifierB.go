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

func solve(a, b int, s string) int {
	cost := 0
	inBlock := false
	seen := false
	zeros := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			if !inBlock {
				if !seen {
					cost += a
					seen = true
				} else {
					if zeros*b < a {
						cost += zeros * b
					} else {
						cost += a
					}
				}
				inBlock = true
			}
			zeros = 0
		} else {
			if inBlock {
				inBlock = false
				zeros = 1
			} else if seen {
				zeros++
			}
		}
	}
	return cost
}

func genTest(rng *rand.Rand) Test {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	expected := make([]int, 0, t)
	for i := 0; i < t; i++ {
		a := rng.Intn(10) + 1
		b := rng.Intn(10) + 1
		n := rng.Intn(20) + 1
		var s strings.Builder
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				s.WriteByte('0')
			} else {
				s.WriteByte('1')
			}
		}
		str := s.String()
		sb.WriteString(fmt.Sprintf("%d %d\n%s\n", a, b, str))
		expected = append(expected, solve(a, b, str))
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genTest(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		nums, err := parseInts(out)
		if err != nil {
			fmt.Printf("test %d bad output: %v\n", i+1, err)
			os.Exit(1)
		}
		if len(nums) != len(tc.expected) {
			fmt.Printf("test %d expected %d numbers got %d\ninput:\n%s", i+1, len(tc.expected), len(nums), tc.input)
			os.Exit(1)
		}
		for j, v := range tc.expected {
			if nums[j] != v {
				fmt.Printf("test %d failed at case %d\ninput:\n%s\nexpected %d got %d\n", i+1, j+1, tc.input, v, nums[j])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
