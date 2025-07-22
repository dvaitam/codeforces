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
	ans := uint64(1)
	n := len(s)
	for i := 0; i < n; {
		if s[i] == '9' {
			i++
			continue
		}
		if i+1 < n && int(s[i]-'0')+int(s[i+1]-'0') == 9 {
			j := i
			for j+1 < n && s[j] != '9' && int(s[j]-'0')+int(s[j+1]-'0') == 9 {
				j++
			}
			length := j - i + 1
			if length%2 == 1 {
				ans *= uint64(length/2 + 1)
			}
			i = j + 1
		} else {
			i++
		}
	}
	return fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := []byte("123456789")
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(50) + 1
		b := make([]byte, n)
		for i := range b {
			b[i] = digits[rng.Intn(len(digits))]
		}
		s := string(b)
		input := s + "\n"
		exp := expected(s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
