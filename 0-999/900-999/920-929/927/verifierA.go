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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		idx++
		parts := strings.Fields(line)
		pos := 0
		if len(parts) < 3 {
			fmt.Printf("test %d: bad format\n", idx)
			os.Exit(1)
		}
		w, _ := strconv.Atoi(parts[pos])
		pos++
		h, _ := strconv.Atoi(parts[pos])
		pos++
		k, _ := strconv.Atoi(parts[pos])
		pos++
		if len(parts) < pos+2*k+1 {
			fmt.Printf("test %d: not enough data\n", idx)
			os.Exit(1)
		}
		cars := make([][2]int, k)
		for i := 0; i < k; i++ {
			x, _ := strconv.Atoi(parts[pos])
			pos++
			y, _ := strconv.Atoi(parts[pos])
			pos++
			cars[i] = [2]int{x, y}
		}
		if len(parts) <= pos {
			fmt.Printf("test %d: missing q\n", idx)
			os.Exit(1)
		}
		q, _ := strconv.Atoi(parts[pos])
		pos++
		if len(parts) < pos+5*q {
			fmt.Printf("test %d: not enough order data\n", idx)
			os.Exit(1)
		}
		orders := make([][5]int, q)
		for i := 0; i < q; i++ {
			t, _ := strconv.Atoi(parts[pos])
			pos++
			sx, _ := strconv.Atoi(parts[pos])
			pos++
			sy, _ := strconv.Atoi(parts[pos])
			pos++
			tx, _ := strconv.Atoi(parts[pos])
			pos++
			ty, _ := strconv.Atoi(parts[pos])
			pos++
			orders[i] = [5]int{t, sx, sy, tx, ty}
		}
		if pos != len(parts) {
			fmt.Printf("test %d: extra data\n", idx)
			os.Exit(1)
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", w, h)
		fmt.Fprintf(&input, "%d\n", k)
		for _, c := range cars {
			fmt.Fprintf(&input, "%d %d\n", c[0], c[1])
		}
		for _, o := range orders {
			fmt.Fprintf(&input, "%d %d %d %d %d\n", o[0], o[1], o[2], o[3], o[4])
		}
		fmt.Fprintf(&input, "-1 0 0 0 0\n")
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(stdout.String())
		expect := strings.TrimSpace(strings.Repeat("0\n", q+2))
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
