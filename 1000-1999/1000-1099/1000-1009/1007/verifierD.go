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

func run(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("g++", "-std=c++17", "solD.cpp", "-O2", "-pipe", "-static", "-s", "-o", ref)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return ref, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scanner.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scanner.Text())
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if !scanner.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			u, _ := strconv.Atoi(scanner.Text())
			if !scanner.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scanner.Text())
			edges[i] = [2]int{u, v}
		}
		pairs := make([]int, 4*m)
		for i := 0; i < 4*m; i++ {
			if !scanner.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			val, _ := strconv.Atoi(scanner.Text())
			pairs[i] = val
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i := 0; i < m; i++ {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", pairs[4*i], pairs[4*i+1], pairs[4*i+2], pairs[4*i+3]))
		}
		input := sb.String()
		want, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("case %d failed\nexpected:%s\ngot:%s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
