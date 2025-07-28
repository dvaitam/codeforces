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

func expected(n int, w, h int64, a, b []int64) string {
	lower := int64(-1 << 60)
	upper := int64(1 << 60)
	for i := 0; i < n; i++ {
		l := (b[i] + h) - (a[i] + w)
		r := (b[i] - h) - (a[i] - w)
		if l > lower {
			lower = l
		}
		if r < upper {
			upper = r
		}
	}
	if lower <= upper {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
	for {
		var line1, line2, line3 string
		if !scanner.Scan() {
			break
		}
		line1 = strings.TrimSpace(scanner.Text())
		if line1 == "" {
			continue
		}
		if !scanner.Scan() {
			break
		}
		line2 = strings.TrimSpace(scanner.Text())
		if !scanner.Scan() {
			break
		}
		line3 = strings.TrimSpace(scanner.Text())
		idx++
		fields := strings.Fields(line1)
		n, _ := strconv.Atoi(fields[0])
		w64, _ := strconv.ParseInt(fields[1], 10, 64)
		h64, _ := strconv.ParseInt(fields[2], 10, 64)
		aVals := make([]int64, n)
		bVals := make([]int64, n)
		aStr := strings.Fields(line2)
		bStr := strings.Fields(line3)
		for i := 0; i < n; i++ {
			av, _ := strconv.ParseInt(aStr[i], 10, 64)
			bv, _ := strconv.ParseInt(bStr[i], 10, 64)
			aVals[i] = av
			bVals[i] = bv
		}
		expect := expected(n, w64, h64, aVals, bVals)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(line1)
		input.WriteByte('\n')
		input.WriteString(line2)
		input.WriteByte('\n')
		input.WriteString(line3)
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
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
		if strings.ToUpper(got) != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
