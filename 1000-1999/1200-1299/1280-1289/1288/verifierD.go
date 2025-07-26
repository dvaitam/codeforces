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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candSrc := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "1288D.go")

	cand, err := buildBinary(candSrc, "candD.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := buildBinary(refSrc, "refD.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	f, err := os.Open(filepath.Join(dir, "testcasesD.txt"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)

	caseNum := 0
	for {
		if !scan.Scan() {
			break
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		vals := make([][]int, n)
		for i := 0; i < n; i++ {
			vals[i] = make([]int, m)
			for j := 0; j < m; j++ {
				if !scan.Scan() {
					fmt.Fprintln(os.Stderr, "bad file")
					os.Exit(1)
				}
				v, _ := strconv.Atoi(scan.Text())
				vals[i][j] = v
			}
		}
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					b.WriteByte(' ')
				}
				b.WriteString(strconv.Itoa(vals[i][j]))
			}
			b.WriteByte('\n')
		}
		input := b.String()
		caseNum++
		exp, err1 := runBin(ref, input)
		got, err2 := runBin(cand, input)
		if err2 != nil {
			fmt.Printf("Test %d: runtime error: %v\n", caseNum, err2)
			os.Exit(1)
		}
		if err1 != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", caseNum, err1)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\ninput:\n%sExpected: %s\nGot: %s\n", caseNum, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", caseNum)
}
