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

type testCase struct {
	input    string
	expected int
}

func solveF(strs []string) int {
	n := len(strs)
	net := make([]int, n)
	minPref := make([]int, n)
	freq := make([]map[int]int, n)
	for i := 0; i < n; i++ {
		sum := 0
		minv := 0
		freq[i] = make(map[int]int)
		for _, ch := range strs[i] {
			if ch == '(' {
				sum++
			} else {
				sum--
			}
			if sum < minv {
				minv = sum
			}
			if sum == minv {
				freq[i][sum]++
			}
		}
		net[i] = sum
		minPref[i] = minv
	}
	size := 1 << n
	bal := make([]int, size)
	for mask := 1; mask < size; mask++ {
		lb := mask & -mask
		idx := 0
		for (lb>>idx)&1 == 0 {
			idx++
		}
		bal[mask] = bal[mask^lb] + net[idx]
	}
	const negInf = -1 << 60
	dp := make([]int, size)
	for i := range dp {
		dp[i] = negInf
	}
	dp[0] = 0
	ans := 0
	for mask := 0; mask < size; mask++ {
		if dp[mask] < 0 {
			continue
		}
		cur := bal[mask]
		if dp[mask] > ans {
			ans = dp[mask]
		}
		for i := 0; i < n; i++ {
			if mask>>i&1 == 0 {
				add := freq[i][-cur]
				if cur+minPref[i] >= 0 {
					nmask := mask | (1 << i)
					if dp[nmask] < dp[mask]+add {
						dp[nmask] = dp[mask] + add
					}
				} else {
					if dp[mask]+add > ans {
						ans = dp[mask] + add
					}
				}
			}
		}
	}
	return ans
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	letters := []rune{'(', ')'}
	for i := 0; i < 100; i++ {
		n := rng.Intn(3) + 1 // 1..3
		strs := make([]string, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			l := rng.Intn(5) + 1
			b := make([]rune, l)
			for k := 0; k < l; k++ {
				b[k] = letters[rng.Intn(2)]
			}
			s := string(b)
			strs[j] = s
			sb.WriteString(fmt.Sprintf("%s\n", s))
		}
		exp := solveF(strs)
		cases[i] = testCase{input: sb.String(), expected: exp}
	}
	return cases
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		var got int
		fmt.Sscan(out, &got)
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
