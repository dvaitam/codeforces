package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type query struct {
	x  int
	ch byte
}

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, s []byte, qs []query) string {
	pairs := 0
	for i := 0; i < n-1; i++ {
		if s[i] == '.' && s[i+1] == '.' {
			pairs++
		}
	}
	var sb strings.Builder
	for idx, q := range qs {
		if s[q.x] != q.ch {
			if s[q.x] == '.' {
				if q.x > 0 && s[q.x-1] == '.' {
					pairs--
				}
				if q.x+1 < n && s[q.x+1] == '.' {
					pairs--
				}
			}
			s[q.x] = q.ch
			if s[q.x] == '.' {
				if q.x > 0 && s[q.x-1] == '.' {
					pairs++
				}
				if q.x+1 < n && s[q.x+1] == '.' {
					pairs++
				}
			}
		}
		if idx > 0 {
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, "%d", pairs)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	letters := []byte("abc.")
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		m := rng.Intn(8) + 1
		s := make([]byte, n)
		for j := 0; j < n; j++ {
			s[j] = letters[rng.Intn(len(letters))]
		}
		qs := make([]query, m)
		for j := 0; j < m; j++ {
			qs[j] = query{rng.Intn(n), letters[rng.Intn(len(letters))]}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		sb.Write(s)
		sb.WriteByte('\n')
		for _, q := range qs {
			fmt.Fprintf(&sb, "%d %c\n", q.x+1, q.ch)
		}
		input := sb.String()
		want := expected(n, append([]byte(nil), s...), qs)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\n", i+1, want, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
