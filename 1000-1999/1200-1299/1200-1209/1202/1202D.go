package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var N int
   if _, err := fmt.Fscan(in, &N); err != nil {
       return
   }
   for i := 0; i < N; i++ {
       var a int
       fmt.Fscan(in, &a)
       mil := a / 1122457
       a %= 1122457
       thou := a / 1000
       one := a % 1000

       // prefix
       out.WriteString("133")
       // ones place as '7'
       for j := 0; j < one; j++ {
           out.WriteByte('7')
       }
       // thousands place
       if thou > 0 {
           for j := 0; j < 994; j++ {
               out.WriteByte('1')
           }
           out.WriteString("33")
           for j := 0; j < thou; j++ {
               out.WriteByte('7')
           }
       }
       // millions place
       if mil > 0 {
           for j := 0; j < 46; j++ {
               out.WriteByte('3')
           }
           for j := 0; j < mil; j++ {
               out.WriteByte('7')
           }
       }
       out.WriteByte('\n')
   }
}
