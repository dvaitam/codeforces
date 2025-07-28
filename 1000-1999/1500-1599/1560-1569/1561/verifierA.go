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

func isSorted(a []int) bool {
	for i := 0; i+1 < len(a); i++ {
		if a[i] > a[i+1] {
			return false
		}
	}
	return true
}

func expectedA(arr []int) int {
	ans := 0
	n := len(arr)
	a := make([]int, n)
	copy(a, arr)
	for {
		if isSorted(a) {
			break
		}
		if ans%2 == 0 {
			for j := 0; j+1 < n; j += 2 {
				if a[j] > a[j+1] {
					a[j], a[j+1] = a[j+1], a[j]
				}
			}
		} else {
			for j := 1; j+1 < n; j += 2 {
				if a[j] > a[j+1] {
					a[j], a[j+1] = a[j+1], a[j]
				}
			}
		}
		ans++
	}
	return ans
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		if len(fields)-1 != n {
			fmt.Fprintf(os.Stderr, "test %d bad length", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[i+1])
			arr[i] = v
		}
		expected := expectedA(arr)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", expected) {
			fmt.Printf("test %d failed: expected %d got %s\n", idx, expected, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
