package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   var m int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // Disjoint set union on positions 0..n-1
   parent := make([]int, n)
   for i := 0; i < n; i++ {
       parent[i] = i
   }
   var find func(int) int
   find = func(x int) int {
       if parent[x] != x {
           parent[x] = find(parent[x])
       }
       return parent[x]
   }
   union := func(a, b int) {
       ra, rb := find(a), find(b)
       if ra != rb {
           parent[rb] = ra
       }
   }
   // enforce palindrome constraint on each substring of length k
   half := k / 2
   for i := 0; i <= n-k; i++ {
       for j := 0; j < half; j++ {
           union(i+j, i+k-1-j)
       }
   }
   // count distinct components
   comps := 0
   for i := 0; i < n; i++ {
       if find(i) == i {
           comps++
       }
   }
   // fast exponentiation: m^comps mod
   const mod = 1000000007
   var res int64 = 1
   base := m % mod
   exp := comps
   for exp > 0 {
       if exp&1 == 1 {
           res = res * base % mod
       }
       base = base * base % mod
       exp >>= 1
   }
   // output result
   fmt.Fprintln(os.Stdout, res)
}
