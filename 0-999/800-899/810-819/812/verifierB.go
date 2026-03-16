package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "812B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesRaw = `5 2 0000 0000 0000 0000 0010
5 6 00000110 00010000 01010010 00100100 01101000
1 3 00010
5 4 001110 001010 011010 000010 011110
5 5 0100100 0000000 0011110 0110100 0101110
4 6 00010110 00011110 01010110 01000110
4 4 010010 010100 000010 001000
5 5 0100000 0110010 0110010 0101010 0010000
4 4 001110 000110 011110 010100
4 5 0000010 0011110 0111010 0110010
4 6 01011110 00000000 01011000 00001010
3 1 010 010 010
1 2 0010
5 6 01110110 01111100 01110100 00011000 01011010
3 6 00001100 00101110 01000100
5 1 000 010 010 010 010
1 5 0101010
4 6 00000100 01000100 01010110 00011100
2 4 010010 001000
4 1 000 000 010 010
3 4 010010 011000 000010
4 2 0010 0000 0000 0010
4 1 010 010 000 010
4 6 01111000 01011010 01000010 01110000
2 1 010 010
5 4 000110 010010 010000 010100 011100
3 6 00100010 00110010 00010000
2 4 011010 000000
3 3 00100 01000 00010
1 6 01111100
1 4 010010
1 5 0110010
5 1 000 000 010 000 010
4 6 00101110 01101100 01111110 00001100
4 3 01000 01100 01010 01100
4 6 01010100 00111110 01100110 01111100
2 3 00010 00000
4 6 01001110 01011010 00001010 01000110
3 5 0101000 0011100 0001010
1 3 00010
5 5 0001110 0101000 0001000 0100010 0111010
3 1 010 000 000
3 6 01100100 01111010 00001110
3 4 011010 010000 000110
1 2 0100
5 1 000 000 010 010 010
5 5 0100000 0111000 0000110 0101010 0001000
2 1 010 000
4 5 0000110 0010110 0011110 0100000
4 3 00010 00110 01100 01100
4 1 010 010 010 010
4 4 010000 010100 011110 011010
5 3 00010 01010 01110 01000 00110
1 4 011110
1 3 01110
2 4 001010 001110
4 6 01101010 00011010 01011110 00110110
2 5 0011110 0011110
4 4 011000 011000 000000 001000
3 2 0010 0010 0010
1 3 01000
5 4 010000 011110 000100 010010 000110
2 6 00101110 00110100
4 3 00010 00010 01000 01110
3 5 0101100 0110110 0110110
5 6 01101010 00100110 00001110 00101110 00111110
3 1 010 000 010
2 5 0101100 0111100
1 6 01110010
4 2 0100 0100 0110 0100
4 4 010000 011000 001000 001010
3 2 0100 0110 0000
2 6 00100110 01000110
2 4 010000 010000
3 1 010 010 010
1 2 0000
2 2 0100 0100
3 1 000 010 010
2 4 001000 011010
4 4 010110 000100 000010 011000
2 3 01000 01010
1 2 0100
3 1 000 000 010
2 6 00010100 00110000
5 4 011110 010010 011110 011000 001000
2 1 000 000
5 6 00100010 00011110 01000010 01101110 00110100
1 1 000
5 5 0101100 0100000 0011000 0100100 0001100
5 6 00001100 01110000 01101100 01110110 01111110
2 3 01010 00010
3 5 0110010 0110110 0111100
2 2 0110 0000
1 1 000
1 5 0110100
3 4 010010 011010 011110
4 1 000 000 000 010
1 6 00110010
2 6 01110110 01111000
1 1 000`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n {
			fmt.Printf("test %d invalid floor count\n", idx)
			os.Exit(1)
		}
		floors := parts[2:]
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			sb.WriteString(floors[i])
			sb.WriteByte('\n')
		}
		input := sb.String()
		expected, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
