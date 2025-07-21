package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var s, t string
   if _, err := fmt.Fscan(in, &s, &t); err != nil {
       return
   }
   n, m := len(s), len(t)
   L := n
   if m > L {
       L = m
   }
   // find first differing position, with virtual leading zeros
   pos := -1
   for i := 0; i < L; i++ {
       var cs, ct byte
       if i < L-n {
           cs = '0'
       } else {
           cs = s[i-(L-n)]
       }
       if i < L-m {
           ct = '0'
       } else {
           ct = t[i-(L-m)]
       }
       if cs != ct {
           pos = i
           break
       }
   }
   if pos == -1 {
       fmt.Println("=")
       return
   }
   // initial sign at first differing bit
   var cs, ct byte
   if pos < L-n {
       cs = '0'
   } else {
       cs = s[pos-(L-n)]
   }
   if pos < L-m {
       ct = '0'
   } else {
       ct = t[pos-(L-m)]
   }
   var sign0 float64
   if cs == '1' {
       sign0 = 1.0
   } else {
       sign0 = -1.0
   }
   invQ := 2.0 / (1.0 + math.Sqrt(5.0))
   tval := sign0
   power := 1.0
   // accumulate tail contributions
   for j := pos + 1; j < L; j++ {
       power *= invQ
       if power < 1e-18 {
           break
       }
       // bits at position j
       if j < L-n {
           cs = '0'
       } else {
           cs = s[j-(L-n)]
       }
       if j < L-m {
           ct = '0'
       } else {
           ct = t[j-(L-m)]
       }
       if cs == ct {
           continue
       }
       if cs == '1' {
           tval += power
       } else {
           tval -= power
       }
   }
   const eps = 1e-9
   switch {
   case tval > eps:
       fmt.Println(">")
   case tval < -eps:
       fmt.Println("<")
   default:
       fmt.Println("=")
   }
}
