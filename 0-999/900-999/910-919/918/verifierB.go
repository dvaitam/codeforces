package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCase(bin string, n, m int, servers [][2]string, cmds [][2]string) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, m))
	ipToName := make(map[string]string, n)
	for _, pair := range servers {
		input.WriteString(fmt.Sprintf("%s %s\n", pair[0], pair[1]))
		ipToName[pair[1]] = pair[0]
	}
	for _, pair := range cmds {
		input.WriteString(fmt.Sprintf("%s %s\n", pair[0], pair[1]))
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}

	outLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(outLines) != m {
		return fmt.Errorf("expected %d lines of output got %d", m, len(outLines))
	}
	for i, pair := range cmds {
		ipWithSemi := pair[1]
		ip := strings.TrimSuffix(ipWithSemi, ";")
		expect := fmt.Sprintf("%s %s #%s", pair[0], ipWithSemi, ipToName[ip])
		got := strings.TrimSpace(outLines[i])
		if got != expect {
			return fmt.Errorf("line %d expected %q got %q", i+1, expect, got)
		}
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
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		fmt.Println("invalid test count")
		os.Exit(1)
	}
	for caseNum := 0; caseNum < t; caseNum++ {
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		servers := make([][2]string, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(in, &servers[i][0], &servers[i][1]); err != nil {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
		}
		cmds := make([][2]string, m)
		for j := 0; j < m; j++ {
			if _, err := fmt.Fscan(in, &cmds[j][0], &cmds[j][1]); err != nil {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
		}
		if err := runCase(bin, n, m, servers, cmds); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
