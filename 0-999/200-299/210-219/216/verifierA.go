package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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

func compute(a, b, c int) int {
	return a*b + b*c + a*c - a - b - c + 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		a := rand.Intn(999) + 2
		b := rand.Intn(999) + 2
		c := rand.Intn(999) + 2
		input := fmt.Sprintf("%d %d %d\n", a, b, c)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", i)
			os.Exit(1)
		}
		if val != compute(a, b, c) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\n", i, compute(a, b, c), val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
