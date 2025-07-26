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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func expected(bits string) string {
	for _, c := range bits {
		if c == '1' {
			return "1"
		}
	}
	return "0"
}

func checkCase(bin string, n int, bits string) error {
	input := fmt.Sprintf("%d\n%s\n", n, bits)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	out = strings.TrimSpace(out)
	exp := expected(bits)
	if out != exp {
		return fmt.Errorf("expected %s got %s", exp, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesB.txt: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Fprintf(os.Stderr, "bad testcase on line %d\n", idx+1)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad N on line %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		bits := fields[1]
		if len(bits) != n {
			fmt.Fprintf(os.Stderr, "bitstring length mismatch on line %d\n", idx+1)
			os.Exit(1)
		}
		idx++
		if err := checkCase(bin, n, bits); err != nil {
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
