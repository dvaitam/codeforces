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

func runCandidate(bin, input string) (string, error) {
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

func expected(n, b, p int) string {
	bottles := (n - 1) * (2*b + 1)
	towels := n * p
	return fmt.Sprintf("%d %d", bottles, towels)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
	for i := 1; i <= t; i++ {
		if !scan.Scan() {
			fmt.Printf("missing n for case %d\n", i)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		b, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		pVal, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d %d %d\n", n, b, pVal)
		expect := expected(n, b, pVal)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
