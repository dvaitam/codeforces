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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func verifyOutput(n int, out string) bool {
	parts := strings.Fields(out)
	if len(parts) < 1 {
		return false
	}
	k, err := strconv.Atoi(parts[0])
	if err != nil || k < 1 || k > n {
		return false
	}
	if len(parts)-1 != k {
		return false
	}
	if k == 0 {
		return false
	}
	digit, err := strconv.Atoi(parts[1])
	if err != nil || digit < 1 || digit > 9 {
		return false
	}
	sum := 0
	for i := 1; i <= k; i++ {
		d, err := strconv.Atoi(parts[i])
		if err != nil || d != digit {
			return false
		}
		sum += d
	}
	return sum == n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		n, _ := strconv.Atoi(line)
		input := fmt.Sprintf("%d\n", n)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%s\noutput:\n%s\n", idx, err, input, out)
			os.Exit(1)
		}
		if !verifyOutput(n, out) {
			fmt.Printf("test %d failed for n=%d\noutput:\n%s\n", idx, n, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
