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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		val := scan.Text()
		input := val + "\n"
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
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", caseNum, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != "1 1 ^" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected '1 1 ^' got %s\n", caseNum, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
