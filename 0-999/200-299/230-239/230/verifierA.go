package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type dragon struct{ x, y int }

func solveA(s int, dr []dragon) string {
	sort.Slice(dr, func(i, j int) bool { return dr[i].x < dr[j].x })
	for _, d := range dr {
		if s <= d.x {
			return "NO"
		}
		s += d.y
	}
	return "YES"
}

func runCase(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("invalid test case %d\n", idx)
			os.Exit(1)
		}
		s, _ := strconv.Atoi(fields[0])
		n, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+2*n {
			fmt.Printf("test %d: expected %d dragon pairs got %d numbers\n", idx, n, len(fields)-2)
			os.Exit(1)
		}
		dr := make([]dragon, n)
		for i := 0; i < n; i++ {
			x, _ := strconv.Atoi(fields[2+2*i])
			y, _ := strconv.Atoi(fields[3+2*i])
			dr[i] = dragon{x, y}
		}
		input := fmt.Sprintf("%d %d\n", s, n)
		for _, d := range dr {
			input += fmt.Sprintf("%d %d\n", d.x, d.y)
		}
		expect := solveA(s, dr)
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("read error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
