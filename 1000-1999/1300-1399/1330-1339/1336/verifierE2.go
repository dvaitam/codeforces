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

const MOD int64 = 998244353

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

func expected(arr []uint64, m int) []int64 {
	n := len(arr)
	res := make([]int64, m+1)
	for mask := 0; mask < (1 << uint(n)); mask++ {
		x := uint64(0)
		for i := 0; i < n; i++ {
			if mask>>uint(i)&1 == 1 {
				x ^= arr[i]
			}
		}
		w := bits.OnesCount64(x)
		if w <= m {
			res[w]++
		}
	}
	for i := range res {
		res[i] %= MOD
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE2.txt")
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
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "case %d invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+n {
			fmt.Fprintf(os.Stderr, "case %d invalid number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]uint64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[2+i])
			arr[i] = uint64(v)
		}
		expect := expected(arr, m)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			input.WriteString(fmt.Sprintf("%d ", arr[i]))
		}
		input.WriteString("\n")
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		outFields := strings.Fields(got)
		if len(outFields) != m+1 {
			fmt.Fprintf(os.Stderr, "case %d wrong number of outputs\n", idx)
			os.Exit(1)
		}
		for i := 0; i <= m; i++ {
			v, err := strconv.ParseInt(outFields[i], 10, 64)
			if err != nil || v%MOD != expect[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at index %d: expected %d got %s\n", idx, i, expect[i], outFields[i])
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
