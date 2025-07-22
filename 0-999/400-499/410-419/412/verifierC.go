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

func expectedC(patterns []string) string {
	if len(patterns) == 0 {
		return ""
	}
	k := len(patterns[0])
	res := make([]rune, k)
	for j := 0; j < k; j++ {
		ch := rune('?')
		conflict := false
		for _, p := range patterns {
			c := rune(p[j])
			if c == '?' {
				continue
			}
			if ch == '?' {
				ch = c
			} else if ch != c {
				conflict = true
				break
			}
		}
		if conflict {
			res[j] = '?'
		} else if ch == '?' {
			res[j] = 'x'
		} else {
			res[j] = ch
		}
	}
	return string(res)
}

func runCase(bin string, patterns []string) error {
	input := fmt.Sprintf("%d\n", len(patterns))
	for _, p := range patterns {
		input += p + "\n"
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(expectedC(patterns))
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func genCase(rng *rand.Rand) []string {
	n := rng.Intn(6) + 1
	l := rng.Intn(10) + 1
	patterns := make([]string, n)
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < n; i++ {
		var sb strings.Builder
		for j := 0; j < l; j++ {
			if rng.Intn(4) == 0 {
				sb.WriteByte('?')
			} else {
				sb.WriteByte(alphabet[rng.Intn(len(alphabet))])
			}
		}
		patterns[i] = sb.String()
	}
	return patterns
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		patterns := genCase(rng)
		if err := runCase(bin, patterns); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%v\n", i+1, err, patterns)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
