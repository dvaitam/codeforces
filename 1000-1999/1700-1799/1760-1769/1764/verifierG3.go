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
		fmt.Println("usage: go run verifierG3.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]int, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 1
		fmt.Fprintln(&input, n)
		perm := rand.Perm(n)
		for j := 0; j < n; j++ {
			val := perm[j] + 1
			fmt.Fprint(&input, val)
			if j+1 < n {
				fmt.Fprint(&input, " ")
			}
			if val == 1 {
				expected[i] = j + 1
			}
		}
		fmt.Fprintln(&input)
	}

	cmd := exec.Command(cand)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("candidate run error:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("missing output for test", i+1)
			os.Exit(1)
		}
		var ans int
		fmt.Sscan(scanner.Text(), &ans)
		if ans != expected[i] {
			fmt.Printf("wrong answer on test %d: expected %d got %d\n", i+1, expected[i], ans)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
