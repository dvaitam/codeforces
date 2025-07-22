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

func expectedAnswer(s string, n int) (int, string, bool) {
	freq := make([]int, 26)
	maxf := 0
	distinct := 0
	for i := 0; i < len(s); i++ {
		c := s[i] - 'a'
		if freq[c] == 0 {
			distinct++
		}
		freq[c]++
		if freq[c] > maxf {
			maxf = freq[c]
		}
	}
	if distinct > n {
		return -1, "", false
	}
	ans := -1
	for k := 1; k <= maxf; k++ {
		need := 0
		for _, f := range freq {
			if f > 0 {
				need += (f + k - 1) / k
			}
		}
		if need <= n {
			ans = k
			break
		}
	}
	if ans == -1 {
		return -1, "", false
	}
	sheet := make([]byte, 0, n)
	used := 0
	for i, f := range freq {
		if f > 0 {
			cnt := (f + ans - 1) / ans
			for j := 0; j < cnt; j++ {
				sheet = append(sheet, byte('a'+i))
			}
			used += cnt
		}
	}
	for used < n {
		sheet = append(sheet, 'a')
		used++
	}
	return ans, string(sheet), true
}

func runCase(bin, s string, n int) error {
	input := fmt.Sprintf("%s\n%d\n", s, n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	expK, expT, ok := expectedAnswer(s, n)
	if !ok {
		if len(lines) != 1 || strings.TrimSpace(lines[0]) != "-1" {
			return fmt.Errorf("expected -1 got %q", out.String())
		}
		return nil
	}
	if len(lines) != 2 {
		return fmt.Errorf("expected two lines, got %d", len(lines))
	}
	gotK, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("failed to parse k: %v", err)
	}
	gotT := strings.TrimSpace(lines[1])
	if gotK != expK || gotT != expT {
		return fmt.Errorf("expected %d %s got %d %s", expK, expT, gotK, gotT)
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int) {
	slen := rng.Intn(10) + 1
	n := rng.Intn(10) + 1
	b := make([]byte, slen)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b), n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// predefined cases
	predefined := []struct {
		s string
		n int
	}{
		{"abc", 2},
		{"aaa", 1},
		{"ab", 1},
		{"a", 5},
	}
	for i, tc := range predefined {
		if err := runCase(bin, tc.s, tc.n); err != nil {
			fmt.Fprintf(os.Stderr, "predefined case %d failed: %v\ninput:\n%s\n%d\n", i+1, err, tc.s, tc.n)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		s, n := generateCase(rng)
		if err := runCase(bin, s, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%d\n", i+1, err, s, n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
