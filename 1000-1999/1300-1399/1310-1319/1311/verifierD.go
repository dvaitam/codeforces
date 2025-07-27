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

const (
	INF  = 1000000000
	MAXA = 10000
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func processB(A, B, c, num int, ans *int, ansA, ansB, ansC *int) {
	jc := c / B
	if jc >= 1 {
		C := B * jc
		total := num + abs(C-c)
		if total < *ans {
			*ans = total
			*ansA = A
			*ansB = B
			*ansC = C
		}
	}
	C2 := B * (jc + 1)
	total2 := num + abs(C2-c)
	if total2 < *ans {
		*ans = total2
		*ansA = A
		*ansB = B
		*ansC = C2
	}
}

func processA(A, a, b, c, costA int, ans *int, ansA, ansB, ansC *int) {
	k := b / A
	if k >= 1 {
		B := A * k
		costB := abs(B - b)
		newNum := costA + costB
		if newNum <= *ans {
			processB(A, B, c, newNum, ans, ansA, ansB, ansC)
		}
	}
	B2 := A * (k + 1)
	costB2 := abs(B2 - b)
	newNum2 := costA + costB2
	if newNum2 <= *ans {
		processB(A, B2, c, newNum2, ans, ansA, ansB, ansC)
	}
}

func expected(a, b, c int) string {
	ans := INF
	ansA, ansB, ansC := 0, 0, 0
	for A := a; A >= 1; A-- {
		da := a - A
		if da > ans {
			break
		}
		processA(A, a, b, c, da, &ans, &ansA, &ansB, &ansC)
	}
	for A := a + 1; A <= MAXA; A++ {
		da := A - a
		if da > ans {
			break
		}
		processA(A, a, b, c, da, &ans, &ansA, &ansB, &ansC)
	}
	return fmt.Sprintf("%d\n%d %d %d", ans, ansA, ansB, ansC)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
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
		fields := strings.Fields(line)
		if len(fields) != 3 {
			fmt.Printf("test %d malformed\n", idx)
			os.Exit(1)
		}
		a, _ := strconv.Atoi(fields[0])
		b, _ := strconv.Atoi(fields[1])
		c, _ := strconv.Atoi(fields[2])
		exp := expected(a, b, c)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
