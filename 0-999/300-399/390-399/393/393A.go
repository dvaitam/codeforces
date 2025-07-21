package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   cntN, cntI, cntE, cntT := 0, 0, 0, 0
   for _, ch := range s {
       switch ch {
       case 'n':
           cntN++
       case 'i':
           cntI++
       case 'e':
           cntE++
       case 't':
           cntT++
       }
   }
   // Each "nineteen" needs 1 i, 1 t, 3 e, and 2*n + 1 n's for k occurrences overlapping
   // Maximum number by n: if cntN < 3, no words; else (cntN-1)/2
   maxByN := 0
   if cntN >= 3 {
       maxByN = (cntN - 1) / 2
   }
   maxByI := cntI
   maxByT := cntT
   maxByE := cntE / 3
   // result is minimum of these
   res := min(min(maxByN, maxByI), min(maxByT, maxByE))
   fmt.Println(res)
}
