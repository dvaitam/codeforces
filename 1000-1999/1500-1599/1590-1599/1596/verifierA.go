package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	tcPath := filepath.Join(dir, "testcasesA.txt")
	f, err := os.Open(tcPath)
	if err != nil {
		panic(err)
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
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[1+i])
			nums[i] = v
		}
		// compute expected sum
		sum := 0
		for _, v := range nums {
			sum += v
		}
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(parts[1:1+n], " "))
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", sum) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, sum, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
