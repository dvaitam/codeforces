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
   const BITS = 10
   // ans0[b][i]: answer for bit b querying group of indices with bit b=0
   // ans1[b][i]: answer for bit b querying group of indices with bit b=1
   ans0 := make([][]int, BITS)
   ans1 := make([][]int, BITS)
   // perform queries
   for b := 0; b < BITS; b++ {
       // group0: indices with bit b = 0
       // group1: indices with bit b = 1
       grp0 := make([]int, 0, n)
       grp1 := make([]int, 0, n)
       for i := 1; i <= n; i++ {
           if (i>>b)&1 == 1 {
               grp1 = append(grp1, i)
           } else {
               grp0 = append(grp0, i)
           }
       }
       // query group0
       if len(grp0) > 0 {
           fmt.Fprintf(writer, "? %d", len(grp0))
           for _, v := range grp0 {
               fmt.Fprintf(writer, " %d", v)
           }
           fmt.Fprintln(writer)
           writer.Flush()
           // read answers
           ans := make([]int, n+1)
           for i := 1; i <= n; i++ {
               fmt.Fscan(reader, &ans[i])
           }
           ans0[b] = ans
       } else {
           // no query, fill with large
           ans0[b] = make([]int, n+1)
           for i := 1; i <= n; i++ {
               ans0[b][i] = 2000000001
           }
       }
       // query group1
       if len(grp1) > 0 {
           fmt.Fprintf(writer, "? %d", len(grp1))
           for _, v := range grp1 {
               fmt.Fprintf(writer, " %d", v)
           }
           fmt.Fprintln(writer)
           writer.Flush()
           ans := make([]int, n+1)
           for i := 1; i <= n; i++ {
               fmt.Fscan(reader, &ans[i])
           }
           ans1[b] = ans
       } else {
           ans1[b] = make([]int, n+1)
           for i := 1; i <= n; i++ {
               ans1[b][i] = 2000000001
           }
       }
   }
   // compute final answers
   res := make([]int, n+1)
   for i := 1; i <= n; i++ {
       best := 2000000001
       for b := 0; b < BITS; b++ {
           if (i>>b)&1 == 1 {
               if ans0[b][i] < best {
                   best = ans0[b][i]
               }
           } else {
               if ans1[b][i] < best {
                   best = ans1[b][i]
               }
           }
       }
       res[i] = best
   }
   // output final
   fmt.Fprintln(writer, -1)
   for i := 1; i <= n; i++ {
       if i > 1 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, res[i])
   }
   fmt.Fprintln(writer)
   writer.Flush()
}
