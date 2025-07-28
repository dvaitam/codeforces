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

func solveB(s string) int {
	days := 1
	present := make(map[byte]bool)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !present[c] {
			if len(present) == 3 {
				days++
				for k := range present {
					delete(present, k)
				}
			}
			present[c] = true
		}
	}
	return days
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(1)
	const T = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	expected := make([]string, T)
	for i := 0; i < T; i++ {
		length := rand.Intn(20) + 1
		b := make([]byte, length)
		for j := 0; j < length; j++ {
			b[j] = byte('a' + rand.Intn(26))
		}
		s := string(b)
		fmt.Fprintln(&input, s)
		expected[i] = fmt.Sprintf("%d", solveB(s))
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
