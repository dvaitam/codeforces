package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)
   cntI, cntF := 0, 0
   for _, c := range s {
       if c == 'I' {
           cntI++
       } else if c == 'F' {
           cntF++
       }
   }
   var res int
   switch {
   case cntI == 0:
       // no one 'IN', all non-folded (A) can show
       res = n - cntF
   case cntI == 1:
       // exactly one 'IN', only that one can show
       res = 1
   default:
       // two or more 'IN', none can show
       res = 0
   }
   fmt.Println(res)
}
