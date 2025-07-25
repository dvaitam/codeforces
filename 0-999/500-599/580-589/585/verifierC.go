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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expectedC(x, y int64) string {
	if gcd(x, y) != 1 {
		return "Impossible\n"
	}
	type step struct {
		k int64
		c byte
	}
	var res []step
	for x > 1 || y > 1 {
		if x > y {
			k := (x - 1) / y
			res = append(res, step{k, 'A'})
			x -= k * y
		} else {
			k := (y - 1) / x
			res = append(res, step{k, 'B'})
			y -= k * x
		}
	}
	var sb strings.Builder
	for _, st := range res {
		sb.WriteString(fmt.Sprintf("%d%c", st.k, st.c))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
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
		var xStr, yStr string
		fmt.Sscan(line, &xStr, &yStr)
		x, _ := strconv.ParseInt(xStr, 10, 64)
		y, _ := strconv.ParseInt(yStr, 10, 64)
		expect := expectedC(x, y)
		input := fmt.Sprintf("%s %s\n", xStr, yStr)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		expectTrim := strings.TrimSpace(expect)
		if got != expectTrim {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expectTrim, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
