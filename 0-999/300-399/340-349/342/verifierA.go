package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var bin string

func isPossible(nums []int) bool {
	n := len(nums)
	cnt := make([]int, 8)
	for _, v := range nums {
		if v >= 0 && v < len(cnt) {
			cnt[v]++
		} else {
			return false
		}
	}
	if cnt[5] > 0 || cnt[7] > 0 {
		return false
	}
	if cnt[1] != n/3 {
		return false
	}
	if cnt[2]-cnt[4] < 0 || cnt[2]-cnt[4] != cnt[6]-cnt[3] || cnt[6] < cnt[3] {
		return false
	}
	return true
}

func verifyCase(input, output string) error {
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return fmt.Errorf("invalid input n: %v", err)
	}
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &nums[i])
	}
	possible := isPossible(nums)
	out := strings.TrimSpace(output)
	if !possible {
		if out != "-1" {
			return fmt.Errorf("expected -1, got %q", out)
		}
		return nil
	}
	lines := strings.Split(out, "\n")
	if len(lines) != n/3 {
		return fmt.Errorf("expected %d lines, got %d", n/3, len(lines))
	}
	cnt := make([]int, 8)
	for _, v := range nums {
		if v >= 0 && v < len(cnt) {
			cnt[v]++
		}
	}
	for idx, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return fmt.Errorf("line %d: expected 3 numbers", idx+1)
		}
		a := make([]int, 3)
		for j := 0; j < 3; j++ {
			x, err := strconv.Atoi(parts[j])
			if err != nil {
				return fmt.Errorf("line %d: parse int: %v", idx+1, err)
			}
			if x < 1 || x > 7 {
				return fmt.Errorf("line %d: number out of range", idx+1)
			}
			a[j] = x
			cnt[x]--
		}
		if !(a[0] < a[1] && a[1] < a[2]) {
			return fmt.Errorf("line %d: numbers not strictly increasing", idx+1)
		}
		if a[1]%a[0] != 0 || a[2]%a[1] != 0 {
			return fmt.Errorf("line %d: divisibility failed", idx+1)
		}
	}
	for i := 1; i < len(cnt); i++ {
		if cnt[i] != 0 {
			return fmt.Errorf("number %d unused balance %d", i, cnt[i])
		}
	}
	return nil
}

func runBinary(input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateCase(rng *rand.Rand) string {
	n := (rng.Intn(10) + 1) * 3
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", rng.Intn(7)+1)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin = os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := generateCase(rng)
		out, err := runBinary(tc)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		if err := verifyCase(tc, out); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%soutput:\n%s\n", t+1, err, tc, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
