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
       var n, k int
       fmt.Fscan(reader, &n, &k)
       if n <= 2*k {
           fmt.Fprintln(writer, -1)
           continue
       }
       res := make([]int, 0, n)
       i, j := 1, n
       for x := 0; x < k; x++ {
           res = append(res, i)
           res = append(res, j)
           i++
           j--
       }
       for l := i; l <= j; l++ {
           res = append(res, l)
       }
       for idx, v := range res {
           if idx > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
