package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	num      string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func formatFinancial(s string) string {
	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}
	intPart := s
	fracPart := ""
	if idx := strings.IndexByte(s, '.'); idx >= 0 {
		intPart = s[:idx]
		fracPart = s[idx+1:]
	}
	var b strings.Builder
	n := len(intPart)
	for i, ch := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(ch)
	}
	formattedInt := b.String()
	if len(fracPart) < 2 {
		fracPart += strings.Repeat("0", 2-len(fracPart))
	} else if len(fracPart) > 2 {
		fracPart = fracPart[:2]
	}
	result := fmt.Sprintf("$%s.%s", formattedInt, fracPart)
	if negative {
		result = fmt.Sprintf("(%s)", result)
	}
	return result
}

func randomNumber(rng *rand.Rand) string {
	negative := rng.Intn(2) == 0
	intLen := rng.Intn(15) + 1
	intPart := make([]byte, intLen)
	intPart[0] = byte('1' + rng.Intn(9))
	for i := 1; i < intLen; i++ {
		intPart[i] = byte('0' + rng.Intn(10))
	}
	if rng.Intn(5) == 0 {
		intPart = []byte("0")
	}
	fracLen := rng.Intn(6)
	fracPart := make([]byte, fracLen)
	for i := 0; i < fracLen; i++ {
		fracPart[i] = byte('0' + rng.Intn(10))
	}
	if string(intPart) == "0" && fracLen == 0 {
		negative = false
	}
	s := ""
	if negative {
		s += "-"
	}
	s += string(intPart)
	if fracLen > 0 {
		s += "." + string(fracPart)
	}
	return s
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := range cases {
		num := randomNumber(rng)
		exp := formatFinancial(num)
		cases[i] = testCase{num: num, expected: exp}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		input := tc.num + "\n"
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:%sexpected:%s\nactual:%s\n", i+1, input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
