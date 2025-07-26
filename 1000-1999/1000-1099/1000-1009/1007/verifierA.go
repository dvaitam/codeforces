package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	ref := "./refA.bin"
	if err := exec.Command("go", "build", "-o", ref, "1007A.go").Run(); err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		want, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput: %s\nexpected: %s\ngot: %s\n", idx, line, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
