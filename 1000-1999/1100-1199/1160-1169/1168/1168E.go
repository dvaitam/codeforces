package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k int
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   n := 1 << k
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   xo := 0
   for _, v := range a {
       xo ^= v
   }
   if xo != 0 {
       fmt.Println("Fou")
       return
   }
   // Try trivial assignment p[i]=i
   p := make([]int, n)
   q := make([]int, n)
   used := make([]bool, n)
   ok := true
   for i := 0; i < n; i++ {
       p[i] = i
       q[i] = p[i] ^ a[i]
       if q[i] < 0 || q[i] >= n || used[q[i]] {
           ok = false
           break
       }
       used[q[i]] = true
   }
   if ok {
       fmt.Println("Shi")
       for i := 0; i < n; i++ {
           if i > 0 {
               fmt.Print(" ")
           }
           fmt.Print(p[i])
       }
       fmt.Println()
       for i := 0; i < n; i++ {
           if i > 0 {
               fmt.Print(" ")
           }
           fmt.Print(q[i])
       }
       fmt.Println()
       return
   }
   fmt.Println("Fou")
}
