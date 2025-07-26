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

func runExe(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1197E.go")
	bin := filepath.Join(dir, "refE.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return bin, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("failed to open testcasesE.txt:", err)
		return
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("bad test file")
		return
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			return
		}
		n, _ := strconv.Atoi(scan.Text())
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			scan.Scan()
			outv := scan.Text()
			scan.Scan()
			inv := scan.Text()
			fmt.Fprintf(&input, "%s %s\n", outv, inv)
		}
		exp, err := runExe(ref, input.String())
		if err != nil {
			fmt.Printf("reference runtime error on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input.String())
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("case %d failed\nInput:\n%sExpected:%s\nGot:%s\n", caseNum, input.String(), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
