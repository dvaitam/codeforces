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
   fmt.Fscan(reader, &t)
   for tt := 0; tt < t; tt++ {
       var x int64
       fmt.Fscan(reader, &x)
       var bin []int
       xx := x
       for xx > 0 {
           bin = append(bin, int(xx%2))
           xx /= 2
       }
       var ans []int64
       ans = append(ans, x)
       at := int64(1)
       sz := len(bin)
       for i := 0; i < sz-1; i++ {
           if bin[i] == 1 {
               x = x - at
               ans = append(ans, x)
           }
           at <<= 1
       }
       for x != 1 {
           x >>= 1
           ans = append(ans, x)
       }
       fmt.Fprintln(writer, len(ans))
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
