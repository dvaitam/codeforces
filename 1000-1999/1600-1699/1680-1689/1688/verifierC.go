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

func expected(list []string, final string) string {
	cnt := make([]int, 26)
	for _, s := range list {
		for i := 0; i < len(s); i++ {
			cnt[s[i]-'a'] ^= 1
		}
	}
	for i := 0; i < len(final); i++ {
		cnt[final[i]-'a'] ^= 1
	}
	for i := 0; i < 26; i++ {
		if cnt[i] == 1 {
			return string('a' + byte(i))
		}
	}
	return ""
}

func runCase(exe string, n int, list []string, final, exp string) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, s := range list {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	sb.WriteString(final)
	sb.WriteByte('\n')
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))

	for caseIdx := 0; caseIdx < t; caseIdx++ {
		for {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			if strings.TrimSpace(scan.Text()) != "" {
				break
			}
		}
		n, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
		list := make([]string, 2*n)
		for i := 0; i < 2*n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			list[i] = strings.TrimSpace(scan.Text())
		}
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		final := strings.TrimSpace(scan.Text())
		exp := expected(list, final)
		if err := runCase(exe, n, list, final, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		if caseIdx != t-1 {
			scan.Scan() // consume blank line
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
