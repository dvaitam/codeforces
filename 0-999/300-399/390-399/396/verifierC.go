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

const MOD = 1000000007

type TestCase struct {
	n       int
	parents []int
	ops     [][]int
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(tc TestCase) []int {
	n := tc.n
	children := make([][]int, n+1)
	for i, p := range tc.parents {
		children[p] = append(children[p], i+2)
	}
	depth := make([]int, n+1)
	stack := []int{1}
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, v := range children[u] {
			depth[v] = depth[u] + 1
			stack = append(stack, v)
		}
	}
	val := make([]int64, n+1)
	res := []int{}
	for _, op := range tc.ops {
		if op[0] == 1 {
			v := op[1]
			x := int64(op[2])
			k := int64(op[3])
			q := []int{v}
			for len(q) > 0 {
				u := q[len(q)-1]
				q = q[:len(q)-1]
				add := x + k*int64(depth[u]-depth[v])
				val[u] = (val[u] + add) % MOD
				for _, w := range children[u] {
					q = append(q, w)
				}
			}
		} else {
			v := op[1]
			res = append(res, int((val[v]%MOD+MOD)%MOD))
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]TestCase, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		parents := make([]int, n-1)
		if n > 1 {
			for j := 0; j < n-1; j++ {
				if !scan.Scan() {
					fmt.Println("bad file")
					os.Exit(1)
				}
				parents[j], _ = strconv.Atoi(scan.Text())
			}
		}
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		q, _ := strconv.Atoi(scan.Text())
		ops := make([][]int, q)
		for j := 0; j < q; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			typ, _ := strconv.Atoi(scan.Text())
			if typ == 1 {
				vals := make([]int, 4)
				vals[0] = typ
				for k := 1; k <= 3; k++ {
					if !scan.Scan() {
						fmt.Println("bad file")
						os.Exit(1)
					}
					vals[k], _ = strconv.Atoi(scan.Text())
				}
				ops[j] = vals
			} else {
				vals := make([]int, 2)
				vals[0] = typ
				if !scan.Scan() {
					fmt.Println("bad file")
					os.Exit(1)
				}
				vals[1], _ = strconv.Atoi(scan.Text())
				ops[j] = vals
			}
		}
		cases[i] = TestCase{n, parents, ops}
	}
	caseNum := 0
	for idx, tc := range cases {
		var input bytes.Buffer
		fmt.Fprintln(&input, tc.n)
		if tc.n > 1 {
			for j, v := range tc.parents {
				if j > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprint(&input, v)
			}
			input.WriteByte('\n')
		} else {
			input.WriteByte('\n')
		}
		fmt.Fprintln(&input, len(tc.ops))
		for _, op := range tc.ops {
			if op[0] == 1 {
				fmt.Fprintf(&input, "1 %d %d %d\n", op[1], op[2], op[3])
			} else {
				fmt.Fprintf(&input, "2 %d\n", op[1])
			}
		}
		expected := solve(tc)
		gotStr, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input.String())
			os.Exit(1)
		}
		outScan := bufio.NewScanner(strings.NewReader(gotStr))
		outScan.Split(bufio.ScanWords)
		for _, exp := range expected {
			if !outScan.Scan() {
				fmt.Fprintf(os.Stderr, "case %d failed: missing output\n", idx+1)
				os.Exit(1)
			}
			gotVal, _ := strconv.Atoi(outScan.Text())
			if gotVal != exp {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", idx+1, exp, gotVal, input.String())
				os.Exit(1)
			}
			caseNum++
		}
		if outScan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d failed: extra output\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
