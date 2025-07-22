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

func solveCase(n int) string {
	k := 1
	for (k*(k-1))/2 <= n {
		k++
	}
	k--
	guests := make([][]int, k)
	id := 1
	for i := 0; i < k; i++ {
		for j := i + 1; j < k; j++ {
			guests[i] = append(guests[i], id)
			guests[j] = append(guests[j], id)
			id++
		}
	}
	var out strings.Builder
	fmt.Fprintln(&out, k)
	for i := 0; i < k; i++ {
		for j, x := range guests[i] {
			if j > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(strconv.Itoa(x))
		}
		if i+1 < k {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
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
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		expected := solveCase(n)
		out, err := run(bin, fmt.Sprintf("%d\n", n))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expected) {
			fmt.Printf("case %d failed:\nexpected:\n%s\n----\ngot:\n%s\n", i+1, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
