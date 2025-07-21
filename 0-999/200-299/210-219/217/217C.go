package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read n, but not used
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // read formula
   var formula string
   if _, err := fmt.Fscan(reader, &formula); err != nil {
       return
   }
   // stacks for two evaluations: val0 = f with ?->0; val1 = f with ?->1
   // operator stack holds bytes: '(', '|', '&', '^'
   op := make([]byte, 0, len(formula))
   val0 := make([]bool, 0, len(formula))
   val1 := make([]bool, 0, len(formula))
   for i := 0; i < len(formula); i++ {
       c := formula[i]
       switch c {
       case '0':
           val0 = append(val0, false)
           val1 = append(val1, false)
       case '1':
           val0 = append(val0, true)
           val1 = append(val1, true)
       case '?':
           val0 = append(val0, false)
           val1 = append(val1, true)
       case '(', '|', '&', '^':
           op = append(op, c)
       case ')':
           // pop v2
           l := len(val0)
           v2_0 := val0[l-1]
           v2_1 := val1[l-1]
           val0 = val0[:l-1]
           val1 = val1[:l-1]
           // pop v1
           l = len(val0)
           v1_0 := val0[l-1]
           v1_1 := val1[l-1]
           val0 = val0[:l-1]
           val1 = val1[:l-1]
           // pop operator
           oper := op[len(op)-1]
           op = op[:len(op)-1]
           // pop matching '('
           if len(op) > 0 && op[len(op)-1] == '(' {
               op = op[:len(op)-1]
           }
           // compute
           var r0, r1 bool
           switch oper {
           case '|':
               r0 = v1_0 || v2_0
               r1 = v1_1 || v2_1
           case '&':
               r0 = v1_0 && v2_0
               r1 = v1_1 && v2_1
           case '^':
               r0 = v1_0 != v2_0
               r1 = v1_1 != v2_1
           }
           val0 = append(val0, r0)
           val1 = append(val1, r1)
       }
   }
   // result is on top of stacks
   res0, res1 := false, false
   if len(val0) > 0 {
       res0 = val0[len(val0)-1]
       res1 = val1[len(val1)-1]
   }
   if res0 != res1 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
