package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // collect positions of ones
   pos := make([]int, 0, n)
   for i, v := range a {
       if v == 1 {
           pos = append(pos, i)
       }
   }
   m := len(pos)
   if m%2 != 0 {
       fmt.Println("NO")
       return
   }
   ops := make([][3]int, 0, m)
   // pair ones symmetrically
   for i := 0; i < m/2; i++ {
       l := pos[i]
       r := pos[m-1-i]
       d := r - l
       if d < 2 {
           // too close to eliminate
           fmt.Println("NO")
           return
       }
       if d%2 == 0 {
           // single operation
           mid := (l + r) / 2
           ops = append(ops, [3]int{l + 1, mid + 1, r + 1})
       } else {
           // two operations
           mid1 := (l + r - 1) / 2
           mid2 := mid1 + 1
           if mid2 > r || mid1 <= l {
               fmt.Println("NO")
               return
           }
           ops = append(ops, [3]int{l + 1, mid1 + 1, mid2 + 1})
           ops = append(ops, [3]int{mid1 + 1, mid2 + 1, r + 1})
       }
   }
   fmt.Println("YES")
   fmt.Println(len(ops))
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for _, t := range ops {
       fmt.Fprintf(w, "%d %d %d\n", t[0], t[1], t[2])
   }
}
