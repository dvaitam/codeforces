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

func isTriangle(x, y, z int64) bool {
	return x+y > z && x+z > y && y+z > x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
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
		parts := strings.Fields(line)
		if len(parts) != 4 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		a, _ := strconv.ParseInt(parts[0], 10, 64)
		b, _ := strconv.ParseInt(parts[1], 10, 64)
		c, _ := strconv.ParseInt(parts[2], 10, 64)
		d, _ := strconv.ParseInt(parts[3], 10, 64)

		input := fmt.Sprintf("1\n%d %d %d %d\n", a, b, c, d)
		outStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(outStr))
		if len(fields) != 3 {
			fmt.Printf("test %d: expected three integers got %q\n", idx, outStr)
			os.Exit(1)
		}
		x, err1 := strconv.ParseInt(fields[0], 10, 64)
		y, err2 := strconv.ParseInt(fields[1], 10, 64)
		z, err3 := strconv.ParseInt(fields[2], 10, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Printf("test %d: invalid integers in output %q\n", idx, outStr)
			os.Exit(1)
		}
		if x < a || x > b {
			fmt.Printf("test %d failed: x=%d not in [%d,%d]\n", idx, x, a, b)
			os.Exit(1)
		}
		if y < b || y > c {
			fmt.Printf("test %d failed: y=%d not in [%d,%d]\n", idx, y, b, c)
			os.Exit(1)
		}
		if z < c || z > d {
			fmt.Printf("test %d failed: z=%d not in [%d,%d]\n", idx, z, c, d)
			os.Exit(1)
		}
		if !isTriangle(x, y, z) {
			fmt.Printf("test %d failed: %d %d %d do not form a triangle\n", idx, x, y, z)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
