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

func solveCase(n int64) int64 {
	var a, b int64
	multA, multB := int64(1), int64(1)
	pos := 0
	for n > 0 {
		digit := n % 10
		if pos%2 == 0 {
			a += digit * multA
			multA *= 10
		} else {
			b += digit * multB
			multB *= 10
		}
		pos++
		n /= 10
	}
	return (a+1)*(b+1) - 2
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) int64 {
	return int64(rng.Intn(1_000_000_000-2) + 2)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := genTest(rng)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input.String())
			os.Exit(1)
		}
		expected := solveCase(n)
		if out != fmt.Sprint(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", i+1, expected, out, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
