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

// solution logic from 1927B.go
func solveCase(nums []int) string {
	counts := make(map[rune]int)
	letters := make([]rune, 0, len(nums))
	var sb strings.Builder
	next := 'a'
	for _, a := range nums {
		if a == 0 {
			counts[next]++
			letters = append(letters, next)
			sb.WriteRune(next)
			next++
		} else {
			for _, r := range letters {
				if counts[r] == a {
					sb.WriteRune(r)
					counts[r]++
					break
				}
			}
		}
	}
	return sb.String()
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

const testcasesData = `
5 0 0 0 0 0
4 0 0 0 1
13 0 0 0 0 0 0 0 0 0 1 1 0 0
9 0 0 0 0 0 0 0 1 2
18 0 0 0 0 0 0 1 0 0 0 0 0 0 1 0 2 1 3
15 0 0 0 0 0 0 0 1 0 1 1 1 0 2 0
17 0 0 0 0 0 1 0 0 1 0 1 0 1 0 0 0 1
14 0 0 0 0 0 0 1 0 1 0 0 2 0 0
6 0 0 0 0 0 0
16 0 0 0 0 0 1 0 0 0 1 0 0 0 0 0 0
18 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 1 1 1
5 0 0 0 0 0
2 0 0
19 0 0 0 0 0 0 1 1 0 1 2 0 0 1 0 0 2 1 1
8 0 0 0 0 1 0 0 1
9 0 0 0 1 0 0 1 0 1
9 0 0 0 0 0 0 0 0 0
6 0 0 0 1 0 1
10 0 0 0 0 1 0 0 0 0 1
14 0 0 0 0 1 0 0 1 0 0 0 0 1 0
5 0 0 0 0 0
17 0 0 0 0 0 0 0 0 0 1 1 1 0 0 1 0 1
11 0 0 0 0 0 0 0 0 1 1 0
3 0 1 0
6 0 0 0 0 0 0
2 0 0
19 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 1 2 0 1
13 0 0 0 0 0 0 0 1 1 0 0 1 0
19 0 0 0 0 0 0 0 0 0 0 0 1 1 0 0 2 0 0 0
2 0 0
6 0 0 0 0 0 0
20 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 1 1 2 0 0
3 0 0 0
5 0 0 0 0 0
13 0 0 0 0 1 1 0 0 0 2 1 0 0
9 0 0 0 0 0 0 0 1 0
14 0 0 0 0 0 1 0 0 0 1 0 1 0 1
6 0 0 0 0 0 0
10 0 0 0 0 0 0 0 0 1 0
1 0
10 0 0 0 0 0 1 1 0 1 2
20 0 0 0 0 0 0 0 0 0 0 0 0 1 0 1 1 0 2 0 1
3 0 0 0
15 0 0 0 1 0 0 0 0 0 1 0 2 0 1 1
8 0 0 0 0 0 0 1 0
8 0 0 0 1 0 0 0 0
3 0 0 0
1 0
12 0 1 0 0 0 0 0 0 0 1 0 0
6 0 0 1 0 0 0
17 0 0 0 0 1 0 0 0 0 0 1 0 0 1 1 0 1
6 0 0 0 0 0 0
8 0 0 0 0 0 0 0 0
9 0 0 1 1 0 0 0 0 0
16 0 0 0 0 0 1 0 0 0 1 0 2 1 2 0 1
13 0 0 0 0 0 0 0 0 1 0 0 1 0
15 0 0 0 0 1 0 0 1 1 2 0 0 1 0 0
9 0 0 0 0 0 0 1 0 0
17 0 0 0 0 1 0 2 0 0 0 1 2 0 0 0 0 0
20 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 1 0 0 0 2
17 0 0 0 0 0 0 0 0 0 0 1 1 1 0 0 0 0
3 0 0 0
14 0 0 0 0 0 0 0 0 0 0 1 0 1 1
10 0 0 0 0 0 0 0 0 0 0
1 0
20 0 0 0 0 0 0 0 1 1 1 0 0 2 0 0 0 1 1 1 0
16 0 0 0 0 0 0 1 0 0 1 1 0 2 1 0 1
18 0 0 0 0 0 0 0 0 0 0 0 1 0 1 1 1 1 2
11 0 0 0 0 1 0 0 0 0 0 0
19 0 0 0 0 1 0 0 0 0 1 0 0 1 2 0 1 0 1 1
14 0 0 0 1 0 0 0 0 0 0 0 1 0 1
13 0 0 0 0 0 1 0 0 1 0 0 2 0
10 0 0 0 0 1 0 2 0 0 0
9 0 0 0 1 0 0 0 0 0
14 0 0 0 0 0 0 0 0 1 0 0 1 1 0
19 0 0 0 1 0 0 0 0 0 0 1 1 2 0 0 2 0 0 0
7 0 0 0 0 0 0 1
17 0 0 0 0 0 0 0 0 0 0 1 0 1 0 1 0 2
8 0 0 1 0 0 0 0 0
11 0 0 0 0 1 1 0 0 1 0 0
14 0 0 0 0 0 0 0 0 0 0 0 1 1 0
20 0 0 0 0 0 0 1 0 1 0 0 1 0 0 0 1 0 2 0 1
3 0 0 0
11 0 0 0 0 1 0 0 0 0 0 0
2 0 0
17 0 0 0 0 0 0 0 1 0 0 1 0 0 1 0 1 1
6 0 1 0 0 0 0
13 0 0 0 0 0 0 0 0 1 0 0 0 0
1 0
11 0 0 0 0 1 0 0 0 0 0 0
1 0
5 0 0 0 0 0
9 0 0 0 0 0 0 0 0 1
17 0 0 0 0 0 0 0 0 0 0 0 1 0 1 0 1 1
15 0 0 0 0 1 1 2 0 1 0 0 0 1 0 0
14 0 0 0 0 0 0 0 1 0 2 0 1 3 1
3 0 0 0
20 0 0 0 0 0 0 0 0 0 0 1 2 0 1 0 0 1 0 1 0
2 0 0
2 0 0
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != n+1 {
			fmt.Printf("test %d expected %d numbers got %d\n", idx, n+1, len(fields))
			os.Exit(1)
		}
		nums := make([]int, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fields[i+1])
			nums[i], _ = strconv.Atoi(fields[i+1])
		}
		sb.WriteByte('\n')
		input := sb.String()

		expect := solveCase(nums)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
