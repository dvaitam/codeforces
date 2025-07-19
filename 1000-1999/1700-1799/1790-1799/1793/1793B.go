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

   var tt int
   if _, err := fmt.Fscan(in, &tt); err != nil {
       return
   }
   for tt > 0 {
       tt--
       var x, y int
       fmt.Fscan(in, &x, &y)
       n := 2 * abs(x-y)
       fmt.Fprintln(out, n)
       // forward from y to x inclusive
       for i := y; i <= x; i++ {
           fmt.Fprint(out, i, " ")
       }
       // backward from x-1 down to y+1
       for i := x - 1; i > y; i-- {
           fmt.Fprint(out, i, " ")
       }
       fmt.Fprintln(out)
   }
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}
