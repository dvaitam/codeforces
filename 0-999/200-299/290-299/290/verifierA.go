package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func factorial(n int) string {
	res := big.NewInt(1)
	for i := 2; i <= n; i++ {
		res.Mul(res, big.NewInt(int64(i)))
	}
	return res.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n := rand.Intn(40) + 1
		expected := factorial(n)

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%d\n", n))
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", i+1, err, string(out))
			os.Exit(1)
		}
		got := string(bytes.TrimSpace(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: input %d expected %s got %s\n", i+1, n, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
