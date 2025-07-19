package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // l[i][j] is the price between attractions i and j
   l := make([][]int, n)
   for i := range l {
       l[i] = make([]int, n)
   }
   // used[c] = true if price c has been assigned
   used := make([]bool, 1001)
   // initialize diagonal to 0 (default), mark initial edges
   if n > 1 {
       l[0][1], l[1][0] = 1, 1
       used[1] = true
   }
   if n > 2 {
       l[0][2], l[2][0] = 2, 2
       used[2] = true
       l[1][2], l[2][1] = 3, 3
       used[3] = true
   }
   // build remaining edges greedily
   for i := 3; i < n; i++ {
       for c := 1; c <= 1000; c++ {
           ok := true
           for j := 0; j < i; j++ {
               if c+l[0][j] > 1000 || used[c+l[0][j]] {
                   ok = false
                   break
               }
           }
           if !ok {
               continue
           }
           for j := 0; j < i; j++ {
               v := c + l[0][j]
               l[j][i], l[i][j] = v, v
               used[v] = true
           }
           break
       }
   }
   // output adjacency matrix
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           fmt.Printf("%d", l[i][j])
           if j+1 < n {
               fmt.Print(" ")
           }
       }
       fmt.Println()
   }
}
