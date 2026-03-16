package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)
const testcasesFRaw = `
3 2 107 011 110
2 0 22 
5 4 941 10100 01100 10010 01010
3 0 850 
4 2 200 0101 0011
4 0 867 
6 2 689 000110 010001
3 2 288 110 110
6 5 905 001010 010100 000110 011000 100001
3 2 267 101 110
4 3 549 0011 1010 1001
2 1 937 11
4 1 366 0011
6 2 651 010010 101000
2 2 236 11 11
3 2 472 101 110
4 2 20 0110 1010
5 4 698 10100 01001 00110 01001
5 4 986 01100 10001 10100 01010
4 2 299 1100 0101
3 0 62 
2 0 756 
3 1 623 110
6 1 331 100001
6 2 794 010100 010100
3 2 422 110 110
5 1 664 01010
5 1 34 10100
4 1 540 1100
5 2 147 10100 10100
6 3 987 100001 000101 100100
3 1 346 011
4 3 542 0110 0110 1001
2 1 986 11
5 4 132 10100 01010 00011 01001
3 0 218 
3 0 626 
4 0 407 
5 1 565 11000
3 2 840 011 011
2 0 37 
6 0 502 
6 6 265 010010 101000 100010 011000 000101 000101
5 1 808 00110
3 3 878 110 101 011
5 5 797 10010 01100 00011 10100 01001
3 1 725 110
5 2 266 10001 11000
4 0 560 
2 1 807 11
6 1 423 001001
3 1 886 101
2 2 931 11 11
6 4 454 010100 001010 001001 010010
2 2 149 11 11
6 6 360 000011 100001 010100 001010 101000 010100
6 0 538 
2 1 828 11
3 2 678 101 101
5 5 909 01100 01010 00011 10100 10001
4 4 643 0101 0101 1010 1010
4 1 920 0011
6 4 847 001100 000110 110000 100001
5 5 944 00101 10100 10001 01010 01010
2 1 137 11
5 4 699 00101 11000 11000 00101
4 0 126 
4 0 901 
4 4 913 0110 0110 1001 1001
2 2 903 11 11
4 3 521 1100 1100 0011
6 5 312 001100 010100 100010 101000 010001
4 1 89 1100
6 6 17 010001 100010 000110 010100 001001 101000
6 4 194 100100 001010 001100 100010
3 2 594 011 011
2 0 916 
4 0 676 
2 2 332 11 11
2 1 533 11
3 0 346 
6 1 715 101000
6 6 846 100010 110000 000101 010001 001010 001100
3 1 370 110
4 3 649 1001 0101 1100
6 6 523 001010 001001 010010 100100 110000 000101
4 2 775 0011 0101
6 1 264 001010
2 0 414 
4 4 296 0110 1010 1001 0101
3 2 784 110 110
6 1 172 110000
6 6 99 101000 001010 010001 110000 000101 000110
4 0 609 
4 2 755 1010 1100
2 0 976 
3 0 691 
4 3 108 1100 0101 0011
6 2 372 000101 100001
3 1 368 110
`


func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "489F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	scanner := bufio.NewScanner(strings.NewReader(testcasesFRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		mod, _ := strconv.Atoi(fields[2])
		if len(fields) != 3+m {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d %d %d\n", n, m, mod))
		for i := 0; i < m; i++ {
			b.WriteString(fields[3+i] + "\n")
		}
		input := b.String()
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
