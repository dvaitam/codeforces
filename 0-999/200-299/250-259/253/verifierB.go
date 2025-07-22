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
	return strings.TrimSpace(out.String()), nil
}

func solveB(nums []int) int {
	sort.Ints(nums)
	n := len(nums)
	maxKeep := 0
	j := 0
	for i := 0; i < n; i++ {
		for j < n && nums[j] <= 2*nums[i] {
			j++
		}
		if j-i > maxKeep {
			maxKeep = j - i
		}
	}
	return n - maxKeep
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 1; tc <= 100; tc++ {
		n := rng.Intn(100) + 1
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i] = rng.Intn(5000) + 1
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
		input := sb.String()
		expect := fmt.Sprintf("%d", solveB(append([]int(nil), nums...)))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", tc, expect, strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
