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

func expected(s string) string {
	n := len(s)
	b := []byte(s)
	res := make([]int, n)
	for i := 1; i < n; i++ {
		rev := false
		if i == n-1 {
			if b[i] == 'a' {
				rev = true
			}
		} else {
			if b[i] != b[i+1] {
				rev = true
			}
		}
		if rev {
			res[i] = 1
			for l, r := 0, i; l < r; l, r = l+1, r-1 {
				b[l], b[r] = b[r], b[l]
			}
		}
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, input, want string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %q got %q", want, got)
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

	tests := []string{"abba", "aaaa"}
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		b := make([]byte, n)
		for j := range b {
			if rng.Intn(2) == 0 {
				b[j] = 'a'
			} else {
				b[j] = 'b'
			}
		}
		tests = append(tests, string(b))
	}

	for idx, s := range tests {
		input := s + "\n"
		want := expected(s)
		if err := runCase(bin, input, want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
