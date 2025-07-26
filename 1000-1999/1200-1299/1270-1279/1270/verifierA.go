package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveA(n, k1, k2 int, a, b []int) string {
	for _, x := range a {
		if x == n {
			return "YES"
		}
	}
	return "NO"
}

func generateTest() (int, int, int, []int, []int) {
	n := rand.Intn(99) + 2 // 2..100
	k1 := rand.Intn(n-1) + 1
	k2 := n - k1
	perm := rand.Perm(n)
	a := make([]int, k1)
	b := make([]int, k2)
	for i := 0; i < k1; i++ {
		a[i] = perm[i] + 1
	}
	for i := 0; i < k2; i++ {
		b[i] = perm[k1+i] + 1
	}
	return n, k1, k2, a, b
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const tests = 100

	var input bytes.Buffer
	fmt.Fprintln(&input, tests)
	expected := make([]string, tests)
	for t := 0; t < tests; t++ {
		n, k1, k2, a, b := generateTest()
		fmt.Fprintf(&input, "%d %d %d\n", n, k1, k2)
		for i, x := range a {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, x)
		}
		input.WriteByte('\n')
		for i, x := range b {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, x)
		}
		input.WriteByte('\n')
		expected[t] = solveA(n, k1, k2, a, b)
	}

	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("binary execution failed: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for i := 0; i < tests; i++ {
		if !scanner.Scan() {
			fmt.Printf("insufficient output, expected %d lines\n", tests)
			os.Exit(1)
		}
		got := strings.TrimSpace(scanner.Text())
		if strings.ToUpper(got) != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("warning: extra output ignored")
	}
	fmt.Println("all tests passed")
}
