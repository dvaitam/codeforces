package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func bestCounts(a, b, c string) (int, int) {
	freqA := [26]int{}
	for i := range a {
		freqA[a[i]-'a']++
	}
	freqB := [26]int{}
	for i := range b {
		freqB[b[i]-'a']++
	}
	freqC := [26]int{}
	for i := range c {
		freqC[c[i]-'a']++
	}
	maxB := len(a) / len(b)
	for i := 0; i < 26; i++ {
		if freqB[i] > 0 {
			if v := freqA[i] / freqB[i]; v < maxB {
				maxB = v
			}
		}
	}
	bestX, bestY, bestSum := 0, 0, 0
	for x := 0; x <= maxB; x++ {
		rem := [26]int{}
		remOK := true
		for i := 0; i < 26; i++ {
			rem[i] = freqA[i] - x*freqB[i]
			if rem[i] < 0 {
				remOK = false
				break
			}
		}
		if !remOK {
			break
		}
		y := len(a) / len(c)
		for i := 0; i < 26; i++ {
			if freqC[i] > 0 {
				if v := rem[i] / freqC[i]; v < y {
					y = v
				}
			}
		}
		if x+y > bestSum {
			bestSum = x + y
			bestX = x
			bestY = y
		}
	}
	return bestX, bestY
}

func countOccurrences(s, b, c string) int {
	n := len(s)
	dp := make([]int, n+1)
	lb, lc := len(b), len(c)
	for i := 0; i < n; i++ {
		if dp[i+1] < dp[i] {
			dp[i+1] = dp[i]
		}
		if i+1 >= lb && s[i-lb+1:i+1] == b {
			if dp[i+1] < dp[i-lb+1]+1 {
				dp[i+1] = dp[i-lb+1] + 1
			}
		}
		if i+1 >= lc && s[i-lc+1:i+1] == c {
			if dp[i+1] < dp[i-lc+1]+1 {
				dp[i+1] = dp[i-lc+1] + 1
			}
		}
	}
	return dp[n]
}

func freqString(s string) [26]int {
	var f [26]int
	for i := 0; i < len(s); i++ {
		f[s[i]-'a']++
	}
	return f
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Printf("case %d: invalid line\n", idx)
			os.Exit(1)
		}
		a, b, c := parts[0], parts[1], parts[2]
		bestX, bestY := bestCounts(a, b, c)
		best := bestX + bestY
		input := fmt.Sprintf("%s\n%s\n%s\n", a, b, c)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		ans := strings.TrimSpace(out.String())
		if len(ans) != len(a) {
			fmt.Printf("case %d: expected length %d got %d\n", idx, len(a), len(ans))
			os.Exit(1)
		}
		if freqString(ans) != freqString(a) {
			fmt.Printf("case %d: output is not a permutation of input\n", idx)
			os.Exit(1)
		}
		got := countOccurrences(ans, b, c)
		if got != best {
			fmt.Printf("case %d failed: expected %d occurrences got %d\n", idx, best, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
