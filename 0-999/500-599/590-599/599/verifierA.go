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

func expected(d1, d2, d3 int) string {
	option1 := d1 + d2 + d3
	option2 := 2 * (d1 + d2)
	option3 := 2 * (d1 + d3)
	option4 := 2 * (d2 + d3)
	ans := option1
	if option2 < ans {
		ans = option2
	}
	if option3 < ans {
		ans = option3
	}
	if option4 < ans {
		ans = option4
	}
	return fmt.Sprintf("%d", ans)
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		d1, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		d2, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		d3, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d %d %d\n", d1, d2, d3)
		exp := expected(d1, d2, d3) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
