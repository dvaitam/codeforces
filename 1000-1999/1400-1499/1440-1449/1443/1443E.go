package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, q int
   fmt.Fscan(in, &n, &q)
   // Suffix length: sufficient to cover all type2 additions
   maxSumX := int64(q) * 100000
   S := 0
   fact := []int64{1}
   for i := 1; i <= n; i++ {
       fact = append(fact, fact[i-1]*int64(i))
       if S == 0 && fact[i] > maxSumX {
           S = i
       }
       if i >= 20 { // safe cap
           break
       }
   }
   if S == 0 || S > n {
       S = n
   }
   // prepare factorials up to S
   fact = fact[:S+1]
   // prefix length stays identity
   prefLen := n - S
   // sorted suffix values
   base := make([]int, S)
   for i := 0; i < S; i++ {
       base[i] = prefLen + 1 + i
   }
   // current suffix and its prefix sums
   b := make([]int, S)
   copy(b, base)
   ps := make([]int64, S+1)
   for i := 0; i < S; i++ {
       ps[i+1] = ps[i] + int64(b[i])
   }
   var kcur int64

   for ; q > 0; q-- {
       var tp int
       fmt.Fscan(in, &tp)
       if tp == 1 {
           var l, r int
           fmt.Fscan(in, &l, &r)
           var ans int64
           if r <= prefLen {
               // sum in identity prefix
               cnt := int64(r - l + 1)
               ans = (int64(l) + int64(r)) * cnt / 2
           } else if l > prefLen {
               // fully in suffix
               l2 := l - prefLen
               r2 := r - prefLen
               ans = ps[r2] - ps[l2-1]
           } else {
               // split
               // prefix part
               cnt := int64(prefLen - l + 1)
               ans = (int64(l) + int64(prefLen)) * cnt / 2
               // suffix part
               r2 := r - prefLen
               ans += ps[r2]
           }
           fmt.Fprintln(out, ans)
       } else {
           var x int64
           fmt.Fscan(in, &x)
           kcur += x
           // rebuild suffix as kcur-th permutation of base
           available := make([]int, S)
           copy(available, base)
           kk := kcur
           for i := 0; i < S; i++ {
               f := fact[S-1-i]
               idx := int(kk / f)
               kk %= f
               b[i] = available[idx]
               // remove element
               available = append(available[:idx], available[idx+1:]...)
           }
           // rebuild prefix sums
           ps[0] = 0
           for i := 0; i < S; i++ {
               ps[i+1] = ps[i] + int64(b[i])
           }
       }
   }
}
