package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

const testcasesE = `1
#define A0 y-1*y
z

3
#define A0 x
#define A1 z-2/y/y
#define A2 8*x*4
y*z

3
#define A0 x-z-9+z
#define A1 2+8+3
#define A2 x
x

3
#define A0 y-y+2
#define A1 1/6/y
#define A2 z/z*4
x*z-z

0
6*8+1

2
#define A0 x*6
#define A1 1
9-5+z

0
x+3/z

2
#define A0 9
#define A1 z*2
y/y

3
#define A0 4
#define A1 z
#define A2 z*z
z-x-z

3
#define A0 x+1+y
#define A1 9
#define A2 y
2

1
#define A0 3
y/5

0
x

0
8-x*y

2
#define A0 7
#define A1 6/1
8/z/7*8

3
#define A0 9
#define A1 y
#define A2 1
y/z

0
1/4/2

1
#define A0 x+z*z
z

1
#define A0 z
3

0
4

3
#define A0 6/y/4-y
#define A1 z/x-2/8
#define A2 x
9-7*z

2
#define A0 7
#define A1 x
x

3
#define A0 y/z+y
#define A1 4+z-x-y
#define A2 5*1+2
y

0
z

2
#define A0 7
#define A1 8
x

2
#define A0 8
#define A1 6/y*z
y

0
5-z+x+8

1
#define A0 y-z+3
9

3
#define A0 7
#define A1 8
#define A2 y/y*x
2/z/4

1
#define A0 4*3*5
y

3
#define A0 y
#define A1 y/z+1
#define A2 9-x*9+y
2+x/2

0
z-1*z

0
y+y

1
#define A0 2/8/8
1

2
#define A0 9
#define A1 1
8+1-7

1
#define A0 8
y

3
#define A0 1
#define A1 8
#define A2 x
4

1
#define A0 2
8/y*6

2
#define A0 x
#define A1 2+y
y

1
#define A0 8+z/z
7

0
y*7*7/8

2
#define A0 z+y/x
#define A1 z
1

0
z

2
#define A0 2+8
#define A1 x-1-5
z-9+x

0
y+4

3
#define A0 8/y*9+5
#define A1 2
#define A2 2-x-y
4`

type Node struct {
	op    rune
	val   string
	left  *Node
	right *Node
}

type parser struct {
	tokens []string
	pos    int
}

func (p *parser) parseExpr() *Node {
	node := p.parseTerm()
	for p.pos < len(p.tokens) {
		t := p.tokens[p.pos]
		if t == "+" || t == "-" {
			op := rune(t[0])
			p.pos++
			right := p.parseTerm()
			node = &Node{op: op, left: node, right: right}
			continue
		}
		break
	}
	return node
}

func (p *parser) parseTerm() *Node {
	node := p.parseFactor()
	for p.pos < len(p.tokens) {
		t := p.tokens[p.pos]
		if t == "*" || t == "/" {
			op := rune(t[0])
			p.pos++
			right := p.parseFactor()
			node = &Node{op: op, left: node, right: right}
			continue
		}
		break
	}
	return node
}

func (p *parser) parseFactor() *Node {
	if p.pos >= len(p.tokens) {
		return nil
	}
	t := p.tokens[p.pos]
	if t == "(" {
		p.pos++
		node := p.parseExpr()
		if p.pos < len(p.tokens) && p.tokens[p.pos] == ")" {
			p.pos++
		}
		return node
	}
	p.pos++
	return &Node{val: t}
}

func equal(a, b *Node) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.op != 0 || b.op != 0 {
		if a.op != b.op {
			return false
		}
		return equal(a.left, b.left) && equal(a.right, b.right)
	}
	return a.val == b.val
}

func normalize(n *Node) *Node {
	if n == nil {
		return nil
	}
	if n.op == 0 {
		return n
	}
	n.left = normalize(n.left)
	n.right = normalize(n.right)
	for {
		if (n.op == '*' || n.op == '/') && n.right != nil {
			r := n.right
			if (r.op == '*' || r.op == '/') && (n.op == '*' || (n.op == '/' && r.op == '*')) {
				A := n.left
				B := r.left
				C := r.right
				var newLeft *Node
				if n.op == '*' {
					newLeft = &Node{op: '*', left: A, right: B}
				} else {
					newLeft = &Node{op: '/', left: A, right: B}
				}
				if n.op == '*' && r.op == '*' {
					n = &Node{op: '*', left: newLeft, right: C}
				} else {
					n = &Node{op: '/', left: newLeft, right: C}
				}
				n.left = normalize(n.left)
				n.right = normalize(n.right)
				continue
			}
		}
		break
	}
	return n
}

func tokenize(s string) []string {
	var toks []string
	i := 0
	for i < len(s) {
		r := rune(s[i])
		if unicode.IsSpace(r) {
			i++
			continue
		}
		if unicode.IsLetter(r) {
			j := i + 1
			for j < len(s) && unicode.IsLetter(rune(s[j])) {
				j++
			}
			toks = append(toks, s[i:j])
			i = j
		} else if unicode.IsDigit(r) {
			j := i + 1
			for j < len(s) && unicode.IsDigit(rune(s[j])) {
				j++
			}
			toks = append(toks, s[i:j])
			i = j
		} else if strings.ContainsRune("+-*/()", r) {
			toks = append(toks, string(r))
			i++
		} else {
			i++
		}
	}
	return toks
}

var macros map[string]string

func expandOrdTokens(s string) []string {
	toks := tokenize(s)
	var res []string
	for _, t := range toks {
		if def, ok := macros[t]; ok {
			inner := expandOrdTokens(def)
			res = append(res, inner...)
		} else {
			res = append(res, t)
		}
	}
	return res
}

func expandSafeTokens(s string) []string {
	toks := tokenize(s)
	var res []string
	for _, t := range toks {
		if def, ok := macros[t]; ok {
			res = append(res, "(")
			inner := expandSafeTokens(def)
			res = append(res, inner...)
			res = append(res, ")")
		} else {
			res = append(res, t)
		}
	}
	return res
}

func compute(lines []string) string {
	n, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	macros = make(map[string]string)
	idx := 1
	for i := 0; i < n; i++ {
		l := strings.TrimSpace(lines[idx])
		idx++
		if l == "" {
			i--
			continue
		}
		fields := strings.Fields(l)
		p := 0
		if fields[0] == "#" {
			p = 2
		} else if fields[0] == "#define" {
			p = 1
		}
		name := fields[p]
		expr := strings.Join(fields[p+1:], "")
		macros[name] = expr
	}
	targ := strings.TrimSpace(lines[idx])
	ord := strings.Join(expandOrdTokens(targ), "")
	safe := strings.Join(expandSafeTokens(targ), "")
	p1 := &parser{tokens: tokenize(ord)}
	ast1 := p1.parseExpr()
	p2 := &parser{tokens: tokenize(safe)}
	ast2 := p2.parseExpr()
	ast1 = normalize(ast1)
	ast2 = normalize(ast2)
	if equal(ast1, ast2) {
		return "OK"
	}
	return "Suspicious"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesE))
	idx := 0
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if len(lines) == 0 {
				continue
			}
			idx++
			expect := compute(lines)
			var input bytes.Buffer
			for _, l := range lines {
				input.WriteString(strings.TrimSpace(l))
				input.WriteByte('\n')
			}
			cmd := exec.Command(bin)
			cmd.Stdin = bytes.NewReader(input.Bytes())
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Test %d: runtime error: %v\n", idx, err)
				os.Exit(1)
			}
			got := strings.TrimSpace(string(out))
			if got != expect {
				fmt.Printf("Test %d failed: expected %s got %s\n", idx, expect, got)
				os.Exit(1)
			}
			lines = lines[:0]
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) > 0 {
		idx++
		expect := compute(lines)
		var input bytes.Buffer
		for _, l := range lines {
			input.WriteString(strings.TrimSpace(l))
			input.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expect {
			fmt.Printf("Test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
