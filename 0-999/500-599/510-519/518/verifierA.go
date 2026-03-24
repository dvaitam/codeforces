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

func expected(s, t string) string {
	b := []byte(s)
	carried := true
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < 'z' {
			b[i]++
			for j := i + 1; j < len(b); j++ {
				b[j] = 'a'
			}
			carried = false
			break
		}
	}
	if carried {
		return "No such string"
	}
	next := string(b)
	if next < t {
		return next
	}
	return "No such string"
}

func generateCase(rng *rand.Rand) (string, string, string) {
	n := rng.Intn(10) + 1
	s := make([]byte, n)
	t := make([]byte, n)
	// generate two random strings of the same length, ensure s < t
	for i := 0; i < n; i++ {
		s[i] = byte('a' + rng.Intn(26))
		t[i] = byte('a' + rng.Intn(26))
	}
	ss, tt := string(s), string(t)
	if ss > tt {
		ss, tt = tt, ss
	}
	if ss == tt {
		// make t strictly greater by setting last char
		t = []byte(tt)
		for i := n - 1; i >= 0; i-- {
			if t[i] < 'z' {
				t[i]++
				tt = string(t)
				break
			}
		}
		if ss >= tt {
			// s is all z's, just pick s="a..a", t="z..z"
			for i := range s {
				s[i] = 'a'
				t[i] = 'z'
			}
			ss, tt = string(s), string(t)
		}
	}
	return ss, tt, expected(ss, tt)
}

func runCase(bin, s, t, exp string) error {
	cmd := exec.Command(bin)
	input := fmt.Sprintf("%s\n%s\n", s, t)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct{ s, t, exp string }{
		{"a", "c", expected("a", "c")},
		{"aaaa", "zzzz", expected("aaaa", "zzzz")},
		{"az", "ba", expected("az", "ba")},
	}
	for i := 0; i < 100; i++ {
		s, t, e := generateCase(rng)
		cases = append(cases, struct{ s, t, exp string }{s, t, e})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.s, tc.t, tc.exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%s\n", i+1, err, tc.s, tc.t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
