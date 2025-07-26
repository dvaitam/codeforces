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
	ref := "./refE.bin"
	if err := exec.Command("go", "build", "-o", ref, "1007E.go").Run(); err != nil {
		return "", err
	}
	return ref, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("could not open testcasesE.txt:", err)
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
		tt, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scanner.Text())
		stations := make([][3]int, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			a, _ := strconv.Atoi(scanner.Text())
			if !scanner.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			b, _ := strconv.Atoi(scanner.Text())
			if !scanner.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			c, _ := strconv.Atoi(scanner.Text())
			stations[i] = [3]int{a, b, c}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, tt, k))
		for _, s := range stations {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", s[0], s[1], s[2]))
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
