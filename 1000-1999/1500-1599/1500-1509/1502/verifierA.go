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

func expected(parts []string) (int, int, []int) {
	n, _ := strconv.Atoi(parts[0])
	arr := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		v, _ := strconv.Atoi(parts[i+1])
		arr[i] = v
		sum += v
	}
	return n, sum, arr
}
func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	idx := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		n, sum, _ := expected(parts)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(parts[i+1])
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err2 := strconv.Atoi(gotStr)
		if err2 != nil || got != sum {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %s\n", idx, sum, gotStr)
			os.Exit(1)
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
