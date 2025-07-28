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

func merge(a, b []int) []int {
	i, j := 0, 0
	res := make([]int, 0, len(a)+len(b))
	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}
	res = append(res, a[i:]...)
	res = append(res, b[j:]...)
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesH.txt")
	if err != nil {
		panic(err)
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
		fields := strings.Fields(line)
		pos := 0
		n, _ := strconv.Atoi(fields[pos])
		pos++
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i], _ = strconv.Atoi(fields[pos])
			pos++
		}
		m, _ := strconv.Atoi(fields[pos])
		pos++
		b := make([]int, m)
		for i := 0; i < m; i++ {
			b[i], _ = strconv.Atoi(fields[pos])
			pos++
		}
		expected := merge(a, b)

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		input.WriteString(fmt.Sprintf("%d\n", m))
		for i, v := range b {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(out.String()))
		if len(outFields) != len(expected) {
			fmt.Printf("Test %d failed: expected %d numbers got %d\n", idx, len(expected), len(outFields))
			os.Exit(1)
		}
		for i, f := range outFields {
			val, _ := strconv.Atoi(f)
			if val != expected[i] {
				fmt.Printf("Test %d failed: expected %v got %v\n", idx, expected, outFields)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
