package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref := "./refB.bin"
	if err := exec.Command("go", "build", "-o", ref, "1067B.go").Run(); err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scanner.Scan() {
			fmt.Printf("test %d: missing header\n", caseNum+1)
			os.Exit(1)
		}
		header := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(header)
		if len(parts) != 2 {
			fmt.Printf("test %d: invalid header\n", caseNum+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n-1; i++ {
			if !scanner.Scan() {
				fmt.Printf("test %d: missing edge\n", caseNum+1)
				os.Exit(1)
			}
			line := strings.TrimSpace(scanner.Text())
			input.WriteString(line + "\n")
		}
		want, err := run(ref, input.String())
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", caseNum+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", caseNum+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", caseNum+1, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
