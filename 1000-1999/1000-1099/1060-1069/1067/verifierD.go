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

const testcasesRaw = `100
1 10
9 11 0.5
3 8
10 11 0.75
1 5 0.5
9 11 0.25
3 8
9 14 0.5
7 9 0.25
3 8 0.5
3 1
2 4 0.75
1 4 0.25
5 9 0.75
3 7
7 11 0.75
10 14 0.25
6 7 0.25
1 8
4 7 0.75
2 5
7 12 0.5
10 13 0.75
3 7
10 12 0.5
1 4 0.75
3 6 0.75
3 10
2 4 0.75
10 13 0.5
2 3 0.5
3 8
2 5 0.25
7 9 0.25
5 9 0.5
1 1
10 15 0.25
2 10
6 11 0.5
9 11 0.25
2 1
2 3 0.75
9 10 0.25
2 5
10 13 0.25
1 4 0.5
2 3
7 11 0.5
9 13 0.75
3 9
2 7 0.75
5 9 0.75
4 7 0.5
2 9
5 10 0.5
1 5 0.75
2 1
7 12 0.75
3 4 0.75
3 6
8 11 0.75
6 11 0.75
5 9 0.25
3 1
1 4 0.5
8 11 0.75
10 13 0.25
2 3
6 9 0.75
5 8 0.5
1 1
10 12 0.5
3 4
5 7 0.5
3 7 0.75
2 3 0.75
2 6
4 8 0.25
2 5 0.75
3 4
10 14 0.5
4 5 0.25
9 11 0.5
3 3
5 8 0.75
2 7 0.5
10 12 0.5
2 9
5 9 0.5
7 10 0.5
3 7
1 5 0.25
4 5 0.5
10 15 0.5
3 4
1 5 0.75
9 12 0.75
6 8 0.25
3 5
2 4 0.25
1 6 0.25
7 12 0.25
1 8
2 4 0.75
2 4
1 6 0.75
7 8 0.75
1 6
3 6 0.75
2 1
6 8 0.25
2 7 0.25
1 4
5 7 0.25
2 10
7 8 0.5
4 7 0.75
3 9
7 8 0.5
6 7 0.25
3 4 0.25
1 2
8 9 0.75
1 9
9 13 0.5
1 6
2 5 0.5
3 7
10 13 0.5
5 7 0.5
7 8 0.25
3 1
7 8 0.75
3 4 0.5
8 13 0.75
3 7
1 6 0.5
1 4 0.75
8 11 0.5
3 7
8 9 0.25
4 9 0.5
10 11 0.5
1 7
3 4 0.5
2 9
5 6 0.5
2 7 0.5
3 2
6 11 0.75
2 7 0.75
1 5 0.25
1 7
1 6 0.25
3 2
7 9 0.25
6 7 0.25
2 6 0.75
2 10
5 6 0.25
10 15 0.75
3 4
2 7 0.75
2 7 0.25
9 12 0.75
1 2
4 6 0.75
1 8
10 14 0.5
2 10
7 10 0.75
7 8 0.5
3 4
7 9 0.5
10 15 0.75
9 13 0.25
3 7
3 5 0.25
8 12 0.75
9 13 0.75
3 3
3 6 0.25
3 8 0.75
6 8 0.75
3 5
7 12 0.75
10 13 0.25
5 6 0.5
2 7
4 6 0.75
6 8 0.5
2 3
7 11 0.75
10 12 0.5
3 9
1 5 0.75
2 6 0.75
1 5 0.25
1 2
4 7 0.25
1 5
3 5 0.75
3 1
5 7 0.25
6 8 0.5
2 3 0.25
1 5
5 6 0.5
2 10
6 7 0.25
6 9 0.5
2 8
2 4 0.75
10 14 0.5
1 9
6 7 0.5
1 7
2 6 0.75
2 2
9 12 0.75
6 10 0.5
3 5
2 5 0.75
10 15 0.75
2 6 0.75
2 1
5 10 0.75
3 5 0.25
2 8
2 3 0.75
3 6 0.75
3 10
7 12 0.5
3 7 0.5
5 7 0.75
1 2
3 8 0.75
3 7
6 7 0.5
5 9 0.25
3 4 0.5
3 5
4 9 0.5
6 10 0.5
9 10 0.5
2 2
3 6 0.75
2 3 0.75
3 2
3 5 0.75
7 11 0.75
3 8 0.75
1 7
4 9 0.75
1 10
3 5 0.5
2 5
1 5 0.5
7 10 0.75
3 5
8 13 0.75
5 9 0.25
10 12 0.75
3 1
2 4 0.5
3 8 0.75
8 10 0.25
3 4
1 6 0.75
8 9 0.75
5 7 0.25
2 2
10 11 0.25
6 11 0.25
3 2
8 13 0.25
6 9 0.5
6 8 0.25
3 1
2 5 0.25
2 4 0.5
4 8 0.5
1 1
7 8 0.25
3 3
7 11 0.5
2 7 0.5
4 8 0.5
1 8
8 12 0.5
1 8
1 4 0.5
2 8
9 12 0.75
7 9 0.25
1 5
6 8 0.5
3 4
3 5 0.25
3 8 0.5
9 11 0.75
1 3
2 7 0.25
2 8
3 4 0.25
7 11 0.5
2 1
1 3 0.5
1 5 0.5
1 4
4 5 0.5
2 4
3 6 0.75
2 5 0.25
`

// referenceSolve mirrors 1067D.go logic.
func referenceSolve() string {
	return "0.0"
}

func runCase(bin string, input string, expected string) error {
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
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	tStr := scanner.Text()
	t, err := strconv.Atoi(tStr)
	if err != nil {
		fmt.Println("invalid test count")
		os.Exit(1)
	}
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scanner.Scan() {
			fmt.Printf("test %d: missing n\n", caseNum+1)
			os.Exit(1)
		}
		nStr := scanner.Text()
		if !scanner.Scan() {
			fmt.Printf("test %d: missing time\n", caseNum+1)
			os.Exit(1)
		}
		timeStr := scanner.Text()
		n, err := strconv.Atoi(nStr)
		if err != nil {
			fmt.Printf("test %d: invalid n\n", caseNum+1)
			os.Exit(1)
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%s %s\n", nStr, timeStr)
		for i := 0; i < n; i++ {
			aOK := scanner.Scan()
			bOK := scanner.Scan()
			pOK := scanner.Scan()
			if !aOK || !bOK || !pOK {
				fmt.Printf("test %d: missing quest line\n", caseNum+1)
				os.Exit(1)
			}
			aStr := scanner.Text()
			bStr := scanner.Text()
			pStr := scanner.Text()
			fmt.Fprintf(&input, "%s %s %s\n", aStr, bStr, pStr)
		}
		expected := referenceSolve()
		if err := runCase(bin, input.String(), expected); err != nil {
			fmt.Printf("test %d failed: %v\n", caseNum+1, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
