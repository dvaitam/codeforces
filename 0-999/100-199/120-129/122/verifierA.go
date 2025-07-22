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

func isAlmostLucky(n int) bool {
	lucky := []int{4, 7, 44, 47, 74, 77, 444, 447, 474, 477, 744, 747, 774, 777}
	for _, v := range lucky {
		if n%v == 0 {
			return true
		}
	}
	return false
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	special := []int{1, 4, 7, 16, 44, 47, 74, 77, 444, 447, 474, 477, 744, 747, 774, 777, 1000}
	for i, n := range special {
		exp := "NO"
		if isAlmostLucky(n) {
			exp = "YES"
		}
		input := fmt.Sprintf("%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("special case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("special case %d failed: n=%d expected %s got %s\n", i+1, n, exp, got)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		n := rng.Intn(1000) + 1
		exp := "NO"
		if isAlmostLucky(n) {
			exp = "YES"
		}
		input := fmt.Sprintf("%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: n=%d expected %s got %s\n", i+1, n, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
