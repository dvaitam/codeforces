package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   b := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &b[j])
   }
   maxRatio := 0
   count := 0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if b[j] % a[i] == 0 {
               ratio := b[j] / a[i]
               if ratio > maxRatio {
                   maxRatio = ratio
                   count = 1
               } else if ratio == maxRatio {
                   count++
               }
           }
       }
   }
   fmt.Println(count)
}
