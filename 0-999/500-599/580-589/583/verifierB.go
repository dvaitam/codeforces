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

func expected(a []int) string {
	n := len(a)
	visited := make([]bool, n)
	pieces := 0
	count := 0
	pos := 0
	dir := 1
	changes := 0
	for count < n {
		if !visited[pos] && pieces >= a[pos] {
			visited[pos] = true
			pieces++
			count++
			if count == n {
				break
			}
		}
		next := pos + dir
		if next < 0 || next >= n {
			dir = -dir
			changes++
			next = pos + dir
		}
		pos = next
	}
	return strconv.Itoa(changes)
}

func loadTests() ([]string, []string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	f, err := os.Open(filepath.Join(dir, "testcasesB.txt"))
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
		arr := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			scan.Scan()
			arr[i], _ = strconv.Atoi(scan.Text())
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		inputs[caseNum] = sb.String()
		expects[caseNum] = expected(arr)
	}
	return inputs, expects, nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
