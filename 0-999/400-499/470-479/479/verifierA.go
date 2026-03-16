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
const testcasesARaw = `
7 7 1
5 9 8
7 5 8
6 10 4
9 3 5
3 2 10
5 9 10
3 5 2
2 6 8
9 2 6
7 6 10
4 9 8
8 9 5
1 9 1
2 7 1
10 8 6
4 6 2
4 10 4
4 3 9
8 2 2
6 9 8
2 5 9
5 2 9
6 9 4
10 9 10
5 8 2
10 7 6
10 4 5
3 4 3
1 10 5
8 2 2
3 3 1
2 9 7
9 5 9
4 4 10
7 10 5
8 8 6
2 6 10
2 8 10
6 4 4
1 5 2
4 6 3
6 7 1
2 3 4
1 10 9
10 2 1
2 4 10
10 2 7
2 6 2
1 10 1
4 3 2
8 4 1
1 9 7
10 2 5
2 4 2
5 6 7
3 1 9
8 1 10
2 7 4
5 6 8
10 3 4
1 3 3
6 9 5
2 10 8
3 1 8
7 10 9
5 6 7
5 3 9
1 8 2
6 1 9
5 3 4
8 6 10
5 6 10
10 3 5
7 7 2
1 10 4
6 3 4
4 8 7
10 7 1
7 10 7
1 3 8
2 5 3
8 9 8
9 10 1
1 8 6
5 8 1
7 4 9
2 3 1
7 7 6
1 4 1
1 9 10
2 4 2
10 4 5
5 3 2
8 7 2
1 5 8
2 5 3
9 6 2
3 5 1
1 1 4
`


func expected(a, b, c int) int {
	vals := []int{
		a + b + c,
		a * b * c,
		(a + b) * c,
		a * (b + c),
		a*b + c,
		a + b*c,
	}
	res := vals[0]
	for _, v := range vals[1:] {
		if v > res {
			res = v
		}
	}
	return res
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
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		c, _ := strconv.Atoi(parts[2])
		expect := strconv.Itoa(expected(a, b, c))
		input := fmt.Sprintf("%d %d %d\n", a, b, c)
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
