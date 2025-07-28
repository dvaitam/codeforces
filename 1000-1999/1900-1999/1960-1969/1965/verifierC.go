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

func run(bin, input string) (string, error) {
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

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func solve(n int, s string) int {
	for L := 1; L <= n; L++ {
		t := s[:L]
		tRev := reverse(t)
		ok := true
		for start := 0; start < n; start += L {
			end := start + L
			if end > n {
				seg := s[start:n]
				prefix := t[:n-start]
				prefRev := tRev[len(tRev)-(n-start):]
				if seg != prefix && seg != prefRev {
					ok = false
					break
				}
			} else {
				seg := s[start:end]
				if seg != t && seg != tRev {
					ok = false
					break
				}
			}
		}
		if ok {
			return L
		}
	}
	return n
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(15) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	var s strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			s.WriteByte('0')
		} else {
			s.WriteByte('1')
		}
	}
	str := s.String()
	sb.WriteString(str)
	sb.WriteByte('\n')
	return sb.String(), solve(n, str)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != fmt.Sprintf("%d", exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
