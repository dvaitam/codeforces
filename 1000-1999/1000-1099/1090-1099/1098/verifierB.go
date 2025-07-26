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
	exe := "refB_bin"
	if _, err := os.Stat(exe); err == nil {
		return "./" + exe
	}
	cmd := exec.Command("go", "build", "-o", exe, "1098B.go")
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
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := buildRef()
	rand.Seed(2)

	letters := []byte{'A', 'G', 'C', 'T'}

	for t := 1; t <= 100; t++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(4) + 2
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				row[j] = letters[rand.Intn(4)]
			}
			grid[i] = row
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			buf.Write(grid[i])
			buf.WriteByte('\n')
		}
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
