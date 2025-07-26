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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin := "./refC.bin"
	if err := exec.Command("go", "build", "-o", refBin, "963C.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(3)
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		used := make(map[[2]int]bool)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			var w, h int
			for {
				w = rand.Intn(20) + 1
				h = rand.Intn(20) + 1
				if !used[[2]int{w, h}] {
					used[[2]int{w, h}] = true
					break
				}
			}
			c := rand.Intn(100) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d\n", w, h, c))
		}
		input := sb.String()

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
