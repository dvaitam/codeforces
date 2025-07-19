package main

import "fmt"

type pair struct { first, second int }

func inRange(x int, p pair) bool {
   return x > p.first && x < p.second
}

func cross(p1, p2 pair) bool {
   return (inRange(p1.first, p2) || inRange(p1.second, p2)) &&
       (inRange(p2.first, p1) || inRange(p2.second, p1))
}

var (
   m    int
   adj  [][]bool
   col  []int
   flag bool
)

func dfs(u, c int) {
   if col[u] != -1 {
       if col[u] != c {
           flag = false
       }
       return
   }
   if !flag {
       return
   }
   col[u] = c
   for v := 0; v < m; v++ {
       if adj[u][v] {
           dfs(v, 1-c)
       }
   }
}

func main() {

   var n int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   p := make([]pair, m)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Scan(&a, &b)
       if a > b {
           a, b = b, a
       }
       p[i] = pair{a, b}
   }
   adj = make([][]bool, m)
   for i := 0; i < m; i++ {
       adj[i] = make([]bool, m)
       for j := 0; j < m; j++ {
           if cross(p[i], p[j]) {
               adj[i][j] = true
           }
       }
   }
   col = make([]int, m)
   for i := range col {
       col[i] = -1
   }
   flag = true
   for i := 0; i < m; i++ {
       if col[i] == -1 {
           dfs(i, 0)
       }
   }
   if !flag {
       fmt.Println("Impossible")
       return
   }
   // output: 0->o, 1->i
   for i := 0; i < m; i++ {
       if col[i] == 0 {
           fmt.Print("o")
       } else {
           fmt.Print("i")
       }
   }
   fmt.Println()
}
