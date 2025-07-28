package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solve(arr []int64) int64 {
	n := len(arr)
	var g int64
	for i := 0; i < n/2; i++ {
		diff := arr[i] - arr[n-1-i]
		if diff < 0 {
			diff = -diff
		}
		g = gcd(g, diff)
	}
	return g
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
		fmt.Println("usage: verifierB <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(0)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(20) + 1
		arr := make([]int64, n)
		for i := range arr {
			arr[i] = rand.Int63n(1000)
		}
		input := fmt.Sprintf("1\n%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		expected := fmt.Sprint(solve(arr))
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s\n", t, err, output)
			os.Exit(1)
		}
		if output != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", t, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
