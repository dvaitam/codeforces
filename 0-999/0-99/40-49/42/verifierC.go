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

func applyOps(nums []int, ops []string) error {
	if len(ops) > 1000 {
		return fmt.Errorf("too many operations")
	}
	for _, op := range ops {
		if len(op) < 2 {
			return fmt.Errorf("bad operation %q", op)
		}
		pos := int(op[1] - '1')
		if pos < 0 || pos > 3 {
			return fmt.Errorf("bad position in %q", op)
		}
		j := (pos + 1) % 4
		switch op[0] {
		case '+':
			nums[pos]++
			nums[j]++
		case '/':
			if nums[pos]%2 != 0 || nums[j]%2 != 0 {
				return fmt.Errorf("invalid divide in %q", op)
			}
			nums[pos] /= 2
			nums[j] /= 2
		default:
			return fmt.Errorf("invalid op %q", op)
		}
	}
	for _, v := range nums {
		if v != 1 {
			return fmt.Errorf("final numbers not all ones: %v", nums)
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, []int) {
	nums := make([]int, 4)
	for i := range nums {
		nums[i] = rng.Intn(20) + 1
	}
	input := fmt.Sprintf("%d %d %d %d\n", nums[0], nums[1], nums[2], nums[3])
	cp := make([]int, 4)
	copy(cp, nums)
	return input, cp
}

func runCase(bin, input string, nums []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 1 && strings.TrimSpace(lines[0]) == "-1" {
		return fmt.Errorf("reported impossible but solution exists")
	}
	ops := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		ops = append(ops, ln)
	}
	if err := applyOps(nums, ops); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, nums := generateCase(rng)
		if err := runCase(bin, in, nums); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
