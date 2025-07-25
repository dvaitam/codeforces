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

func expectedC(s string) int {
	n := len(s)
	freqTotal := make(map[byte]bool)
	for i := 0; i < n; i++ {
		freqTotal[s[i]] = true
	}
	totalTypes := len(freqTotal)
	freq := make(map[byte]int)
	have := 0
	left := 0
	best := n
	for right := 0; right < n; right++ {
		c := s[right]
		freq[c]++
		if freq[c] == 1 {
			have++
		}
		for have == totalTypes {
			if right-left+1 < best {
				best = right - left + 1
			}
			lch := s[left]
			freq[lch]--
			if freq[lch] == 0 {
				have--
			}
			left++
		}
	}
	return best
}

func generateCaseC(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func runCaseC(bin string, s string) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n%s\n", len(s), s))
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		return fmt.Errorf("invalid output: %s", out.String())
	}
	expect := expectedC(s)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
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
		s := generateCaseC(rng)
		if err := runCaseC(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
