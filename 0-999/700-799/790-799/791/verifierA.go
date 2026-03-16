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
13 19
2 10
17 20
13 17
16 18
19 19
17 18
10 12
4 12
18 20
20 20
10 11
3 13
16 20
4 15
14 16
20 20
18 19
15 19
9 9
18 18
3 15
1 20
16 18
8 19
11 12
7 16
8 11
5 19
3 5
11 19
16 16
10 18
10 11
18 19
18 18
20 20
15 15
20 20
11 20
8 12
6 9
6 6
20 20
16 16
3 7
5 6
3 20
13 17
17 18
7 17
19 20
19 20
15 18
12 13
11 20
4 19
19 20
7 10
1 9
4 11
12 14
11 17
2 5
5 12
2 20
18 20
3 3
4 10
20 20
13 14
12 13
2 2
7 9
4 19
7 18
2 2
18 19
20 20
9 10
8 9
10 15
14 15
2 18
15 15
20 20
13 16
9 14
16 20
6 17
7 19
2 7
6 11
17 19
4 18
6 6
16 19
19 20
12 18
9 11
`

func expected(a, b int) string {
	years := 0
	for a <= b {
		a *= 3
		b *= 2
		years++
	}
	return fmt.Sprintf("%d\n", years)
}

func runCase(bin string, a, b int) error {
	input := fmt.Sprintf("%d %d\n", a, b)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected(a, b))
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
	bin := os.Args[1]


	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("invalid test count")
		os.Exit(1)
	}

	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		a, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		b, _ := strconv.Atoi(scanner.Text())
		if err := runCase(bin, a, b); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
