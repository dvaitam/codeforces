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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   s := make([]byte, n)
   // read string, skip whitespace
   var tmp string
   fmt.Fscan(reader, &tmp)
   copy(s, tmp)

   // compute Z-array
   z := make([]int, n)
   z[0] = n
   l, r := 0, 0
   for i := 1; i < n; i++ {
       if i <= r {
           k0 := r - i + 1
           if z[i-l] < k0 {
               z[i] = z[i-l]
           } else {
               z[i] = k0
           }
       }
       for i+z[i] < n && s[z[i]] == s[i+z[i]] {
           z[i]++
       }
       if i+z[i]-1 > r {
           l = i
           r = i + z[i] - 1
       }
   }

   // diff array for prefix lengths 1..n
   diff := make([]int, n+2)
   for L := 1; L < n; L++ {
       // general case: C = A+B length L, need C repeat k times
       l64 := int64(k) * int64(L)
       if l64 <= int64(n) && int64(z[L]) >= l64 {
           // mark lengths from k*L to min(z[L]+L, (k+1)*L)
           r64 := int64(z[L]) + int64(L)
           maxR := int64(k+1) * int64(L)
           if r64 > maxR {
               r64 = maxR
           }
           if r64 > int64(n) {
               r64 = int64(n)
           }
           lpos := int(l64)
           rpos := int(r64)
           if lpos <= rpos {
               diff[lpos]++
               diff[rpos+1]--
           }
       }
       // A empty case: B only, need B repeat k times
       // pattern length L for B, need z[L] >= (k-1)*L
       if int64(k-1)*int64(L) <= int64(z[L]) {
           pos := k * L
           if pos <= n {
               diff[pos]++
               diff[pos+1]--
           }
       }
   }
   // build result
   res := make([]byte, n)
   cur := 0
   for i := 1; i <= n; i++ {
       cur += diff[i]
       if cur > 0 {
           res[i-1] = '1'
       } else {
           res[i-1] = '0'
       }
   }
   writer.Write(res)
}
