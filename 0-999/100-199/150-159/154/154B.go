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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   spf := make([]int, n+1)
   for i := 2; i <= n; i++ {
       if spf[i] == 0 {
           for j := i; j <= n; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   owner := make([]int, n+1)
   on := make([]bool, n+1)

   for k := 0; k < m; k++ {
       var op string
       var x int
       fmt.Fscan(reader, &op, &x)
       if op == "+" {
           if on[x] {
               writer.WriteString("Already on\n")
               continue
           }
           conflict := false
           conflictWith := 0
           j := x
           for j > 1 {
               p := spf[j]
               if owner[p] != 0 {
                   conflict = true
                   conflictWith = owner[p]
                   break
               }
               for j%p == 0 {
                   j /= p
               }
           }
           if conflict {
               writer.WriteString(fmt.Sprintf("Conflict with %d\n", conflictWith))
           } else {
               writer.WriteString("Success\n")
               on[x] = true
               j = x
               for j > 1 {
                   p := spf[j]
                   owner[p] = x
                   for j%p == 0 {
                       j /= p
                   }
               }
           }
       } else {
           // op == "-"
           if !on[x] {
               writer.WriteString("Already off\n")
               continue
           }
           // turn off
           j := x
           for j > 1 {
               p := spf[j]
               owner[p] = 0
               for j%p == 0 {
                   j /= p
               }
           }
           on[x] = false
           writer.WriteString("Success\n")
       }
   }
}
