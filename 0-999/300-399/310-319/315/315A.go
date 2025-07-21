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
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
   }
   count := 0
   for i := 0; i < n; i++ {
       canOpen := false
       for j := 0; j < n; j++ {
           if i != j && b[j] == a[i] {
               canOpen = true
               break
           }
       }
       if !canOpen {
           count++
       }
   }
   fmt.Println(count)
}
