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

func runProg(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func solve(nums [5]int) int {
	freq := make(map[int]int)
	sum := 0
	for _, v := range nums {
		sum += v
		freq[v]++
	}
	best := 0
	for v, c := range freq {
		if c >= 2 {
			cand := v * 2
			if c >= 3 {
				cand = v * 3
			}
			if cand > best {
				best = cand
			}
		}
	}
	return sum - best
}

func generateCase(rng *rand.Rand) string {
	nums := [5]int{}
	for i := 0; i < 5; i++ {
		nums[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", nums[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInput(input string) ([5]int, error) {
	var nums [5]int
	parts := strings.Fields(input)
	if len(parts) != 5 {
		return nums, fmt.Errorf("invalid input")
	}
	for i := 0; i < 5; i++ {
		fmt.Sscan(parts[i], &nums[i])
	}
	return nums, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		input := generateCase(rng)
		expectedNums, err := parseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse generated input: %v", err)
			os.Exit(1)
		}
		expected := fmt.Sprint(solve(expectedNums))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
