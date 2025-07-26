package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Node represents an AST node: either leaf or binary op
type Node struct {
	op    rune   // '+', '-', '*', '/' for op nodes
	val   string // for leaf: variable or number
	left  *Node
	right *Node
}

// parse expression into AST
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
	// leaf
	p.pos++
	return &Node{val: t}
}

// equal checks if two ASTs are identical
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
	// both leaf
	return a.val == b.val
}

// normalize flattens safe multiplication/division nodes by associativity where allowed
func normalize(n *Node) *Node {
	if n == nil {
		return nil
	}
	if n.op == 0 {
		return n
	}
	// normalize children first
	n.left = normalize(n.left)
	n.right = normalize(n.right)
	// attempt to flatten patterns
	for {
		if (n.op == '*' || n.op == '/') && n.right != nil {
			r := n.right
			// check safe flattening condition: right child op is '*' or '/' and (n.op=='*' or (n.op=='/' and r.op=='*'))
			if (r.op == '*' || r.op == '/') && (n.op == '*' || (n.op == '/' && r.op == '*')) {
				A := n.left
				B := r.left
				C := r.right
				// build new left node
				var newLeft *Node
				if n.op == '*' {
					newLeft = &Node{op: '*', left: A, right: B}
				} else {
					newLeft = &Node{op: '/', left: A, right: B}
				}
				// determine new root op
				if n.op == '*' && r.op == '*' {
					n = &Node{op: '*', left: newLeft, right: C}
				} else {
					// cases: * followed by /, or / followed by *
					n = &Node{op: '/', left: newLeft, right: C}
				}
				// normalize the new subtree
				n.left = normalize(n.left)
				n.right = normalize(n.right)
				continue
			}
		}
		break
	}
	return n
}

// tokenize splits into identifiers, numbers, operators, parentheses
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
			// skip unknown
			i++
		}
	}
	return toks
}

var macros map[string]string

// expandOrdTokens does ordinary macro expansion
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

// expandSafeTokens expands macros with parentheses
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	n, _ := strconv.Atoi(line)
	macros = make(map[string]string)
	for i := 0; i < n; i++ {
		l, _ := reader.ReadString('\n')
		l = strings.TrimSpace(l)
		if l == "" {
			i--
			continue
		}
		fields := strings.Fields(l)
		idx := 0
		if fields[0] == "#" {
			idx = 2 // # define
		} else if fields[0] == "#define" {
			idx = 1
		} else {
			// unexpected, skip
		}
		name := fields[idx]
		expr := strings.Join(fields[idx+1:], "")
		macros[name] = expr
	}
	// read target expression
	targ, _ := reader.ReadString('\n')
	targ = strings.TrimSpace(targ)
	// ordinary expansion
	ordToks := expandOrdTokens(targ)
	ord := strings.Join(ordToks, "")
	// safe expansion
	safeToks := expandSafeTokens(targ)
	safe := strings.Join(safeToks, "")
	// parse both
	p1 := &parser{tokens: tokenize(ord)}
	ast1 := p1.parseExpr()
	p2 := &parser{tokens: tokenize(safe)}
	ast2 := p2.parseExpr()
	// normalize multiplication/division associativity where safe
	ast1 = normalize(ast1)
	ast2 = normalize(ast2)
	if equal(ast1, ast2) {
		fmt.Println("OK")
	} else {
		fmt.Println("Suspicious")
	}
}
