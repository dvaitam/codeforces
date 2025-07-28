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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expected(arr []int) string {
	n := len(arr)
	sum := 0
	for _, v := range arr {
		sum += v
	}
	half := sum / 2
	dp := make([]bool, half+1)
	dp[0] = true
	for _, x := range arr {
		for s := half; s >= x; s-- {
			if dp[s-x] {
				dp[s] = true
			}
		}
	}
	if sum%2 == 1 || !dp[half] {
		return "0"
	}
	g := arr[0]
	for i := 1; i < n; i++ {
		g = gcd(g, arr[i])
	}
	best := 31
	idx := -1
	for i, v := range arr {
		x := v / g
		cnt := 0
		for x%2 == 0 {
			x /= 2
			cnt++
		}
		if cnt < best {
			best = cnt
			idx = i
		}
	}
	return fmt.Sprintf("1\n%d", idx+1)
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Fprintf(os.Stderr, "bad testcase line %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+n {
			fmt.Fprintf(os.Stderr, "bad testcase line %d\n", idx+1)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(fields[1+i])
		}
		idx++

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()

		exp := expected(arr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx, exp, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
