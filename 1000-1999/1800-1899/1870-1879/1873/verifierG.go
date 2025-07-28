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

func solve(s string) int {
	bpos := []int{}
	for i := 0; i < len(s); i++ {
		if s[i] == 'B' {
			bpos = append(bpos, i)
		}
	}
	if len(bpos) == 0 {
		return 0
	}
	k := len(bpos)
	segments := make([]int, k+1)
	segments[0] = bpos[0]
	for i := 0; i < k-1; i++ {
		segments[i+1] = bpos[i+1] - bpos[i] - 1
	}
	segments[k] = len(s) - 1 - bpos[k-1]
	edges := make([]int, 0, 2*k)
	for i := 0; i < k; i++ {
		edges = append(edges, segments[i])
		edges = append(edges, segments[i+1])
	}
	if len(edges) == 0 {
		return 0
	}
	dpPrev2 := 0
	dpPrev1 := edges[0]
	if dpPrev1 < 0 {
		dpPrev1 = 0
	}
	for i := 2; i <= len(edges); i++ {
		val := dpPrev1
		if dpPrev2+edges[i-1] > val {
			val = dpPrev2 + edges[i-1]
		}
		dpPrev2, dpPrev1 = dpPrev1, val
	}
	return dpPrev1
}

func expected(s string) string {
	return fmt.Sprintf("%d\n", solve(s))
}

func runCase(exe, input, expect string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesG.txt")
	if err != nil {
		fmt.Println("could not read testcasesG.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		s := scan.Text()
		input := fmt.Sprintf("1\n%s\n", s)
		exp := expected(s)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
