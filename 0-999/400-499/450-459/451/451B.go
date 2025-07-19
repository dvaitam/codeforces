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
   // find the first position where order breaks
   l := -1
   for i := 0; i < n-1; i++ {
       if a[i] > a[i+1] {
           l = i
           break
       }
   }
   // already sorted
   if l == -1 {
       fmt.Println("yes")
       fmt.Println(1, 1)
       return
   }
   // find the last position where order breaks
   r := -1
   for j := n - 1; j > 0; j-- {
       if a[j-1] > a[j] {
           r = j
           break
       }
   }
   // reverse the segment [l, r]
   for i, j := l, r; i < j; i, j = i+1, j-1 {
       a[i], a[j] = a[j], a[i]
   }
   // check if sorted
   for i := 0; i < n-1; i++ {
       if a[i] > a[i+1] {
           fmt.Println("no")
           return
       }
   }
   fmt.Println("yes")
   fmt.Println(l+1, r+1)
}
