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

type Pair struct{ u, v int }

func expectedE(n, m int, edges []Pair) (int, []Pair) {
	v := make([][]int, n)
	for i := range v {
		v[i] = make([]int, n)
	}
	for _, e := range edges {
		a, b := e.u-1, e.v-1
		v[a][b] = 1
		v[b][a] = 1
	}
	vv := make([]int, n)
	deg := make([]int, n)
	cp := 0
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < n; j++ {
			sum += v[i][j]
		}
		deg[i] = sum
		vv[i] = 1 - (sum % 2)
		cp += vv[i]
	}
	for i := 0; i < n; i++ {
		if vv[i] == 0 {
			continue
		}
		jFlag := 0
		for deg[i] > 0 {
			a := i
			for deg[a] > 0 && vv[a] > 0 {
				for k := 0; k < n; k++ {
					if v[a][k] == 1 {
						if jFlag == 0 {
							v[a][k] = 2
							v[k][a] = 0
						} else {
							v[a][k] = 0
							v[k][a] = 2
						}
						deg[a]--
						deg[k]--
						a = k
						break
					}
				}
			}
			if deg[i]%2 == 1 {
				jFlag = 1
			} else {
				jFlag = 0
			}
		}
	}
	res := make([]Pair, 0, m)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if v[i][j] == 2 {
				res = append(res, Pair{i + 1, j + 1})
			}
			if v[i][j] == 1 {
				res = append(res, Pair{i + 1, j + 1})
				v[j][i] = 0
			}
		}
	}
	return cp, res
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		parts := strings.Fields(scan.Text())
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		edges := make([]Pair, m)
		for i := 0; i < m; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			p := strings.Fields(scan.Text())
			u, _ := strconv.Atoi(p[0])
			v, _ := strconv.Atoi(p[1])
			edges[i] = Pair{u, v}
		}
		cp, orient := expectedE(n, m, edges)
		input := fmt.Sprintf("1\n%d %d\n", n, m)
		for _, e := range edges {
			input += fmt.Sprintf("%d %d\n", e.u, e.v)
		}
		exp := fmt.Sprintf("%d\n", cp)
		for _, e := range orient {
			exp += fmt.Sprintf("%d %d\n", e.u, e.v)
		}
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
