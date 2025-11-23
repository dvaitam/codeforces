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

// Naive implementation to serve as reference
func solveNaive(t, sBegin, sEnd string) int {
	n := len(t)
	nb := len(sBegin)
	ne := len(sEnd)
	
	if nb > n || ne > n {
		return 0
	}
	
	minLen := nb
	if ne > minLen {
		minLen = ne
	}
	if nb + ne <= n && nb + ne > minLen {
		// They might overlap, but the total length must at least accommodate the max(nb, ne).
		// Actually, the condition is that the substring starts with sBegin and ends with sEnd.
		// The length of the substring must be at least max(len(sBegin), len(sEnd)).
		// Overlap is allowed. e.g. sBegin="aba", sEnd="bab", substring "abab" is valid.
		// length is 4. max(3,3) is 3. 4 >= 3. Correct.
	}


unique := make(map[string]struct{})
	
	for i := 0; i <= n-nb; i++ {
		if t[i:i+nb] == sBegin {
			for j := i + minLen; j <= n; j++ {
				// Check end
				// The substring is t[i:j]. It ends with sEnd if t[j-ne:j] == sEnd.
				// BUT we must ensure that the length j-i is sufficient.
				// We already ensured j >= i + minLen, where minLen = max(len(sBegin), len(sEnd)).
				// So j-i >= len(sEnd) is guaranteed.
				if t[j-ne:j] == sEnd {
				
unique[t[i:j]] = struct{}{}
				}
			}
		}
	}
	return len(unique)
}

func randStr(rng *rand.Rand, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(4)))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 5
	s := randStr(rng, n)
	b := randStr(rng, rng.Intn(3)+1)
	e := randStr(rng, rng.Intn(3)+1)
	// The problem input format is "t \n sBegin \n sEnd \n" (on separate lines usually, 
	// or space separated if using Fscan for all).
	// The solution provided uses reader.ReadString('\n') so it expects newlines.
	// Let's format with newlines to be safe for bufio.ReadString or Fscan.
	input := fmt.Sprintf("%s\n%s\n%s\n", s, b, e)
	exp := fmt.Sprintf("%d", solveNaive(s, b, e))
	return input, exp
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Edge cases from problem description or logic
	edge := []struct{ s, b, e string }{{"abc", "a", "c"},{"aaaa", "aa", "aa"},{"ababab", "ab", "ab"},{"abcde", "x", "y"},{"zzz", "z", "z"}}
	for i, tc := range edge {
		input := fmt.Sprintf("%s\n%s\n%s\n", tc.s, tc.b, tc.e)
		exp := fmt.Sprintf("%d", solveNaive(tc.s, tc.b, tc.e))
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\ninput:\n%s", i+1, err, strings.ReplaceAll(input, "\n", " "))
			os.Exit(1)
		}
	}
	
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, strings.ReplaceAll(in, "\n", " "))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}