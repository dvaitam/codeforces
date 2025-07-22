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

func isPalindrome(t string) bool {
	i, j := 0, len(t)-1
	for i < j {
		if t[i] != t[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func isSubsequence(s, t string) bool {
	j := 0
	for i := 0; i < len(s) && j < len(t); i++ {
		if s[i] == t[j] {
			j++
		}
	}
	return j == len(t)
}

func lpsString(s string) string {
	n := len(s)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for i := n - 1; i >= 0; i-- {
		dp[i][i] = 1
		for j := i + 1; j < n; j++ {
			if s[i] == s[j] {
				dp[i][j] = dp[i+1][j-1] + 2
			} else {
				if dp[i+1][j] >= dp[i][j-1] {
					dp[i][j] = dp[i+1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}
	// reconstruct
	i, j := 0, n-1
	left := make([]byte, 0, dp[0][n-1])
	right := make([]byte, 0, dp[0][n-1])
	for i <= j {
		if i == j {
			left = append(left, s[i])
			break
		}
		if s[i] == s[j] {
			left = append(left, s[i])
			right = append(right, s[j])
			i++
			j--
		} else if dp[i+1][j] >= dp[i][j-1] {
			i++
		} else {
			j--
		}
	}
	for k := len(right) - 1; k >= 0; k-- {
		left = append(left, right[k])
	}
	return string(left)
}

func expectedLength(s string) int {
	cnt := [26]int{}
	for i := 0; i < len(s); i++ {
		cnt[s[i]-'a']++
	}
	for _, c := range cnt {
		if c >= 100 {
			return 100
		}
	}
	if len(s) == 0 {
		return 0
	}
	lps := lpsString(s)
	if len(lps) > 100 {
		return 100
	}
	return len(lps)
}

func runCase(bin string, s string) error {
	input := fmt.Sprintf("%s\n", s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expLen := expectedLength(s)
	if len(got) != expLen {
		return fmt.Errorf("expected length %d got %d", expLen, len(got))
	}
	if !isPalindrome(got) {
		return fmt.Errorf("output is not a palindrome: %s", got)
	}
	if !isSubsequence(s, got) {
		return fmt.Errorf("output is not a subsequence: %s", got)
	}
	return nil
}

func generateCase(rng *rand.Rand) string {
	if rng.Intn(5) == 0 {
		// case with many 'a'
		n := 120
		return strings.Repeat("a", n)
	}
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		s := generateCase(rng)
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
