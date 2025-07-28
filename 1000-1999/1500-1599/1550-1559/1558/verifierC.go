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

func runCase(bin string, arr []int) error {
	n := len(arr)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	tokens := strings.Fields(out.String())
	impossible := false
	for i := 0; i < n; i++ {
		if (arr[i]-(i+1))%2 != 0 {
			impossible = true
			break
		}
	}
	if impossible {
		if len(tokens) != 1 || tokens[0] != "-1" {
			return fmt.Errorf("expected -1")
		}
		return nil
	}
	if len(tokens) < 1 {
		return fmt.Errorf("no output")
	}
	m, err := strconv.Atoi(tokens[0])
	if err != nil {
		return fmt.Errorf("invalid count")
	}
	if m > 5*n/2 {
		return fmt.Errorf("too many operations")
	}
	if len(tokens)-1 != m {
		return fmt.Errorf("expected %d ops got %d", m, len(tokens)-1)
	}
	ops := make([]int, m)
	for i := 0; i < m; i++ {
		v, err := strconv.Atoi(tokens[1+i])
		if err != nil {
			return fmt.Errorf("invalid op %q", tokens[1+i])
		}
		if v < 1 || v > n || v%2 == 0 {
			return fmt.Errorf("invalid op value %d", v)
		}
		ops[i] = v
	}
	arr2 := append([]int(nil), arr...)
	for _, p := range ops {
		for i, j := 0, p-1; i < j; i, j = i+1, j-1 {
			arr2[i], arr2[j] = arr2[j], arr2[i]
		}
	}
	if !isSorted(arr2) {
		return fmt.Errorf("array not sorted after operations")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("could not open testcasesC.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			arr[i], _ = strconv.Atoi(scanner.Text())
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
