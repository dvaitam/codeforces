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

func expected(n int, x int64, a []int64) string {
	var maxA int64
	for _, v := range a {
		if v > maxA {
			maxA = v
		}
	}
	lo := int64(1)
	hi := maxA + x + 1
	for lo+1 < hi {
		mid := (lo + hi) / 2
		var need int64
		for _, v := range a {
			if mid > v {
				need += mid - v
				if need > x {
					break
				}
			}
		}
		if need <= x {
			lo = mid
		} else {
			hi = mid
		}
	}
	return fmt.Sprintf("%d\n", lo)
}

func runCase(exe, input, expect string) error {
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
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		parts := strings.Fields(scan.Text())
		if len(parts) != 2 {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		x, _ := strconv.ParseInt(parts[1], 10, 64)
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		arrParts := strings.Fields(scan.Text())
		a := make([]int64, n)
		for i := 0; i < n && i < len(arrParts); i++ {
			val, _ := strconv.ParseInt(arrParts[i], 10, 64)
			a[i] = val
		}
		input := fmt.Sprintf("1\n%d %d\n%s\n", n, x, strings.Join(arrParts, " "))
		exp := expected(n, x, a)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
