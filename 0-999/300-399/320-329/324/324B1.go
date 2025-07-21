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

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n+1)
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
       pos[a[i]] = i
   }

   var q int
   fmt.Fscan(reader, &q)
   for qi := 0; qi < q; qi++ {
       var typ, x, y int
       fmt.Fscan(reader, &typ, &x, &y)
       if typ == 1 {
           // query: minimum sessions to shave ids from x to y
           sessions := 1
           for id := x; id < y; id++ {
               if pos[id+1] < pos[id] {
                   sessions++
               }
           }
           fmt.Fprintln(writer, sessions)
       } else {
           // swap beavers at positions x and y
           v1 := a[x]
           v2 := a[y]
           // swap in a
           a[x], a[y] = v2, v1
           // update positions
           pos[v1], pos[v2] = y, x
       }
   }
}
