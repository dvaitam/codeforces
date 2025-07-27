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

func expected(n int) (string, string) {
	if n <= 30 {
		return "NO", ""
	}
	switch n {
	case 36:
		return "YES", "5 6 10 15"
	case 40:
		return "YES", "5 6 14 15"
	case 44:
		return "YES", "6 7 10 21"
	default:
		return "YES", fmt.Sprintf("6 10 14 %d", n-30)
	}
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	outLines := strings.Fields(strings.TrimSpace(out.String()))
	exp1, exp2 := expected(n)
	if len(outLines) == 0 {
		return fmt.Errorf("no output")
	}
	if outLines[0] != exp1 {
		return fmt.Errorf("expected %s got %s", exp1, outLines[0])
	}
	if exp1 == "YES" {
		got := strings.Join(outLines[1:], " ")
		if got != exp2 {
			return fmt.Errorf("expected numbers %q got %q", exp2, got)
		}
	} else {
		if len(outLines) > 1 {
			return fmt.Errorf("extra output after NO")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("bad number on line %d\n", idx+1)
			os.Exit(1)
		}
		idx++
		if err := runCase(bin, n); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
