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

func balance(s []byte) string {
	n := len(s)
	target := n / 3
	a, b, c := 0, 0, 0
	for _, ch := range s {
		switch ch {
		case '0':
			a++
		case '1':
			b++
		case '2':
			c++
		}
	}
	// increase zeros
	for i := 0; i < n && a < target; i++ {
		if s[i] == '1' && b > target {
			s[i], a, b = '0', a+1, b-1
		} else if s[i] == '2' && c > target {
			s[i], a, c = '0', a+1, c-1
		}
	}
	for i := n - 1; i >= 0 && a > target; i-- {
		if s[i] == '0' {
			s[i] = '3'
			a--
		}
	}
	for i := 0; i < n && b < target; i++ {
		if s[i] == '3' {
			s[i] = '1'
			b++
		}
	}
	for i := 0; i < n && b < target; i++ {
		if s[i] == '2' {
			s[i], b, c = '1', b+1, c-1
		}
	}
	for i := n - 1; i >= 0 && b > target; i-- {
		if s[i] == '1' {
			s[i] = '3'
			b--
		}
	}
	for i := 0; i < n; i++ {
		if s[i] == '3' {
			s[i] = '2'
		}
	}
	return string(s)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := (rng.Intn(30) + 1) * 3
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + rng.Intn(3))
	}
	input := fmt.Sprintf("%d\n%s\n", n, string(b))
	ans := balance(append([]byte(nil), b...))
	return input, ans
}

func runCase(bin, input, expected string) error {
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
