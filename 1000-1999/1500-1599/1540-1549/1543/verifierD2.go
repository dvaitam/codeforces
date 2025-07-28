package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

func xorK(a, b, k int) int {
	res := 0
	mul := 1
	for a > 0 || b > 0 {
		res += ((a%k + b%k) % k) * mul
		a /= k
		b /= k
		mul *= k
	}
	return res
}

type testD2 struct {
	n int
	k int
	x int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(42)
	const T = 100
	tests := make([]testD2, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(10) + 1
		k := rand.Intn(99) + 2
		x := rand.Intn(n)
		tests[i] = testD2{n: n, k: k, x: x}
	}
	cmd := exec.Command(binary)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "start error:", err)
		os.Exit(1)
	}
	writer := bufio.NewWriter(stdin)
	reader := bufio.NewReader(stdout)
	fmt.Fprintln(writer, T)
	for _, tc := range tests {
		fmt.Fprintf(writer, "%d %d\n", tc.n, tc.k)
	}
	writer.Flush()
	for _, tc := range tests {
		x := tc.x
		for i := 0; i < tc.n; i++ {
			var q int
			if _, err := fmt.Fscan(reader, &q); err != nil {
				fmt.Fprintln(os.Stderr, "failed to read query:", err)
				cmd.Process.Kill()
				os.Exit(1)
			}
			if q == x {
				fmt.Fprintln(writer, 1)
				writer.Flush()
				break
			} else {
				fmt.Fprintln(writer, 0)
				writer.Flush()
				x = xorK(x, q, tc.k)
			}
		}
	}
	stdin.Close()
	writer.Flush()
	if err := cmd.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "binary exited with error:", err)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
