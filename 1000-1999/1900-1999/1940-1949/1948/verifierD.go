package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

func maxTandem(s string) int {
	b := []byte(s)
	n := len(b)
	for k := n / 2; k >= 1; k-- {
		run := 0
		for i := n - k - 1; i >= 0; i-- {
			c1 := b[i]
			c2 := b[i+k]
			if c1 == c2 || c1 == '?' || c2 == '?' {
				run++
				if run >= k {
					return 2 * k
				}
			} else {
				run = 0
			}
		}
	}
	return 0
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(42)
	const t = 100
	strs := make([]string, t)
	letters := []byte("abcde?")
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			b[j] = letters[rand.Intn(len(letters))]
		}
		strs[i] = string(b)
	}

	var input bytes.Buffer
	fmt.Fprintf(&input, "%d\n", t)
	for i := 0; i < t; i++ {
		fmt.Fprintf(&input, "%s\n", strs[i])
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", i+1)
			os.Exit(1)
		}
		got := scanner.Text()
		want := fmt.Sprintf("%d", maxTandem(strs[i]))
		if got != want {
			fmt.Printf("case %d: expected %s, got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("warning: extra output detected")
	}
	fmt.Println("All tests passed!")
}
