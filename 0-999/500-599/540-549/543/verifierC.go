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

// solveC computes the minimum cost to learn at least one character from each
// string. For any column we may either learn a single character from some
// string or learn the same character for a group of strings in that column with
// the most expensive one being free. The straightforward solution is a DP over
// subsets of strings.
func solveC(n, m int, s []string, a [][]int) int {
	const inf = int(1e9)
	dp := make([]int, 1<<n)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0
	for mask := 0; mask < (1 << n); mask++ {
		if dp[mask] == inf {
			continue
		}
		// find the first string that is not yet covered
		var l int
		for l = 0; l < n; l++ {
			if mask&(1<<l) == 0 {
				break
			}
		}
		if l == n {
			continue
		}
		for c := 0; c < m; c++ {
			// option 1: learn character only for string l
			nmask := mask | (1 << l)
			cost := dp[mask] + a[l][c]
			if cost < dp[nmask] {
				dp[nmask] = cost
			}

			// option 2: learn this character for all strings with the
			// same letter in column c, keeping the most expensive free
			nmask = mask
			sum, mx := 0, 0
			ch := s[l][c]
			for i := 0; i < n; i++ {
				if s[i][c] == ch {
					nmask |= 1 << i
					sum += a[i][c]
					if a[i][c] > mx {
						mx = a[i][c]
					}
				}
			}
			cost = dp[mask] + sum - mx
			if cost < dp[nmask] {
				dp[nmask] = cost
			}
		}
	}
	return dp[(1<<n)-1]
}

func genRandomString(m int) string {
	b := make([]byte, m)
	for i := 0; i < m; i++ {
		b[i] = byte('a' + rand.Intn(3))
	}
	return string(b)
}

func genTestC() (string, int) {
	n := rand.Intn(3) + 1
	m := rand.Intn(3) + 1
	s := make([]string, n)
	for i := 0; i < n; i++ {
		s[i] = genRandomString(m)
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			a[i][j] = rand.Intn(10)
		}
	}
	input := fmt.Sprintf("%d %d\n", n, m)
	for i := 0; i < n; i++ {
		input += s[i] + "\n"
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", a[i][j])
		}
		input += "\n"
	}
	expected := solveC(n, m, s, a)
	return input, expected
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 1; t <= 100; t++ {
		input, expected := genTestC()
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected: %d\nGot: %s\n", t, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
