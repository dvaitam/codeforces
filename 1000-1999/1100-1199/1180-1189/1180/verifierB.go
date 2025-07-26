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

func run(bin string, input string) (string, error) {
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

func solve(a []int) []int {
	n := len(a)
	allZero := true
	for i := 0; i < n; i++ {
		if a[i] >= 0 {
			a[i] = ^a[i]
		}
		if ^a[i] != 0 {
			allZero = false
		}
	}
	if allZero {
		if n%2 == 1 {
			for i := 0; i < n; i++ {
				a[i] = 0
			}
		}
	} else {
		if n%2 == 1 {
			minVal := a[0]
			minIdx := 0
			for i := 1; i < n; i++ {
				if a[i] < minVal {
					minVal = a[i]
					minIdx = i
				}
			}
			a[minIdx] = ^a[minIdx]
		}
	}
	return a
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(scan.Text())
	for tc := 0; tc < T; tc++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test %d\n", tc+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "bad test %d\n", tc+1)
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scan.Text())
			arr[i] = v
		}
		input := fmt.Sprintf("%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		expectedArr := solve(append([]int(nil), arr...))
		var want strings.Builder
		if n > 0 {
			want.WriteString(fmt.Sprintf("%d", expectedArr[0]))
			for i := 1; i < n; i++ {
				want.WriteString(" ")
				want.WriteString(fmt.Sprintf("%d", expectedArr[i]))
			}
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want.String() {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", tc+1, input, want.String(), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", T)
}
