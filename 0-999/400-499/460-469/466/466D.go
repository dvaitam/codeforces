package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, h int
   if _, err := fmt.Fscan(reader, &n, &h); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // compute required increments d[i] = h - a[i]
   d := make([]int, n+2)
   for i := 1; i <= n; i++ {
       if a[i-1] > h {
           fmt.Println(0)
           return
       }
       d[i] = h - a[i-1]
   }
   const mod = 1000000007
   // boundary
   d[0], d[n+1] = 0, 0
   ans := int64(1)
   // transitions from d[i] to d[i+1] for i = 0..n
   for i := 0; i <= n; i++ {
       cur := d[i]
       next := d[i+1]
       delta := next - cur
       switch {
       case delta == 1:
           // start a new segment
       case delta == 0:
           // either do nothing or start+end: contributes (cur+1) ways
           ans = ans * int64(cur+1) % mod
       case delta == -1:
           // must end one of cur open segments
           ans = ans * int64(cur) % mod
       default:
           fmt.Println(0)
           return
       }
   }
   fmt.Println(ans)
}
