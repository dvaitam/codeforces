package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(n int64) string {
	if n == 1 || n == 2 {
		return "-1"
	}
	rec := 0
	for n%2 == 0 {
		n /= 2
		rec++
	}
	if n == 1 {
		x, y := int64(3), int64(5)
		for i := 0; i < rec-2; i++ {
			x *= 2
			y *= 2
		}
		return fmt.Sprintf("%d %d", x, y)
	}
	x := n / 2
	y := x + 1
	ans1 := 2 * x * y
	ans2 := x*x + y*y
	for i := 0; i < rec; i++ {
		ans1 *= 2
		ans2 *= 2
	}
	return fmt.Sprintf("%d %d", ans1, ans2)
}

func generateTest() string {
	n := rand.Int63n(1_000_000_000) + 1
	return fmt.Sprintf("%d\n", n)
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	for t := 0; t < 100; t++ {
		input := generateTest()
		var n int64
		fmt.Sscanf(strings.TrimSpace(input), "%d", &n)
		want := expected(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "test", t+1, "error running binary:", err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput: %sexpected: %s\nactual: %s\n", t+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
