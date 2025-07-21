package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

// Parser for Perse-script string expressions
type Parser struct {
   s   string
   pos int
}

func (p *Parser) skipSpaces() {
   for p.pos < len(p.s) && (p.s[p.pos] == ' ' || p.s[p.pos] == '\t' || p.s[p.pos] == '\r' || p.s[p.pos] == '\n') {
       p.pos++
   }
}

func (p *Parser) peek() byte {
   if p.pos < len(p.s) {
       return p.s[p.pos]
   }
   return 0
}

func (p *Parser) eat(ch byte) bool {
   p.skipSpaces()
   if p.pos < len(p.s) && p.s[p.pos] == ch {
       p.pos++
       return true
   }
   return false
}

func (p *Parser) parseString() string {
   // assume current char is '"'
   p.pos++ // skip '"'
   start := p.pos
   for p.pos < len(p.s) && p.s[p.pos] != '"' {
       p.pos++
   }
   res := p.s[start:p.pos]
   // skip ending '"'
   if p.pos < len(p.s) && p.s[p.pos] == '"' {
       p.pos++
   }
   return res
}

func (p *Parser) parseInt() int {
   p.skipSpaces()
   start := p.pos
   if p.pos < len(p.s) && (p.s[p.pos] == '-' || p.s[p.pos] == '+') {
       p.pos++
   }
   for p.pos < len(p.s) && p.s[p.pos] >= '0' && p.s[p.pos] <= '9' {
       p.pos++
   }
   val, err := strconv.Atoi(p.s[start:p.pos])
   if err != nil {
       // invalid number, but per problem input is valid
       return 0
   }
   return val
}

func (p *Parser) parseIdent() string {
   p.skipSpaces()
   start := p.pos
   for p.pos < len(p.s) && ((p.s[p.pos] >= 'a' && p.s[p.pos] <= 'z') || (p.s[p.pos] >= 'A' && p.s[p.pos] <= 'Z')) {
       p.pos++
   }
   return p.s[start:p.pos]
}

func (p *Parser) parseExpr() string {
   p.skipSpaces()
   if p.peek() == '"' {
       return p.parseString()
   }
   // function call
   name := strings.ToLower(p.parseIdent())
   p.skipSpaces()
   // expect '('
   if !p.eat('(') {
       return ""
   }
   var res string
   switch name {
   case "concat":
       // concat(x,y)
       x := p.parseExpr()
       p.eat(',')
       y := p.parseExpr()
       p.eat(')')
       res = x + y
   case "reverse":
       // reverse(x)
       x := p.parseExpr()
       p.eat(')')
       // reverse string x
       rs := []rune(x)
       for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
           rs[i], rs[j] = rs[j], rs[i]
       }
       res = string(rs)
   case "substr":
       // substr(x,a,b) or substr(x,a,b,c)
       x := p.parseExpr()
       p.eat(',')
       a := p.parseInt()
       p.eat(',')
       b := p.parseInt()
       // check if step c exists
       step := 1
       if p.eat(',') {
           step = p.parseInt()
       }
       p.eat(')')
       // 1-based inclusive [a,b]
       // ensure indices
       if a < 1 {
           a = 1
       }
       if b > len(x) {
           b = len(x)
       }
       if step <= 1 {
           // simple substring
           if a > b {
               res = ""
           } else {
               res = x[a-1 : b]
           }
       } else {
           // step substring
           var sb strings.Builder
           for i := a; i <= b; i += step {
               sb.WriteByte(x[i-1])
           }
           res = sb.String()
       }
   default:
       // unknown, but input valid
       res = ""
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, _ := reader.ReadString('\n')
   line = strings.TrimSpace(line)
   parser := Parser{s: line, pos: 0}
   result := parser.parseExpr()
   fmt.Print(result)
}
