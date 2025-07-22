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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(10) + 1
	nums := make([]int, 2*n)
	for i := range nums {
		nums[i] = rng.Intn(5) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), nums
}

func expectedPossible(nums []int) bool {
	count := make(map[int]int)
	for _, v := range nums {
		count[v]++
	}
	for _, c := range count {
		if c%2 == 1 {
			return false
		}
	}
	return true
}

func verifyOutput(out string, nums []int) error {
	fields := strings.Fields(out)
	if expectedPossible(nums) {
		n := len(nums) / 2
		if len(fields) != 2*n {
			return fmt.Errorf("expected %d numbers, got %d", 2*n, len(fields))
		}
		used := make([]bool, len(nums))
		for i := 0; i < n; i++ {
			idx1, err1 := strconv.Atoi(fields[2*i])
			idx2, err2 := strconv.Atoi(fields[2*i+1])
			if err1 != nil || err2 != nil {
				return fmt.Errorf("failed to parse integers")
			}
			if idx1 < 1 || idx1 > len(nums) || idx2 < 1 || idx2 > len(nums) {
				return fmt.Errorf("index out of range")
			}
			if idx1 == idx2 {
				return fmt.Errorf("pair uses same index")
			}
			if used[idx1-1] || used[idx2-1] {
				return fmt.Errorf("index repeated")
			}
			used[idx1-1], used[idx2-1] = true, true
			if nums[idx1-1] != nums[idx2-1] {
				return fmt.Errorf("values mismatch for pair")
			}
		}
		for _, u := range used {
			if !u {
				return fmt.Errorf("some indices not used")
			}
		}
	} else {
		if len(fields) != 1 || fields[0] != "-1" {
			return fmt.Errorf("expected -1")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, nums := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verifyOutput(out, nums); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
