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

func solveExpected(n, k int, arr []int) string {
	freq := make(map[int]int)
	for _, v := range arr {
		freq[v]++
	}
	nums := make([]int, 0)
	for v, c := range freq {
		if c >= k {
			nums = append(nums, v)
		}
	}
	if len(nums) == 0 {
		return "-1"
	}
	sort.Ints(nums)
	bestL, bestR := nums[0], nums[0]
	curL, curR := nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[i-1]+1 {
			curR = nums[i]
		} else {
			if curR-curL > bestR-bestL {
				bestL, bestR = curL, curR
			}
			curL, curR = nums[i], nums[i]
		}
	}
	if curR-curL > bestR-bestL {
		bestL, bestR = curL, curR
	}
	return fmt.Sprintf("%d %d", bestL, bestR)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(5) + 1
	arr := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(10) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	exp := solveExpected(n, k, arr)
	return sb.String(), exp
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
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
