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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, k int, w []int) int {
	total := 0
	for _, x := range w {
		total += (x + k - 1) / k
	}
	return (total + 1) / 2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "empty testcases")
		os.Exit(1)
	}
	cases, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	for idx := 1; idx <= cases; idx++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing first line\n", idx)
			os.Exit(1)
		}
		header := strings.Fields(scanner.Text())
		if len(header) != 2 {
			fmt.Fprintf(os.Stderr, "case %d malformed header\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(header[0])
		k, _ := strconv.Atoi(header[1])
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing weights\n", idx)
			os.Exit(1)
		}
		wfields := strings.Fields(scanner.Text())
		if len(wfields) != n {
			fmt.Fprintf(os.Stderr, "case %d expected %d weights got %d\n", idx, n, len(wfields))
			os.Exit(1)
		}
		weights := make([]int, n)
		for i, v := range wfields {
			weights[i], _ = strconv.Atoi(v)
		}
		input := fmt.Sprintf("%d %d\n%s\n", n, k, strings.Join(wfields, " "))
		want := fmt.Sprintf("%d", expected(n, k, weights))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", cases)
}
