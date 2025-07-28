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

func solveD(n int, b []int) (int, []int) {
	v := make([][]int, n+2)
	k := 0
	for i := 1; i <= n; i++ {
		if b[i] > i {
			k++
		}
		if b[i] >= 0 && b[i] <= n+1 {
			v[b[i]] = append(v[b[i]], i)
		}
	}
	a := make([]int, 0, n)
	cur := 0
	if len(v[n+1]) > 0 {
		cur = n + 1
	}
	cnt := 0
	for cnt < n {
		cnt += len(v[cur])
		last := len(v[cur]) - 1
		good := -1
		for j := 0; j <= last; j++ {
			nxt := v[cur][j]
			if nxt >= 0 && nxt < len(v) && len(v[nxt]) > 0 {
				good = j
			}
		}
		if good != -1 && good != last {
			v[cur][good], v[cur][last] = v[cur][last], v[cur][good]
		}
		a = append(a, v[cur]...)
		if len(v[cur]) > 0 {
			cur = v[cur][len(v[cur])-1]
		} else {
			break
		}
	}
	return k, a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
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
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != n+1 {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		bvals := make([]int, n+1)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[1+i])
			bvals[i+1] = v
		}
		expectK, expectA := solveD(n, bvals)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 1; i <= n; i++ {
			if i > 1 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", bvals[i]))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if len(lines) < 2 {
			fmt.Printf("test %d: output has fewer lines\n", idx)
			os.Exit(1)
		}
		gotK, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
		gotParts := strings.Fields(lines[1])
		if len(gotParts) != n {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, n, len(gotParts))
			os.Exit(1)
		}
		gotA := make([]int, n)
		for i := 0; i < n; i++ {
			val, _ := strconv.Atoi(gotParts[i])
			gotA[i] = val
		}
		if gotK != expectK {
			fmt.Printf("test %d: expected k=%d got %d\n", idx, expectK, gotK)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			if gotA[i] != expectA[i] {
				fmt.Printf("test %d: permutation mismatch\nexpected %v\ngot %v\n", idx, expectA, gotA)
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
