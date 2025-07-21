package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var p1, p2, p3, p4 int
   var a, b int
   fmt.Fscan(reader, &p1, &p2, &p3, &p4, &a, &b)
   // find minimum of p1, p2, p3, p4
   m := p1
   if p2 < m {
       m = p2
   }
   if p3 < m {
       m = p3
   }
   if p4 < m {
       m = p4
   }
   // count x in [a, b] such that x < m
   // valid x are in [a, min(b, m-1)]
   r := b
   if r > m-1 {
       r = m - 1
   }
   if r >= a {
       fmt.Println(r - a + 1)
   } else {
       fmt.Println(0)
   }
}
