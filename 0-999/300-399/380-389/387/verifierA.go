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

func expectedBedtime(s, t string) string {
	shh, _ := strconv.Atoi(s[:2])
	smm, _ := strconv.Atoi(s[3:])
	thh, _ := strconv.Atoi(t[:2])
	tmm, _ := strconv.Atoi(t[3:])
	curr := shh*60 + smm
	dur := thh*60 + tmm
	p := curr - dur
	p = ((p % (24 * 60)) + (24 * 60)) % (24 * 60)
	return fmt.Sprintf("%02d:%02d", p/60, p%60)
}

func runCase(bin, s, t string) error {
	input := fmt.Sprintf("%s\n%s\n", s, t)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedBedtime(s, t)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "bad test file")
		os.Exit(1)
	}
	tcases, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < tcases; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "test %d: missing s\n", i+1)
			os.Exit(1)
		}
		s := scanner.Text()
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "test %d: missing t\n", i+1)
			os.Exit(1)
		}
		tt := scanner.Text()
		if err := runCase(bin, s, tt); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra data in test file")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", tcases)
}
