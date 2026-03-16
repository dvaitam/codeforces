package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type nodeType int

const (
	nConst nodeType = iota
	nVar
	nAbs
	nAdd
	nSub
	nMul
)

type node struct {
	typ         nodeType
	val         int64
	left, right *node
}

func parseExpr(s string) (*node, error) {
	if s == "t" {
		return &node{typ: nVar}, nil
	}
	if len(s) > 0 && s[0] == '-' {
		// no negatives allowed
	}
	if s != "" && (s[0] >= '0' && s[0] <= '9') {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return &node{typ: nConst, val: v}, nil
	}
	if strings.HasPrefix(s, "abs(") && strings.HasSuffix(s, ")") {
		inside := s[4 : len(s)-1]
		sub, err := parseExpr(inside)
		if err != nil {
			return nil, err
		}
		return &node{typ: nAbs, left: sub}, nil
	}
	if len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' {
		inner := s[1 : len(s)-1]
		level := 0
		for i, ch := range inner {
			switch ch {
			case '(':
				level++
			case ')':
				level--
			case '+', '-', '*':
				if level == 0 {
					left := inner[:i]
					right := inner[i+1:]
					lnode, err := parseExpr(left)
					if err != nil {
						return nil, err
					}
					rnode, err := parseExpr(right)
					if err != nil {
						return nil, err
					}
					switch ch {
					case '+':
						return &node{typ: nAdd, left: lnode, right: rnode}, nil
					case '-':
						return &node{typ: nSub, left: lnode, right: rnode}, nil
					case '*':
						return &node{typ: nMul, left: lnode, right: rnode}, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("cannot parse: %s", s)
}

func eval(n *node, t int64) int64 {
	switch n.typ {
	case nConst:
		return n.val
	case nVar:
		return t
	case nAbs:
		v := eval(n.left, t)
		if v < 0 {
			return -v
		}
		return v
	case nAdd:
		return eval(n.left, t) + eval(n.right, t)
	case nSub:
		return eval(n.left, t) - eval(n.right, t)
	case nMul:
		return eval(n.left, t) * eval(n.right, t)
	}
	return 0
}

func parseInput(data []byte) ([][][3]int, error) {
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, _ := strconv.Atoi(scan.Text())
	tests := make([][][3]int, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		arr := make([][3]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			x, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			y, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			r, _ := strconv.Atoi(scan.Text())
			arr[j] = [3]int{x, y, r}
		}
		tests[i] = arr
	}
	return tests, nil
}

const testcasesCRaw = `100
3
3 1 5
2 12 5
15 4 3
1
3 4 2
2
19 3 4
2 18 2
1
9 12 3
3
14 8 4
20 5 5
8 5 3
3
2 6 4
14 2 4
16 13 4
2
1 4 4
5 15 3
1
2 4 5
2
3 11 3
3 10 3
1
7 19 4
3
15 10 5
16 8 4
0 19 3
2
6 16 2
9 14 3
1
1 14 5
2
10 19 4
7 10 5
1
16 0 5
1
0 2 2
1
3 1 4
2
6 2 2
7 5 4
2
20 16 4
13 6 5
3
1 19 3
18 0 5
0 13 3
1
2 18 3
2
8 9 5
10 16 5
2
10 8 4
14 17 4
2
10 12 4
0 6 5
2
12 13 3
5 15 3
1
19 6 2
2
16 16 3
18 5 5
1
15 7 5
1
3 1 2
3
2 3 3
4 18 3
20 11 4
1
3 12 5
1
10 18 3
3
7 19 5
19 9 2
4 13 2
3
6 4 3
8 8 4
5 14 4
3
3 2 5
0 1 3
10 13 3
2
13 18 2
10 2 5
3
15 9 3
12 1 4
10 5 2
3
8 14 2
19 20 3
4 17 4
3
14 2 5
8 14 4
0 12 2
2
3 4 4
10 9 3
2
14 0 4
4 13 2
2
5 18 3
2 13 5
3
20 4 5
10 4 4
18 5 2
1
7 3 3
2
20 6 2
12 10 4
2
18 8 5
16 1 4
2
13 7 5
0 16 5
3
10 6 4
9 5 3
8 6 2
2
12 4 4
8 19 3
2
0 2 3
12 13 2
1
18 7 2
1
10 14 2
1
5 19 2
2
8 6 2
15 2 2
3
1 19 3
5 15 3
16 9 2
1
11 17 2
3
10 20 4
17 7 4
0 18 5
2
2 10 4
11 3 4
3
4 11 2
0 20 5
0 13 3
2
18 14 4
0 5 5
3
2 20 2
1 2 3
12 12 3
3
7 20 2
0 1 5
12 9 3
2
0 10 5
16 6 4
3
7 1 5
15 16 2
6 10 5
3
5 18 3
6 1 3
3 9 4
2
9 4 3
18 13 4
1
12 3 5
1
12 7 5
2
4 17 5
19 11 2
3
19 16 3
6 16 5
15 7 2
3
3 4 3
9 0 4
3 4 3
3
13 8 2
11 4 2
16 6 5
3
19 8 5
1 7 4
6 10 2
3
8 15 4
9 1 4
3 11 4
1
10 7 3
1
6 12 2
3
8 5 4
12 0 2
17 13 2
2
6 16 4
2 3 5
3
17 11 3
4 3 4
9 5 5
2
7 17 3
2 15 3
2
13 8 4
3 4 5
3
11 16 3
2 14 5
15 10 4
2
4 9 2
13 2 5
3
5 2 2
12 7 4
20 3 5
3
18 9 2
18 1 5
13 16 4
1
7 5 2
1
17 19 3
2
11 15 3
13 4 3
2
5 0 4
4 3 4
3
7 3 2
15 16 5
9 10 2
1
1 20 5
2
16 10 3
18 10 4
3
2 14 3
1 3 5
8 10 4
2
0 10 4
20 7 4
3
11 11 5
6 0 5
20 18 4
2
19 0 2
10 0 5
1
2 4 2
3
17 9 3
20 11 3
18 2 3
3
13 10 4
0 2 5
8 20 5
1
5 14 2
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := []byte(testcasesCRaw)
	tests, err := parseInput(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for idx, circles := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(circles)))
		for _, c := range circles {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", c[0], c[1], c[2]))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		outReader := bufio.NewReader(bytes.NewReader(out))
		fLine, err := outReader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Printf("missing output for test %d\n", idx+1)
			os.Exit(1)
		}
		gLine, err := outReader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Printf("missing output for test %d\n", idx+1)
			os.Exit(1)
		}
		fLine = strings.TrimSpace(fLine)
		gLine = strings.TrimSpace(gLine)
		if fLine == "" || gLine == "" {
			fmt.Printf("missing output for test %d\n", idx+1)
			os.Exit(1)
		}
		fExpr, err := parseExpr(fLine)
		if err != nil {
			fmt.Printf("invalid function f in test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		gExpr, err := parseExpr(gLine)
		if err != nil {
			fmt.Printf("invalid function g in test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		ok := true
		for _, c := range circles {
			hit := false
			for t := int64(0); t <= 50 && !hit; t++ {
				x := eval(fExpr, t)
				y := eval(gExpr, t)
				dx := float64(x - int64(c[0]))
				dy := float64(y - int64(c[1]))
				if dx*dx+dy*dy <= float64(c[2]*c[2])+1e-9 {
					hit = true
				}
			}
			if !hit {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Printf("test %d failed\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
