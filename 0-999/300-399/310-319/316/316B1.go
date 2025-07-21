package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, x int
   if _, err := fmt.Fscan(reader, &n, &x); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // build next pointers
   next := make([]int, n+1)
   for i := 1; i <= n; i++ {
       if a[i] != 0 {
           next[a[i]] = i
       }
   }
   // find chains and build sizes of other chains, and pos in x's chain
   var other []int
   sumOther := 0
   posInChain := 0
   // identify heads
   for i := 1; i <= n; i++ {
       if a[i] != 0 {
           continue
       }
       // traverse this chain
       cur := i
       size := 0
       containsX := false
       for cur != 0 {
           if cur == x {
               containsX = true
               posInChain = size
           }
           size++
           cur = next[cur]
       }
       if !containsX {
           other = append(other, size)
           sumOther += size
       }
   }
   // subset sums dp
   dp := make([]bool, sumOther+1)
   dp[0] = true
   for _, sz := range other {
       for s := sumOther; s >= sz; s-- {
           if dp[s-sz] {
               dp[s] = true
           }
       }
   }
   // collect positions
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   first := true
   for s := 0; s <= sumOther; s++ {
       if dp[s] {
           pos := s + posInChain + 1
           if !first {
               writer.WriteByte(' ')
           }
           first = false
           fmt.Fprint(writer, pos)
       }
   }
   fmt.Fprintln(writer)
}
