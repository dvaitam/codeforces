package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   b := make([]int64, n)
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &c[i])
   }
   var sumBC int64
   for i := 0; i < n; i++ {
       sumBC += b[i] + c[i]
   }
   denom := int64(2 * n)
   if sumBC%denom != 0 {
       fmt.Println(-1)
       return
   }
   S := sumBC / denom
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       tmp := b[i] + c[i] - S
       if tmp%int64(n) != 0 {
           fmt.Println(-1)
           return
       }
       ai := tmp / int64(n)
       if ai < 0 {
           fmt.Println(-1)
           return
       }
       a[i] = ai
   }
   const maxbit = 31
   cnt := make([]int64, maxbit+1)
   for _, ai := range a {
       for k := 0; k <= maxbit; k++ {
           if (ai>>k)&1 == 1 {
               cnt[k]++
           }
       }
   }
   // verify b and c
   for i, ai := range a {
       var bi2, ci2 int64
       for k := 0; k <= maxbit; k++ {
           bit := int64(1) << k
           if (ai>>k)&1 == 1 {
               bi2 += cnt[k] * bit
               ci2 += int64(n) * bit
           } else {
               ci2 += cnt[k] * bit
           }
       }
       if bi2 != b[i] || ci2 != c[i] {
           fmt.Println(-1)
           return
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i, ai := range a {
       if i > 0 {
           fmt.Fprint(w, " ")
       }
       fmt.Fprint(w, ai)
   }
   fmt.Fprintln(w)
}
