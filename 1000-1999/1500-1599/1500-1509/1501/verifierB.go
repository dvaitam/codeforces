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

func expected(n int, a []int) string {
	ans := make([]int, n)
	cream := 0
	for i := n - 1; i >= 0; i-- {
		if a[i] > cream {
			cream = a[i]
		}
		if cream > 0 {
			ans[i] = 1
			cream--
		} else {
			ans[i] = 0
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ans[i]))
	}
	return sb.String()
}

func runCase(bin string, n int, a []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	input := sb.String()

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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	want := expected(n, a)
	if got != want {
		return fmt.Errorf("expected %s got %s\ninput:\n%s", want, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "bad test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "bad test file")
			os.Exit(1)
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) < 1 {
			fmt.Fprintln(os.Stderr, "bad test case")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields)-1 != n {
			fmt.Fprintln(os.Stderr, "wrong number of values")
			os.Exit(1)
		}
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j], _ = strconv.Atoi(fields[j+1])
		}
		if err := runCase(bin, n, a); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
