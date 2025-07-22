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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expectedAnswer(nums []int64) string {
	g := nums[0]
	maxV := nums[0]
	for _, v := range nums[1:] {
		g = gcd(g, v)
		if v > maxV {
			maxV = v
		}
	}
	total := maxV / g
	moves := total - int64(len(nums))
	if moves%2 == 1 {
		return "Alice"
	}
	return "Bob"
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(9) + 2 // 2..10
	used := make(map[int64]bool)
	nums := make([]int64, n)
	for i := 0; i < n; i++ {
		var v int64
		for {
			v = int64(rng.Intn(1000) + 1)
			if !used[v] {
				used[v] = true
				break
			}
		}
		nums[i] = v
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseNums(input string) []int64 {
	parts := strings.Fields(input)
	n := 0
	fmt.Sscan(parts[0], &n)
	nums := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Sscan(parts[i+1], &nums[i])
	}
	return nums
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{
		"2\n2 3\n",
		"3\n1 2 3\n",
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for idx, tc := range cases {
		nums := parseNums(tc)
		expect := expectedAnswer(nums)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
