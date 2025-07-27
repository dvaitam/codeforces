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

func solve(n int, a []int64, b []int) (int64, []int) {
	f := make([]int64, n)
	copy(f, a)
	indeg := make([]int, n)
	for i := 0; i < n; i++ {
		if b[i] >= 0 {
			indeg[b[i]]++
		}
	}
	graph := make([][]int, n)
	deg := make([]int, n)
	q := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	for head := 0; head < len(q); head++ {
		x := q[head]
		if b[x] >= 0 {
			p := b[x]
			if f[x] < 0 {
				graph[p] = append(graph[p], x)
				deg[x]++
			} else {
				f[p] += f[x]
				graph[x] = append(graph[x], p)
				deg[p]++
			}
			indeg[p]--
			if indeg[p] == 0 {
				q = append(q, p)
			}
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		ans += f[i]
	}
	q2 := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if deg[i] == 0 {
			q2 = append(q2, i)
		}
	}
	order := make([]int, 0, n)
	for head := 0; head < len(q2); head++ {
		x := q2[head]
		order = append(order, x)
		for _, y := range graph[x] {
			deg[y]--
			if deg[y] == 0 {
				q2 = append(q2, y)
			}
		}
	}
	for i := range order {
		order[i]++
	}
	return ans, order
}

func runCase(bin string, n int, a []int64, b []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(a[i], 10))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if b[i] < 0 {
			sb.WriteString("-1")
		} else {
			sb.WriteString(strconv.Itoa(b[i] + 1))
		}
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	outFields := strings.Fields(strings.TrimSpace(out.String()))
	if len(outFields) < 1+n {
		return fmt.Errorf("not enough output")
	}
	ansGot, err := strconv.ParseInt(outFields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("bad sum output")
	}
	orderGot := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(outFields[1+i])
		if err != nil {
			return fmt.Errorf("bad order value")
		}
		orderGot[i] = v
	}
	ansExp, orderExp := solve(n, a, b)
	if ansGot != ansExp {
		return fmt.Errorf("expected sum %d got %d", ansExp, ansGot)
	}
	for i := 0; i < n; i++ {
		if orderGot[i] != orderExp[i] {
			return fmt.Errorf("expected order %v got %v", orderExp, orderGot)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("bad test line %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+2*n {
			fmt.Printf("bad test line %d field count\n", idx+1)
			os.Exit(1)
		}
		a := make([]int64, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[1+i], 10, 64)
			a[i] = v
		}
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[1+n+i])
			if v == -1 {
				b[i] = -1
			} else {
				b[i] = v - 1
			}
		}
		idx++
		if err := runCase(bin, n, a, b); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
