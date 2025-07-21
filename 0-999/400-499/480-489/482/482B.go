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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // constraints
   type cons struct{ l, r, q int }
   consList := make([]cons, m)
   // diff array for bits
   const B = 30
   diff := make([][]int, B)
   for b := 0; b < B; b++ {
       diff[b] = make([]int, n+1)
   }
   for i := 0; i < m; i++ {
       var l, r, q int
       fmt.Fscan(reader, &l, &r, &q)
       consList[i] = cons{l - 1, r - 1, q}
       for b := 0; b < B; b++ {
           if (q>>b)&1 == 1 {
               diff[b][l-1]++
               diff[b][r]--
           }
       }
   }
   // build array
   a := make([]int, n)
   for b := 0; b < B; b++ {
       cnt := 0
       for i := 0; i < n; i++ {
           cnt += diff[b][i]
           if cnt > 0 {
               a[i] |= 1 << b
           }
       }
   }
   // segment tree for AND
   size := 1
   for size < n {
       size <<= 1
   }
   seg := make([]int, 2*size)
   // initialize leaves
   for i := 0; i < n; i++ {
       seg[size+i] = a[i]
   }
   // fill remaining leaves with all bits set
   allOnes := (1 << B) - 1
   for i := n; i < size; i++ {
       seg[size+i] = allOnes
   }
   // build internal nodes
   for i := size - 1; i > 0; i-- {
       seg[i] = seg[2*i] & seg[2*i+1]
   }
   // function to query AND on [l,r]
   var query func(l, r int) int
   query = func(l, r int) int {
       l += size
       r += size
       res := allOnes
       for l <= r {
           if l&1 == 1 {
               res &= seg[l]
               l++
           }
           if r&1 == 0 {
               res &= seg[r]
               r--
           }
           l >>= 1
           r >>= 1
       }
       return res
   }
   // verify constraints
   for _, c := range consList {
       if q := query(c.l, c.r); q != c.q {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   // output
   fmt.Fprintln(writer, "YES")
   for i, v := range a {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
