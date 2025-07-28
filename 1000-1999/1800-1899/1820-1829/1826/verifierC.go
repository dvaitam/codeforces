package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func spf(n int) int {
	if n%2 == 0 {
		return 2
	}
	limit := int(math.Sqrt(float64(n)))
	for i := 3; i <= limit; i += 2 {
		if n%i == 0 {
			return i
		}
	}
	return n
}

func solve(n, m int) string {
	if n == 1 || m == 1 {
		return "YES"
	}
	d := spf(n)
	if m >= d {
		return "NO"
	}
	return "YES"
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierC <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(0)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(1000) + 1
		m := rand.Intn(1000) + 1
		input := fmt.Sprintf("1\n%d %d\n", n, m)
		expected := solve(n, m)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s\n", t, err, output)
			os.Exit(1)
		}
		output = strings.ToUpper(output)
		if output != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", t, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
