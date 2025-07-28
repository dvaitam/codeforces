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

func runCandidate(bin, input string) (string, error) {
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

func isGood(arr []int) bool {
	n := len(arr)
	for l := 0; l < n; l++ {
		freq := make(map[int]int)
		for r := l; r < n; r++ {
			freq[arr[r]]++
			unique := false
			for _, c := range freq {
				if c == 1 {
					unique = true
					break
				}
			}
			if !unique {
				return false
			}
		}
	}
	return true
}

func expected(arr []int) int {
	n := len(arr)
	for k := 0; k <= n; k++ {
		for mask := 0; mask < (1 << n); mask++ {
			if bits.OnesCount(uint(mask)) != k {
				continue
			}
			b := append([]int(nil), arr...)
			val := n + 1
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					b[i] = val
					val++
				}
			}
			if isGood(b) {
				return k
			}
		}
	}
	return n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("failed to open testcasesE.txt:", err)
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
		n, _ := strconv.Atoi(fields[0])
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[1+i])
			arr[i] = v
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		want := strconv.Itoa(expected(arr))
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", idx, want, strings.TrimSpace(got), sb.String())
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
