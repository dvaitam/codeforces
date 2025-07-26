package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (
   in  = bufio.NewReader(os.Stdin)
   out = bufio.NewWriter(os.Stdout)
)

func main() {
   defer out.Flush()
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       solve()
   }
}

func solve() {
   var n, m int
   fmt.Fscan(in, &n, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   b := make([]bool, n-1)
   for i := 0; i < m; i++ {
       var p int
       fmt.Fscan(in, &p)
       if p >= 1 && p < n {
           b[p-1] = true
       }
   }
   // For each consecutive segment of allowed swaps, sort the subarray
   i := 0
   for i < n-1 {
       if !b[i] {
           i++
           continue
       }
       l := i
       for i < n-1 && b[i] {
           i++
       }
       // segment covers elements from l to i (inclusive)
       sort.Ints(a[l : i+1])
   }
   // Check if sorted
   ok := true
   for i := 0; i < n-1; i++ {
       if a[i] > a[i+1] {
           ok = false
           break
       }
   }
   if ok {
       fmt.Fprintln(out, "YES")
   } else {
       fmt.Fprintln(out, "NO")
   }
}
