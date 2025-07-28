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

var memo [105][105][2][2]int8

func dfs(e, o, turn, parity int) bool {
	if e == 0 && o == 0 {
		return parity == 0
	}
	m := &memo[e][o][turn][parity]
	if *m != 0 {
		return *m == 1
	}
	var res bool
	if turn == 0 {
		if e > 0 && dfs(e-1, o, 1, parity) {
			res = true
		} else if o > 0 && dfs(e, o-1, 1, parity^1) {
			res = true
		}
	} else {
		res = true
		if e > 0 && !dfs(e-1, o, 0, parity) {
			res = false
		}
		if res && o > 0 && !dfs(e, o-1, 0, parity) {
			res = false
		}
	}
	if res {
		*m = 1
	} else {
		*m = 2
	}
	return res
}

func solveC(arr []int) string {
	e, o := 0, 0
	for _, x := range arr {
		if x%2 == 0 {
			e++
		} else {
			o++
		}
	}
	for i := 0; i <= e; i++ {
		for j := 0; j <= o; j++ {
			memo[i][j][0][0] = 0
			memo[i][j][0][1] = 0
			memo[i][j][1][0] = 0
			memo[i][j][1][1] = 0
		}
	}
	if dfs(e, o, 0, 0) {
		return "Alice"
	}
	return "Bob"
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
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != n+1 {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[1+i])
			arr[i] = v
		}
		expect := solveC(arr)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", arr[i]))
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
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
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
