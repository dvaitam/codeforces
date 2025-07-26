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

func buildBinary(src, tag string) (string, error) {
	if strings.HasSuffix(src, ".go") {
		out := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", out, src)
		if outb, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", src, err, string(outb))
		}
		return out, nil
	}
	return src, nil
}

func runBin(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	candSrc := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "1288C.go")

	cand, err := buildBinary(candSrc, "candC.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := buildBinary(refSrc, "refC.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	f, err := os.Open(filepath.Join(dir, "testcasesC.txt"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot open testcasesC.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	cases := make([][2]int, 0, 100)
	for scan.Scan() {
		a, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		b, _ := strconv.Atoi(scan.Text())
		cases = append(cases, [2]int{a, b})
	}

	for i, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc[0], tc[1])
		exp, err1 := runBin(ref, input)
		got, err2 := runBin(cand, input)
		if err2 != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err2)
			os.Exit(1)
		}
		if err1 != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err1)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed: input %d %d\nexpected %s got %s\n", i+1, tc[0], tc[1], exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
