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
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign = "-"
		s = s[1:]
	}
	i := 0
	for i < len(s) && s[i] == '0' {
		i++
	}
	s = s[i:]
	if len(s) == 0 {
		return "0"
	}
	bs := []byte(s)
	for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
		bs[i], bs[j] = bs[j], bs[i]
	}
	j := 0
	for j < len(bs) && bs[j] == '0' {
		j++
	}
	bs = bs[j:]
	if len(bs) == 0 {
		return "0"
	}
	if sign != "" {
		return sign + string(bs)
	}
	return string(bs)
}

func generateCase(rng *rand.Rand) (string, string) {
	sign := ""
	if rng.Intn(2) == 0 {
		sign = "-"
	}
	l := rng.Intn(8) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('0' + rng.Intn(10))
	}
	s := sign + string(b)
	return fmt.Sprintf("%s\n", s), expected(s)
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected '%s' got '%s'", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
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
