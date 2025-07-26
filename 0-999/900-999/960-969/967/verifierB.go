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

func expected(n, A, B int, sizes []int) string {
	total := 0
	for _, v := range sizes {
		total += v
	}
	others := append([]int(nil), sizes[1:]...)
	// sort descending
	for i := 0; i < len(others); i++ {
		for j := i + 1; j < len(others); j++ {
			if others[j] > others[i] {
				others[i], others[j] = others[j], others[i]
			}
		}
	}
	S := int64(total)
	target := int64(sizes[0]) * int64(A)
	if target >= int64(B)*S {
		return "0\n"
	}
	blocked := 0
	for _, v := range others {
		S -= int64(v)
		blocked++
		if target >= int64(B)*S {
			return fmt.Sprintf("%d\n", blocked)
		}
	}
	return fmt.Sprintf("%d\n", blocked)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		return
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		return
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			return
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		A, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		Bv, _ := strconv.Atoi(scan.Text())
		sizes := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			sizes[i] = v
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, A, Bv))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(sizes[i]))
		}
		sb.WriteByte('\n')
		in := sb.String()
		exp := expected(n, A, Bv, sizes)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, in)
			return
		}
	}
	fmt.Println("All tests passed")
}
