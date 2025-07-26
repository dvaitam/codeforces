package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveC(a []int, s string) string {
	n := len(a)
	i := 0
	for i < n-1 {
		if s[i] == '1' {
			start := i
			for i < n-1 && s[i] == '1' {
				i++
			}
			end := i
			// sort a[start:end+1]
			for x := start; x <= end; x++ {
				for y := x + 1; y <= end; y++ {
					if a[y] < a[x] {
						a[x], a[y] = a[y], a[x]
					}
				}
			}
		} else {
			i++
		}
	}
	for i := 0; i < n; i++ {
		if a[i] != i+1 {
			return "NO"
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(3)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(10) + 2
		a := rand.Perm(n)
		for i := range a {
			a[i]++
		}
		b := make([]int, n)
		copy(b, a)
		sbytes := make([]byte, n-1)
		for i := range sbytes {
			if rand.Intn(2) == 0 {
				sbytes[i] = '0'
			} else {
				sbytes[i] = '1'
			}
		}
		s := string(sbytes)
		expected := solveC(b, s)
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for i, v := range a {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, v)
		}
		fmt.Fprintln(&input)
		fmt.Fprintln(&input, s)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Failed to run binary:", err)
			os.Exit(1)
		}
		answer := strings.TrimSpace(string(out))
		if answer != expected {
			fmt.Printf("Test %d failed: expected %s got %s\n", t+1, expected, answer)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
