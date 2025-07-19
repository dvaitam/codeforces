package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       v := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &v[i])
       }
       c := make([]int, n+1)
       for _, x := range v {
           c[x]++
       }
       // collect missing values
       pq := make([]int, 0, n)
       for i := 1; i <= n; i++ {
           if c[i] == 0 {
               pq = append(pq, i)
           }
       }
       sort.Ints(pq)
       pair := make([]int, n+1)
       f := false
       // pair values from largest to smallest
       for i := n; i >= 1; i-- {
           switch c[i] {
           case 1:
               pair[i] = i
           case 2:
               if len(pq) == 0 {
                   f = true
                   break
               }
               m := pq[len(pq)-1]
               if m < i {
                   // pair missing m with i
                   pq = pq[:len(pq)-1]
                   pair[i] = m
                   pair[m] = i
               } else {
                   f = true
               }
           default:
               // c[i] == 0 handled in pq; >2 invalid
               if c[i] > 2 {
                   f = true
               }
           }
           if f {
               break
           }
       }
       if f {
           fmt.Fprintln(writer, "NO")
           continue
       }
       fmt.Fprintln(writer, "YES")
       // build permutations
       vis := make([]bool, n+1)
       p := make([]int, n)
       q := make([]int, n)
       for i, x := range v {
           if !vis[x] {
               p[i] = x
               q[i] = pair[x]
               vis[x] = true
           } else {
               p[i] = pair[x]
               q[i] = x
           }
       }
       // print p
       for i := 0; i < n; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, p[i])
       }
       writer.WriteByte('\n')
       // print q
       for i := 0; i < n; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, q[i])
       }
       writer.WriteByte('\n')
   }
}
