package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func segCost(sa, ca, sb, cb int64) int64 {
	dr := sb - sa
	dc := cb - ca
	if dr < 0 || dc < 0 || dc > dr {
		return 0
	}
	if (sa+ca)%2 == 0 {
		if dc == dr {
			return dr
		}
		return (dr - dc) / 2
	}
	return (dr - dc + 1) / 2
}

func solveCase(r, c []int64) int64 {
	n := len(r)
	pts := make([][2]int64, n)
	for i := 0; i < n; i++ {
		pts[i] = [2]int64{r[i], c[i]}
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i][0] == pts[j][0] {
			return pts[i][1] < pts[j][1]
		}
		return pts[i][0] < pts[j][0]
	})
	var ans int64
	curR, curC := int64(1), int64(1)
	for _, p := range pts {
		ans += segCost(curR, curC, p[0], p[1])
		curR, curC = p[0], p[1]
	}
	return ans
}

func main() {
	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		arg = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := arg
	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesF.txt: %v\n", err)
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
		if len(parts) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts)-1 != 2*n {
			fmt.Printf("test %d: wrong number count\n", idx)
			os.Exit(1)
		}
		r := make([]int64, n)
		c := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(parts[1+i], 10, 64)
			r[i] = v
		}
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(parts[1+n+i], 10, 64)
			c[i] = v
		}
		expect := fmt.Sprintf("%d", solveCase(r, c))
		input := fmt.Sprintf("1\n%d\n", n)
		for i, v := range r {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		for i, v := range c {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
