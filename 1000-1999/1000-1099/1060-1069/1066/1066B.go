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
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, r int
   if _, err := fmt.Fscan(reader, &n, &r); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   cnt := 0
   now := 1
   for now <= n {
       // furthest heater to cover position now
       temp := min(now+r-1, n)
       lower := now - r + 1
       if lower < 1 {
           lower = 1
       }
       found := false
       for pos := temp; pos >= lower; pos-- {
           if a[pos] == 1 {
               temp = pos
               found = true
               break
           }
       }
       if !found {
           fmt.Fprintln(writer, -1)
           return
       }
       cnt++
       now = temp + r
   }
   fmt.Fprintln(writer, cnt)
}
