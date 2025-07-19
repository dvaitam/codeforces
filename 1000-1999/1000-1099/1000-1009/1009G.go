package main

import "fmt"

func main() {

   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   n := len(s)
   const A = 6
   maxMask := 1 << A

   // count occurrences of each character
   cou := make([]int, A)
   for i := 0; i < n; i++ {
       c := s[i] - 'a'
       if c >= 0 && c < A {
           cou[c]++
       }
   }

   // initial left counts per mask
   leftArr := make([]int, maxMask)
   for mask := 0; mask < maxMask; mask++ {
       sum := 0
       for j := 0; j < A; j++ {
           if mask&(1<<j) != 0 {
               sum += cou[j]
           }
       }
       leftArr[mask] = sum
   }

   // allowed masks per position, default all letters allowed
   allowed := make([]int, n)
   fullMask := maxMask - 1
   for i := range allowed {
       allowed[i] = fullMask
   }

   // read constraints
   var m int
   fmt.Scan(&m)
   for k := 0; k < m; k++ {
       var x int
       var t string
       fmt.Scan(&x, &t)
       x-- // zero-based index
       mask := 0
       for i := 0; i < len(t); i++ {
           c := t[i] - 'a'
           if c >= 0 && c < A {
               mask |= 1 << c
           }
       }
       if x >= 0 && x < n {
           allowed[x] = mask
       }
   }

   // compute matching counts per mask
   matching := make([]int, maxMask)
   for i := 0; i < n; i++ {
       a := allowed[i]
       for mask := 0; mask < maxMask; mask++ {
           if mask&a != 0 {
               matching[mask]++
           }
       }
   }

   // feasibility check
   for mask := 0; mask < maxMask; mask++ {
       if leftArr[mask] > matching[mask] {
           fmt.Println("Impossible")
           return
       }
   }

   // build answer greedily
   ans := make([]byte, n)
   for i := 0; i < n; i++ {
       poss := allowed[i]
       for mask := 1; mask < maxMask; mask++ {
           if leftArr[mask] == matching[mask] && (mask&allowed[i]) != 0 {
               poss &= mask
           }
       }
       if poss == 0 {
           fmt.Println("Impossible")
           return
       }
       // pick lowest set bit
       take := 0
       for j := 0; j < A; j++ {
           if poss&(1<<j) != 0 {
               take = j
               break
           }
       }
       ans[i] = byte('a' + take)
       // update counts
       for mask := 0; mask < maxMask; mask++ {
           if (mask>>take)&1 != 0 {
               leftArr[mask]--
           }
           if mask&allowed[i] != 0 {
               matching[mask]--
           }
       }
   }

   fmt.Println(string(ans))
}
