package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const t = 100
	for i := 0; i < t; i++ {
		n := rand.Intn(50) + 1
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		cmd := exec.Command(cand)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("candidate run error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(bytes.NewReader(out))
		if !scanner.Scan() {
			fmt.Println("no output on test", i+1)
			os.Exit(1)
		}
		var val int
		fmt.Sscan(scanner.Text(), &val)
		if val != 1 {
			fmt.Printf("wrong answer on test %d: expected 1 got %d\n", i+1, val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
