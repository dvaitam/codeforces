package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

const (
	pAdd  = 0 // + or -
	pMul  = 1 // * or /
	pAtom = 2 // atom (number, variable, parenthesized expression)
)

type macroInfo struct {
	prio       int
	suspicious bool
}

type testCase struct {
	n      int
	macros []string
	target string
}

// Embedded testcases (from testcasesE.txt) so the verifier is self contained.
const rawTestcases = `1
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
z+z*y

2
#define A0 9/9
#define A1 6+3/6
z/x

0
z

0
y-x+y

0
9

3
#define A0 9*1
#define A1 z*z
#define A2 5+7
9/8

0
1

2
#define A0 x/3+3
#define A1 3*z
8+z*8*z

0
9+7*z/9

0
8-4

2
#define A0 z/6*x-x
#define A1 9/1/8
8

0
5-1*7*4

3
#define A0 z
#define A1 4
#define A2 2
4

2
#define A0 z*y*x/4
#define A1 2+5
y-6

2
#define A0 6+5/6
#define A1 x+z*x
7*z

1
#define A0 5
2+z+4

3
#define A0 8-6+4
#define A1 x
#define A2 z*1*6
4+9+3

1
#define A0 y-6
3*z

2
#define A0 3*z-9*z
#define A1 z
2/8/2

2
#define A0 2*z+8*y
#define A1 z/4*y-1
y

0
y*3*1

1
#define A0 8*4+5/9
y/1*z+4

1
#define A0 z
z

2
#define A0 6/9+1
#define A1 2*z+z
4

1
#define A0 y/6+8
z/z/y*z+z

2
#define A0 2-2/8
#define A1 5
z*z

1
#define A0 5*y*z
z*z*z

1
#define A0 y-9
y*x*x

3
#define A0 z*5*z
#define A1 z
#define A2 4*y
8-3

2
#define A0 z
#define A1 7*z*9/2
z

0
y

0
6+x

1
#define A0 7+5
x*z*z+z

2
#define A0 z
#define A1 z/5*x
x

0
8*x*8*z

0
6*z-z+x

2
#define A0 6
#define A1 7/9+2-5
2/6

3
#define A0 y+4+7
#define A1 x
#define A2 z+x
y+8

2
#define A0 8
#define A1 y*6*7
z*y*y+y

1
#define A0 y
5/1+1

0
9*9

3
#define A0 9*z-8+x
#define A1 3
#define A2 5+y
9+8

0
6

1
#define A0 1
y*y

2
#define A0 5
#define A1 x*8
4/1/3*z*z

0
5/9*2/2/6

2
#define A0 y*z/6/5
#define A1 x
2+5

3
#define A0 z/1/7
#define A1 y-7/2*x
#define A2 8*x*1
1+9*z/z

3
#define A0 y*z*z+8
#define A1 z/6*z
#define A2 z
z

2
#define A0 6+8
#define A1 z*x
6*3

1
#define A0 z*1
y*z

2
#define A0 8
#define A1 x*1
x*x*5

3
#define A0 1+7
#define A1 4
#define A2 z+z
z/1

2
#define A0 x
#define A1 3*x
x+8/x

2
#define A0 x
#define A1 y
x/x

2
#define A0 8*x
#define A1 6*z+5
y/1

0
x/6+x/8

3
#define A0 5*y
#define A1 2
#define A2 x*y
x*x*6*z

1
#define A0 x*y*z+4
z/4

1
#define A0 z
z+8/1*z*z

1
#define A0 9
y

0
5*z

2
#define A0 x-8*6
#define A1 x*y
9

2
#define A0 5
#define A1 3*x+z
1+6/5*z

0
x*y*4*x

2
#define A0 5*z*8
#define A1 z
9*z

0
z/2/2

1
#define A0 y
7

3
#define A0 x/2*y
#define A1 z
#define A2 4/z/6
8

1
#define A0 4
6/7

2
#define A0 z*7-3
#define A1 2
z*z*z

2
#define A0 5+3
#define A1 y-2+z
y-3

2
#define A0 z
#define A1 x
7+x/2

0
x-4*z-z

0
9

3
#define A0 9+z
#define A1 3*z
#define A2 7
4

0
y*z*x/2

2
#define A0 y*z
#define A1 9*z-7-2
7

0
3

2
#define A0 y
#define A1 z+1*z
5-8/6*z*x

1
#define A0 x*z
8/8*5*z

0
y*7

3
#define A0 x+6
#define A1 2*x*z
#define A2 x*z*z
7+7*x

1
#define A0 z*4
3/4

3
#define A0 x
#define A1 3*x*z
#define A2 4
3/5+5*z*8

3
#define A0 9
#define A1 8*y
#define A2 y
x*x

2
#define A0 4*x*x
#define A1 4
x*6-5

0
2*x*x*x

1
#define A0 x*y*y
2+x/3

0
8*z*9

1
#define A0 y*y
x*6

0
4*x+y*z

2
#define A0 8*1*y
#define A1 8
3/x

3
#define A0 5*z*8
#define A1 9
#define A2 x-6*y
2

0
8/7

0
z*z*9

0
6

0
5*z

0
6

0
1

3
#define A0 y
#define A1 7
#define A2 4
4*y

3
#define A0 9
#define A1 x-2+2+z
#define A2 z-7-9
7

1
#define A0 1
z

1
#define A0 x*x
2*x

3
#define A0 x
#define A1 z-4
#define A2 z*z*z
z

1
#define A0 y+z/5
z/x

0
9/5*z

1
#define A0 9+z
1

1
#define A0 9*z
5

1
#define A0 x
z

1
#define A0 6
z*z

2
#define A0 z*y
#define A1 7
x

0
5

1
#define A0 x
y-9

1
#define A0 z*z
1

1
#define A0 y+6
8

2
#define A0 z*y
#define A1 7
z

2
#define A0 7*z
#define A1 y*z+3
7

0
1

2
#define A0 z*z*y
#define A1 5
4

3
#define A0 z
#define A1 3
#define A2 8*x
z

2
#define A0 z/9-5
#define A1 y
y*z

2
#define A0 y
#define A1 8
x*z*y

3
#define A0 8/7*z
#define A1 x
#define A2 6
3

3
#define A0 3
#define A1 x*x
#define A2 3
6

1
#define A0 z
x/6

0
z

3
#define A0 x
#define A1 5*z*z*z
#define A2 3+2*x
1*z

0
5*z

0
z

0
z
`

type parser struct {
	toks []string
	pos  int
	err  bool
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
		if strings.ContainsRune("+-*/()", r) {
			toks = append(toks, string(r))
			i++
			continue
		}
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
	return toks
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
func (p *parser) parseExpr(macros map[string]macroInfo) int {
	lhs := p.parseTerm(macros)
	if p.err {
		return 0
	}
	for {
		op := p.peek()
		if op == "+" || op == "-" {
			p.next()
			rhs := p.parseTerm(macros)
			if p.err {
				return 0
			}
			if op == "-" && rhs < pMul {
				p.err = true
				return 0
			}
			lhs = pAdd
		} else {
			break
		}
	}
	return lhs
}

// term: factor { ('*'|'/') factor }
func (p *parser) parseTerm(macros map[string]macroInfo) int {
	lhs := p.parseFactor(macros)
	if p.err {
		return 0
	}
	for {
		op := p.peek()
		if op == "*" || op == "/" {
			p.next()
			rhs := p.parseFactor(macros)
			if p.err {
				return 0
			}
			if lhs < pMul {
				p.err = true
				return 0
			}
			if op == "*" {
				if rhs < pMul {
					p.err = true
					return 0
				}
			} else {
				if rhs < pAtom {
					p.err = true
					return 0
				}
			}
			lhs = pMul
		} else {
			break
		}
	}
	return lhs
}

func (p *parser) parseFactor(macros map[string]macroInfo) int {
	t := p.next()
	if t == "(" {
		p.parseExpr(macros)
		if p.peek() == ")" {
			p.next()
		}
		if p.err {
			return 0
		}
		return pAtom
	}
	if val, ok := macros[t]; ok {
		if val.suspicious {
			p.err = true
		}
		return val.prio
	}
	return pAtom
}

func parseMacro(line string) (name string, exprTokens []string) {
	toks := tokenize(line)
	for i := 0; i < len(toks); i++ {
		if toks[i] == "define" || toks[i] == "#define" {
			if i+1 < len(toks) {
				name = toks[i+1]
				exprTokens = toks[i+2:]
			}
			break
		}
	}
	return
}

func buildMacros(defs []string) map[string]macroInfo {
	mac := make(map[string]macroInfo)
	for _, line := range defs {
		name, exprToks := parseMacro(line)
		if name == "" {
			continue
		}
		p := &parser{toks: exprToks}
		prio := p.parseExpr(mac)
		mac[name] = macroInfo{prio: prio, suspicious: p.err}
	}
	return mac
}

func solve(tc testCase) string {
	mac := buildMacros(tc.macros)
	p := &parser{toks: tokenize(tc.target)}
	p.parseExpr(mac)
	if p.err {
		return "Suspicious"
	}
	return "OK"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(rawTestcases, "\n")
	var cases []testCase
	for idx := 0; idx < len(lines); {
		line := strings.TrimSpace(lines[idx])
		if line == "" {
			idx++
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		idx++
		if idx+n > len(lines) {
			return nil, fmt.Errorf("unexpected EOF reading macros")
		}
		macros := make([]string, n)
		for i := 0; i < n; i++ {
			macros[i] = lines[idx]
			idx++
		}
		if idx >= len(lines) {
			return nil, fmt.Errorf("unexpected EOF reading target")
		}
		target := lines[idx]
		idx++
		cases = append(cases, testCase{n: n, macros: macros, target: target})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, line := range tc.macros {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	sb.WriteString(tc.target)
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected := solve(tc)
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
