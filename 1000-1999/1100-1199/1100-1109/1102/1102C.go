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

   var n, x, y int
   fmt.Fscan(reader, &n, &x, &y)
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }

   // If break power greater than repair, all doors can be broken
   if x > y {
       fmt.Fprintln(writer, n)
       return
   }

   // Count doors breakable in one move
   cnt := 0
   for _, v := range arr {
       if v <= x {
           cnt++
       }
   }

   // Doors broken is ceil(cnt/2)
   ans := (cnt + 1) / 2
   fmt.Fprintln(writer, ans)
}
