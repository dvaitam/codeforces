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

func solve(arr []int) int {
	comp := make(map[int]int, len(arr))
	vals := make([]int, 0, len(arr))
	for _, v := range arr {
		if _, ok := comp[v]; !ok {
			comp[v] = len(vals)
			vals = append(vals, v)
		}
	}
	m := len(vals)
	pos := make([][]int, m)
	for i, v := range arr {
		id := comp[v]
		pos[id] = append(pos[id], i)
	}
	ans := 0
	for _, lst := range pos {
		if len(lst) > ans {
			ans = len(lst)
		}
	}
	for i := 0; i < m; i++ {
		pi0 := pos[i]
		for j := 0; j < m; j++ {
			if j == i {
				continue
			}
			pj0 := pos[j]
			pi, pj := 0, 0
			lastPos := -1
			lastIsI := true
			cnt := 0
			for {
				if lastIsI {
					for pi < len(pi0) && pi0[pi] <= lastPos {
						pi++
					}
					if pi >= len(pi0) {
						break
					}
					lastPos = pi0[pi]
					pi++
					cnt++
					lastIsI = false
				} else {
					for pj < len(pj0) && pj0[pj] <= lastPos {
						pj++
					}
					if pj >= len(pj0) {
						break
					}
					lastPos = pj0[pj]
					pj++
					cnt++
					lastIsI = true
				}
			}
			if cnt > ans {
				ans = cnt
			}
		}
	}
	return ans
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
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		if len(fields)-1 != n {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(fields[i+1])
		}
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(fields[1:], " "))
		want := solve(arr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		gotVal, _ := strconv.Atoi(strings.TrimSpace(got))
		if gotVal != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, gotVal)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
