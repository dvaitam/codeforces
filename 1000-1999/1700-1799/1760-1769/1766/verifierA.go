package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveA(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return ""
	}
	var out strings.Builder
	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(in, &n)
		pow := 1
		for pow*10 <= n {
			pow *= 10
		}
		digits := 0
		for tmp := n; tmp > 0; tmp /= 10 {
			digits++
		}
		ans := 9*(digits-1) + n/pow
		fmt.Fprintf(&out, "%d\n", ans)
	}
	return strings.TrimSpace(out.String())
}

func genTestA(rng *rand.Rand) string {
	t := rng.Intn(10) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(999999) + 1
		fmt.Fprintf(&buf, "%d\n", n)
	}
	return buf.String()
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestA(rng)
		expect := solveA(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
