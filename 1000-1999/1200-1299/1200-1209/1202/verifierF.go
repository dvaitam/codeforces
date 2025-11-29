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

// Embedded source for the reference solution (currently a stub) (was 1202F.go).
const solutionSource = `package main

// TODO: implement solution for Problem F.
func main() {}
`

const testcasesRaw = `2 3
8 3
8 6
8 3
8 5
8 6
4 1
6 5
7 2
5 3
3 3
2 6
4 1
1 8
8 2
1 3
8 3
5 2
5 6
3 8
6 1
6 2
2 2
5 5
6 3
1 8
7 7
4 2
7 5
8 7
8 7
3 7
4 7
4 5
2 3
8 6
6 5
8 6
4 3
8 3
6 4
8 3
2 8
6 3
2 6
3 5
6 2
2 7
6 6
5 6
1 4
5 3
3 3
8 4
2 5
2 1
7 2
7 8
8 3
3 1
2 3
2 7
2 5
2 2
4 5
7 3
8 5
6 3
7 5
3 3
2 4
6 8
2 6
7 5
8 1
3 6
3 8
4 4
2 3
8 7
2 6
4 2
6 3
4 4
6 8
6 8
5 5
1 3
5 4
2 5
8 3
2 1
8 7
7 6
6 4
6 8
3 1
8 3
1 2
1 3`

var _ = solutionSource

// Expected logic derived from current verifier (counts valid k partitions).
func expected(a, b int) int {
	n := a + b
	count := 0
	for k := 1; k <= n; k++ {
		m := n / k
		if m == 0 {
			continue
		}
		r := n % k
		for x := 0; x <= r; x++ {
			rem := a - x*(m+1)
			if rem < 0 {
				continue
			}
			if rem%m != 0 {
				continue
			}
			y := rem / m
			if y >= 0 && y <= k-r {
				count++
				break
			}
		}
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("test %d invalid line\n", idx)
			os.Exit(1)
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		input := fmt.Sprintf("%d %d\n", a, b)
		expectedOut := fmt.Sprintf("%d", expected(a, b))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expectedOut, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
