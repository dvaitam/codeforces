package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "317B.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(2)
	for i := 0; i < 100; i++ {
		n := rand.Intn(50)
		q := rand.Intn(5) + 1
		input := fmt.Sprintf("%d %d\n", n, q)
		for j := 0; j < q; j++ {
			x := rand.Intn(11) - 5
			y := rand.Intn(11) - 5
			input += fmt.Sprintf("%d %d\n", x, y)
		}
		exp, err := runRef(input)
		if err != nil {
			fmt.Println("reference run error:", err)
			return
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("binary run error on test %d: %v\n", i+1, err)
			return
		}
		if exp != got {
			fmt.Printf("mismatch on test %d\ninput:\n%sexpected:\n%s\n got:\n%s\n", i+1, input, exp, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
