package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   totalStress := 0
   removed := make([]int, k+1)
   added := make([]int, k+1)
   // count original stresses and removals
   for i := 1; i < n; i++ {
       if a[i] != a[i-1] {
           totalStress++
           removed[a[i]]++
           removed[a[i-1]]++
       }
   }
   // count added stresses when removing segments
   for i := 0; i < n; {
       j := i
       for j+1 < n && a[j+1] == a[i] {
           j++
       }
       // segment [i..j] of genre g
       g := a[i]
       if i > 0 && j < n-1 {
           left := a[i-1]
           right := a[j+1]
           if left != right {
               added[g]++
           }
       }
       i = j + 1
   }
   // find best genre to remove
   best := 1
   minStress := totalStress - removed[1] + added[1]
   for x := 2; x <= k; x++ {
       stress := totalStress - removed[x] + added[x]
       if stress < minStress {
           minStress = stress
           best = x
       }
   }
   fmt.Println(best)
}
