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

func isKLucky(ticket string, k int) bool {
	if len(ticket) != 8 {
		return false
	}
	digits := make([]int, 8)
	for i := 0; i < 8; i++ {
		if ticket[i] < '0' || ticket[i] > '9' {
			return false
		}
		digits[i] = int(ticket[i] - '0')
	}
	dp := make([][]map[int]struct{}, 9)
	for i := range dp {
		dp[i] = make([]map[int]struct{}, 9)
	}
	for i := 0; i < 8; i++ {
		dp[i][i+1] = map[int]struct{}{digits[i]: {}}
	}
	for length := 2; length <= 8; length++ {
		for l := 0; l+length <= 8; l++ {
			r := l + length
			m := make(map[int]struct{})
			for mid := l + 1; mid < r; mid++ {
				left := dp[l][mid]
				right := dp[mid][r]
				for a := range left {
					for b := range right {
						m[a+b] = struct{}{}
						m[a-b] = struct{}{}
						m[a*b] = struct{}{}
					}
				}
			}
			dp[l][r] = m
		}
	}
	_, ok := dp[0][8][k]
	return ok
}

func runCase(bin string, k, m int) error {
	input := fmt.Sprintf("%d %d\n", k, m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != m {
		return fmt.Errorf("expected %d lines, got %d", m, len(lines))
	}
	seen := make(map[string]struct{})
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if _, ok := seen[t]; ok {
			return fmt.Errorf("duplicate ticket %s", t)
		}
		seen[t] = struct{}{}
		if !isKLucky(t, k) {
			return fmt.Errorf("ticket %s is not %d-lucky", t, k)
		}
	}
	return nil
}

func genCase(rng *rand.Rand) (int, int) {
	k := rng.Intn(1000)
	m := rng.Intn(3) + 1
	return k, m
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := [][2]int{{0, 1}, {1, 1}}
	for i := 0; i < 100; i++ {
		k, m := genCase(rng)
		cases = append(cases, [2]int{k, m})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc[0], tc[1]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
