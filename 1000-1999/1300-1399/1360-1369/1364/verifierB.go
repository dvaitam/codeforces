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

func run(bin string, input string) (string, error) {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	ref := "./refB.bin"
	if err := exec.Command("go", "build", "-o", ref, "1364B.go").Run(); err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		fmt.Println("empty testcasesB.txt")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	for i := 1; i <= t; i++ {
		if !scanner.Scan() {
			fmt.Printf("missing test case %d\n", i)
			os.Exit(1)
		}
                line := strings.TrimSpace(scanner.Text())
                fields := strings.Fields(line)
                if len(fields) == 0 {
                        fmt.Printf("empty test case %d\n", i)
                        os.Exit(1)
                }
                n := fields[0]
                arr := strings.Join(fields[1:], " ")
                input := fmt.Sprintf("1\n%s\n%s\n", n, arr)
		want, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput: %s\nexpected: %s\ngot: %s\n", i, line, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
