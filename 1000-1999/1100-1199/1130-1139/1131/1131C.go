package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Ints(a)

   type pair struct { val, cnt int }
   var b []pair
   cnt := 1
   for i := 1; i < n; i++ {
       if a[i] == a[i-1] {
           cnt++
       } else {
           b = append(b, pair{a[i-1], cnt})
           cnt = 1
       }
   }
   if n > 0 {
       b = append(b, pair{a[n-1], cnt})
   }
   bb := len(b)

   var l, r []int
   flag := 0
   for i := bb - 2; i >= 0; i-- {
       x := b[i].val
       c := b[i].cnt
       half := c / 2
       for j := 0; j < half; j++ {
           l = append(l, x)
           r = append(r, x)
       }
       if c % 2 == 1 {
           if flag == 0 {
               l = append(l, x)
               flag = 1
           } else {
               r = append(r, x)
               flag = 0
           }
       }
   }

   var res []int
   // left part in reverse
   for i := len(l) - 1; i >= 0; i-- {
       res = append(res, l[i])
   }
   // middle: all of the largest element
   if bb > 0 {
       last := b[bb-1]
       for i := 0; i < last.cnt; i++ {
           res = append(res, last.val)
       }
   }
   // right part
   for _, v := range r {
       res = append(res, v)
   }

   for i, v := range res {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
