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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(arr []int) string {
	n := len(arr)
	if n == 0 {
		return "YES"
	}
	L := make([]int, n)
	R := make([]int, n)
	for i := 0; i < n; i++ {
		if i == 0 || arr[i] < L[i-1] {
			L[i] = arr[i]
		} else {
			L[i] = L[i-1]
		}
	}
	for i := n - 1; i >= 0; i-- {
		if i == n-1 || arr[i] < R[i+1] {
			R[i] = arr[i]
		} else {
			R[i] = R[i+1]
		}
	}
	for i := 0; i < n; i++ {
		if L[i]+R[i] < arr[i] {
			return "NO"
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
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
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scan.Text())
			arr[j] = v
		}
		var input bytes.Buffer
		fmt.Fprintln(&input, 1)
		fmt.Fprintln(&input, n)
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		expected := solveCase(arr)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
