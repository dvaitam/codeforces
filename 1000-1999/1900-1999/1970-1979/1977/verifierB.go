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
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func check(x int, out string) error {
	scan := bufio.NewScanner(strings.NewReader(out))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return fmt.Errorf("missing n")
	}
	n, err := strconv.Atoi(scan.Text())
	if err != nil || n < 1 || n > 32 {
		return fmt.Errorf("invalid n")
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if !scan.Scan() {
			return fmt.Errorf("missing array element")
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil || (v != -1 && v != 0 && v != 1) {
			return fmt.Errorf("invalid element")
		}
		arr[i] = v
		if i > 0 && arr[i] != 0 && arr[i-1] != 0 {
			return fmt.Errorf("consecutive non-zero elements")
		}
	}
	if scan.Scan() {
		return fmt.Errorf("extra output")
	}
	sum := int64(0)
	pow := int64(1)
	for i := 0; i < n; i++ {
		sum += int64(arr[i]) * pow
		pow <<= 1
	}
	if sum != int64(x) {
		return fmt.Errorf("representation does not match")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesB.txt")
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
		x, _ := strconv.Atoi(line)
		input := fmt.Sprintf("1\n%d\n", x)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if err := check(x, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
