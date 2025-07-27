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

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func smallestPrimeFactor(n int64) int64 {
	if n%2 == 0 {
		return 2
	}
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return i
		}
	}
	return n
}

func expectedPair(n int64) (int64, int64) {
	if n%2 == 0 {
		return n / 2, n / 2
	}
	sp := smallestPrimeFactor(n)
	b := n / sp
	a := n - b
	return a, b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	ns := make([]int64, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		v, _ := strconv.ParseInt(scan.Text(), 10, 64)
		ns[i] = v
	}
	out, err := runCandidate(bin, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(strings.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for idx, n := range ns {
		var a, b int64
		if !outScan.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		a, _ = strconv.ParseInt(outScan.Text(), 10, 64)
		if !outScan.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		b, _ = strconv.ParseInt(outScan.Text(), 10, 64)
		if a+b != n {
			fmt.Printf("case %d failed: sum mismatch\n", idx+1)
			os.Exit(1)
		}
		ea, eb := expectedPair(n)
		if !((a == ea && b == eb) || (a == eb && b == ea)) {
			fmt.Printf("case %d failed: expected %d %d (in any order) got %d %d\n", idx+1, ea, eb, a, b)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
