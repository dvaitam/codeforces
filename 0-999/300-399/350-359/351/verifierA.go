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

func expected(n int, nums []string) string {
	k := 0
	sum := 0
	for _, s := range nums {
		parts := strings.SplitN(s, ".", 2)
		frac := 0
		if len(parts) == 2 {
			fracPart := parts[1]
			if len(fracPart) > 3 {
				fracPart = fracPart[:3]
			}
			fracPart = fmt.Sprintf("%03s", fracPart)
			fracPart = strings.ReplaceAll(fracPart, " ", "0")
			frac, _ = strconv.Atoi(fracPart)
		}
		if frac != 0 {
			k++
			sum += frac
		}
	}
	L := k - n
	if L < 0 {
		L = 0
	}
	R := k
	if R > n {
		R = n
	}
	c0 := (sum + 500) / 1000
	c := c0
	if c < L {
		c = L
	} else if c > R {
		c = R
	}
	diff := c*1000 - sum
	if diff < 0 {
		diff = -diff
	}
	return fmt.Sprintf("%d.%03d", diff/1000, diff%1000)
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "351A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

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
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+2*n {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		nums := parts[1:]
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		sb.WriteString(strings.Join(nums, " "))
		sb.WriteByte('\n')
		input := sb.String()

		// run oracle
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
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
