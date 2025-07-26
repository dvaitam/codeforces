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

func solveA(a []int) string {
	n := len(a)
	if n == 1 || (n == 2 && a[0] == a[1]) {
		return "-1"
	}
	sum := 0
	for _, v := range a {
		sum += v
	}
	if sum-a[0] != a[0] {
		return "1\n1"
	}
	return "2\n1 2"
}

func runCandidate(bin, input string) (string, error) {
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
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			a[i], _ = strconv.Atoi(scan.Text())
		}
		var inSB strings.Builder
		inSB.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				inSB.WriteByte(' ')
			}
			inSB.WriteString(strconv.Itoa(v))
		}
		inSB.WriteByte('\n')
		input := inSB.String()
		expected := solveA(a)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("case %d failed: expected %q got %q\ninput:\n%s", caseIdx+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
