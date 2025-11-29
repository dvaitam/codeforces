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

const testcasesRaw = `13 98 54 0 1 2 1 1 1 1 1 2 0 2 0 1
5 97 13 2 1 2 2 2
5 40 13 2 0 2 1 1
18 13 46 1 1 2 2 0 2 1 1 2 1 0 2 0 0 2 1 2 2
1 79 64 1
8 94 42 2 0 0 2 0 0 0 2
15 12 11 1 2 1 0 1 2 1 2 0 2 1 2 0 2 2
19 37 57 0 2 1 1 2 0 1 0 0 0 0 2 2 1 1 0 0 2 0
5 5 11 2 2 2 1 2
17 36 67 0 0 2 2 1 2 1 1 1 2 2 2 1 0 1 2 0
16 76 81 1 0 0 0 2 1 0 2 0 1 0 1 1 0 0 0
8 6 74 2 2 2 2 0 0 0 2
7 78 74 0 1 0 1 0 0 2
1 25 24 2
4 62 27 2 0 2 0
18 55 80 0 1 0 0 0 2 1 1 1 0 0 2 1 0 2 0 2 1
7 34 46 2 1 2 0 2 2 0
2 87 21 0 1
17 33 16 2 1 2 0 0 1 2 1 2 2 1 2 1 1 2 1 0
18 89 2 1 2 0 1 2 0 2 1 0 0 1 1 2 1 2 1 2 2
20 17 92 1 1 2 1 2 0 0 2 0 2 1 0 0 0 2 1 1 2 2 2
14 5 52 2 2 1 2 2 0 0 1 0 1 2 0 1 2
16 72 78 0 0 1 1 1 1 0 1 0 2 2 0 2 0 0 1
14 41 1 0 0 2 0 2 2 2 0 0 0 2 2 0 1
9 89 24 0 1 1 2 0 0 1 1 0
9 18 84 2 2 2 1 0 0 1 0 0
2 27 88 1 2
11 47 73 0 2 2 2 2 1 2 2 1 2 1
12 69 23 0 1 2 1 0 0 0 1 1 1 1 2
3 44 100 2 0 0
9 21 20 2 1 1 1 2 0 1 0 1
8 7 40 0 2 2 0 1 1 1 1
14 14 13 2 1 1 1 1 0 1 0 2 1 1 0 1 1
5 22 81 2 1 2 0 0
3 26 96 0 0 1
1 13 51 2
17 38 58 1 2 2 2 0 1 0 1 0 1 2 0 1 0 1 0 0
1 68 58 2
7 16 64 1 1 0 2 0 0 2
5 14 26 1 1 1 2 0
4 77 63 0 2 1 2
14 67 64 2 1 1 1 2 2 0 2 2 0 0 1 2 2
11 42 5 2 0 1 2 0 1 2 1 2 2 1
3 11 67 0 0 0
5 6 39 0 1 1 0 0
15 48 65 1 2 2 0 2 0 2 2 2 0 2 1 0 1 2
20 54 62 1 2 2 0 0 2 0 2 0 1 2 2 1 1 0 1 1 1 1 1
13 8 21 2 0 0 1 2 1 0 0 1 1 0 1 2
3 87 90 0 1 1
2 79 60 1 1
2 13 61 0 0
2 77 80 0 2
11 14 90 2 2 1 0 1 1 0 0 2 2 1
20 81 44 2 0 2 2 2 1 0 1 1 2 2 0 2 0 0 1 1 1 0 1
12 81 10 0 0 1 1 0 2 2 2 2 0 0 0
17 90 68 1 2 1 0 0 1 2 1 0 0 1 0 0 1 0 2 0
15 56 88 1 0 1 1 2 1 0 1 0 0 1 2 0 1 1
6 2 30 1 0 2 0 0 0
7 85 87 2 0 2 0 1 1 2
1 78 30 0
6 59 15 1 1 2 1 0 0
7 47 43 1 1 1 2 2 1 0
19 11 14 2 2 1 0 1 0 0 0 1 2 0 0 0 1 1 1 2 0 0
20 3 51 0 2 0 0 1 1 2 1 2 0 2 1 1 2 1 0 0 1 2 1
17 8 49 1 0 1 2 1 1 0 1 1 1 0 0 0 1 0 2 2
5 90 58 1 0 1 1 0
8 59 44 2 0 1 1 2 2 0 1
7 38 1 2 1 2 1 0 0 1
4 99 81 1 2 2 0
14 91 97 1 0 2 1 0 2 1 1 2 0 0 1 2 0
1 33 51 2
19 91 51 1 0 2 1 1 1 2 0 2 0 0 0 1 1 2 1 0 2 0
6 9 54 1 1 2 0 2 2
7 69 14 1 2 2 1 2 1 1
15 48 73 2 0 0 0 2 0 1 1 2 1 0 2 2 1 1
16 96 54 0 1 1 2 2 1 1 2 0 1 1 0 2 1 0 2
7 4 46 1 1 0 2 0 2 0
13 1 47 0 0 2 0 1 2 2 1 2 0 0 2 1
7 14 56 1 2 1 1 0 1 1
14 19 58 2 0 2 1 0 0 0 1 1 1 1 1 1 2
8 26 57 0 2 2 0 1 0 0 2
3 24 47 0 2 2
6 30 79 1 2 0 2 2 1
12 53 59 0 2 2 2 2 2 2 2 1 2 1 1
9 91 61 0 1 1 0 0 0 0 1 0
10 84 1 0 0 1 2 0 2 1 2 0 1
7 44 78 0 2 0 1 1 2 1
11 33 4 2 0 0 1 0 0 2 1 0 0 1
10 40 67 1 1 1 1 2 0 0 1 2 0
1 59 64 2
15 7 53 1 1 1 0 0 0 0 0 0 1 0 1 2 0 1
18 97 51 0 0 0 1 0 0 1 1 1 1 0 2 1 2 2 0 2 1
15 66 78 1 2 2 1 1 0 0 0 2 2 0 0 2 1 0
7 37 95 2 0 2 2 2 1 0
4 50 83 1 0 2 2
12 30 87 2 2 2 2 1 0 2 0 0 2 1 2
11 30 48 2 1 1 2 0 0 0 2 2 1 1
19 82 4 0 1 0 0 2 0 0 0 1 2 2 0 0 2 0 0 1 1 2
19 17 81 1 0 2 0 2 2 1 1 1 1 0 0 2 2 0 2 1 2 1
18 41 91 0 1 0 2 1 2 0 1 1 1 0 0 0 1 2 2 0 1`

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Embedded reference logic from 1040A.go.
func referenceCost(n, d0, d1 int, colors []int) int64 {
	D := []int{d0, d1}
	var res int64
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		x := colors[i]
		y := colors[j]
		if x != 2 && y != 2 {
			if x != y {
				return -1
			}
			continue
		}
		if x == 2 && y == 2 {
			res += int64(2 * min(d0, d1))
			continue
		}
		if x == 2 {
			res += int64(D[y])
			continue
		}
		// y == 2
		res += int64(D[x])
	}
	if n%2 == 1 && colors[n/2] == 2 {
		res += int64(min(d0, d1))
	}
	return res
}

func parseCase(line string) (int, int, int, []int, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return 0, 0, 0, nil, fmt.Errorf("invalid case")
	}
	n, _ := strconv.Atoi(parts[0])
	a, _ := strconv.Atoi(parts[1])
	b, _ := strconv.Atoi(parts[2])
	if len(parts) != 3+n {
		return 0, 0, 0, nil, fmt.Errorf("expected %d colors, got %d", n, len(parts)-3)
	}
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		colors[i], _ = strconv.Atoi(parts[3+i])
	}
	return n, a, b, colors, nil
}

func runCase(bin string, input string, exp int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse output %q", gotStr)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, a, b, colors, err := parseCase(line)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", n, a, b)
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(colors[i]))
		}
		input.WriteByte('\n')
		expect := referenceCost(n, a, b, colors)
		if err := runCase(bin, input.String(), expect); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
