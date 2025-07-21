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

	lines := strings.Split(strings.TrimSpace(input), "\n")
	var n, k int
	fmt.Sscan(lines[0], &n, &k)
	p := strings.TrimSpace(lines[1])
	mask := strings.TrimSpace(lines[2])
	ls := len(p)
	required := strings.Count(mask, "1")

	check := func(s string) bool {
		if len(s) != n {
			return false
		}
		for i := 0; i < n; i++ {
			if s[i] < 'a' || s[i] >= 'a'+byte(k) {
				return false
			}
		}
		cnt := 0
		for i := 0; i+ls <= n; i++ {
			if s[i:i+ls] == p {
				cnt++
			}
		}
		if cnt != required {
			return false
		}
		for i := 0; i < len(mask); i++ {
			if mask[i] == '1' && s[i:i+ls] != p {
				return false
			}
			if mask[i] == '0' && s[i:i+ls] == p {
				return false
			}
		}
		return true
	}

	if expected == "" {
		if got != "No solution" {
			return fmt.Errorf("expected No solution got %s", got)
		}
	} else {
		if got == "No solution" {
			return fmt.Errorf("expected solution got No solution")
		}
		if !check(got) {
			return fmt.Errorf("invalid solution: %s", got)
		}
	}
	return nil
}

func construct(n, k int, p, mask string) string {
	ls := len(p)
	st := make([]byte, n)
	for i := range st {
		st[i] = '!'
	}
	expected := 0
	for i := 0; i < len(mask); i++ {
		if mask[i] == '1' {
			expected++
			for j := 0; j < ls; j++ {
				idx := i + j
				if st[idx] != '!' && st[idx] != p[j] {
					return ""
				}
				st[idx] = p[j]
			}
		}
	}
	check := func(s []byte) bool {
		count := 0
		for i := 0; i+ls <= n; i++ {
			match := true
			for j := 0; j < ls; j++ {
				if s[i+j] != p[j] {
					match = false
					break
				}
			}
			if match {
				count++
			}
		}
		return count == expected
	}
	for ch := byte('a'); ch < byte('a')+byte(k); ch++ {
		s := make([]byte, n)
		for i := 0; i < n; i++ {
			if st[i] == '!' {
				s[i] = ch
			} else {
				s[i] = st[i]
			}
		}
		if check(s) {
			return string(s)
		}
	}
	return ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		k := rng.Intn(5) + 2
		if k > 26 {
			k = 26
		}
		pLen := rng.Intn(min(n, 5)) + 1
		letters := []rune("abcdefghijklmnopqrstuvwxyz")[:k]
		var p strings.Builder
		for j := 0; j < pLen; j++ {
			p.WriteRune(letters[rng.Intn(len(letters))])
		}
		maskLen := n - pLen + 1
		var maskBuilder strings.Builder
		for j := 0; j < maskLen; j++ {
			if rng.Intn(2) == 0 {
				maskBuilder.WriteByte('0')
			} else {
				maskBuilder.WriteByte('1')
			}
		}
		pStr := p.String()
		mask := maskBuilder.String()
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		sb.WriteString(pStr + "\n")
		sb.WriteString(mask + "\n")
		input := sb.String()
		expected := construct(n, k, pStr, mask)
		if err := runCase(exe, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
