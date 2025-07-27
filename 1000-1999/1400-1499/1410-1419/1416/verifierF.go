package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runProg(prog string, args []string, input string) (string, error) {
	cmd := exec.Command(prog, args...)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(6)
	for tcase := 1; tcase <= 100; tcase++ {
		n := rand.Intn(3) + 1
		m := rand.Intn(3) + 1
		size := n * m
		a := make([]int, size)
		for i := range a {
			a[i] = rand.Intn(50) + 2
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", a[i*m+j]))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		candOut, err := runProg(binary, nil, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", tcase, err)
			os.Exit(1)
		}
		refOut, err := runProg("go", []string{"run", "1416F.go"}, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", tcase, err)
			os.Exit(1)
		}
		if candOut != refOut {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\nexpected: %s\ngot: %s\n", tcase, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
