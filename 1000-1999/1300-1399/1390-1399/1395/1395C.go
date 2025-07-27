package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &b[j])
   }
   // brute ans from 0 to 2^9-1
   const maxMask = 1<<9 - 1
   for ans := 0; ans <= maxMask; ans++ {
       ok := true
       for i := 0; i < n; i++ {
           found := false
           ai := a[i]
           for j := 0; j < m; j++ {
               ci := ai & b[j]
               // ci must not have bits outside ans
               if ci|ans == ans {
                   found = true
                   break
               }
           }
           if !found {
               ok = false
               break
           }
       }
       if ok {
           fmt.Println(ans)
           return
       }
   }
}
