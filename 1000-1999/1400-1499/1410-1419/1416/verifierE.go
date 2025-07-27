package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveE(a []int) int {
	M := 0
	r := 0
	for _, x := range a {
		if r > 0 && r < x {
			v := x - r
			if v == r {
				M += 2
			} else {
				M++
			}
			r = v
		} else {
			if x%2 == 0 {
				M++
				r = x / 2
			} else {
				r = 1
			}
		}
	}
	return 2*len(a) - M
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(5)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(20) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(100) + 2
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		out, err := runBinary(binary, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", t, err)
			os.Exit(1)
		}
		expected := fmt.Sprintf("%d", solveE(a))
		if out != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\nexpected: %s\ngot: %s\n", t, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
