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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refB.bin"
	if err := exec.Command("go", "build", "-o", ref, "821B.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("failed to open testcases:", err)
		os.Exit(1)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(scan.Text())
	for tc := 0; tc < T; tc++ {
		if !scan.Scan() {
			fmt.Printf("bad test %d\n", tc+1)
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Printf("bad test %d\n", tc+1)
			os.Exit(1)
		}
		b, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d %d\n", m, b)
		want, err := run(ref, []byte(input))
		if err != nil {
			fmt.Println("reference runtime error:", err)
			os.Exit(1)
		}
		got, err := run(cand, []byte(input))
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tc+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", tc+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", T)
}
