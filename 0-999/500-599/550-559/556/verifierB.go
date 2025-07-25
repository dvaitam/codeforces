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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, a []int) string {
	for t := 0; t < n; t++ {
		ok := true
		for i := 0; i < n; i++ {
			val := a[i]
			if i%2 == 0 {
				val = (val + t) % n
			} else {
				val = ((val-t)%n + n) % n
			}
			if val != i {
				ok = false
				break
			}
		}
		if ok {
			return "Yes"
		}
	}
	return "No"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("failed to open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad test format on line %d\n", idx+1)
			os.Exit(1)
		}
		idx++
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Printf("bad n on line %d: %v\n", idx, err)
			os.Exit(1)
		}
		if len(fields)-1 != n {
			fmt.Printf("line %d: expected %d numbers, got %d\n", idx, n, len(fields)-1)
			os.Exit(1)
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[i+1])
			if err != nil {
				fmt.Printf("bad value on line %d: %v\n", idx, err)
				os.Exit(1)
			}
			a[i] = v
		}
		var sb strings.Builder
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		input := fmt.Sprintf("%d\n%s\n", n, sb.String())
		exp := expected(n, a)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != exp {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", idx, exp, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
