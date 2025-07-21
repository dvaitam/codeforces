package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // grid to track occupied sections
   grid := make([][]bool, n+1)
   for i := range grid {
       grid[i] = make([]bool, m+1)
   }
   // map id to its position [shelf, section]
   pos := make(map[string][2]int)
   for i := 0; i < k; i++ {
       var op string
       fmt.Fscan(in, &op)
       if op == "+1" {
           var x, y int
           var id string
           fmt.Fscan(in, &x, &y, &id)
           placed := false
           // try to place starting at (x,y), then to the right and down shelves
           for sx := x; sx <= n && !placed; sx++ {
               start := 1
               if sx == x {
                   start = y
               }
               for sy := start; sy <= m; sy++ {
                   if !grid[sx][sy] {
                       grid[sx][sy] = true
                       pos[id] = [2]int{sx, sy}
                       placed = true
                       break
                   }
               }
           }
       } else if op == "-1" {
           var id string
           fmt.Fscan(in, &id)
           if p, ok := pos[id]; ok {
               fmt.Println(p[0], p[1])
               grid[p[0]][p[1]] = false
               delete(pos, id)
           } else {
               fmt.Println(-1, -1)
           }
       }
   }
}
