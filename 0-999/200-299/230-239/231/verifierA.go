package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"7\n1 0 1\n1 1 1\n1 1 0\n0 1 0\n0 1 0\n1 0 0\n1 1 0",
	"6\n1 1 0\n1 1 1\n0 0 0\n1 0 1\n1 0 1\n0 0 0",
	"4\n0 1 0\n0 1 1\n0 1 1\n0 1 0",
	"10\n1 1 0\n1 1 0\n1 0 0\n0 0 1\n1 0 0\n0 0 0\n0 1 1\n0 0 1\n1 1 1\n1 0 1",
	"10\n0 1 1\n0 0 0\n1 0 0\n1 0 1\n1 0 0\n0 0 0\n0 0 0\n0 0 1\n0 1 0\n0 0 0",
	"3\n0 1 0\n0 0 1\n0 1 0",
	"4\n0 1 1\n1 0 0\n1 0 0\n1 0 1",
	"6\n1 0 0\n0 0 0\n1 1 0\n1 0 0\n1 1 1\n1 1 1",
	"3\n0 1 0\n1 0 1\n0 0 1",
	"6\n1 1 0\n1 1 1\n0 0 0\n1 0 0\n0 1 1\n1 0 1",
	"10\n1 0 0\n1 0 1\n0 1 1\n0 0 1\n1 1 1\n0 1 0\n0 0 0\n1 1 1\n0 0 0\n0 0 0",
	"2\n0 1 1\n0 0 1",
	"7\n0 0 1\n1 0 1\n0 1 0\n0 1 0\n0 0 0\n1 1 1\n0 1 1",
	"7\n1 0 0\n1 1 0\n0 0 1\n1 1 1\n0 1 0\n0 1 0\n0 1 1",
	"7\n0 1 0\n1 0 0\n1 0 0\n1 1 1\n1 1 0\n0 1 1\n1 1 0",
	"8\n0 1 1\n0 1 1\n0 0 1\n0 0 0\n0 0 0\n1 0 0\n1 1 1\n1 0 1",
	"2\n1 0 1\n0 1 0",
	"6\n0 0 0\n1 0 0\n1 1 1\n0 0 0\n0 0 0\n1 1 1",
	"9\n0 0 1\n0 1 1\n1 1 1\n1 0 0\n0 1 1\n1 0 0\n1 0 1\n1 1 0\n0 0 0",
	"4\n0 0 1\n0 1 1\n0 0 1\n1 1 0",
	"10\n0 0 1\n0 1 1\n1 1 0\n0 0 0\n1 1 1\n0 1 1\n1 1 1\n1 0 0\n0 0 1\n1 0 0",
	"8\n1 0 1\n0 0 1\n1 0 1\n1 1 0\n0 1 0\n0 0 0\n1 0 1\n0 1 1",
	"2\n0 1 1\n0 1 0",
	"7\n1 0 0\n0 1 1\n1 0 1\n1 0 0\n0 1 1\n0 0 0\n0 1 1",
	"2\n0 1 1\n0 0 1",
	"2\n0 1 0\n0 1 1",
	"7\n0 1 1\n1 0 1\n0 0 1\n0 1 1\n0 0 0\n1 0 0\n0 0 0",
	"2\n0 1 1\n0 0 0",
	"3\n1 0 1\n1 1 0\n0 0 1",
	"6\n1 1 1\n1 0 0\n0 1 0\n1 0 0\n0 1 0\n0 0 1",
	"6\n1 0 0\n0 1 0\n0 0 1\n1 1 0\n1 1 1\n0 0 1",
	"6\n0 1 1\n0 1 1\n1 0 1\n1 1 0\n0 0 1\n0 0 1",
	"7\n0 1 1\n0 0 1\n1 0 1\n1 0 1\n0 1 0\n1 1 0\n0 1 0",
	"5\n0 1 1\n0 1 0\n1 1 0\n0 1 0\n0 1 1",
	"9\n1 1 0\n1 1 1\n0 0 0\n0 1 1\n1 0 0\n0 0 1\n1 1 0\n0 0 1\n1 1 1",
	"8\n1 0 0\n0 0 1\n1 1 0\n1 1 1\n1 0 1\n1 1 1\n0 1 1\n0 1 0",
	"10\n0 0 1\n1 1 0\n0 0 1\n0 1 0\n0 0 1\n1 0 0\n1 0 0\n1 1 1\n1 0 1\n1 1 0",
	"8\n0 1 0\n0 0 1\n1 1 1\n1 1 0\n0 1 0\n0 1 0\n0 0 0\n1 0 0",
	"4\n1 0 1\n1 1 1\n0 1 1\n1 1 1",
	"4\n1 1 0\n0 0 0\n1 0 1\n0 0 0",
	"7\n0 1 0\n1 0 1\n0 0 1\n1 1 1\n1 0 0\n0 1 0\n0 1 0",
	"4\n1 1 1\n1 1 1\n1 0 0\n1 0 0",
	"8\n1 1 0\n1 1 1\n1 0 0\n0 0 0\n0 1 0\n1 0 1\n1 0 0\n0 1 0",
	"3\n1 1 1\n1 0 1\n0 1 1",
	"9\n1 1 1\n0 0 0\n0 0 1\n0 0 1\n0 1 0\n0 1 1\n0 1 0\n1 0 0\n0 1 1",
	"4\n1 1 1\n0 0 0\n1 1 0\n0 1 0",
	"3\n0 0 0\n1 0 0\n0 0 1",
	"6\n0 1 0\n0 1 1\n1 1 0\n0 0 1\n1 1 0\n1 0 1",
	"4\n1 1 1\n0 0 0\n1 0 1\n1 0 1",
	"6\n0 1 0\n1 0 0\n1 0 1\n0 0 0\n0 0 0\n1 0 0",
	"6\n1 0 0\n1 1 0\n1 1 0\n1 1 1\n0 1 1\n1 0 0",
	"7\n1 0 0\n0 0 0\n1 0 0\n0 0 0\n0 0 1\n0 1 0\n0 1 1",
	"7\n0 0 1\n0 1 1\n0 1 0\n1 0 1\n1 1 1\n1 0 1\n1 1 0",
	"9\n1 0 1\n1 0 1\n1 1 0\n0 1 0\n1 1 0\n1 1 1\n0 0 1\n0 1 0\n0 1 0",
	"10\n0 0 1\n1 1 1\n1 1 1\n1 0 0\n0 0 1\n0 1 0\n1 1 0\n0 1 1\n0 0 0\n0 0 1",
	"7\n1 1 1\n0 1 1\n1 1 1\n0 1 0\n0 0 1\n1 1 0\n1 0 0",
	"6\n0 0 1\n1 0 0\n0 1 0\n1 0 1\n0 1 0\n0 0 1",
	"2\n1 1 1\n1 1 1",
	"10\n1 0 1\n0 1 0\n1 0 1\n1 1 1\n0 1 1\n0 1 1\n0 1 1\n0 0 0\n1 1 0\n0 0 0",
	"9\n1 0 1\n0 1 1\n0 1 1\n0 0 1\n1 0 0\n0 0 1\n0 0 1\n0 1 0\n0 0 1",
	"4\n1 1 1\n1 1 0\n1 0 0\n0 0 1",
	"2\n0 1 0\n1 0 1",
	"3\n1 0 0\n1 0 0\n0 1 1",
	"10\n1 1 1\n1 0 0\n1 0 0\n0 0 1\n0 0 1\n0 1 1\n0 0 1\n0 1 0\n1 1 1\n1 1 1",
	"3\n0 1 1\n0 1 0\n1 0 0",
	"1\n1 0 0",
	"7\n0 1 1\n1 1 0\n0 1 1\n0 0 0\n1 1 0\n1 0 0\n1 0 0",
	"7\n1 0 0\n0 0 0\n0 1 0\n1 1 0\n0 0 1\n1 0 0\n1 0 0",
	"7\n1 1 0\n1 0 0\n0 1 0\n0 0 0\n1 0 0\n0 0 0\n0 1 0",
	"7\n1 1 0\n0 1 0\n1 0 0\n1 1 0\n1 1 0\n1 1 0\n1 0 1",
	"7\n0 1 1\n1 0 1\n0 1 0\n1 1 0\n0 0 1\n0 1 1\n1 0 1",
	"10\n1 0 0\n1 1 1\n0 1 0\n1 1 0\n1 1 0\n0 1 0\n1 1 1\n0 1 1\n0 1 1\n1 0 1",
	"3\n0 0 1\n0 0 1\n1 1 0",
	"2\n1 0 0\n0 0 0",
	"4\n0 0 1\n1 1 1\n0 1 1\n0 1 1",
	"9\n0 1 1\n0 1 1\n1 1 1\n1 1 0\n0 0 0\n0 0 0\n1 0 1\n0 0 0\n0 1 1",
	"6\n0 1 1\n1 1 0\n0 0 0\n0 0 1\n0 1 1\n1 0 0",
	"7\n0 1 0\n1 0 1\n1 0 1\n0 1 0\n1 1 1\n1 1 0\n0 1 0",
	"5\n0 1 0\n0 1 1\n1 0 1\n1 0 0\n1 0 1",
	"2\n0 0 0\n0 0 0",
	"7\n1 0 1\n1 1 0\n0 0 0\n1 1 1\n0 0 1\n1 1 0\n0 1 1",
	"4\n1 1 0\n0 1 0\n1 1 1\n1 0 0",
	"7\n1 1 0\n1 1 1\n1 1 1\n1 0 1\n1 0 1\n0 1 0\n0 1 1",
	"4\n0 1 1\n1 1 0\n1 1 1\n0 1 0",
	"2\n1 1 1\n0 1 0",
	"10\n0 1 0\n1 0 1\n0 0 0\n1 0 0\n0 0 1\n0 0 0\n0 0 1\n0 0 1\n0 0 0\n0 1 1",
	"8\n0 1 0\n1 0 1\n0 1 1\n1 1 0\n0 1 0\n0 1 0\n0 0 1\n0 1 0",
	"7\n0 0 1\n0 1 1\n0 1 1\n0 1 1\n1 0 1\n0 1 0\n1 1 1",
	"1\n1 0 1",
	"8\n0 1 0\n1 1 1\n0 0 1\n1 1 1\n0 0 0\n1 0 0\n0 1 0\n1 0 1",
	"7\n0 1 1\n0 0 1\n0 1 0\n1 0 1\n0 1 0\n1 0 0\n1 1 0",
	"1\n0 0 0",
	"4\n1 0 0\n1 1 1\n1 0 0\n1 1 1",
	"7\n1 1 1\n0 0 0\n0 0 1\n1 0 1\n0 0 1\n0 0 1\n1 0 0",
	"3\n1 1 1\n0 0 1\n0 0 1",
	"10\n1 1 0\n0 0 0\n1 0 1\n1 0 1\n1 0 1\n1 1 0\n1 0 0\n0 1 1\n0 1 1\n1 1 0",
	"2\n1 0 0\n0 0 1",
	"4\n0 1 1\n1 0 1\n1 1 0\n0 1 0",
	"3\n1 0 1\n1 0 1\n0 1 0",
	"5\n0 0 0\n1 0 0\n1 0 0\n0 0 1\n0 1 1",
}

func solveCaseA(n int, tri [][3]int) int {
	cnt := 0
	for _, t := range tri {
		if t[0]+t[1]+t[2] >= 2 {
			cnt++
		}
	}
	return cnt
}

func parseCase(input string) (int, [][3]int, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return 0, nil, fmt.Errorf("empty case")
	}
	idx := 0
	n, err := strconv.Atoi(fields[idx])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid n: %v", err)
	}
	idx++
	if len(fields) != 1+3*n {
		return 0, nil, fmt.Errorf("expected %d numbers, got %d", 1+3*n, len(fields))
	}
	tri := make([][3]int, n)
	for i := 0; i < n; i++ {
		a, _ := strconv.Atoi(fields[idx])
		b, _ := strconv.Atoi(fields[idx+1])
		c, _ := strconv.Atoi(fields[idx+2])
		idx += 3
		tri[i] = [3]int{a, b, c}
	}
	return n, tri, nil
}

func runCase(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tcStr := range rawTestcases {
		n, tri, err := parseCase(tcStr)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := fmt.Sprintf("%d", solveCaseA(n, tri))
		got, err := runCase(bin, tcStr+"\n")
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
