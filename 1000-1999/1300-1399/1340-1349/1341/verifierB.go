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

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func expected(n, k int, arr []int) (int, int) {
	isPeak := make([]int, n)
	for i := 1; i < n-1; i++ {
		if arr[i] > arr[i-1] && arr[i] > arr[i+1] {
			isPeak[i] = 1
		}
	}
	pref := make([]int, n)
	for i := 1; i < n; i++ {
		pref[i] = pref[i-1] + isPeak[i]
	}
	bestPeaks := -1
	bestL := 0
	for l := 0; l+k-1 < n; l++ {
		peaks := pref[l+k-2] - pref[l]
		if peaks > bestPeaks {
			bestPeaks = peaks
			bestL = l
		}
	}
	return bestPeaks + 1, bestL + 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
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
		if len(fields) < 3 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+n {
			fmt.Printf("test %d wrong length\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[2+i])
			arr[i] = v
		}
		exp1, exp2 := expected(n, k, arr)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(arr[i]))
		}
		input.WriteByte('\n')

		cmd := exec.Command(cand)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		parts := strings.Fields(strings.TrimSpace(out.String()))
		if len(parts) != 2 {
			fmt.Printf("test %d: expected two numbers got %q\n", idx, out.String())
			os.Exit(1)
		}
		got1, _ := strconv.Atoi(parts[0])
		got2, _ := strconv.Atoi(parts[1])
		if got1 != exp1 || got2 != exp2 {
			fmt.Printf("test %d failed: expected %d %d got %d %d\n", idx, exp1, exp2, got1, got2)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
