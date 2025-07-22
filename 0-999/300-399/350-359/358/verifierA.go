package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func expected(xs []int) string {
	type seg struct{ a, b int }
	segs := make([]seg, 0, len(xs)-1)
	for i := 0; i+1 < len(xs); i++ {
		a, b := xs[i], xs[i+1]
		if a > b {
			a, b = b, a
		}
		segs = append(segs, seg{a, b})
	}
	for i := 0; i < len(segs); i++ {
		a1, b1 := segs[i].a, segs[i].b
		for j := i + 1; j < len(segs); j++ {
			a2, b2 := segs[j].a, segs[j].b
			if (a1 < a2 && a2 < b1 && b1 < b2) || (a2 < a1 && a1 < b2 && b2 < b1) {
				return "yes"
			}
		}
	}
	return "no"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		n := 0
		fmt.Sscan(fields[0], &n)
		if len(fields) != n+1 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		xs := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(fields[i+1], &xs[i])
		}
		want := expected(xs)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range xs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
