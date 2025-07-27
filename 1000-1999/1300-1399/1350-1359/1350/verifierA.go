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

type Test struct {
	n int
	k int
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildRef() (string, error) {
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1350A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v: %s", err, string(out))
	}
	return ref, nil
}

func loadTests() ([]Test, error) {
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("invalid test file")
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid test count")
	}
	tests := make([]Test, 0, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("unexpected EOF")
		}
		n, _ := strconv.Atoi(sc.Text())
		if !sc.Scan() {
			return nil, fmt.Errorf("unexpected EOF")
		}
		k, _ := strconv.Atoi(sc.Text())
		tests = append(tests, Test{n: n, k: k})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests, err := loadTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "load tests:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)
		want, err := runProg(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("case %d failed\ninput: %sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
