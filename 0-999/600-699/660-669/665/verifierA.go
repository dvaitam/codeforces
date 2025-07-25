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

func expected(a, ta, b, tb, hh, mm int) string {
	start := hh*60 + mm
	end := start + ta
	dayStart := 5 * 60
	dayEnd := 23*60 + 59
	count := 0
	for t := dayStart; t <= dayEnd; t += b {
		tEnd := t + tb
		if t < end && tEnd > start && t != end && tEnd != start {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
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
	for i := 0; i < t; i++ {
		vals := make([]int, 6)
		for j := 0; j < 6; j++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			vals[j], _ = strconv.Atoi(scan.Text())
		}
		a, ta, b, tb, hh, mm := vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]
		input := fmt.Sprintf("%d %d\n%d %d\n%02d:%02d\n", a, ta, b, tb, hh, mm)
		exp := expected(a, ta, b, tb, hh, mm) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
