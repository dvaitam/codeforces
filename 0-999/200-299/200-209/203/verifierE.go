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
        oracle := filepath.Join(dir, "oracleE")
        cmd := exec.Command("rustc", "203E.rs", "-O", "-o", oracle)
        if out, err := cmd.CombinedOutput(); err != nil {
                return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
        }
        return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesE.txt")
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
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Fprintf(os.Stderr, "test %d: invalid test case\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid n: %v\n", idx, err)
			os.Exit(1)
		}
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("%s %s %s\n", fields[0], fields[1], fields[2]))
		pos := 3
		for i := 0; i < n; i++ {
			if pos+2 >= len(fields) {
				fmt.Fprintf(os.Stderr, "test %d: not enough robot data\n", idx)
				os.Exit(1)
			}
			builder.WriteString(fmt.Sprintf("%s %s %s\n", fields[pos], fields[pos+1], fields[pos+2]))
			pos += 3
		}
		input := builder.String()
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(input)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
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
