package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expected(nums []int64) string {
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	n := len(nums)
	ans := int64(0)
	j := 0
	for i := 0; i < n; i++ {
		if j < i {
			j = i
		}
		for j+1 < n && nums[j+1] <= 2*nums[i] {
			j++
		}
		ans += int64(j - i)
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	nums := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		val := int64(rng.Intn(2000) - 1000)
		nums[i] = abs(val)
		sb.WriteString(fmt.Sprintf("%d", val))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	exp := expected(nums)
	return sb.String(), exp
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
