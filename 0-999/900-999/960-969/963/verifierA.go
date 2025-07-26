package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCmd(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	// build reference solution
	refBin := "./refA.bin"
	if err := exec.Command("go", "build", "-o", refBin, "963A.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(1)
	for t := 0; t < 100; t++ {
		n := rand.Int63n(1_000_000_000)
		a := rand.Int63n(1_000_000_000) + 1
		b := rand.Int63n(1_000_000_000) + 1
		k := rand.Intn(10) + 1
		var sb strings.Builder
		for i := 0; i < k; i++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('+')
			} else {
				sb.WriteByte('-')
			}
		}
		s := sb.String()
		input := fmt.Sprintf("%d %d %d %d\n%s\n", n, a, b, k, s)

		exp, err := runCmd(refBin, input)
		if err != nil {
			fmt.Println("reference solution failed:", err)
			os.Exit(1)
		}

		got, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("test %d: candidate runtime error: %v\n", t+1, err)
			os.Exit(1)
		}

		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
