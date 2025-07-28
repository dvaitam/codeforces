package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func generateInput() []byte {
	r := rand.New(rand.NewSource(42))
	t := 100
	var buf bytes.Buffer
	fmt.Fprintln(&buf, t)
	for i := 0; i < t; i++ {
		n := r.Intn(11) + 2 // 2..12
		fmt.Fprintln(&buf, n)
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if r.Intn(2) == 0 {
				sb.WriteByte('A')
			} else {
				sb.WriteByte('B')
			}
		}
		fmt.Fprintln(&buf, sb.String())
	}
	return buf.Bytes()
}

func run(cmd *exec.Cmd, input []byte) ([]byte, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	input := generateInput()
	refOut, err := run(exec.Command("go", "run", "1896B.go"), input)
	if err != nil {
		fmt.Println("reference solution error:", err)
		fmt.Print(string(refOut))
		os.Exit(1)
	}
	out, err := run(exec.Command(os.Args[1]), input)
	if err != nil {
		fmt.Println("solution runtime error:", err)
		fmt.Print(string(out))
		os.Exit(1)
	}
	refLines := strings.Split(strings.TrimSpace(string(refOut)), "\n")
	outLines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(refLines) != len(outLines) {
		fmt.Printf("mismatched lines: expected %d got %d\n", len(refLines), len(outLines))
		os.Exit(1)
	}
	for i := range refLines {
		if strings.TrimSpace(refLines[i]) != strings.TrimSpace(outLines[i]) {
			fmt.Printf("mismatch on line %d: expected %q got %q\n", i+1, refLines[i], outLines[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
