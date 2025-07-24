package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

var s []byte
var n int
var match []int
var digitEnd []int

// build matching parentheses
func buildMatch() {
   match = make([]int, n)
   stack := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if s[i] == '(' {
           stack = append(stack, i)
       } else if s[i] == ')' {
           j := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           match[i] = j
           match[j] = i
       }
   }
}

// build digit end positions
func buildDigitEnd() {
   digitEnd = make([]int, n)
   for i := n - 1; i >= 0; i-- {
       if s[i] >= '0' && s[i] <= '9' {
           if i+1 < n && s[i+1] >= '0' && s[i+1] <= '9' {
               digitEnd[i] = digitEnd[i+1]
           } else {
               digitEnd[i] = i
           }
       } else {
           digitEnd[i] = i
       }
   }
}

// parseE from pos up to bound r, returns value, next pos, ok
func parseE(pos, r int) (int, int, bool) {
   val, p, ok := parseT(pos, r)
   if !ok {
       return 0, pos, false
   }
   for p <= r && s[p] == '+' {
       p++
       t, np, ok2 := parseT(p, r)
       if !ok2 {
           return 0, pos, false
       }
       val += t
       if val >= MOD {
           val -= MOD
       }
       p = np
   }
   return val, p, true
}

// parseT: term
func parseT(pos, r int) (int, int, bool) {
   val, p, ok := parseF(pos, r)
   if !ok {
       return 0, pos, false
   }
   for p <= r && s[p] == '*' {
       p++
       f, np, ok2 := parseF(p, r)
       if !ok2 {
           return 0, pos, false
       }
       val = int((int64(val) * int64(f)) % MOD)
       p = np
   }
   return val, p, true
}

// parseF: factor
func parseF(pos, r int) (int, int, bool) {
   if pos > r {
       return 0, pos, false
   }
   if s[pos] == '(' {
       j := match[pos]
       if j > r {
           return 0, pos, false
       }
       v, np, ok := parseE(pos+1, j-1)
       if !ok || np != j {
           return 0, pos, false
       }
       return v, j + 1, true
   }
   // number
   if s[pos] >= '0' && s[pos] <= '9' {
       end := digitEnd[pos]
       var v int64
       for i := pos; i <= end; i++ {
           v = (v*10 + int64(s[i]-'0')) % MOD
       }
       return int(v), end + 1, true
   }
   return 0, pos, false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   str, _ := reader.ReadString('\n')
   s = []byte(str)
   // remove newline
   if len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
       s = s[:len(s)-1]
   }
   n = len(s)
   buildMatch()
   buildDigitEnd()
   var m int
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       // to 0-index
       l--
       r--
       val, p, ok := parseE(l, r)
       if !ok || p != r+1 {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, val)
       }
   }
}
