package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   N := 1 << uint(n)
   a := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // inv0[k]: inversions at level k when not reversed
   inv0 := make([]int64, n+1)
   tmp := make([]int, N)
   // Compute inv0 by bottom-up merge sort
   for k := 1; k <= n; k++ {
       width := 1 << uint(k-1)
       block := width << 1
       var cnt int64
       for i := 0; i < N; i += block {
           l, r := i, i+width
           endL, endR := r, i+block
           if endR > N {
               endR = N
           }
           ti := i
           for l < endL && r < endR {
               if a[l] <= a[r] {
                   tmp[ti] = a[l]
                   l++
               } else {
                   tmp[ti] = a[r]
                   cnt += int64(endL - l)
                   r++
               }
               ti++
           }
           for l < endL {
               tmp[ti] = a[l]
               l++
               ti++
           }
           for r < endR {
               tmp[ti] = a[r]
               r++
               ti++
           }
           // copy back
           for j := i; j < endR; j++ {
               a[j] = tmp[j]
           }
       }
       inv0[k] = cnt
   }
   // total pairs at each level
   total := make([]int64, n+1)
   for k := 1; k <= n; k++ {
       // total pairs: 2^(n+k-2)
       total[k] = int64(1) << uint(n+k-2)
   }
   // current flips and current inversion count
   flip := make([]bool, n+1)
   var curInv int64
   for k := 1; k <= n; k++ {
       curInv += inv0[k]
   }
   // process queries
   var m int
   fmt.Fscan(reader, &m)
   for qi := 0; qi < m; qi++ {
       var q int
       fmt.Fscan(reader, &q)
       // flip levels 1..q
       for k := 1; k <= q; k++ {
           if flip[k] {
               // turning off: subtract (total[k] - 2*inv0[k])
               curInv -= (total[k] - 2*inv0[k])
           } else {
               // turning on: add (total[k] - 2*inv0[k])
               curInv += (total[k] - 2*inv0[k])
           }
           flip[k] = !flip[k]
       }
       // output
       writer.WriteString(strconv.FormatInt(curInv, 10))
       writer.WriteByte('\n')
   }
}
