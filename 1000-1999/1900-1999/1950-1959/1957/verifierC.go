package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod int64 = 1000000007
const maxN = 300000

var f [maxN + 1]int64

const testcasesC = `4 2 2 2 1 4
1 1 1 1
2 2 2 2 1 1
2 1 2 2
4 0
4 2 2 4 1 1
1 1 1 1
1 0
6 3 1 6 5 2 4 3
1 0
3 2 3 1 2 2
2 0
4 2 3 4 1 2
4 2 2 4 1 3
2 0
2 1 2 2
4 2 4 2 1 3
2 1 1 2
1 0
6 0
2 1 1 1
2 1 2 1
6 3 5 3 4 6 2 1
6 4 5 6 4 4 2 2 3 1
3 2 1 1 2 2
3 0
3 0
2 0
4 1 4 2
4 1 3 1
5 0
3 2 3 2 1 1
1 0
2 1 1 1
5 2 2 2 4 4
3 2 2 2 3 1
1 1 1 1
4 3 4 1 3 3 2 2
1 0
1 0
1 0
4 1 4 1
4 3 4 1 2 2 3 3
2 1 2 2
4 0
5 2 5 5 2 3
4 2 4 3 2 1
1 0
5 3 1 2 5 4 3 3
3 0
6 4 5 4 1 1 3 6 2 2
6 4 3 2 5 4 1 1 6 6
2 1 2 2
5 1 2 1
5 3 3 3 1 2 5 4
2 2 1 1 2 2
5 4 2 2 3 3 4 4 1 5
1 0
5 4 3 4 2 2 1 1 5 5
6 3 5 3 1 4 2 6
6 1 3 3
4 0
2 0
3 1 2 2
6 4 6 6 2 3 4 5 1 1
1 1 1 1
3 0
4 0
1 0
1 0
2 1 2 1
1 1 1 1
2 1 2 1
3 2 1 3 2 2
4 0
4 2 4 1 2 3
4 1 1 2
5 3 5 2 1 3 4 4
2 0
2 2 1 1 2 2
5 3 4 2 5 5 1 3
3 2 1 3 2 2
1 1 1 1
4 1 3 2
6 3 4 3 5 6 2 1
3 2 3 2 1 1
6 0
4 3 1 4 3 3 2 2
3 0
3 2 3 1 2 2
3 1 1 3
3 1 3 3
6 4 6 3 2 2 4 5 1 1
5 1 3 3
3 2 2 2 3 3
2 0
2 1 1 2
2 1 1 2
6 4 6 3 2 1 5 5 4 4
1 0`

func init() {
	f[0] = 1
	if maxN >= 1 {
		f[1] = 1
	}
	for i := 2; i <= maxN; i++ {
		f[i] = (f[i-1] + 2*int64(i-1)%mod*f[i-2]) % mod
	}
}

func expected(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", fmt.Errorf("bad line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	if len(fields) != 2+2*k {
		return "", fmt.Errorf("expected %d numbers, got %d", 2+2*k, len(fields))
	}
	used := make([]bool, n+1)
	usedCount := 0
	for i := 0; i < k; i++ {
		r, err := strconv.Atoi(fields[2+2*i])
		if err != nil {
			return "", err
		}
		c, err := strconv.Atoi(fields[3+2*i])
		if err != nil {
			return "", err
		}
		if !used[r] {
			used[r] = true
			usedCount++
		}
		if !used[c] {
			used[c] = true
			usedCount++
		}
	}
	remaining := n - usedCount
	if remaining < 0 || remaining > maxN {
		return "", fmt.Errorf("remaining out of range")
	}
	return strconv.FormatInt(f[remaining], 10), nil
}

func buildInput(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", fmt.Errorf("bad line")
	}
	n := fields[0]
	k := fields[1]
	if len(fields) < 2 {
		return "", fmt.Errorf("bad line")
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%s %s\n", n, k))
	idx := 2
	for idx < len(fields) {
		sb.WriteString(fmt.Sprintf("%s %s\n", fields[idx], fields[idx+1]))
		idx += 2
	}
	return sb.String(), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	lines := strings.Split(testcasesC, "\n")
	idx := 0
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		idx++
		exp, err := expected(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", idx, err)
			os.Exit(1)
		}
		inp, err := buildInput(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d input build error: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := runCandidate(os.Args[1], inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
