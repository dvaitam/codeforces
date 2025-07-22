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

func reduce(s string) string {
	b := []byte(s)
	for {
		n := len(b)
		found := false
		for L := 1; 2*L <= n && !found; L++ {
			for i := 0; i+2*L <= n; i++ {
				eq := true
				for k := 0; k < L; k++ {
					if b[i+k] != b[i+L+k] {
						eq = false
						break
					}
				}
				if eq {
					nb := make([]byte, 0, n-L)
					nb = append(nb, b[:i+L]...)
					nb = append(nb, b[i+2*L:]...)
					b = nb
					found = true
					break
				}
			}
		}
		if !found {
			break
		}
	}
	return string(b)
}

func runCase(bin string, s string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(s + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := reduce(s)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"a", "aa", "ab", "aba", "abab", "abcabc"}
	for i := 0; i < 94; i++ {
		n := rng.Intn(10) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = byte('a' + rng.Intn(3))
		}
		cases = append(cases, string(b))
	}
	for i, s := range cases {
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %q\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
