package main

import (
   "bufio"
   "fmt"
   "io/ioutil"
   "os"
   "strings"
)

const mod = 32768

// Expr represents an arithmetic expression
type Expr interface {
   Eval(n int, fvals []int) int
}

// Const is a constant number
type Const struct{ v int }
func (c *Const) Eval(_ int, _ []int) int { return c.v }

// Var is the variable n
type Var struct{}
func (v *Var) Eval(n int, _ []int) int { return n }

// Call is a recursive call f(expr)
type Call struct{ arg Expr }
func (c *Call) Eval(n int, fvals []int) int {
   idx := c.arg.Eval(n, fvals)
   return fvals[idx]
}

// BinOp is a binary operation
type BinOp struct{
   op    byte
   left, right Expr
}
func (b *BinOp) Eval(n int, fvals []int) int {
   a := b.left.Eval(n, fvals)
   c := b.right.Eval(n, fvals)
   switch b.op {
   case '+':
       x := a + c
       if x >= mod || x < 0 {
           x %= mod
           if x < 0 {
               x += mod
           }
       }
       return x
   case '-':
       x := a - c
       if x < 0 || x >= mod {
           x %= mod
           if x < 0 {
               x += mod
           }
       }
       return x
   case '*':
       x := (a * c) % mod
       if x < 0 {
           x += mod
       }
       return x
   case '/':
       if c != 0 {
           return a / c
       }
       return 0
   }
   return 0
}

// Condition for if
type Cond struct{
   left Expr
   op   string
   right Expr
}
func (c *Cond) Eval(n int, fvals []int) bool {
   a := c.left.Eval(n, fvals)
   b := c.right.Eval(n, fvals)
   switch c.op {
   case ">":
       return a > b
   case "<":
       return a < b
   case "==":
       return a == b
   }
   return false
}

// Operator represents one return or conditional return
type Operator struct{
   cond *Cond
   expr Expr
}

// Parser for the function description
type Parser struct{
   s   string
   pos int
   n   int
}

func (p *Parser) skipSpace() {
   for p.pos < p.n {
       c := p.s[p.pos]
       if c == ' ' || c == '\n' || c == '\r' || c == '\t' {
           p.pos++
           continue
       }
       break
   }
}

func (p *Parser) expect(str string) {
   p.skipSpace()
   if strings.HasPrefix(p.s[p.pos:], str) {
       p.pos += len(str)
   }
}

func (p *Parser) parseOperators() []Operator {
   ops := []Operator{}
   // find opening brace
   if idx := strings.Index(p.s, "{"); idx >= 0 {
       p.pos = idx + 1
   }
   for p.pos < p.n {
       p.skipSpace()
       if p.pos < p.n && p.s[p.pos] == '}' {
           break
       }
       // parse operator
       if strings.HasPrefix(p.s[p.pos:], "if") {
           p.pos += 2
           p.skipSpace()
           p.expect("(")
           cond := p.parseCond()
           p.expect(")")
           p.skipSpace()
           p.expect("return")
           expr := p.parseExpr()
           p.expect(";")
           ops = append(ops, Operator{cond: &cond, expr: expr})
       } else if strings.HasPrefix(p.s[p.pos:], "return") {
           p.pos += 6
           expr := p.parseExpr()
           p.expect(";")
           ops = append(ops, Operator{cond: nil, expr: expr})
       } else {
           // skip unexpected
           p.pos++
       }
   }
   return ops
}

func (p *Parser) parseCond() Cond {
   p.skipSpace()
   left := p.parseExpr()
   p.skipSpace()
   op := ""
   if strings.HasPrefix(p.s[p.pos:], "==") {
       op = "=="
       p.pos += 2
   } else if p.s[p.pos] == '>' || p.s[p.pos] == '<' {
       op = string(p.s[p.pos])
       p.pos++
   }
   p.skipSpace()
   right := p.parseExpr()
   return Cond{left: left, op: op, right: right}
}

func (p *Parser) parseExpr() Expr {
   return p.parseSum()
}

func (p *Parser) parseSum() Expr {
   p.skipSpace()
   left := p.parseProduct()
   for {
       p.skipSpace()
       if p.pos < p.n && (p.s[p.pos] == '+' || p.s[p.pos] == '-') {
           op := p.s[p.pos]
           p.pos++
           right := p.parseProduct()
           left = &BinOp{op: op, left: left, right: right}
       } else {
           break
       }
   }
   return left
}

func (p *Parser) parseProduct() Expr {
   p.skipSpace()
   left := p.parseMultiplier()
   for {
       p.skipSpace()
       if p.pos < p.n && (p.s[p.pos] == '*' || p.s[p.pos] == '/') {
           op := p.s[p.pos]
           p.pos++
           right := p.parseMultiplier()
           left = &BinOp{op: op, left: left, right: right}
       } else {
           break
       }
   }
   return left
}

func (p *Parser) parseMultiplier() Expr {
   p.skipSpace()
   if p.pos < p.n {
       c := p.s[p.pos]
       if c >= '0' && c <= '9' {
           v := 0
           for p.pos < p.n && p.s[p.pos] >= '0' && p.s[p.pos] <= '9' {
               v = v*10 + int(p.s[p.pos]-'0')
               p.pos++
           }
           return &Const{v: v % mod}
       }
       if c == 'n' {
           p.pos++
           return &Var{}
       }
       if c == 'f' {
           p.pos++
           p.skipSpace()
           p.expect("(")
           arg := p.parseExpr()
           p.expect(")")
           return &Call{arg: arg}
       }
   }
   // fallback
   p.pos++
   return &Const{v: 0}
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var target int
   if _, err := fmt.Fscan(in, &target); err != nil {
       return
   }
   codeBytes, _ := ioutil.ReadAll(in)
   p := Parser{s: string(codeBytes), n: len(codeBytes)}
   ops := p.parseOperators()
   // compute f values
   const N = mod
   fvals := make([]int, N)
   ans := -1
   for i := 0; i < N; i++ {
       for _, op := range ops {
           if op.cond != nil {
               if !op.cond.Eval(i, fvals) {
                   continue
               }
           }
           v := op.expr.Eval(i, fvals)
           fvals[i] = v
           break
       }
       if fvals[i] == target {
           ans = i
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, ans)
}
