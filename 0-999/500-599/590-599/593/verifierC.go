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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
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
		if err != nil {
			fmt.Printf("missing output for test %d\n", idx+1)
			os.Exit(1)
		}
		gLine, err := outReader.ReadString('\n')
		if err != nil {
			fmt.Printf("missing output for test %d\n", idx+1)
			os.Exit(1)
		}
		fLine = strings.TrimSpace(fLine)
		gLine = strings.TrimSpace(gLine)
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
