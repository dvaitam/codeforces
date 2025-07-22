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

func expectedURL(s string) string {
	x := -1
	for i := 1; i < len(s); i++ {
		if s[i-1] == 'r' && s[i] == 'u' {
			x = i - 1
		}
	}
	var res strings.Builder
	start := 0
	if strings.HasPrefix(s, "http") {
		res.WriteString("http://")
		start = 4
	} else {
		res.WriteString("ftp://")
		start = 3
	}
	if x > start {
		res.WriteString(s[start:x])
	}
	res.WriteString(".ru")
	if x+2 < len(s) {
		res.WriteString("/")
		res.WriteString(s[x+2:])
	}
	return res.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	protocol := "http"
	if rng.Intn(2) == 0 {
		protocol = "ftp"
	}
	domainLen := rng.Intn(10) + 1
	contextLen := rng.Intn(5)
	build := func(n int) string {
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = byte('a' + rng.Intn(26))
		}
		return string(b)
	}
	domain := build(domainLen)
	context := build(contextLen)
	raw := protocol + domain + "ru" + context
	return raw + "\n", expectedURL(raw)
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
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
