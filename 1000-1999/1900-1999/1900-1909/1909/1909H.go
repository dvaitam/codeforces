package main

import (
   "bufio"
   "fmt"
   "os"
)

type op struct{ l, r int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // fast int reader
   readInt := func() int {
       var x int
       var c byte
       var neg bool
       for {
           b, err := reader.ReadByte()
           if err != nil {
               break
           }
           c = b
           if (c >= '0' && c <= '9') || c == '-' {
               break
           }
       }
       if c == '-' {
           neg = true
       } else {
           x = int(c - '0')
       }
       for {
           b, err := reader.ReadByte()
           if err != nil {
               break
           }
           c = b
           if c < '0' || c > '9' {
               break
           }
           x = x*10 + int(c-'0')
       }
       if neg {
           x = -x
       }
       return x
   }

   n := readInt()
   p := make([]int, n)
   id := make([]int, n)
   xl1 := make([]int, n)
   for i := 0; i < n; i++ {
       pi := readInt() - 1
       p[i] = pi
       id[pi] = i
       xl1[i] = i
   }
   // process to transform xl1 and p
   for i := 0; i < n; i++ {
       zz := i
       for p[p[zz]] != zz {
           a := p[zz]
           b := p[a]
           c := p[b]
           // swap id[a], id[c]
           id[a], id[c] = id[c], id[a]
           // swap xl1[zz], xl1[c]
           xl1[zz], xl1[c] = xl1[c], xl1[zz]
           // swap p[zz], p[c]
           p[zz], p[c] = p[c], p[zz]
           zz = id[zz]
       }
   }
   xl2 := make([]int, n)
   copy(xl2, p)

   var ans []op

   var sol func(a []int, opt bool)
   sol = func(a []int, opt bool) {
       dy := make([][]int, n)
       for i := 0; i < n; i++ {
           if a[i] != i {
               d1, d2 := i, a[i]
               if opt {
                   d1 = n - 1 - d1
                   d2 = n - 1 - d2
               }
               if d1 > d2 {
                   d1, d2 = d2, d1
               }
               var d3 int
               if (d2-d1)%2 != 0 {
                   d3 = (d2 - d1 + 1) / 2
               } else {
                   d3 = (2*n - d2 - d1) / 2
               }
               if d1%2 != 0 {
                   d3 = n + 1 - d3
               }
               d4 := d3
               if d1%2 == 0 {
                   d4 += d1
               } else {
                   d4 -= d1 + 1
               }
               // convert to 0-based indices for storage
               dy[d3-1] = append(dy[d3-1], d4-1)
               // swap a[i], a[a[i]]
               a[i], a[a[i]] = a[a[i]], a[i]
           }
       }
       for i := 0; i < n; i++ {
           d1 := i % 2
           d2 := n
           if (d2-d1)%2 != 0 {
               d2--
           }
           if d1 < d2 {
               ans = append(ans, op{l: d1, r: d2})
           }
           for _, dq := range dy[i] {
               ans = append(ans, op{l: dq, r: dq + 2})
           }
       }
   }

   sol(xl1, false)
   sol(xl2, true)

   // output
   fmt.Fprintln(writer, len(ans))
   for _, o := range ans {
       // print 1-based l and r
       fmt.Fprintln(writer, o.l+1, o.r)
   }
}
