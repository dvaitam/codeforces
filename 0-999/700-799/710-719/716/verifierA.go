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

const testcasesARaw = `7 7 1 6 15 23 30 35 43
6 10 4 13 16 21 24 26
10 5 9 19 22 27 29 31 37 45 54 56
6 7 6 16 20 29 37 45
9 5 1 10 11 13 20 21 31 39 45
4 6 2 6 16 20
4 3 9 17 19 21
6 9 8 10 15 24 29 31
9 6 9 13 23 32 42 47 55 57 67
7 6 10 14 19 22 26 29 30
10 5 8 10 12 15 18 19 21 30 37 46
5 9 4 8 18 25 35
5 8 8 14 16 22 32
2 8 10 16
4 4 1 6 8 12
6 3 6 13 14 16 19 23
1 10 9
10 2 1 3 7 17 27 29 36 38 44 46
1 10 1
4 3 2 10 14 15
1 9 7
10 2 5 7 11 13 18 24 31 34 35 44
8 1 10 12 19 23 28 34 42 52
3 4 1 4 7
6 9 5 7 17 25 28 29
8 7 10 19 24 30 37 42 45 54
1 8 2
6 1 9 14 17 21 29 35
10 5 6 16 26 29 34 41 48 50 51 61
4 6 3 7 11 19
7 10 7 8 15 25 32 33 36
8 2 5 8 16 25 33 42 52 53
1 8 6
5 8 1 8 12 21 23
3 1 7 14 20
1 4 1
1 9 10
2 4 2 12
4 5 5 8 10 18
7 2 1 6 14 16 21 24 33
6 2 3 8 9 10 11 15
5 9 6 12 22 23 33
8 8 7 13 22 25 29 36 46 51
1 3 3
5 6 6 12 14 20 30
1 1 5
3 3 10 15 21
7 9 3 8 10 18 22 23 28
3 9 2 7 14
6 5 7 9 11 20 28 36
6 6 2 10 12 20 27 28
5 6 3 6 16 23 25
2 2 4 8
1 7 1
2 7 9 18
5 8 8 18 22 29 31
6 4 5 15 18 25 29 35
2 2 1 10
8 4 2 10 17 22 26 27 31 41
3 2 4 12 19
6 9 3 5 15 23 26 36
7 7 9 17 23 31 39 43 52
10 4 1 7 13 19 20 29 32 37 47 50
7 10 5 13 15 17 26 27 29
4 3 1 6 7 15
6 3 3 11 17 26 33 42
9 1 10 12 21 31 33 40 44 49 58
10 7 8 15 25 35 39 40 41 44 49 58
10 5 6 8 16 21 26 33 40 47 48 51
3 4 5 11 12
1 8 7
3 8 10 12 15
6 7 1 11 19 26 34 35
2 8 3 4
1 10 10
3 6 2 11 17
4 7 8 10 11 21
8 10 6 8 18 23 26 33 38 40
9 4 1 8 16 22 26 34 40 42 43
1 8 5
1 9 10
10 4 4 6 15 24 31 40 45 47 50 57
10 7 2 4 11 13 15 22 25 26 34 41
7 1 8 14 19 21 27 29 31
6 1 6 12 15 16 20 26
2 10 3 7
1 4 2
1 5 6
1 10 4
3 3 8 10 18
6 5 3 4 8 14 20 28
5 5 9 15 18 28 30
2 9 10 15
3 7 3 6 10
6 9 4 8 11 16 22 29
1 3 10
1 7 2
2 3 7 12
9 7 3 13 20 25 31 33 37 45 51
9 1 7 14 15 22 28 36 40 46 51
`

func expectedA(times []int64, c int64) int {
	count := 0
	var last int64
	for i, t := range times {
		if i == 0 || t-last <= c {
			count++
		} else {
			count = 1
		}
		last = t
	}
	return count
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "test %d: invalid format\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		c, _ := strconv.ParseInt(fields[1], 10, 64)
		if len(fields) != 2+n {
			fmt.Fprintf(os.Stderr, "test %d: expected %d timestamps got %d\n", idx, n, len(fields)-2)
			os.Exit(1)
		}
		times := make([]int64, n)
		for i := 0; i < n; i++ {
			val, _ := strconv.ParseInt(fields[2+i], 10, 64)
			times[i] = val
		}
		input := fmt.Sprintf("%d %d\n", n, c)
		for i, t := range times {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", t)
		}
		input += "\n"
		expected := expectedA(times, c)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to parse output %q\n", idx, out)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
