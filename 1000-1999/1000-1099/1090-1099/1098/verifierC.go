package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() string {
	exe := "refC_bin"
	if _, err := os.Stat(exe); err == nil {
		return "./" + exe
	}
	cmd := exec.Command("go", "build", "-o", exe, "1098C.go")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	return "./" + exe
}

func run(path string, input []byte) ([]byte, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return out, err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := buildRef()
	rand.Seed(3)

	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 2
		maxSum := n * (n + 1) / 2
		s := rand.Intn(maxSum-n+1) + n
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", n, s)
		input := buf.Bytes()

		exp, err := run(ref, input)
		if err != nil {
			fmt.Println("reference run error:", err)
			os.Exit(1)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(string(exp)) != strings.TrimSpace(string(out)) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:\n%sGot:\n%s", t, string(input), string(exp), string(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
