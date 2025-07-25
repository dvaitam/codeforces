package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   h1 := make([]int64, n)
   h2 := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &h1[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &h2[i])
   }

   var pmax0, pmax1 int64 // pmax0: max dp ending in row1 so far; pmax1: row2
   for i := 0; i < n; i++ {
       dp0 := h1[i] + max(pmax1, 0)
       dp1 := h2[i] + max(pmax0, 0)
       pmax0 = max(pmax0, dp0)
       pmax1 = max(pmax1, dp1)
   }
   ans := max(pmax0, pmax1)
   fmt.Fprintln(writer, ans)
}
