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
       var n int
       var k int64
       fmt.Fscan(reader, &n, &k)
       var a int64
       var b int64
       found := false
       for ai := int64(0); ai <= int64(n); ai++ {
           bi := int64(n) - ai
           cur := ai*(ai-1)/2 + bi*(bi-1)/2
           if cur == k {
               a = ai
               b = bi
               found = true
               break
           }
       }
       if found {
           fmt.Fprintln(writer, "YES")
           // print a ones
           for i := int64(0); i < a; i++ {
               writer.WriteString("1 ")
           }
           // print b minus ones
           for i := int64(0); i < b; i++ {
               writer.WriteString("-1 ")
           }
           writer.WriteByte('\n')
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
