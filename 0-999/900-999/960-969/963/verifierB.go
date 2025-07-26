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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin := "./refB.bin"
	if err := exec.Command("go", "build", "-o", refBin, "963B.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(2)
	for t := 0; t < 100; t++ {
		n := rand.Intn(50) + 1
		parents := make([]int, n)
		for i := 1; i < n; i++ {
			parents[i] = rand.Intn(i) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
				sb.WriteString(fmt.Sprintf("%d", parents[i]))
			} else {
				sb.WriteString("0")
			}
		}
		sb.WriteByte('\n')
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
