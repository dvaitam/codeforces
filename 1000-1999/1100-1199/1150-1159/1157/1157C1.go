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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   l, r := 0, n-1
   now := 0
   var res []byte
   for l <= r {
       if a[l] <= now && a[r] <= now {
           break
       }
       if a[l] < a[r] {
           if now < a[l] {
               res = append(res, 'L')
               now = a[l]
               l++
           } else {
               res = append(res, 'R')
               now = a[r]
               r--
           }
       } else {
           if now < a[r] {
               res = append(res, 'R')
               now = a[r]
               r--
           } else {
               res = append(res, 'L')
               now = a[l]
               l++
           }
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(res))
   w.Write(res)
   fmt.Fprintln(w)
}
