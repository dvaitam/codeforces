package main

import (
	"bytes"
	"fmt"
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
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	refBin := "./_refC"
	if err := exec.Command("go", "build", "-o", refBin, "392C.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer os.Remove(refBin)

	rand.Seed(1)
	for tc := 1; tc <= 100; tc++ {
		n := rand.Int63n(100000000000000000) + 1
		k := rand.Intn(40) + 1
		input := fmt.Sprintf("%d %d\n", n, k)
		exp, err := runCmd(refBin, input)
		if err != nil {
			fmt.Println("reference solution error:", err)
			return
		}
		got, err := runCmd(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", tc, err)
			return
		}
		if got != exp {
			fmt.Printf("wrong answer on test %d\ninput:%sexpected:%s\ngot:%s\n", tc, input, exp, got)
			return
		}
	}
	fmt.Println("All tests passed!")
}
