package main

import (
   "fmt"
   "sort"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   // single node case
   if n == 1 {
       if m == a[0] {
           fmt.Println("Yes")
           for i := 0; i < a[0]; i++ {
               if i > 0 {
                   fmt.Print(" ")
               }
               fmt.Print("1 1")
           }
           fmt.Println()
       } else {
           fmt.Println("No")
       }
       return
   }
   // sort indices by degree
   t := make([]int, n)
   for i := 0; i < n; i++ {
       t[i] = i
   }
   sort.Slice(t, func(i, j int) bool {
       return a[t[i]] < a[t[j]]
   })
   // impossible checks
   if m < 1 || (m == 1 && a[t[0]] > 1) {
       fmt.Println("No")
       return
   }
   // total degree
   sum := 0
   for _, v := range a {
       sum += v
   }
   if sum < m {
       fmt.Println("No")
       return
   }
   fmt.Println("Yes")
   // build sequence
   if a[t[0]] <= m {
       // central node is smallest
       fmt.Print(t[0]+1)
       // reduce extra visits from others
       for i := 1; i < n && sum > m; i++ {
           for a[t[i]] > 0 && sum > m {
               sum--
               a[t[i]]--
               fmt.Print(" ", t[i]+1, " ", t[i]+1)
           }
       }
       // connect back to central
       fmt.Print(" ", t[0]+1)
       a[t[0]]--
       // remaining visits
       for i := 0; i < n; i++ {
           for j := 0; j < a[t[i]]; j++ {
               fmt.Print(" ", t[i]+1, " ", t[i]+1)
           }
       }
       fmt.Println()
   } else {
       // central node degree > m, use second smallest start
       fmt.Print(t[1]+1)
       // initial loop from m to degree of central
       for i := m; i <= a[t[0]]; i++ {
           fmt.Print(" ", t[0]+1, " ", t[0]+1)
       }
       m--
       // link second to central
       fmt.Print(" ", t[1]+1, " ", t[0]+1)
       // decrement second node
       a[t[1]]--
       // remaining visits of others
       for idx := 1; idx < n; idx++ {
           for j := 0; j < a[t[idx]]; j++ {
               fmt.Print(" ", t[idx]+1, " ", t[idx]+1)
           }
       }
       // final visits to central
       fmt.Print(" ", t[0]+1)
       for i := 1; i < m; i++ {
           fmt.Print(" ", t[0]+1, " ", t[0]+1)
       }
       fmt.Println()
   }
}
