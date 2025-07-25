package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, pairs [][2]int) string {
	horiz := make([]bool, n+1)
	vert := make([]bool, n+1)
	var days []int
	for i, p := range pairs {
		h, v := p[0], p[1]
		if !horiz[h] && !vert[v] {
			horiz[h] = true
			vert[v] = true
			days = append(days, i+1)
		}
	}
	var sb strings.Builder
	for i, d := range days {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(d))
	}
	return sb.String()
}

func loadTests() ([]string, []string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	f, err := os.Open(filepath.Join(dir, "testcasesA.txt"))
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, nil, fmt.Errorf("empty test file")
	}
	t, _ := strconv.Atoi(scan.Text())
	inputs := make([]string, t)
	expects := make([]string, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scan.Scan() {
			return nil, nil, fmt.Errorf("missing n for case %d", caseNum+1)
		}
		n, _ := strconv.Atoi(scan.Text())
		pairs := make([][2]int, n*n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n*n; i++ {
			scan.Scan()
			h, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			pairs[i] = [2]int{h, v}
			sb.WriteString(fmt.Sprintf("%d %d\n", h, v))
		}
		inputs[caseNum] = sb.String()
		expects[caseNum] = expected(n, pairs)
	}
	return inputs, expects, nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputs, expects, err := loadTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for i := range inputs {
		got, err := runCandidate(bin, inputs[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expects[i] {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, expects[i], got, inputs[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
