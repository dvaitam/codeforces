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
   i, j, matched := 0, 0, 0
   for i < n && j < m {
       if b[j] >= a[i] {
           matched++
           i++
           j++
       } else {
           j++
       }
   }
   // need to come up with problems for unmatched requirements
   fmt.Println(n - matched)
}
