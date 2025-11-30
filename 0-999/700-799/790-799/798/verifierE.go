package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt.
const embeddedTestcasesE = `1 -1
6 3 1 1 1 0 4
1 1
3 0 1 3
1 -1
3 0 1 2
3 2 2 3
5 -1 -1 1 2 2
4 3 2 3 2
5 -1 0 0 4 2
5 1 0 -1 0 0
1 0
6 0 0 3 3 4 6
4 3 1 4 4
5 2 3 1 0 5
5 -1 4 2 5 1
3 3 -1 0
5 -1 2 -1 1 2
5 0 3 -1 3 3
1 -1
5 1 -1 1 0 1
3 2 2 2
1 1
2 2 1
6 1 -1 6 1 5 3
6 6 2 -1 6 1 6
3 3 1 2
5 5 5 5 4 0
5 2 2 0 1 2
3 0 0 1
2 1 2
2 1 0
3 0 0 1
5 0 3 4 4 3
1 0
4 0 4 0 3
3 3 0 -1
5 2 0 1 2 1
5 0 1 5 2 1
6 -1 -1 2 2 1 6
1 0
5 1 4 2 2 5
3 2 -1 -1
1 0
2 1 2
3 3 2 2
3 0 -1 -1
6 0 -1 4 0 2 -1
1 -1
6 3 0 5 1 2 3
3 0 0 3
2 0 0
2 1 2
3 0 3 1
2 -1 0
6 5 1 2 -1 4 6
6 1 3 0 3 2 -1
2 0 -1
4 0 0 2 3
1 1
3 0 0 2
3 2 0 1
6 2 0 4 6 3 4
5 3 1 0 1 3
4 -1 0 2 -1
5 4 -1 3 5 0
6 6 1 5 3 2 1
6 0 3 6 4 5 -1
3 2 2 3
2 0 -1
3 3 2 3
5 2 0 2 0 2
1 -1
5 2 -1 3 4 2
6 1 -1 2 2 2 0
5 2 2 -1 5 3
2 0 0
4 -1 0 -1 3
2 0 2
5 2 5 4 1 4
4 4 4 3 -1
3 1 0 -1
2 2 0
5 4 0 5 -1 0
2 2 2
1 1
3 1 2 1
6 3 0 0 6 4 0
2 0 2
2 -1 2
3 0 3 1
6 1 0 4 4 0 4
2 0 1
5 4 2 4 4 4
4 2 -1 0 2
1 0
6 1 2 6 4 1 1
5 2 1 4 0 -1
2 0 2
5 2 0 5 -1 1`

func solve798E(a, b, c, d int) int {
	return a ^ b ^ c ^ d
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesE), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 4 {
			fmt.Fprintf(os.Stderr, "case %d malformed\n", idx+1)
			os.Exit(1)
		}
		a, err1 := strconv.Atoi(fields[0])
		b, err2 := strconv.Atoi(fields[1])
		c, err3 := strconv.Atoi(fields[2])
		d, err4 := strconv.Atoi(fields[3])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error\n", idx+1)
			os.Exit(1)
		}
		want := strconv.Itoa(solve798E(a, b, c, d))
		input := fmt.Sprintf("%d %d %d %d\n", a, b, c, d)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
