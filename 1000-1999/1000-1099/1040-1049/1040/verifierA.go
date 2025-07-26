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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(n, a, b int, colors []int) int {
	D := []int{a, b}
	res := 0
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
			res += 2 * min(a, b)
			continue
		}
		if x == 2 {
			res += D[y]
			continue
		}
		// y == 2
		res += D[x]
	}
	if n%2 == 1 && colors[n/2] == 2 {
		res += min(a, b)
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

func runCase(bin string, input string, exp int) error {
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
	got, err := strconv.Atoi(gotStr)
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

	file, err := os.Open("testcasesA.txt")
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
		expect := expected(n, a, b, colors)
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
