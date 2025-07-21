package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   // find first position where p[i] != i
   l := 1
   for l <= n && p[l] == l {
       l++
   }
   if l > n {
       // already identity, no valid single reversal
       fmt.Println(0, 0)
       return
   }
   // find r such that p[r] == l
   r := n
   for r >= l && p[r] != l {
       r--
   }
   if r <= l {
       fmt.Println(0, 0)
       return
   }
   // check segment [l, r] is reversed identity
   ok := true
   for i := l; i <= r; i++ {
       if p[i] != (l + r - i) {
           ok = false
           break
       }
   }
   // check outside segment
   if ok {
       for i := 1; i < l; i++ {
           if p[i] != i {
               ok = false
               break
           }
       }
   }
   if ok {
       for i := r + 1; i <= n; i++ {
           if p[i] != i {
               ok = false
               break
           }
       }
   }
   if !ok {
       fmt.Println(0, 0)
   } else {
       fmt.Println(l, r)
   }
}
