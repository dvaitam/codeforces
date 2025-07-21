package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&arr[i])
   }
   if n == 1 {
       fmt.Println("Exemplary pages.")
       return
   }
   // Check if all equal
   allEq := true
   for i := 1; i < n; i++ {
       if arr[i] != arr[0] {
           allEq = false
           break
       }
   }
   if allEq {
       fmt.Println("Exemplary pages.")
       return
   }
   // Count values
   cnt := make(map[int]int)
   for _, v := range arr {
       cnt[v]++
   }
   // Collect unique values
   uniq := make([]int, 0, len(cnt))
   for v := range cnt {
       uniq = append(uniq, v)
   }
   // sort uniq
   for i := 0; i < len(uniq); i++ {
       for j := i + 1; j < len(uniq); j++ {
           if uniq[j] < uniq[i] {
               uniq[i], uniq[j] = uniq[j], uniq[i]
           }
       }
   }
   u := len(uniq)
   // Helper to print unrecoverable
   unrecover := func() {
       fmt.Println("Unrecoverable configuration.")
   }
   if u > 3 {
       unrecover()
       return
   }
   // Case with three unique: s-v, s, s+v
   if u == 3 {
       x, y, z := uniq[0], uniq[1], uniq[2]
       // differences must match
       if z-y != y-x || y-x <= 0 {
           unrecover()
           return
       }
       v := y - x
       if cnt[x] != 1 || cnt[z] != 1 || cnt[y] != n-2 {
           unrecover()
           return
       }
       var ia, ib int
       for i, val := range arr {
           if val == x {
               ia = i + 1
           } else if val == z {
               ib = i + 1
           }
       }
       fmt.Printf("%d ml. from cup #%d to cup #%d.\n", v, ia, ib)
       return
   }
   // Case with two unique
   if u == 2 {
       if n != 2 {
           unrecover()
           return
       }
       x, y := uniq[0], uniq[1]
       sum := x + y
       if sum%2 != 0 {
           unrecover()
           return
       }
       s := sum / 2
       v := y - s
       if v <= 0 {
           unrecover()
           return
       }
       // find positions
       var ia, ib int
       for i, val := range arr {
           if val == x {
               ia = i + 1
           } else if val == y {
               ib = i + 1
           }
       }
       fmt.Printf("%d ml. from cup #%d to cup #%d.\n", v, ia, ib)
       return
   }
   // Should not reach here
   unrecover()
}
