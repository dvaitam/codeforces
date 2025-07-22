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

func expected(r, g, b int64) int64 {
	base := r/3 + g/3 + b/3
	best := base
	for k := int64(1); k <= 2; k++ {
		if r < k || g < k || b < k {
			break
		}
		cur := k + (r-k)/3 + (g-k)/3 + (b-k)/3
		if cur > best {
			best = cur
		}
	}
	return best
}

func runCase(bin string, r, g, b int64) error {
	input := fmt.Sprintf("%d %d %d\n", r, g, b)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected(r, g, b) {
		return fmt.Errorf("expected %d got %d", expected(r, g, b), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing r\n", i+1)
			os.Exit(1)
		}
		rVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing g\n", i+1)
			os.Exit(1)
		}
		gVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing b\n", i+1)
			os.Exit(1)
		}
		bVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if err := runCase(bin, rVal, gVal, bVal); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
