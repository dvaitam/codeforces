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

const testcasesARaw = `0 0 1
0 0 2
0 0 3
0 0 4
0 0 5
0 1 0
0 1 1
0 1 2
0 1 3
0 1 4
0 1 5
0 2 0
0 2 1
0 2 2
0 2 3
0 2 4
0 2 5
0 3 0
0 3 1
0 3 2
0 3 3
0 3 4
0 3 5
0 4 0
0 4 1
0 4 2
0 4 3
0 4 4
0 4 5
0 5 0
0 5 1
0 5 2
0 5 3
0 5 4
0 5 5
1 0 0
1 0 1
1 0 2
1 0 3
1 0 4
1 0 5
1 1 0
1 1 1
1 1 2
1 1 3
1 1 4
1 1 5
1 2 0
1 2 1
1 2 2
1 2 3
1 2 4
1 2 5
1 3 0
1 3 1
1 3 2
1 3 3
1 3 4
1 3 5
1 4 0
1 4 1
1 4 2
1 4 3
1 4 4
1 4 5
1 5 0
1 5 1
1 5 2
1 5 3
1 5 4
1 5 5
2 0 0
2 0 1
2 0 2
2 0 3
2 0 4
2 0 5
2 1 0
2 1 1
2 1 2
2 1 3
2 1 4
2 1 5
2 2 0
2 2 1
2 2 2
2 2 3
2 2 4
2 2 5
2 3 0
2 3 1
2 3 2
2 3 3
2 3 4
2 3 5
2 4 0
2 4 1
2 4 2
2 4 3
2 4 4
2 4 5
2 5 0
2 5 1
2 5 2
2 5 3
2 5 4
2 5 5
3 0 0
3 0 1
3 0 2
3 0 3
3 0 4
3 0 5
3 1 0
3 1 1
3 1 2
3 1 3
3 1 4
3 1 5
3 2 0
3 2 1
3 2 2
3 2 3
3 2 4
3 2 5
3 3 0
3 3 1
3 3 2
3 3 3
3 3 4
3 3 5
3 4 0
3 4 1
3 4 2
3 4 3
3 4 4
3 4 5
3 5 0
3 5 1
3 5 2
3 5 3
3 5 4
3 5 5
4 0 0
4 0 1
4 0 2
4 0 3
4 0 4
4 0 5
4 1 0
`

func expected(r, g, b int) int {
	ans := 0
	if r > 0 {
		cars := (r + 1) / 2
		t := 0 + (cars-1)*3 + 30
		if t > ans {
			ans = t
		}
	}
	if g > 0 {
		cars := (g + 1) / 2
		t := 1 + (cars-1)*3 + 30
		if t > ans {
			ans = t
		}
	}
	if b > 0 {
		cars := (b + 1) / 2
		t := 2 + (cars-1)*3 + 30
		if t > ans {
			ans = t
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		r, _ := strconv.Atoi(parts[0])
		g, _ := strconv.Atoi(parts[1])
		b, _ := strconv.Atoi(parts[2])
		expect := expected(r, g, b)
		input := fmt.Sprintf("%d %d %d\n", r, g, b)
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
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
