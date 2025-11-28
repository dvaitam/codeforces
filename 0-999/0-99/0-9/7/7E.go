package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

const (
	P_ADD  = 0 // +, -
	P_MUL  = 1 // *, /
	P_ATOM = 2 // ( ), var, num
)

type MacroInfo struct {
	prio       int
	suspicious bool
}

var macros map[string]MacroInfo

func tokenize(s string) []string {
	var toks []string
	i := 0
	for i < len(s) {
		r := rune(s[i])
		if unicode.IsSpace(r) {
			i++
			continue
		}
		if strings.ContainsRune("+-*/()", r) {
			toks = append(toks, string(r))
			i++
		} else {
			j := i
			for j < len(s) {
				rj := rune(s[j])
				if strings.ContainsRune("+-*/()", rj) || unicode.IsSpace(rj) {
					break
				}
				j++
			}
			toks = append(toks, s[i:j])
			i = j
		}
	}
	return toks
}

type parser struct {
	toks []string
	pos  int
	err  bool
}

func (p *parser) peek() string {
	if p.pos < len(p.toks) {
		return p.toks[p.pos]
	}
	return ""
}

func (p *parser) next() string {
	t := p.peek()
	if t != "" {
		p.pos++
	}
	return t
}

// expr: term { ('+'|'-') term }
func (p *parser) parseExpr() int {
	lhs := p.parseTerm()
	if p.err {
		return 0
	}

	for {
		op := p.peek()
		if op == "+" || op == "-" {
			p.next()
			rhs := p.parseTerm()
			if p.err {
				return 0
			}

			// Validity Check
			if op == "-" {
				// For A - B, B must be >= P_MUL (cannot have exposed + or -)
				if rhs < P_MUL {
					p.err = true
					return 0
				}
			}
			// For +, rhs >= P_ADD is sufficient (always true)

			lhs = P_ADD
		} else {
			break
		}
	}
	return lhs
}

// term: factor { ('*'|'/') factor }
func (p *parser) parseTerm() int {
	lhs := p.parseFactor()
	if p.err {
		return 0
	}

	for {
		op := p.peek()
		if op == "*" || op == "/" {
			p.next()
			rhs := p.parseFactor()
			if p.err {
				return 0
			}

			// Check LHS
			if lhs < P_MUL {
				p.err = true
				return 0
			}

			// Check RHS
			if op == "*" {
				if rhs < P_MUL {
					p.err = true
					return 0
				}
			} else { // op == "/"
				if rhs < P_ATOM {
					p.err = true
					return 0
				}
			}

			lhs = P_MUL
		} else {
			break
		}
	}
	return lhs
}

func (p *parser) parseFactor() int {
	t := p.next()
	if t == "(" {
		p.parseExpr()
		if p.peek() == ")" {
			p.next()
		}
		if p.err {
			return 0
		}
		return P_ATOM
	}

	// Identifier or Number
	if val, ok := macros[t]; ok {
		if val.suspicious {
			p.err = true
		}
		return val.prio
	}
	return P_ATOM
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	var n int
	fmt.Sscan(scanner.Text(), &n)

	macros = make(map[string]MacroInfo)

	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		fullToks := tokenize(line)
		
		// Identify macro name
		// Format: # define NAME ... or #define NAME ...
		name := ""
		exprStart := 0
		
		// Scan tokens to find "define" or "#define" then next is NAME
		for k := 0; k < len(fullToks); k++ {
			if fullToks[k] == "define" || fullToks[k] == "#define" {
				if k+1 < len(fullToks) {
					name = fullToks[k+1]
					exprStart = k + 2
				}
				break
			}
		}
		
		if name == "" {
			continue
		}

		exprToks := fullToks[exprStart:]
		
		p := &parser{toks: exprToks, pos: 0}
		prio := p.parseExpr()
		// Don't exit on error, just mark suspicious
		macros[name] = MacroInfo{prio: prio, suspicious: p.err}
	}

	if !scanner.Scan() {
		return
	}
	target := scanner.Text()
	p := &parser{toks: tokenize(target), pos: 0}
	p.parseExpr()
	if p.err {
		fmt.Println("Suspicious")
	} else {
		fmt.Println("OK")
	}
}
