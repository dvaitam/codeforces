package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func expected(arr []int64) int {
	n := len(arr)
	best := 0
	for mask := 1; mask < 1<<n; mask++ {
		l := int64(1)
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				l = lcm(l, arr[i])
			}
		}
		found := false
		for _, v := range arr {
			if int64(v) == l {
				found = true
				break
			}
		}
		if !found {
			if c := bits.OnesCount(uint(mask)); c > best {
				best = c
			}
		}
	}
	// empty subsequence gives length 0 which is always allowed
	return best
}

func parseCase(line string) []int64 {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	n, _ := strconv.Atoi(fields[0])
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		v, _ := strconv.Atoi(fields[i+1])
		arr[i] = int64(v)
	}
	return arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesC.txt")
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
		arr := parseCase(line)
		if arr == nil {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		input := line + "\n"
		want := expected(arr)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got, _ := strconv.Atoi(strings.TrimSpace(gotStr))
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
