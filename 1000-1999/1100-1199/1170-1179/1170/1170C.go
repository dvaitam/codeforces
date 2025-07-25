package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k int
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   // process each test case
   for ; k > 0; k-- {
       var s, t string
       fmt.Fscan(reader, &s)
       fmt.Fscan(reader, &t)
       if solve(s, t) {
           fmt.Println("YES")
       } else {
           fmt.Println("NO")
       }
   }
}

// solve returns true if s can be transformed to t via operations
func solve(s, t string) bool {
   // t must not contain two consecutive '-' since leftover minus runs max length 1
   if len(t) > 1 && containsDoubleMinus(t) {
       return false
   }
   sp, tp := 0, 0
   cm, cp := 0, 0
   ns, nt := len(s), len(t)
   // match t chars
   for tp < nt {
       need := t[tp]
       if need == '-' {
           // get one minus
           for sp < ns && cm < 1 {
               if s[sp] == '-' {
                   cm++
               } else {
                   cp++
               }
               sp++
           }
           if cm < 1 {
               return false
           }
           cm--
           tp++
       } else { // need '+'
           for sp < ns && cp < 1 && cm < 2 {
               if s[sp] == '-' {
                   cm++
               } else {
                   cp++
               }
               sp++
           }
           if cp > 0 {
               cp--
               tp++
           } else if cm >= 2 {
               cm -= 2
               tp++
           } else {
               return false
           }
       }
   }
   // process remaining s
   for sp < ns {
       if s[sp] == '-' {
           cm++
       } else {
           cp++
       }
       sp++
   }
   // leftover must be zero; cannot have extra '-' or '+'
   if cm != 0 || cp != 0 {
       return false
   }
   return true
}

func containsDoubleMinus(t string) bool {
   // fast scan for "--"
   prev := byte(0)
   for i := 0; i < len(t); i++ {
       if t[i] == '-' && prev == '-' {
           return true
       }
       prev = t[i]
   }
   return false
}
