package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expected(n int) (int, int) {
	for a := n / 2; a >= 1; a-- {
		b := n - a
		if a < b && gcd(a, b) == 1 {
			return a, b
		}
	}
	return 0, 0
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, n := range func() []int {
		arr := make([]int, 0, 100)
		for x := 3; x < 103; x++ {
			arr = append(arr, x)
		}
		return arr
	}() {
		expA, expB := expected(n)
		got, err := run(bin, fmt.Sprintf("%d\n", n))
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var a, b int
		if _, err := fmt.Sscanf(got, "%d %d", &a, &b); err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", i+1, got)
			os.Exit(1)
		}
		if a != expA || b != expB {
			fmt.Printf("test %d failed: n=%d expected %d %d got %d %d\n", i+1, n, expA, expB, a, b)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
