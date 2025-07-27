package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = int(1e9 + 7)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // prev and next occurrence
   prev := make([]int, n+1)
   next := make([]int, n+1)
   last := make(map[int]int)
   for i := 1; i <= n; i++ {
       if p, ok := last[a[i]]; ok {
           prev[i] = p
       } else {
           prev[i] = 0
       }
       last[a[i]] = i
   }
   // reset map for next
   last = make(map[int]int)
   for i := n; i >= 1; i-- {
       if p, ok := last[a[i]]; ok {
           next[i] = p
       } else {
           next[i] = n + 1
       }
       last[a[i]] = i
   }

   ans := 0
   for x := 1; x <= n; x++ {
       curPl := 0
       curNr := n + 1
       for y := x; y <= n; y++ {
           if prev[y] >= x {
               break
           }
           if prev[y] > curPl {
               curPl = prev[y]
           }
           if next[y] < curNr {
               curNr = next[y]
           }
           left := x - curPl - 1
           right := curNr - y - 1
           if left <= 0 || right <= 0 {
               continue
           }
           ans = (ans + left*right) % MOD
       }
   }
   fmt.Fprintln(writer, ans)
}
