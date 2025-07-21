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

func solveA(x, y int) string {
	turn := true
	for {
		if turn {
			if x >= 2 && y >= 2 {
				x -= 2
				y -= 2
			} else if x >= 1 && y >= 12 {
				x -= 1
				y -= 12
			} else if y >= 22 {
				y -= 22
			} else {
				return "Hanako"
			}
		} else {
			if y >= 22 {
				y -= 22
			} else if x >= 1 && y >= 12 {
				x -= 1
				y -= 12
			} else if x >= 2 && y >= 2 {
				x -= 2
				y -= 2
			} else {
				return "Ciel"
			}
		}
		turn = !turn
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([][2]int, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		x, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		y, _ := strconv.Atoi(scan.Text())
		cases[i] = [2]int{x, y}
		expected[i] = solveA(x, y)
	}
	for i, c := range cases {
		in := fmt.Sprintf("%d %d\n", c[0], c[1])
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = strings.NewReader(in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("execution failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
