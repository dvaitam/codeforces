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

func solveF(a, b []int) string {
	freq := make(map[int]int)
	for _, x := range a {
		for x%2 == 0 {
			x /= 2
		}
		freq[x]++
	}
	for _, y := range b {
		for y%2 == 0 {
			y /= 2
		}
		for y > 0 && freq[y] == 0 {
			y /= 2
		}
		if y == 0 {
			return "NO"
		}
		freq[y]--
	}
	return "YES"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(1)
	const T = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	expected := make([]string, T)

	for i := 0; i < T; i++ {
		n := rand.Intn(10) + 1
		fmt.Fprintln(&input, n)
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(100) + 1
			fmt.Fprintf(&input, "%d ", a[j])
		}
		fmt.Fprintln(&input)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			b[j] = rand.Intn(100) + 1
			fmt.Fprintf(&input, "%d ", b[j])
		}
		fmt.Fprintln(&input)
		expected[i] = solveF(a, b)
	}

	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	i := 0
	for scanner.Scan() {
		if i >= T {
			fmt.Println("binary produced extra output")
			os.Exit(1)
		}
		got := strings.TrimSpace(scanner.Text())
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("reading output failed:", err)
		os.Exit(1)
	}
	if i < T {
		fmt.Println("binary produced insufficient output")
		os.Exit(1)
	}

	fmt.Println("all tests passed")
}
