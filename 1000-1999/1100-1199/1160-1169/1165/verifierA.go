package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const numTestsA = 100

func solveA(n, x, y int, s string) int {
	ans := 0
	for i := 0; i < x; i++ {
		target := byte('0')
		if i == y {
			target = '1'
		}
		if s[n-1-i] != target {
			ans++
		}
	}
	return ans
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= numTestsA; t++ {
		n := rand.Intn(100) + 2 // >=2
		x := rand.Intn(n-1) + 1
		y := rand.Intn(x)
		b := make([]byte, n)
		b[0] = '1'
		for i := 1; i < n; i++ {
			if rand.Intn(2) == 0 {
				b[i] = '0'
			} else {
				b[i] = '1'
			}
		}
		s := string(b)
		input := fmt.Sprintf("%d %d %d\n%s\n", n, x, y, s)
		expect := solveA(n, x, y, s)
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d failed to run: %v\noutput:%s\n", t, err, out)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fmt.Printf("test %d: no output\n", t)
			os.Exit(1)
		}
		var got int
		fmt.Sscanf(fields[0], "%d", &got)
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\ninput:%s\n", t, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
