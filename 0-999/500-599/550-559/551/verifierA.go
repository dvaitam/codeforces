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

func expected(a []int) []int {
	freq := make([]int, 2002)
	for _, v := range a {
		if v >= 0 && v <= 2000 {
			freq[v]++
		}
	}
	suffix := make([]int, 2002)
	for r := 2000; r >= 1; r-- {
		suffix[r] = suffix[r+1] + freq[r+1]
	}
	res := make([]int, len(a))
	for i, v := range a {
		res[i] = suffix[v] + 1
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcases: %v\n", err)
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
		if len(parts) < 1 {
			fmt.Printf("case %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+n {
			fmt.Printf("case %d: wrong number of ratings\n", idx)
			os.Exit(1)
		}
		ratings := make([]int, n)
		for i := 0; i < n; i++ {
			ratings[i], _ = strconv.Atoi(parts[1+i])
		}
		expect := expected(ratings)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(ratings[i]))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		outParts := strings.Fields(strings.TrimSpace(out.String()))
		if len(outParts) != n {
			fmt.Printf("case %d: expected %d numbers got %d\n", idx, n, len(outParts))
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			got, err := strconv.Atoi(outParts[i])
			if err != nil || got != expect[i] {
				fmt.Printf("case %d failed: expected %v got %v\n", idx, expect, outParts)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
