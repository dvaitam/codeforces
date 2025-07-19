package main

import (
   "fmt"
   "os"
)

func main() {
   var n, d, h int
   if _, err := fmt.Fscan(os.Stdin, &n, &d, &h); err != nil {
       return
   }
   // Impossible if diameter greater than twice the height
   if d > 2*h {
       fmt.Println(-1)
       return
   }
   // Special case: only two nodes
   if d == 1 {
       if n == 2 {
           fmt.Println("1 2")
       } else {
           fmt.Println(-1)
       }
       return
   }
   // Build tree edges
   edges := make([][2]int, 0, n-1)
   used := 1
   // Build the main chain of height h
   for i := 0; i < h; i++ {
       edges = append(edges, [2]int{used, used + 1})
       used++
   }
   if h < d {
       // Build the second chain to reach diameter d
       // Connect from root (1)
       edges = append(edges, [2]int{1, used + 1})
       used++
       for i := 1; i < d-h; i++ {
           edges = append(edges, [2]int{used, used + 1})
           used++
       }
       // Attach any remaining nodes to root
       for i := used + 1; i <= n; i++ {
           edges = append(edges, [2]int{1, i})
       }
   } else {
       // Attach any remaining nodes to the second node of main chain
       for i := used + 1; i <= n; i++ {
           edges = append(edges, [2]int{2, i})
       }
   }
   // Output edges
   for _, e := range edges {
       fmt.Println(e[0], e[1])
   }
}
