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

// Embedded testcases (same format as original file).
const embeddedTestcases = `100
3 5 1 3 1 4 4
8 6 4 2 1 4 1 4 4 5 1 6 4 3 6 2 5
2 3 1 1 1
9 1 4 6 2 4 6 1 5 2 4 4 5 2 3 2 6 2 4
5 1 4 5 6 1 2 6 6 3 1
6 6 6 5 4 5 6 2 3 3 5 4 5
7 5 1 4 2 6 4 4 6 2 3 5 6 6 6
6 1 4 6 5 1 2 5 4 3 4 6 1
8 1 3 6 5 5 5 4 6 2 2 5 2 1 2 5 5
4 4 5 3 5 3 4 3 6
9 5 6 1 4 6 5 2 5 5 2 4 1 4 3 5 5 2 5
7 4 3 4 3 1 5 5 5 5 3 4 5 1 2
3 5 5 2 1 5 3
1 6 1
2 1 4 1 3
4 3 1 5 2 3 3 1 2
3 3 5 2 6 3 6
5 4 6 3 4 4 1 1 3 4 3
7 2 3 1 3 6 5 2 5 4 1 2 1 4 2
1 6 2
8 6 5 6 4 5 2 6 6 5 4 2 5 6 1 4 6
10 3 6 6 4 1 6 3 2 2 1 3 1 1 3 3 6 2 4 5 3
3 1 5 1 5 2 5
8 2 6 5 5 1 4 2 3 1 2 5 6 4 5 2 4
2 6 4 3 5
8 1 3 5 4 3 1 2 2 3 5 2 3 4 2 3 6
2 4 5 3 6
9 4 5 2 1 6 1 1 2 2 2 5 2 3 3 5 5 3 3
6 3 1 3 2 5 6 4 2 5 5 1 3
1 4 1
7 2 2 3 1 5 5 4 1 5 5 2 5 1 3
6 3 5 5 1 4 3 1 1 3 1 5 6
1 1 4
2 1 2 2 5
7 2 1 4 2 6 2 2 6 1 4 4 5 3 5
5 6 4 3 1 2 6 3 1 1 1
5 6 5 3 4 4 3 4 1 1 3
10 4 1 3 2 5 5 6 4 6 3 3 2 5 2 3 2 2 3 1 3
2 4 1 6 5
6 2 4 3 1 3 2 3 5 3 2 3 1
9 5 5 5 1 2 2 1 2 4 1 3 5 1 6 1 1 6 1
5 3 4 4 2 1 5 3 1 5 6
3 2 2 2 3 3 1
9 5 3 2 2 2 5 6 1 3 5 6 5 6 6 2 2 3 4
9 2 1 6 6 2 3 1 6 4 4 5 3 5 4 5 4 1 4
6 2 3 4 1 6 4 5 1 1 6 3 5
3 5 2 2 3 3 4
10 4 2 5 1 2 4 1 2 5 3 5 6 4 6 6 6 2 2 3 4
8 2 6 4 3 5 5 6 6 3 6 2 1 1 5 6 3
3 5 2 3 3 6 3
9 3 2 6 6 6 4 5 1 1 5 5 5 4 2 2 3 4 2
10 6 1 4 6 4 6 6 3 4 5 2 5 6 1 5 1 3 6 1 3
2 2 5 6 6
2 4 2 4 4
7 2 3 4 2 5 4 2 1 4 5 5 4 1 6
5 3 2 4 6 5 1 2 5 4 5
1 1 6
10 2 3 2 2 3 2 5 2 3 3 5 3 6 4 2 5 3 4 4 1
4 5 4 2 3 1 1 1 5
1 5 3
3 1 5 3 5 3 4
9 6 3 5 3 1 1 4 6 4 3 3 5 4 3 6 6 5 4
2 6 4 4 2
9 1 3 6 5 6 6 6 5 2 4 5 5 4 6 6 3 6 2
8 5 6 5 2 3 5 1 6 4 5 4 4 3 5 5 6
2 4 6 2 6
5 6 1 4 6 6 2 6 4 3 2
2 5 1 3 3
7 6 5 3 2 4 3 4 2 4 5 1 3 5 1
10 4 1 3 1 6 4 1 2 5 6 2 6 1 4 6 6 3 5 3 2
9 2 2 3 3 1 1 6 5 6 3 4 5 5 6 1 2 3 6
9 3 3 5 6 2 4 5 4 2 4 3 5 3 6 2 3 5 6
4 6 1 5 4 3 4 2 3
4 1 6 6 2 5 4 5 6
3 5 3 4 5 2 2
3 6 4 3 3 4 2
2 6 2 6 6
5 1 1 2 4 3 4 1 2 1 1
10 1 2 6 1 4 6 5 6 5 4 3 6 3 1 5 6 2 1 2 4
4 4 4 4 2 2 2 3 4
9 5 4 2 4 6 3 3 4 5 1 2 1 1 1 1 4 3 4
10 3 2 4 2 6 2 1 1 4 2 6 5 1 5 4 3 2 1 4 6
5 1 1 5 1 5 2 1 3 1 4
2 2 1 4 6
3 6 3 6 2 6 4
7 3 6 3 3 6 6 2 2 1 5 5 2 3 4
10 6 5 6 5 1 3 5 4 5 2 6 5 4 6 1 6 3 6 5 6
2 3 2 1 2
1 2 4
1 1 6
2 5 4 5 3
2 3 1 2 5
1 4 6
3 4 6 4 1 6 5
5 1 3 3 1 3 1 4 1 6 3
6 6 2 3 4 1 6 3 1 4 2 5 5
4 3 3 5 4 5 4 1 2
8 5 5 6 5 6 5 5 1 3 6 2 2 3 4 5 3
2 4 3 2 5
2 1 3 6 5`

// Embedded reference solution (from 703A.go).
func expected(n int, pairs [][2]int) string {
	mishka, chris := 0, 0
	for _, p := range pairs {
		m, c := p[0], p[1]
		if m > c {
			mishka++
		} else if c > m {
			chris++
		}
	}
	if mishka > chris {
		return "Mishka"
	} else if chris > mishka {
		return "Chris"
	}
	return "Friendship is magic!^^"
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	scan := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		pairs := make([][2]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			m, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			c, _ := strconv.Atoi(scan.Text())
			pairs[j] = [2]int{m, c}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, p := range pairs {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		input := sb.String()
		exp := expected(n, pairs) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
