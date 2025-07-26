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
	words []string
}

func solveCase(words []string) string {
	roots := make(map[string]struct{})
	for _, w := range words {
		seen := [26]bool{}
		for _, ch := range w {
			if ch >= 'a' && ch <= 'z' {
				seen[ch-'a'] = true
			}
		}
		var root []byte
		for i := 0; i < 26; i++ {
			if seen[i] {
				root = append(root, byte('a'+i))
			}
		}
		roots[string(root)] = struct{}{}
	}
	return fmt.Sprintf("%d\n", len(roots))
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	words := make([]string, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(8) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(5))
		}
		words[i] = string(b)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, w := range words {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(w)
	}
	sb.WriteByte('\n')
	return sb.String(), solveCase(words)
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s", exp, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
