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

func expectedE(s string) int64 {
	cnt := make([]int, 10)
	for _, ch := range s {
		cnt[ch-'0']++
	}
	digits := []int{}
	for d := 0; d < 10; d++ {
		if cnt[d] > 0 {
			digits = append(digits, d)
		}
	}
	maxLen := len(s)
	fact := make([]int64, maxLen+1)
	fact[0] = 1
	for i := 1; i <= maxLen; i++ {
		fact[i] = fact[i-1] * int64(i)
	}
	cur := make([]int, 10)
	var ans int64
	var dfs func(int, int)
	dfs = func(pos, length int) {
		if pos == len(digits) {
			L := length
			total := fact[L]
			for _, d := range digits {
				total /= fact[cur[d]]
			}
			if cur[0] > 0 {
				t := fact[L-1] / fact[cur[0]-1]
				for _, d := range digits {
					if d == 0 {
						continue
					}
					t /= fact[cur[d]]
				}
				total -= t
			}
			ans += total
			return
		}
		d := digits[pos]
		for c := 1; c <= cnt[d]; c++ {
			cur[d] = c
			dfs(pos+1, length+c)
		}
		cur[d] = 0
	}
	dfs(0, 0)
	return ans
}

func generateCaseE(rng *rand.Rand) (string, int64) {
	length := rng.Intn(18) + 1
	var sb strings.Builder
	sb.WriteByte(byte('1' + rng.Intn(9)))
	for i := 1; i < length; i++ {
		sb.WriteByte(byte('0' + rng.Intn(10)))
	}
	input := sb.String() + "\n"
	expected := expectedE(sb.String())
	return input, expected
}

func runCaseE(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		if err := runCaseE(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
