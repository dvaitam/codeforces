package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

func expected(n int) (int, int) {
	a := int(math.Sqrt(float64(n)))
	for a > 0 {
		if n%a == 0 {
			return a, n / a
		}
		a--
	}
	return 1, n
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(bytes.TrimSpace(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for n := 1; n <= 120; n++ {
		input := fmt.Sprintf("%d\n", n)
		expA, expB := expected(n)
		outStr, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", n, err)
			return
		}
		var a, b int
		if _, err := fmt.Sscanf(outStr, "%d %d", &a, &b); err != nil {
			fmt.Printf("Could not parse output on test %d: %v\nOutput: %s\n", n, err, outStr)
			return
		}
		if a != expA || b != expB {
			fmt.Printf("Wrong answer on test %d: expected %d %d got %d %d\n", n, expA, expB, a, b)
			return
		}
	}
	fmt.Println("All tests passed")
}
