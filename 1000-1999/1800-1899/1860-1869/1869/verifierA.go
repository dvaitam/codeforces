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
	return out.String(), nil
}

func applyOps(arr []int, ops [][2]int) error {
	n := len(arr)
	for _, op := range ops {
		l := op[0] - 1
		r := op[1] - 1
		if l < 0 || r < l || r >= n {
			return fmt.Errorf("invalid range %d %d", op[0], op[1])
		}
		s := 0
		for i := l; i <= r; i++ {
			s ^= arr[i]
		}
		for i := l; i <= r; i++ {
			arr[i] = s
		}
	}
	for _, v := range arr {
		if v != 0 {
			return fmt.Errorf("array not zero after ops")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesA.txt:", err)
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
		if len(fields) < 1 {
			fmt.Fprintf(os.Stderr, "case %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != n+1 {
			fmt.Fprintf(os.Stderr, "case %d: invalid data\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[i+1])
			arr[i] = v
		}
		input := fmt.Sprintf("1\n%d\n%s\n", n, strings.Join(fields[1:], " "))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		sc := bufio.NewScanner(strings.NewReader(out))
		sc.Split(bufio.ScanWords)
		if !sc.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: no output\n", idx)
			os.Exit(1)
		}
		k, err := strconv.Atoi(sc.Text())
		if err != nil || k < 0 || k > 8 {
			fmt.Fprintf(os.Stderr, "case %d: invalid k\n", idx)
			os.Exit(1)
		}
		ops := make([][2]int, k)
		for i := 0; i < k; i++ {
			if !sc.Scan() {
				fmt.Fprintf(os.Stderr, "case %d: insufficient output\n", idx)
				os.Exit(1)
			}
			l, err := strconv.Atoi(sc.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: invalid l\n", idx)
				os.Exit(1)
			}
			if !sc.Scan() {
				fmt.Fprintf(os.Stderr, "case %d: insufficient output\n", idx)
				os.Exit(1)
			}
			r, err := strconv.Atoi(sc.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: invalid r\n", idx)
				os.Exit(1)
			}
			ops[i] = [2]int{l, r}
		}
		arrCopy := make([]int, n)
		copy(arrCopy, arr)
		if err := applyOps(arrCopy, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
