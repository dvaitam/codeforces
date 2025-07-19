package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type vl struct { a, b, c int }
type gs struct { ps, sg int }

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   x := make([]vl, 0, k)
   y := make([]vl, 0, k)
   for i := 0; i < k; i++ {
       var a, b, c, d int
       fmt.Fscan(reader, &a, &b, &c, &d)
       if a > c {
           a, c = c, a
       }
       if b > d {
           b, d = d, b
       }
       if b == d {
           x = append(x, vl{a, c, b})
       } else {
           y = append(y, vl{b, d, a})
       }
   }
   sort.Slice(x, func(i, j int) bool {
       if x[i].c != x[j].c {
           return x[i].c < x[j].c
       }
       return x[i].a < x[j].a
   })
   sort.Slice(y, func(i, j int) bool {
       if y[i].c != y[j].c {
           return y[i].c < y[j].c
       }
       return y[i].a < y[j].a
   })
   var sg int
   var a1, a2 int
   if n%2 != 0 {
       a1 = 0
   } else {
       a1 = m
   }
   if m%2 != 0 {
       a2 = 0
   } else {
       a2 = n
   }
   sg = a1 ^ a2

   var (
       px  []gs
       py  []gs
       t1  int
       t2  int
       lst int
       mi  int
       til int
       xx  int
       yy  int
   )
   // process vertical segments x
   lst = 0
   xx = 0
   for i := 0; i <= len(x); i++ {
       atEnd := i == len(x)
       cval := 0
       if !atEnd {
           cval = x[i].c
       }
       if (atEnd && lst != m-1) || (!atEnd && cval > lst+1) {
           xx = lst + 1
       }
       if atEnd || cval != lst {
           if i > 0 {
               sg ^= n - mi
               px = append(px, gs{ps: x[i-1].c, sg: n - mi})
               t1++
           }
           til = 0
           mi = 0
           if atEnd {
               break
           }
           sg ^= n
       }
       if !atEnd && x[i].b > til {
           mi += x[i].b - max(til, x[i].a)
           til = x[i].b
       }
       if !atEnd {
           lst = x[i].c
       }
   }
   // process horizontal segments y
   lst = 0
   yy = 0
   for i := 0; i <= len(y); i++ {
       atEnd := i == len(y)
       cval := 0
       if !atEnd {
           cval = y[i].c
       }
       if (atEnd && lst != n-1) || (!atEnd && cval > lst+1) {
           yy = lst + 1
       }
       if atEnd || cval != lst {
           if i > 0 {
               sg ^= m - mi
               py = append(py, gs{ps: y[i-1].c, sg: m - mi})
               t2++
           }
           til = 0
           mi = 0
           if atEnd {
               break
           }
           sg ^= m
       }
       if !atEnd && y[i].b > til {
           mi += y[i].b - max(til, y[i].a)
           til = y[i].b
       }
       if !atEnd {
           lst = y[i].c
       }
   }
   if xx != 0 {
       px = append(px, gs{ps: xx, sg: n})
       t1++
   }
   if yy != 0 {
       py = append(py, gs{ps: yy, sg: m})
       t2++
   }
   if sg == 0 {
       fmt.Fprintln(writer, "SECOND")
       return
   }
   fmt.Fprintln(writer, "FIRST")
   // find winning move
   var i int
   for i = 0; i < t1; i++ {
       if (sg ^ px[i].sg) < px[i].sg {
           break
       }
   }
   if i == t1 {
       // horizontal move
       var j int
       for j = 0; j < t2; j++ {
           if (sg ^ py[j].sg) < py[j].sg {
               break
           }
       }
       yy = py[j].ps
       lst = py[j].sg - (sg ^ py[j].sg)
       fmt.Fprintf(writer, "%d 0 ", yy)
       til = 0
       for ii := 0; ii < len(y); ii++ {
           if y[ii].c == yy {
               if y[ii].a > til {
                   if lst <= y[ii].a-til {
                       break
                   } else {
                       lst -= y[ii].a - til
                   }
               }
               if y[ii].b > til {
                   til = y[ii].b
               }
           }
       }
       fmt.Fprintf(writer, "%d %d\n", yy, til+lst)
   } else {
       // vertical move
       xx = px[i].ps
       lst = px[i].sg - (sg ^ px[i].sg)
       fmt.Fprintf(writer, "0 %d ", xx)
       til = 0
       for ii := 0; ii < len(x); ii++ {
           if x[ii].c == xx {
               if x[ii].a > til {
                   if lst <= x[ii].a-til {
                       break
                   } else {
                       lst -= x[ii].a - til
                   }
               }
               if x[ii].b > til {
                   til = x[ii].b
               }
           }
       }
       fmt.Fprintf(writer, "%d %d\n", til+lst, xx)
   }
}
