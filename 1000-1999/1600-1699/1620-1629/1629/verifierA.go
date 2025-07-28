package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type pair struct{ req, gain int }

func expected(n, k int, a, b []int) string {
	arr := make([]pair, n)
	for i := 0; i < n; i++ {
		arr[i] = pair{a[i], b[i]}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].req < arr[j].req })
	ram := k
	for _, p := range arr {
		if ram < p.req {
			break
		}
		ram += p.gain
	}
	return strconv.Itoa(ram)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	f, err := os.Open(filepath.Join(filepath.Dir(file), "testcasesA.txt"))
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
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "bad testcase line %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+2*n {
			fmt.Fprintf(os.Stderr, "bad testcase line %d\n", idx+1)
			os.Exit(1)
		}
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			a[i], _ = strconv.Atoi(fields[2+i])
		}
		for i := 0; i < n; i++ {
			b[i], _ = strconv.Atoi(fields[2+n+i])
		}
		idx++
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(n, k, a, b)
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
