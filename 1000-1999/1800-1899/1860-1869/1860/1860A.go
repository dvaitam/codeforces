package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var s string
       fmt.Fscan(reader, &s)
       n := len(s)
       type1 := false
       for i := 1; i < n; i++ {
           if s[i-1] == ')' && s[i] == '(' {
               type1 = true
               break
           }
       }
       if type1 {
           writer.WriteString("YES\n")
           for i := 0; i < n; i++ {
               writer.WriteByte('(')
           }
           for i := 0; i < n; i++ {
               writer.WriteByte(')')
           }
           writer.WriteByte('\n')
           continue
       }
       type2 := false
       for i := 1; i < n; i++ {
           if s[i-1] == s[i] {
               type2 = true
               break
           }
       }
       if type2 {
           writer.WriteString("YES\n")
           for i := 0; i < n; i++ {
               writer.WriteString("()")
           }
           writer.WriteByte('\n')
           continue
       }
       writer.WriteString("NO\n")
   }
}
