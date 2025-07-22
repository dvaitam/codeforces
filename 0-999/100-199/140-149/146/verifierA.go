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

type Test struct {
	n      int
	s      string
	expect string
}

func expected(n int, s string) string {
	if len(s) != n {
		return "NO"
	}
	half := n / 2
	sum1, sum2 := 0, 0
	for i := 0; i < n; i++ {
		ch := s[i]
		if ch != '4' && ch != '7' {
			return "NO"
		}
		d := int(ch - '0')
		if i < half {
			sum1 += d
		} else {
			sum2 += d
		}
	}
	if sum1 == sum2 {
		return "YES"
	}
	return "NO"
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(25)*2 + 2 // even between 2 and 50
	digits := make([]byte, n)
	for i := 0; i < n; i++ {
		digits[i] = byte(rng.Intn(10)) + '0'
	}
	ticket := string(digits)
	input := fmt.Sprintf("%d\n%s\n", n, ticket)
	exp := expected(n, ticket)
	return input, exp
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s", strings.TrimSpace(exp), got)
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
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
