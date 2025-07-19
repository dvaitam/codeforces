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
   var t int
   fmt.Fscan(reader, &t)
   // precompute powers
   const K = 51
   p := make([]int64, K)
   p[1], p[2] = 1, 1
   for i := 3; i < K; i++ {
       p[i] = p[i-1] * 2
   }
   for ; t > 0; t-- {
       var a, b, m int64
       fmt.Fscan(reader, &a, &b, &m)
       if a == b {
           fmt.Fprintf(writer, "1 %d\n", a)
           continue
       }
       found := false
       for k := 2; k < K; k++ {
           // minimal and maximal possible last term: p[k]*a + [p[k], p[k]*m]
           if p[k]*a > b {
               break
           }
           r, ok := tryMake(a, b, m, k, p)
           if !ok {
               continue
           }
           // build sequence
           fmt.Fprintf(writer, "%d ", k)
           seq := make([]int64, k)
           x := a
           var sum int64
           for i := 0; i < k; i++ {
               seq[i] = x
               sum += x
               // r is 1-based, r[1..k], r[k] == 0
               x = sum + r[i+1]
           }
           for i, v := range seq {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprintf(writer, "%d", v)
           }
           writer.WriteByte('\n')
           found = true
           break
       }
       if !found {
           fmt.Fprintln(writer, -1)
       }
   }
}

// tryMake attempts to find r[1..k] (1-based) for length k sequence
// returns slice r of length k+1, r[0] unused, r[k] should be 0, and ok
func tryMake(a, b, m int64, k int, p []int64) ([]int64, bool) {
   // reduce b by p[k]*a
   b2 := b - p[k]*a
   if b2 < p[k] || (b2+p[k]-1)/p[k] > m {
       return nil, false
   }
   r := make([]int64, k+1)
   // ensure at least 1 for r[1..k-1]
   for i := 1; i < k; i++ {
       r[i] = 1
       b2 -= p[i]
   }
   // distribute remaining b2
   for i := 1; i < k; i++ {
       if b2 <= 0 {
           break
       }
       // p[k-i]
       pi := p[k-i]
       maxAdd := (b2 / pi)
       if maxAdd > m-1 {
           maxAdd = m-1
       }
       r[i] += maxAdd
       b2 -= pi * maxAdd
   }
   if b2 != 0 {
       return nil, false
   }
   return r, true
}
