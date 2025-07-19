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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   parent := make([]int, n+1)
   for i := 2; i <= n; i++ {
       fmt.Fscan(reader, &parent[i])
   }
   const inf = int64(1e18)
   s := make([]int64, n+1)
   orig := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &orig[i])
   }
   // Initialize s and propagate minima to parents
   for i := 1; i <= n; i++ {
       if orig[i] < 0 {
           s[i] = inf
       } else {
           s[i] = orig[i]
           if i > 1 && s[i] < s[parent[i]] {
               s[parent[i]] = s[i]
           }
       }
   }
   var result int64
   // Compute answer by differences
   for i := n; i >= 1; i-- {
       if s[i] == inf {
           if i > 1 {
               s[i] = s[parent[i]]
           } else {
               s[i] = 0
           }
       }
       var pval int64
       if i > 1 {
           pval = s[parent[i]]
       }
       delta := s[i] - pval
       if delta < 0 {
           fmt.Fprintln(writer, -1)
           return
       }
       result += delta
   }
   fmt.Fprintln(writer, result)
