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

const mod = 1000003

var codes = map[rune]string{
	'>': "1000",
	'<': "1001",
	'+': "1010",
	'-': "1011",
	'.': "1100",
	',': "1101",
	'[': "1110",
	']': "1111",
}

var charset = []rune{'<', '>', '+', '-', '.', ',', '[', ']'}

func expectedOutput(s string) string {
	result := 0
	for _, ch := range s {
		code, ok := codes[ch]
		if !ok {
			continue
		}
		for _, bit := range code {
			result = (result*2 + int(bit-'0')) % mod
		}
	}
	return fmt.Sprint(result)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(100) + 1
	b := make([]rune, n)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

func runCase(bin, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedOutput(input)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
