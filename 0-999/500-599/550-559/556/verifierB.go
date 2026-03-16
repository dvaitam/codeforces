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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, a []int) string {
	for t := 0; t < n; t++ {
		ok := true
		for i := 0; i < n; i++ {
			val := a[i]
			if i%2 == 0 {
				val = (val + t) % n
			} else {
				val = ((val-t)%n + n) % n
			}
			if val != i {
				ok = false
				break
			}
		}
		if ok {
			return "Yes"
		}
	}
	return "No"
}

const testcasesBRaw = `3 2 0 1
2 1 1
8 6 3 1 7 0 6 6 0
8 4 3 1 5 0 0 0 0
7 5 1 3 5 0 4 1
8 7 3 5 3 3 7 4 0
7 6 4 5 0 1 5 5
5 0 2 4 3 4
4 2 2 3 3
10 0 7 3 6 6 2 5 8 5 1
8 1 2 6 5 7 0 7 0
5 4 4 4 3 1
3 2 0 0
4 1 3 2 2
8 4 0 6 2 3 6 0 7
6 4 4 1 4 3 3
6 3 2 0 4 4 4
10 5 7 9 0 3 2 8 9 2 1
9 4 0 1 1 0 7 0 4 3
5 0 4 1 2 2
2 0 0
5 4 1 2 2 3
6 3 3 0 0 2 3
6 3 1 2 0 2 5
9 3 6 0 3 0 6 2 0 2
8 6 3 7 3 0 6 5 6
1 0
3 0 0 1
2 0 1
5 1 3 4 2 1
1 0
10 3 9 7 2 9 8 0 6 3 5
2 0 1
10 3 7 1 6 4 8 7 0 5 9
7 2 0 1 1 6 2 6
10 2 5 6 3 4 1 6 8 5 8
8 3 1 0 1 2 2 2 3
5 2 4 4 2 2
6 2 0 2 1 4 5
8 2 1 5 0 6 1 6 2
3 1 0 2
10 6 1 9 8 3 9 1 4 5 4
10 8 1 7 4 1 0 4 0 9 0
2 1 0
1 0
4 3 1 0 3
3 2 0 0
2 1 1
9 4 8 4 7 5 1 3 5 0
1 0
5 4 2 3 3 2
7 0 0 2 4 3 0 2
4 3 2 2 1
9 3 4 3 3 5 1 4 1 7
2 1 0
7 2 0 2 1 2 6 6
10 4 3 5 1 8 9 9 9 1 3
4 0 1 3 0
5 4 0 0 0 0
5 2 3 3 1 0
9 5 1 8 2 2 2 2 5 4
2 1 0
4 1 0 2 1
3 1 1 2
3 0 2 2
4 2 0 3 3
9 4 8 7 8 7 0 6 5 2
5 3 0 3 4 0
1 0
10 2 9 2 2 4 4 6 9 6 2
10 1 3 7 0 2 8 5 8 7 3
4 2 3 3 1
7 2 4 4 5 5 2 5
4 0 0 2 1
9 3 4 4 4 8 5 2 7 1
2 1 0
3 1 1 0
10 0 7 6 5 6 8 2 8 0 8
2 1 0
5 0 1 4 0 3
4 3 3 3 1
6 3 1 4 3 1 0
7 4 4 3 0 5 2 2
4 3 0 1 3
10 0 0 9 3 4 3 2 4 2 8
4 2 2 2 3
3 2 1 1
7 6 0 6 1 4 3 1
5 0 0 0 4 0
9 4 2 1 8 5 4 6 8 5
9 5 0 1 7 7 5 4 8 6
6 5 5 4 3 0 5
7 3 1 4 0 2 5 4
9 3 7 8 6 4 2 7 8 3
6 4 0 5 3 4 3
7 2 6 4 4 5 5 5
2 1 0
5 0 3 1 3 2
3 0 2 0
6 2 5 3 5 4 2
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad test format on line %d\n", idx+1)
			os.Exit(1)
		}
		idx++
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Printf("bad n on line %d: %v\n", idx, err)
			os.Exit(1)
		}
		if len(fields)-1 != n {
			fmt.Printf("line %d: expected %d numbers, got %d\n", idx, n, len(fields)-1)
			os.Exit(1)
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[i+1])
			if err != nil {
				fmt.Printf("bad value on line %d: %v\n", idx, err)
				os.Exit(1)
			}
			a[i] = v
		}
		var sb strings.Builder
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		input := fmt.Sprintf("%d\n%s\n", n, sb.String())
		exp := expected(n, a)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != exp {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", idx, exp, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
