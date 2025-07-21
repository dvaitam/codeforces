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
   var n int
   var d int64 // d is unused in this generator
   var x int64
   if _, err := fmt.Fscan(in, &n, &d, &x); err != nil {
       return
   }
   const mod = 1000000007
   // generate permutation a of 1..n
   a := make([]int, n)
   for i := 0; i < n; i++ {
       a[i] = i + 1
   }
   for i := 0; i < n; i++ {
       x = (x*7 + 13) % mod
       j := int(x % int64(i+1))
       a[i], a[j] = a[j], a[i]
   }
   // generate binary mask b
   b := make([]bool, n)
   for i := 0; i < n; i++ {
       x = (x*7 + 13) % mod
       b[i] = (x & 1) == 1
   }
   // collect offsets where b[k] == 1
   offs := make([]int, 0, n)
   for k, v := range b {
       if v {
           offs = append(offs, k)
       }
   }
   // DSU skip for assigning c
   parent := make([]int, n+1)
   for i := 0; i <= n; i++ {
       parent[i] = i
   }
   var find func(int) int
   find = func(u int) int {
       if parent[u] != u {
           parent[u] = find(parent[u])
       }
       return parent[u]
   }
   // map value to position
   pos := make([]int, n+1)
   for i, v := range a {
       pos[v] = i
   }
   // result c
   c := make([]int, n)
   // assign in descending order
   for v := n; v >= 1; v-- {
       p := pos[v]
       for _, k := range offs {
           t := p + k
           if t >= n {
               break
           }
           u := find(t)
           if u != t {
               continue
           }
           c[t] = v
           parent[t] = t + 1
       }
   }
   // output c[0..n-1]
   for i := 0; i < n; i++ {
       fmt.Fprintln(out, c[i])
   }
}
