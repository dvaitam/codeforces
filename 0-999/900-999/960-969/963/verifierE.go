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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin := "./refE.bin"
	if err := exec.Command("go", "build", "-o", refBin, "963E.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(5)
	for t := 0; t < 100; t++ {
		R := rand.Intn(10)
		a1 := rand.Intn(1000) + 1
		a2 := rand.Intn(1000) + 1
		a3 := rand.Intn(1000) + 1
		a4 := rand.Intn(1000) + 1
		input := fmt.Sprintf("%d %d %d %d %d\n", R, a1, a2, a3, a4)

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
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
