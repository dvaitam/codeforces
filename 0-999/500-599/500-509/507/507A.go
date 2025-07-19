package main

import "fmt"

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   a := make([]int, n)
   idx := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
       idx[i] = i + 1
   }
   // sort by value using selection sort (n <= 100)
   for i := 0; i < n; i++ {
       minj := i
       for j := i + 1; j < n; j++ {
           if a[j] < a[minj] {
               minj = j
           }
       }
       a[i], a[minj] = a[minj], a[i]
       idx[i], idx[minj] = idx[minj], idx[i]
   }
   // pick as many as possible
   ans := 0
   for i := 0; i < n; i++ {
       if k >= a[i] {
           k -= a[i]
           ans++
       } else {
           break
       }
   }
   fmt.Println(ans)
   for i := 0; i < ans; i++ {
       fmt.Print(idx[i], " ")
   }
}
