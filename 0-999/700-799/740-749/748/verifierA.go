package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(n, m, k int64) (int64, int64, string) {
	idx := k - 1
	seatsPerLane := 2 * m
	lane := idx/seatsPerLane + 1
	pos := idx % seatsPerLane
	desk := pos/2 + 1
	side := "L"
	if pos%2 == 1 {
		side = "R"
	}
	return lane, desk, side
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		panic(err)
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
		var n, m, k int64
		fmt.Sscan(line, &n, &m, &k)
		lane, desk, side := expected(n, m, k)
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		var gotLane, gotDesk int64
		var gotSide string
		fmt.Sscan(strings.TrimSpace(string(out)), &gotLane, &gotDesk, &gotSide)
		if gotLane != lane || gotDesk != desk || strings.ToUpper(gotSide) != side {
			fmt.Printf("Test %d failed: expected %d %d %s got %s\n", idx, lane, desk, side, strings.TrimSpace(string(out)))
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
