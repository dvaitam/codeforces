package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func expected(a, b, c int) string {
	if a+b >= 10 || a+c >= 10 || b+c >= 10 {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cases := [][3]int{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var a, b, c int
		fmt.Sscan(line, &a, &b, &c)
		cases = append(cases, [3]int{a, b, c})
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	exps := make([]string, len(cases))
	for i, tc := range cases {
		fmt.Fprintf(&sb, "%d %d %d\n", tc[0], tc[1], tc[2])
		exps[i] = expected(tc[0], tc[1], tc[2])
	}

	out, err := runBinary(bin, sb.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s", err, out)
		os.Exit(1)
	}
	scannerOut := bufio.NewScanner(strings.NewReader(out))
	got := []string{}
	for scannerOut.Scan() {
		line := strings.TrimSpace(scannerOut.Text())
		if line != "" {
			got = append(got, strings.ToUpper(line))
		}
	}
	if err := scannerOut.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "output scan error: %v\n", err)
		os.Exit(1)
	}
	if len(got) != len(exps) {
		fmt.Fprintf(os.Stderr, "wrong number of lines: got %d expected %d\ninput:\n%soutput:\n%s", len(got), len(exps), sb.String(), out)
		os.Exit(1)
	}
	for i := range exps {
		if got[i] != exps[i] {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), exps[i], got[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
