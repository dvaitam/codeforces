package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesF = `4 3 1 2 1 1 3 1 2 4 1
2 1 1 2 1
4 5 1 2 1 1 3 1 2 4 1 1 2 0 2 4 1
3 2 1 2 1 1 3 1
3 2 1 2 1 1 3 1
2 1 1 2 1
3 3 1 2 1 1 3 1 2 1 0
3 3 1 2 1 2 3 1 1 2 1
3 2 1 2 1 2 3 1
4 5 1 2 1 1 3 1 3 4 1 3 1 1 3 3 1
4 5 1 2 1 2 3 1 2 4 1 2 1 1 1 3 1
2 2 1 2 1 2 1 1
2 3 1 2 1 1 1 0 2 2 1
4 5 1 2 1 1 3 1 3 4 1 3 3 0 4 1 0
4 4 1 2 1 1 3 1 2 4 1 1 3 1
4 3 1 2 1 1 3 1 3 4 1
4 3 1 2 1 1 3 1 3 4 1
2 1 1 2 1
3 3 1 2 1 2 3 1 3 2 0
2 2 1 2 1 2 1 0
3 2 1 2 1 1 3 1
5 5 1 2 1 1 3 1 1 4 1 4 5 1 2 4 0
5 5 1 2 1 2 3 1 2 4 1 3 5 1 4 4 0
2 3 1 2 1 1 2 1 1 1 1
4 5 1 2 1 1 3 1 2 4 1 1 3 1 3 4 1
4 3 1 2 1 1 3 1 3 4 1
5 5 1 2 1 2 3 1 3 4 1 3 5 1 3 5 1
4 4 1 2 1 2 3 1 3 4 1 3 2 1
4 3 1 2 1 1 3 1 3 4 1
3 3 1 2 1 2 3 1 3 1 1
3 4 1 2 1 2 3 1 3 3 1 3 1 0
3 2 1 2 1 1 3 1
3 2 1 2 1 2 3 1
3 2 1 2 1 1 3 1
4 5 1 2 1 2 3 1 1 4 1 1 4 1 3 4 0
3 2 1 2 1 2 3 1
4 3 1 2 1 2 3 1 2 4 1
4 4 1 2 1 2 3 1 1 4 1 3 3 1
3 4 1 2 1 2 3 1 3 3 1 2 1 0
5 5 1 2 1 1 3 1 1 4 1 1 5 1 1 4 0
4 5 1 2 1 2 3 1 3 4 1 4 4 1 2 4 1
5 6 1 2 1 2 3 1 1 4 1 3 5 1 3 4 1 4 1 1
4 4 1 2 1 2 3 1 1 4 1 1 1 1
3 2 1 2 1 1 3 1
5 4 1 2 1 2 3 1 3 4 1 3 5 1
5 4 1 2 1 1 3 1 2 4 1 3 5 1
4 4 1 2 1 2 3 1 3 4 1 1 1 0
4 4 1 2 1 2 3 1 3 4 1 1 4 1
4 4 1 2 1 1 3 1 1 4 1 4 1 0
2 2 1 2 1 2 2 1
5 5 1 2 1 2 3 1 2 4 1 1 5 1 1 4 1
5 6 1 2 1 1 3 1 1 4 1 4 5 1 5 4 1 5 2 1
2 3 1 2 1 1 2 0 1 1 1
5 4 1 2 1 2 3 1 3 4 1 4 5 1
2 1 1 2 1
2 2 1 2 1 2 1 0
3 2 1 2 1 1 3 1
3 4 1 2 1 1 3 1 2 3 0 2 1 0
5 5 1 2 1 1 3 1 2 4 1 1 5 1 5 4 1
4 5 1 2 1 1 3 1 1 4 1 3 2 1 2 3 0
3 4 1 2 1 2 3 1 1 1 1 2 3 1
3 4 1 2 1 1 3 1 2 1 0 1 1 1
3 3 1 2 1 2 3 1 2 1 0
5 6 1 2 1 1 3 1 2 4 1 2 5 1 1 1 1 2 2 1
4 5 1 2 1 2 3 1 2 4 1 3 4 1 2 3 1
5 6 1 2 1 1 3 1 3 4 1 2 5 1 2 4 0 1 1 0
3 2 1 2 1 1 3 1
5 5 1 2 1 2 3 1 1 4 1 4 5 1 3 5 1
4 3 1 2 1 2 3 1 1 4 1
3 3 1 2 1 1 3 1 2 2 0
4 4 1 2 1 2 3 1 1 4 1 2 1 1
2 2 1 2 1 2 1 0
5 5 1 2 1 1 3 1 1 4 1 3 5 1 4 2 0
2 3 1 2 1 1 1 0 2 2 1
3 3 1 2 1 2 3 1 1 3 1
3 4 1 2 1 1 3 1 1 1 1 3 2 0
5 4 1 2 1 2 3 1 3 4 1 4 5 1
2 3 1 2 1 2 1 0 1 2 0
5 5 1 2 1 1 3 1 3 4 1 1 5 1 4 1 0
4 5 1 2 1 2 3 1 2 4 1 1 1 0 2 4 0
3 2 1 2 1 2 3 1
2 3 1 2 1 1 1 0 2 2 1
5 5 1 2 1 1 3 1 3 4 1 1 5 1 4 2 0
3 2 1 2 1 2 3 1
4 3 1 2 1 2 3 1 1 4 1
2 1 1 2 1
2 3 1 2 1 1 2 0 2 1 0
3 2 1 2 1 2 3 1
2 2 1 2 1 2 1 1
3 2 1 2 1 2 3 1
2 3 1 2 1 2 2 0 2 1 1
2 1 1 2 1
3 4 1 2 1 2 3 1 2 2 0 1 3 0
2 1 1 2 1
3 3 1 2 1 2 3 1 3 1 0
5 5 1 2 1 2 3 1 1 4 1 3 5 1 1 2 0
2 3 1 2 1 1 1 1 2 2 0
4 4 1 2 1 2 3 1 1 4 1 4 2 1
2 3 1 2 1 2 2 0 1 1 1
2 1 1 2 1`

type edge struct {
	v, id int
}

type fastReader struct {
	r *bufio.Reader
}

func (fr *fastReader) nextInt() int {
	var (
		b   byte
		err error
	)
	for {
		b, err = fr.r.ReadByte()
		if err != nil {
			return 0
		}
		if (b >= '0' && b <= '9') || b == '-' {
			break
		}
	}
	sign := 1
	if b == '-' {
		sign = -1
		b, err = fr.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	x := 0
	for (b >= '0' && b <= '9') && err == nil {
		x = x*10 + int(b-'0')
		b, err = fr.r.ReadByte()
	}
	return x * sign
}

func solve1994F(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var out bytes.Buffer
	writer := bufio.NewWriter(&out)
	fr := fastReader{r: reader}

	t := fr.nextInt()
	for tc := 0; tc < t; tc++ {
		n := fr.nextInt()
		m := fr.nextInt()
		finalG := make([][]edge, n)
		optionalG := make([][]edge, n)
		for i := 0; i < m; i++ {
			u := fr.nextInt() - 1
			v := fr.nextInt() - 1
			c := fr.nextInt()
			if c == 1 {
				finalG[u] = append(finalG[u], edge{v, i})
				finalG[v] = append(finalG[v], edge{u, i})
			} else {
				optionalG[u] = append(optionalG[u], edge{v, i})
				optionalG[v] = append(optionalG[v], edge{u, i})
			}
		}
		trav := make([]bool, n)
		parent := make([]int, n)
		parentEdge := make([]int, n)
		iter := make([]int, n)
		odd := make([]int, n)
		for i := 0; i < n; i++ {
			odd[i] = len(finalG[i]) & 1
			parent[i] = -1
		}
		order := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if trav[i] {
				continue
			}
			trav[i] = true
			stack := []int{i}
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				if iter[u] < len(optionalG[u]) {
					e := optionalG[u][iter[u]]
					iter[u]++
					if !trav[e.v] {
						trav[e.v] = true
						parent[e.v] = u
						parentEdge[e.v] = e.id
						stack = append(stack, e.v)
					}
				} else {
					order = append(order, u)
					stack = stack[:len(stack)-1]
				}
			}
		}

		hasSol := true
		for _, u := range order {
			p := parent[u]
			if p == -1 {
				if odd[u] == 1 {
					hasSol = false
				}
				continue
			}
			if odd[u] == 1 {
				odd[p] ^= 1
				finalG[p] = append(finalG[p], edge{u, parentEdge[u]})
				finalG[u] = append(finalG[u], edge{p, parentEdge[u]})
			}
		}
		if !hasSol {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
		used := make([]bool, m)
		ptr := make([]int, n)
		var ans []int
		stack := []int{0}
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			for ptr[u] < len(finalG[u]) && used[finalG[u][ptr[u]].id] {
				ptr[u]++
			}
			if ptr[u] == len(finalG[u]) {
				ans = append(ans, u+1)
				stack = stack[:len(stack)-1]
				continue
			}
			e := finalG[u][ptr[u]]
			ptr[u]++
			used[e.id] = true
			stack = append(stack, e.v)
		}
		fmt.Fprintln(writer, len(ans)-1)
		for i, v := range ans {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		writer.WriteByte('\n')
	}
	if err := writer.Flush(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func parseTestcases() []string {
	lines := strings.Split(testcasesF, "\n")
	cases := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		cases = append(cases, trimmed)
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for idx, line := range testcases {
		input := "1\n" + line + "\n"
		expected, err := solve1994F(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solver error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected = strings.TrimSpace(expected)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed. Expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
