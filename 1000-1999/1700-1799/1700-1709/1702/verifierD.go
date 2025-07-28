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

func solveD(w string, p int) string {
	freq := make([]int, 26)
	total := 0
	for _, ch := range w {
		idx := int(ch - 'a')
		freq[idx]++
		total += idx + 1
	}
	keep := make([]int, 26)
	copy(keep, freq)
	for c := 25; c >= 0 && total > p; c-- {
		for keep[c] > 0 && total > p {
			keep[c]--
			total -= c + 1
		}
	}
	var res []byte
	for _, ch := range w {
		idx := int(ch - 'a')
		if keep[idx] > 0 {
			res = append(res, byte(ch))
			keep[idx]--
		}
	}
	return string(res)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		w := string(b)
		p := rand.Intn(26*length + 1)
		fmt.Fprintf(&input, "%s\n%d\n", w, p)
		expected[i] = solveD(w, p)
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
