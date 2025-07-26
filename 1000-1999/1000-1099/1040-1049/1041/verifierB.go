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
	return a
}

func solveB(a, b, x, y int64) int64 {
	g := gcd(x, y)
	x /= g
	y /= g
	c := a / x
	d := b / y
	if c < d {
		return c
	}
	return d
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	for t := 1; t <= 100; t++ {
		a := rand.Int63n(1000000000) + 1
		b := rand.Int63n(1000000000) + 1
		x := rand.Int63n(1000000000) + 1
		y := rand.Int63n(1000000000) + 1
		input := fmt.Sprintf("%d %d %d %d\n", a, b, x, y)
		expect := fmt.Sprintf("%d", solveB(a, b, x, y))
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", t, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
