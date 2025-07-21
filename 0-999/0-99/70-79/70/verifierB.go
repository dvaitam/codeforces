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

func solve(input string) string {
	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")
	var n int
	fmt.Sscan(lines[0], &n)
	text := lines[1]
	var sentences []string
	start := 0
	for i := 0; i < len(text); i++ {
		c := text[i]
		if (c == '.' || c == '!' || c == '?') && i+1 < len(text) && text[i+1] == ' ' {
			sentences = append(sentences, text[start:i+1])
			start = i + 2
			i++
		}
	}
	if start < len(text) {
		sentences = append(sentences, text[start:])
	}
	msgs := 0
	curr := 0
	for _, s := range sentences {
		L := len(s)
		if L > n {
			return "Impossible"
		}
		if curr == 0 {
			curr = L
		} else if curr+1+L <= n {
			curr += 1 + L
		} else {
			msgs++
			curr = L
		}
	}
	if curr > 0 {
		msgs++
	}
	return fmt.Sprintf("%d", msgs)
}

var punct = []byte{'.', '!', '?'}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(254) + 2 // 2..255
	sentences := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < sentences; i++ {
		L := rng.Intn(n*2) + 1
		for j := 0; j < L; j++ {
			sb.WriteByte(byte('a' + rng.Intn(3)))
		}
		sb.WriteByte(punct[rng.Intn(len(punct))])
		if i+1 != sentences {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
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
		input := genCase(rng)
		expected := solve(input)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
