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
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
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
