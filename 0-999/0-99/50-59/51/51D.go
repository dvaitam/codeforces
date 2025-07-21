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
   // Check if already geometric progression
   if isGP(a, -1) {
       fmt.Println(0)
       return
   }
   // Collect candidate indices for deletion
   cand := make(map[int]bool)
   haveCand := false
   // Special initial case: first is zero but second non-zero
   if n >= 2 && a[0] == 0 && a[1] != 0 {
       cand[0] = true
       cand[1] = true
       haveCand = true
   }
   // Check triples for general failures
   for i := 2; i < n; i++ {
       x0, x1, x2 := a[i-2], a[i-1], a[i]
       fail := false
       if x0 == 0 {
           if x1 != 0 {
               fail = true
           }
       } else {
           if int64(x1)*int64(x1) != int64(x0)*int64(x2) {
               fail = true
           }
       }
       if fail {
           // indices that could be removed
           tri := []int{i - 2, i - 1, i}
           if !haveCand {
               for _, j := range tri {
                   if j >= 0 && j < n {
                       cand[j] = true
                   }
               }
               haveCand = true
           } else {
               newCand := make(map[int]bool)
               for _, j := range tri {
                   if cand[j] {
                       newCand[j] = true
                   }
               }
               cand = newCand
           }
           if len(cand) == 0 {
               fmt.Println(2)
               return
           }
       }
   }
   // Try deleting each candidate index
   for j := range cand {
       if isGP(a, j) {
           fmt.Println(1)
           return
       }
   }
   // Impossible
   fmt.Println(2)
}

// isGP checks if sequence a without element at skip is a geometric progression
func isGP(a []int, skip int) bool {
   n := len(a)
   k := 0
   var prevprev, prev int
   for i := 0; i < n; i++ {
       if i == skip {
           continue
       }
       x := a[i]
       if k == 0 {
           prevprev = x
       } else if k == 1 {
           prev = x
           if prevprev == 0 && prev != 0 {
               return false
           }
       } else {
           if prevprev == 0 {
               if x != 0 {
                   return false
               }
           } else {
               if int64(prev)*int64(prev) != int64(prevprev)*int64(x) {
                   return false
               }
           }
           prevprev = prev
           prev = x
       }
       k++
   }
   return true
}
