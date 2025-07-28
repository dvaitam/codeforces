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

func expected(n int, a, b, tm []int) int {
	prevDepart := 0
	prevB := 0
	for i := 0; i < n; i++ {
		travel := a[i] - prevB + tm[i]
		arrival := prevDepart + travel
		if i == n-1 {
			return arrival
		}
		stay := (b[i] - a[i] + 1) / 2
		depart := arrival + stay
		if depart < b[i] {
			depart = b[i]
		}
		prevDepart = depart
		prevB = b[i]
	}
	return 0
}

func runCase(bin string, n int, a, b, tm []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", a[i], b[i]))
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(tm[i]))
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
	want := strconv.Itoa(expected(n, a, b, tm))
	if got != want {
		return fmt.Errorf("expected %s got %s\ninput:\n%s", want, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "bad test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "missing n")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Fprintln(os.Stderr, "bad test file")
				os.Exit(1)
			}
			a[j], _ = strconv.Atoi(scan.Text())
			if !scan.Scan() {
				fmt.Fprintln(os.Stderr, "bad test file")
				os.Exit(1)
			}
			b[j], _ = strconv.Atoi(scan.Text())
		}
		tm := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Fprintln(os.Stderr, "bad test file")
				os.Exit(1)
			}
			tm[j], _ = strconv.Atoi(scan.Text())
		}
		if err := runCase(bin, n, a, b, tm); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
