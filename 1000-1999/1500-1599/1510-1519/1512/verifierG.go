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

const N = 1000000

var ans []int32

func precompute() {
	ans = make([]int32, N+1)
	lp := make([]int32, N+1)
	pw := make([]int32, N+1)
	sumP := make([]int32, N+1)
	sigma := make([]int32, N+1)
	primes := make([]int32, 0)

	pw[1] = 1
	sumP[1] = 1
	sigma[1] = 1
	ans[1] = 1

	for i := 2; i <= N; i++ {
		if lp[i] == 0 {
			lp[i] = int32(i)
			primes = append(primes, int32(i))
			pw[i] = int32(i)
			sumP[i] = 1 + int32(i)
			sigma[i] = sumP[i]
		}
		li := lp[i]
		for _, p32 := range primes {
			p := int(p32)
			if p > int(li) || i*p > N {
				break
			}
			idx := i * p
			lp[idx] = p32
			if p32 == li {
				pw[idx] = pw[i] * p32
				sumP[idx] = sumP[i] + pw[idx]
				sigma[idx] = sigma[i/int(pw[i])] * sumP[idx]
			} else {
				pw[idx] = p32
				sumP[idx] = 1 + p32
				sigma[idx] = sigma[i] * sumP[idx]
			}
		}
		s := sigma[i]
		if int(s) <= N && ans[s] == 0 {
			ans[s] = int32(i)
		}
	}
}

func parseTestcasesG(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cases []int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		v, _ := strconv.Atoi(line)
		cases = append(cases, v)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func solveG(c int) string {
	if c <= N && ans[c] != 0 {
		return strconv.Itoa(int(ans[c]))
	}
	return "-1"
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	precompute()
	bin := os.Args[1]
	cases, err := parseTestcasesG("testcasesG.txt")
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	expected := make([]string, len(cases))
	for i, c := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", c))
		expected[i] = solveG(c)
	}

	got, err := run(bin, sb.String())
	if err != nil {
		fmt.Printf("failed: %v\n", err)
		os.Exit(1)
	}
	outputs := strings.Fields(got)
	if len(outputs) != len(expected) {
		fmt.Printf("expected %d lines of output, got %d\n", len(expected), len(outputs))
		os.Exit(1)
	}
	for i, exp := range expected {
		if outputs[i] != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, exp, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
