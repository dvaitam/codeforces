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

func solveA(n int64) string {
	names := []string{"Sheldon", "Leonard", "Penny", "Rajesh", "Howard"}
	group := int64(1)
	count := group * int64(len(names))
	for n > count {
		n -= count
		group *= 2
		count = group * int64(len(names))
	}
	idx := (n - 1) / group
	return names[idx]
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
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		n := rand.Int63n(1_000_000_000) + 1
		expected := solveA(n)
		input := fmt.Sprintf("%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: n=%d expected %s got %s\n", i+1, n, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
