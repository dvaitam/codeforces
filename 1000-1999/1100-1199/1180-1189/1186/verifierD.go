package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func validate(a []float64, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	var b []int
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		if t == "" {
			continue
		}
		v, err := strconv.Atoi(t)
		if err != nil {
			return fmt.Errorf("invalid integer %q", t)
		}
		b = append(b, v)
	}
	if len(b) != len(a) {
		return fmt.Errorf("expected %d numbers got %d", len(a), len(b))
	}
	sum := 0
	for i, v := range b {
		fl := math.Floor(a[i])
		ce := math.Ceil(a[i])
		if float64(v) != fl && float64(v) != ce {
			return fmt.Errorf("value %d not floor/ceil of %.5f", v, a[i])
		}
		sum += v
	}
	if sum != 0 {
		return fmt.Errorf("sum is %d, expected 0", sum)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != n+1 {
			fmt.Printf("bad testcase %d\n", idx)
			os.Exit(1)
		}
		a := make([]float64, n)
		for i := 0; i < n; i++ {
			v, err := parseFloat(fields[i+1])
			if err != nil {
				fmt.Printf("bad number on test %d: %v\n", idx, err)
				os.Exit(1)
			}
			a[i] = v
		}
		// build input string
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%.5f\n", a[i])
		}
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if err := validate(a, out); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
