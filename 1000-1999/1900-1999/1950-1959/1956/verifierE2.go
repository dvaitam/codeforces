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

func solveCase(n int, a []int) string {
	return "0"
}

func genCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var in bytes.Buffer
	var out bytes.Buffer
	fmt.Fprintf(&in, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		fmt.Fprintf(&in, "%d\n", n)
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(11)
			if j > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", a[j])
		}
		in.WriteByte('\n')
		out.WriteString(solveCase(n, a))
		if i+1 < t {
			out.WriteByte('\n')
		}
	}
	return in.String(), strings.TrimSpace(out.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
