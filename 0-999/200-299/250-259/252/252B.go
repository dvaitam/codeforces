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
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   if n <= 2 {
       fmt.Println(-1)
       return
   }
   for i := 0; i < n-1; i++ {
       if a[i] != a[i+1] {
           a[i], a[i+1] = a[i+1], a[i]
           if !isSorted(a) {
               fmt.Println(i+1, i+2)
               return
           }
           a[i], a[i+1] = a[i+1], a[i]
       }
   }
   fmt.Println(-1)
}

func isSorted(a []int) bool {
   asc, desc := true, true
   for i := 1; i < len(a); i++ {
       if a[i] < a[i-1] {
           asc = false
       }
       if a[i] > a[i-1] {
           desc = false
       }
       if !asc && !desc {
           return false
       }
   }
   return true
}
