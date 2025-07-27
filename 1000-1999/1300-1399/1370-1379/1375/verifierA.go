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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		panic(err)
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
		n, err := strconv.Atoi(parts[0])
		if err != nil || len(parts) != n+1 {
			fmt.Printf("invalid test case %d\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[i+1])
			arr[i] = v
		}
		// build input
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(outBuf.String()))
		if len(outFields) != n {
			fmt.Printf("Test %d: expected %d numbers got %d\n", idx, n, len(outFields))
			os.Exit(1)
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(outFields[i])
			if err != nil {
				fmt.Printf("Test %d: invalid output %q\n", idx, outFields[i])
				os.Exit(1)
			}
			b[i] = val
			if b[i] != arr[i] && b[i] != -arr[i] {
				fmt.Printf("Test %d: b[%d] not +/- a[%d]\n", idx, i, i)
				os.Exit(1)
			}
		}
		nonNeg, nonPos := 0, 0
		for i := 0; i < n-1; i++ {
			diff := b[i+1] - b[i]
			if diff >= 0 {
				nonNeg++
			}
			if diff <= 0 {
				nonPos++
			}
		}
		need := (n - 1) / 2
		if nonNeg < need || nonPos < need {
			fmt.Printf("Test %d failed: condition not satisfied\n", idx)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
