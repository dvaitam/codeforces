package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x, y, a, b int
   if _, err := fmt.Fscan(reader, &x, &y, &a, &b); err != nil {
       return
   }
   var res [][2]int
   for c := a; c <= x; c++ {
       for d := b; d <= y && d < c; d++ {
           res = append(res, [2]int{c, d})
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, len(res))
   for _, p := range res {
       fmt.Fprintln(writer, p[0], p[1])
   }
}
