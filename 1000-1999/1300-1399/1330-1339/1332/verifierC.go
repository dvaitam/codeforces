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

func solveCase(n, k int, s string) int {
	freq := make([][26]int, k)
	for i := 0; i < n; i++ {
		c := s[i] - 'a'
		idx := i % k
		freq[idx][c]++
	}
	m := n / k
	ans := 0
	for i := 0; i <= (k-1)/2; i++ {
		j := k - 1 - i
		counts := [26]int{}
		groupSize := m
		if i != j {
			groupSize = 2 * m
			for ch := 0; ch < 26; ch++ {
				counts[ch] = freq[i][ch] + freq[j][ch]
			}
		} else {
			for ch := 0; ch < 26; ch++ {
				counts[ch] = freq[i][ch]
			}
		}
		maxFreq := 0
		for ch := 0; ch < 26; ch++ {
			if counts[ch] > maxFreq {
				maxFreq = counts[ch]
			}
		}
		ans += groupSize - maxFreq
	}
	return ans
}

func generateCase(rng *rand.Rand) (int, int, string) {
	k := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	n := k * m
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return n, k, string(b)
}

func runCase(bin string, n, k int, s string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d %d\n%s\n", n, k, s))

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := solveCase(n, k, s)
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, s := generateCase(rng)
		if err := runCase(bin, n, k, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n%s\n", i+1, err, n, k, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
