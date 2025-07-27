package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveC(x string) (string, string) {
	n := len(x)
	a := make([]byte, n)
	b := make([]byte, n)
	a[0], b[0] = '1', '1'
	broken := false
	for i := 1; i < n; i++ {
		switch x[i] {
		case '0':
			a[i], b[i] = '0', '0'
		case '1':
			if !broken {
				a[i], b[i] = '1', '0'
				broken = true
			} else {
				a[i], b[i] = '0', '1'
			}
		case '2':
			if !broken {
				a[i], b[i] = '1', '1'
			} else {
				a[i], b[i] = '0', '2'
			}
		}
	}
	return string(a), string(b)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(30) + 1
		sb := make([]byte, n)
		sb[0] = '2'
		for i := 1; i < n; i++ {
			sb[i] = byte('0' + rand.Intn(3))
		}
		x := string(sb)
		input := fmt.Sprintf("1\n%d\n%s\n", n, x)
		a, b := solveC(x)
		expect := strings.TrimSpace(a + "\n" + b)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", t, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected\n%s\ngot\n%s\ninput:\n%s", t, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
