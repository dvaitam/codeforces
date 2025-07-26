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

func runCase(bin string, a, b, x int) error {
	input := fmt.Sprintf("%d %d %d\n", a, b, x)
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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	s := strings.TrimSpace(out.String())
	if len(s) != a+b {
		return fmt.Errorf("expected length %d got %d", a+b, len(s))
	}
	ca := 0
	cb := 0
	trans := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			ca++
		} else if s[i] == '1' {
			cb++
		} else {
			return fmt.Errorf("invalid character %q", s[i])
		}
		if i > 0 && s[i] != s[i-1] {
			trans++
		}
	}
	if ca != a || cb != b || trans != x {
		return fmt.Errorf("wrong string properties: have %d zeros %d ones %d transitions", ca, cb, trans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		a, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		b, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		x, _ := strconv.Atoi(scan.Text())
		if err := runCase(bin, a, b, x); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
