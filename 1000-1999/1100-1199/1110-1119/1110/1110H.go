package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

// compare two decimal strings of equal length lex
func ge(a, b string) bool { // a >= b
   return a >= b
}
func le(a, b string) bool { // a <= b
   return a <= b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var l, r string
   var n int
   fmt.Fscan(reader, &l)
   fmt.Fscan(reader, &r)
   fmt.Fscan(reader, &n)
   L := len(l)
   R := len(r)
   // interior lengths d in (L,R)
   total := 0
   // sum interior
   for d := L + 1; d <= R-1; d++ {
       total += n - d + 1
   }
   ones := strings.Repeat("1", n)
   // count L-length
   if L == R {
       if n >= L {
           pat := strings.Repeat("1", L)
           if ge(pat, l) && le(pat, r) {
               total += n - L + 1
           }
       }
   } else {
       if n >= L {
           patL := strings.Repeat("1", L)
           if ge(patL, l) {
               total += n - L + 1
           }
       }
       if n >= R {
           patR := strings.Repeat("1", R)
           if le(patR, r) {
               total += n - R + 1
           }
       }
   }
   // print result
   fmt.Println(total)
   fmt.Println(ones)
}
