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

func runExe(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "786D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func readTests() ([][]byte, error) {
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("empty test file")
	}
	t, _ := strconv.Atoi(scan.Text())
	tests := make([][]byte, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("bad test %d", i+1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		qStr := scan.Text()
		q, _ := strconv.Atoi(qStr)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for j := 0; j < n-1; j++ {
			scan.Scan()
			vStr := scan.Text()
			scan.Scan()
			uStr := scan.Text()
			scan.Scan()
			ch := scan.Text()
			sb.WriteString(fmt.Sprintf("%s %s %s\n", vStr, uStr, ch))
		}
		for j := 0; j < q; j++ {
			scan.Scan()
			xStr := scan.Text()
			scan.Scan()
			yStr := scan.Text()
			sb.WriteString(fmt.Sprintf("%s %s\n", xStr, yStr))
		}
		tests = append(tests, []byte(sb.String()))
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	tests, err := readTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, input := range tests {
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, string(input), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
