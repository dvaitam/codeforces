package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded correct solver for 39G
const (
	numKind  = iota
	varKind
	callKind
	binKind
)

type Expr struct {
	kind        int
	val         int
	op          byte
	left, right *Expr
}

type Cond struct {
	kind        int
	left, right *Expr
}

type Stmt struct {
	cond *Cond
	expr *Expr
}

type Parser struct {
	s string
	i int
}

func (p *Parser) expect(t string) {
	if !strings.HasPrefix(p.s[p.i:], t) {
		panic("parse")
	}
	p.i += len(t)
}

func (p *Parser) parseFunction() []Stmt {
	p.expect("intf(intn){")
	var stmts []Stmt
	for p.i < len(p.s) && p.s[p.i] != '}' {
		stmts = append(stmts, p.parseStmt())
	}
	p.expect("}")
	return stmts
}

func (p *Parser) parseStmt() Stmt {
	if strings.HasPrefix(p.s[p.i:], "return") {
		p.i += len("return")
		e := p.parseExpr()
		p.expect(";")
		return Stmt{expr: e}
	}
	p.expect("if(")
	c := p.parseLogical()
	p.expect(")")
	p.expect("return")
	e := p.parseExpr()
	p.expect(";")
	return Stmt{cond: c, expr: e}
}

func (p *Parser) parseLogical() *Cond {
	l := p.parseExpr()
	if strings.HasPrefix(p.s[p.i:], "==") {
		p.i += 2
		r := p.parseExpr()
		return &Cond{kind: 2, left: l, right: r}
	}
	if p.s[p.i] == '<' {
		p.i++
		r := p.parseExpr()
		return &Cond{kind: 0, left: l, right: r}
	}
	p.i++
	r := p.parseExpr()
	return &Cond{kind: 1, left: l, right: r}
}

func (p *Parser) parseExpr() *Expr {
	return p.parseSum()
}

func (p *Parser) parseSum() *Expr {
	node := p.parseProduct()
	for p.i < len(p.s) {
		c := p.s[p.i]
		if c != '+' && c != '-' {
			break
		}
		p.i++
		r := p.parseProduct()
		node = &Expr{kind: binKind, op: c, left: node, right: r}
	}
	return node
}

func (p *Parser) parseProduct() *Expr {
	node := p.parseMultiplier()
	for p.i < len(p.s) {
		c := p.s[p.i]
		if c != '*' && c != '/' {
			break
		}
		p.i++
		r := p.parseMultiplier()
		node = &Expr{kind: binKind, op: c, left: node, right: r}
	}
	return node
}

func (p *Parser) parseMultiplier() *Expr {
	if strings.HasPrefix(p.s[p.i:], "f(") {
		p.i += 2
		arg := p.parseExpr()
		p.expect(")")
		return &Expr{kind: callKind, left: arg}
	}
	if p.s[p.i] == 'n' {
		p.i++
		return &Expr{kind: varKind}
	}
	start := p.i
	for p.i < len(p.s) && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	v, _ := strconv.Atoi(p.s[start:p.i])
	return &Expr{kind: numKind, val: v}
}

func evalExpr(e *Expr, n int, memo []int) int {
	switch e.kind {
	case numKind:
		return e.val
	case varKind:
		return n
	case callKind:
		return memo[evalExpr(e.left, n, memo)]
	default:
		l := evalExpr(e.left, n, memo)
		r := evalExpr(e.right, n, memo)
		switch e.op {
		case '+':
			return (l + r) & 32767
		case '-':
			return (l - r) & 32767
		case '*':
			return (l * r) & 32767
		default:
			return l / r
		}
	}
}

func evalCond(c *Cond, n int, memo []int) bool {
	l := evalExpr(c.left, n, memo)
	r := evalExpr(c.right, n, memo)
	switch c.kind {
	case 0:
		return l < r
	case 1:
		return l > r
	default:
		return l == r
	}
}

func evalFunc(stmts []Stmt, n int, memo []int) int {
	for _, st := range stmts {
		if st.cond == nil || evalCond(st.cond, n, memo) {
			return evalExpr(st.expr, n, memo)
		}
	}
	return 0
}

func isSpace(c byte) bool {
	switch c {
	case ' ', '\n', '\r', '\t', '\v', '\f':
		return true
	}
	return false
}

func solveG(input string) string {
	data := []byte(input)
	i := 0
	for i < len(data) && isSpace(data[i]) {
		i++
	}
	j := i
	for j < len(data) && data[j] >= '0' && data[j] <= '9' {
		j++
	}
	target, _ := strconv.Atoi(string(data[i:j]))
	rest := data[j:]

	buf := make([]byte, 0, len(rest))
	for _, c := range rest {
		if !isSpace(c) {
			buf = append(buf, c)
		}
	}

	p := Parser{s: string(buf)}
	stmts := p.parseFunction()

	memo := make([]int, 32768)
	ans := -1
	for n := 0; n < 32768; n++ {
		memo[n] = evalFunc(stmts, n, memo)
		if memo[n] == target {
			ans = n
		}
	}

	return strconv.Itoa(ans)
}

func generateCaseG(rng *rand.Rand) string {
	typ := rng.Intn(3)
	switch typ {
	case 0:
		target := rng.Intn(32768)
		return fmt.Sprintf("%d\nreturn n;\n", target)
	case 1:
		c := rng.Intn(10) + 1
		n := rng.Intn(32768)
		var target int
		if n > c {
			target = (n - c) % 32768
		} else {
			target = n
		}
		code := fmt.Sprintf("if (n>%d) return n-%d; return n;", c, c)
		return fmt.Sprintf("%d\n%s\n", target, code)
	default:
		c := rng.Intn(10) + 1
		n := rng.Intn(32768)
		var target int
		if n < c {
			target = n
		} else {
			target = n - c
		}
		code := fmt.Sprintf("if (n<%d) return n; return n-%d;", c, c)
		return fmt.Sprintf("%d\n%s\n", target, code)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// Keep io import alive
var _ = io.Discard

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseG(rng)
	}
	for i, tc := range cases {
		expect := solveG(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
