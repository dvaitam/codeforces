package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   names := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &names[i])
   }
   // read separator
   var sep string
   fmt.Fscan(reader, &sep)
   // group names by length
   groups := make(map[int][]string)
   lengths := make([]int, 0)
   for _, s := range names {
       l := len(s)
       if _, ok := groups[l]; !ok {
           lengths = append(lengths, l)
       }
       groups[l] = append(groups[l], s)
   }
   sort.Ints(lengths)
   // determine target sum of lengths
   minL := lengths[0]
   maxL := lengths[len(lengths)-1]
   K := minL + maxL
   // sort each group lexographically
   for _, l := range lengths {
       sort.Strings(groups[l])
   }
   // build pairs
   res := make([]string, 0, n/2)
   i, j := 0, len(lengths)-1
   for i <= j {
       li := lengths[i]
       lj := lengths[j]
       sum := li + lj
       if sum < K {
           i++
       } else if sum > K {
           j--
       } else {
           if i < j {
               A := groups[li]
               B := groups[lj]
               // assume len(A) == len(B)
               for idx := 0; idx < len(A); idx++ {
                   a, b := A[idx], B[idx]
                   // choose lex smaller order with separator
                   s1 := a + sep + b
                   s2 := b + sep + a
                   if s1 < s2 {
                       res = append(res, s1)
                   } else {
                       res = append(res, s2)
                   }
               }
           } else {
               A := groups[li]
               // li == lj: pair within group
               for idx := 0; idx+1 < len(A); idx += 2 {
                   a, b := A[idx], A[idx+1]
                   s1 := a + sep + b
                   s2 := b + sep + a
                   if s1 < s2 {
                       res = append(res, s1)
                   } else {
                       res = append(res, s2)
                   }
               }
           }
           i++
           j--
       }
   }
   // sort lines lex to get minimal concatenation
   sort.Strings(res)
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for _, line := range res {
       fmt.Fprintln(w, line)
   }
}
