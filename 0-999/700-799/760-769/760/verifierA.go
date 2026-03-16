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

const testcasesARaw = `100
7 7
7 1
5 5
8 4
5 4
6 5
4 5
3 3
3 7
2 5
5 5
12 7
10 2
5 1
12 1
11 3
8 5
2 3
7 3
10 6
4 5
8 4
9 3
1 7
9 1
2 6
7 6
11 6
1 5
8 7
6 2
12 3
12 7
2 2
10 2
4 7
3 7
9 4
2 1
6 5
8 1
5 5
5 6
2 5
6 7
9 2
10 5
10 3
8 1
10 7
7 3
10 2
5 2
4 7
3 1
10 6
5 4
2 1
11 7
3 2
1 7
2 6
9 6
7 7
12 5
5 5
4 7
4 6
10 7
7 5
5 4
8 6
11 6
6 1
6 5
2 4
10 6
6 7
4 2
1 6
5 1
12 2
6 7
3 3
7 7
1 1
3 7
12 2
1 7
10 6
9 5
11 1
1 1
11 2
10 7
10 1
7 1
6 7
2 1
10 1
`

func expected(m, d int) string {
	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	total := days[m-1] + d - 1
	columns := (total + 6) / 7
	return fmt.Sprintf("%d", columns)
}

func runCase(exe, input, exp string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
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
	scan := bufio.NewScanner(strings.NewReader(testcasesARaw))
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
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		d, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d %d\n", m, d)
		exp := expected(m, d) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
