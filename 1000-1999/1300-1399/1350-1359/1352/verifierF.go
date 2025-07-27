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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func checkString(s string, n0, n1, n2 int) bool {
	if len(s) != n0+n1+n2+1 {
		return false
	}
	for i := range s {
		if s[i] != '0' && s[i] != '1' {
			return false
		}
	}
	c0 := 0
	c1 := 0
	c2 := 0
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '0' && s[i+1] == '0' {
			c0++
		} else if s[i] == '1' && s[i+1] == '1' {
			c2++
		} else {
			c1++
		}
	}
	return c0 == n0 && c1 == n1 && c2 == n2
}

func runCase(bin string, n0, n1, n2 int) error {
	input := fmt.Sprintf("1\n%d %d %d\n", n0, n1, n2)
	out, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	s := strings.TrimSpace(out)
	if !checkString(s, n0, n1, n2) {
		return fmt.Errorf("output string invalid: %q", s)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(sc.Text())
	for i := 0; i < t; i++ {
		sc.Scan()
		n0, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		n1, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		n2, _ := strconv.Atoi(sc.Text())
		if err := runCase(bin, n0, n1, n2); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
