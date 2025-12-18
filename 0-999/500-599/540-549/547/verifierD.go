package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `4 10 9 3 6 10 8 10 2
10 1 8 5 9 4 4 8 9 9 8 7 3 4 3 9 7 1 2 3 10
1 5 1
5 8 10 7 7 7 10 8 3 6 2
1 3 8
4 5 7 5 7 9 7 10 6
9 10 7 10 4 6 1 5 10 3 6 9 10 10 2 4 10 5 5
2 2 8 8 2
6 2 7 3 1 5 7 7 2 1 10 10 1
7 10 6 9 5 9 4 1 5 1 2 2 10 9 1
4 7 5 10 5 3 1 6 6
6 3 7 7 8 9 7 10 9 2 10 9 5
7 4 5 7 5 9 5 9 6 1 7 10 6 1 7
10 10 3 1 6 8 6 6 10 5 8 1 10 1 1 6 5 8 5 10 10
6 3 6 3 6 6 10 5 5 7 2 1 10
3 5 9 4 5 4 6
3 7 2 2 10 6 6
4 8 3 2 6 4 10 8 5
4 2 1 9 4 6 10 3 5
6 2 10 6 10 3 7 5 9 5 8 6 7
5 7 10 7 1 7 3 4 1 8 10
9 7 9 4 1 8 9 5 9 6 4 2 10 5 2 4 1 1 9
4 7 10 1 1 8 2 3 9
5 4 1 9 9 7 1 10 2 6 3
5 9 8 1 6 4 4 2 9 2 3
4 5 3 1 8 10 7 1 5
4 5 10 9 9 7 1 8 6
1 1 3
1 2 1
2 8 1 2 9
9 8 6 3 6 2 6 7 7 10 5 6 5 4 6 7 2 3 9
1 7 2
10 3 1 6 8 10 9 7 1 10 7 1 6 8 6 7 7 8 1 4 4
9 5 10 2 7 4 7 3 1 6 6 9 5 2 8 2 9 7 2
6 10 9 2 10 1 8 3 4 7 1 9 2
10 2 7 3 1 6 2 1 2 8 5 10 5 2 1 10 9 9 4 2 9
2 9 1 9 6
10 3 2 4 3 4 8 10 7 5 6 10 7 6 9 7 2 7 9 4 7
3 7 10 10 9 8 3
7 3 3 2 8 8 9 8 10 3 3 5 4 3 10
9 6 4 9 5 7 10 10 10 5 4 5 1 5 8 7 4 3 10
6 4 6 8 3 7 8 10 4 8 10 9 1
8 2 7 1 8 4 4 2 4 5 4 4 5 3 3 10 1
5 3 1 6 3 7 2 2 2 2 5
5 1 6 8 10 6 1 1 6 6 7
7 8 2 4 10 8 7 3 9 6 2 5 2 7 2
8 9 5 2 9 6 6 8 5 5 2 6 10 9 9 2 8
9 6 1 5 10 3 3 3 6 8 2 2 9 3 6 10 7 9 5
3 8 8 5 3 2 2
3 9 9 10 7 6 2
5 5 7 1 3 1 8 9 5 4 9
6 6 7 8 9 2 6 8 2 3 5 10 2
2 10 2 3 4
10 7 7 3 10 10 3 7 4 9 9 3 10 3 4 5 6 5 1 8 7
7 6 9 10 5 8 9 5 8 1 10 4 1 2 4
8 3 9 8 4 4 9 4 1 9 8 2 10 5 3 3 8
2 10 1 1 6
10 4 9 2 8 9 1 6 6 6 6 3 2 10 1 2 6 4 2 4 7
4 8 6 2 1 7 2 4 3
7 8 8 2 9 7 4 8 5 1 8 8 7 8 3
8 1 5 6 6 8 9 6 10 7 4 1 4 5 6 3 8
9 4 3 4 1 3 10 7 9 3 1 3 2 10 3 8 8 3 1
1 7 8
6 7 1 1 4 7 1 7 8 1 4 4 2
7 8 4 3 6 10 2 6 2 10 1 5 5 8 5
8 4 9 5 1 6 6 6 2 1 7 2 10 10 1 2 1
2 1 3 9 1
8 1 4 9 6 4 8 6 8 6 1 7 5 10 7 2 5
3 7 2 9 7 9 6
9 7 3 7 9 6 3 6 7 8 4 8 8 6 5 3 9 10 7
8 1 3 3 1 8 2 2 6 4 10 1 10 1 8 8 6
6 1 2 4 7 2 6 10 5 2 8 2 4
4 1 3 3 10 1 2 4 5
4 4 9 9 7 9 10 6 9
4 8 3 10 2 1 2 10 1
2 4 5 2 2
8 7 4 10 2 8 6 7 10 8 2 5 10 8 7 4 2
9 1 8 5 2 6 6 4 8 2 9 6 7 2 10 9 4 4 6
1 6 4
7 8 2 5 4 6 3 4 4 10 8 9 7 6 4
10 7 8 7 8 10 1 5 1 3 2 1 3 5 9 9 1 8 1 4 4
5 8 7 1 6 8 4 5 3 2 8
5 7 8 2 4 3 8 5 7 6 3
7 5 8 8 9 9 4 6 5 5 1 8 6 6 5
4 9 1 1 3 9 3 9 1
3 1 1 4 8 6 6
9 1 8 3 4 1 5 7 6 1 10 9 2 8 5 5 4 8 7
5 6 1 1 7 1 3 10 4 3 7
9 6 9 3 5 1 3 1 1 8 1 8 8 9 10 9 7 6 9
3 5 3 2 3 9 2
7 6 8 8 5 5 8 5 9 3 10 6 3 9 1
7 8 4 8 10 10 5 1 6 10 10 9 2 8 3
5 5 2 7 2 6 1 9 8 8 4
5 6 3 7 7 6 1 5 4 1 6
6 10 7 9 5 1 3 7 5 7 2 8 4
4 2 9 2 2 1 5 2 7
5 8 8 5 5 9 9 1 3 4 8
3 3 3 3 8 7 1
3 7 1 3 3 5 4
3 3 1 9 3 9 4`

func runProgram(bin, input string) (string, error) {
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

func check(input, output string) error {
	inScanner := bufio.NewScanner(strings.NewReader(input))
	inScanner.Split(bufio.ScanWords)
	
	if !inScanner.Scan() {
		return fmt.Errorf("empty input")
	}
	n, err := strconv.Atoi(inScanner.Text())
	if err != nil {
		return fmt.Errorf("bad n: %v", err)
	}

	type Point struct {
		x, y int
	}
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		if !inScanner.Scan() {
			return fmt.Errorf("missing x for point %d", i)
		}
		x, _ := strconv.Atoi(inScanner.Text())
		if !inScanner.Scan() {
			return fmt.Errorf("missing y for point %d", i)
		}
		y, _ := strconv.Atoi(inScanner.Text())
		points[i] = Point{x, y}
	}

	ans := strings.TrimSpace(output)
	if len(ans) != n {
		return fmt.Errorf("expected length %d, got %d", n, len(ans))
	}

	// Maps coordinate -> balance (red - blue)
	xBal := make(map[int]int)
	yBal := make(map[int]int)

	for i, ch := range ans {
		val := 0
		if ch == 'r' {
			val = 1
		} else if ch == 'b' {
			val = -1
		} else {
			return fmt.Errorf("invalid char '%c' at index %d", ch, i)
		}
		xBal[points[i].x] += val
		yBal[points[i].y] += val
	}

	for x, bal := range xBal {
		if math.Abs(float64(bal)) > 1 {
			return fmt.Errorf("x=%d has balance %d", x, bal)
		}
	}
	for y, bal := range yBal {
		if math.Abs(float64(bal)) > 1 {
			return fmt.Errorf("y=%d has balance %d", y, bal)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcases))
	// Buffer for long lines
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		got, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed to run: %v\n", idx, err)
			os.Exit(1)
		}
		if err := check(input, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
