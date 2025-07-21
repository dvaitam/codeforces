package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   d2     int
   dx, dy []int
   offset int
   vis    map[int]bool
   winMap map[int]bool
)

// state key encoding: x and y in [-offset, offset], rA, rB, turn bits
func encode(x, y, rA, rB, turn int) int {
   // shift x and y to non-negative
   xi := x + offset
   yi := y + offset
   return (((xi << 9) | yi) << 3) | (rA<<2) | (rB<<1) | turn
}

// dfs returns true if current player to move wins
func dfs(x, y, rA, rB, turn int) bool {
   key := encode(x, y, rA, rB, turn)
   if vis[key] {
       return winMap[key]
   }
   vis[key] = true
   // assume losing
   win := false
   // try all vector moves
   for i := range dx {
       nx := x + dx[i]
       ny := y + dy[i]
       if nx*nx+ny*ny > d2 {
           // move leads to loss for mover, skip
           continue
       }
       if !dfs(nx, ny, rA, rB, 1-turn) {
           win = true
           break
       }
   }
   // try reflection if not yet used and still losing
   if !win {
       if turn == 0 && rA == 1 {
           // Anton reflect
           if !dfs(y, x, 0, rB, 1-turn) {
               win = true
           }
       } else if turn == 1 && rB == 1 {
           // Dasha reflect
           if !dfs(y, x, rA, 0, 1-turn) {
               win = true
           }
       }
   }
   winMap[key] = win
   return win
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x, y, n, d int
   if _, err := fmt.Fscan(reader, &x, &y, &n, &d); err != nil {
       return
   }
   dx = make([]int, n)
   dy = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &dx[i], &dy[i])
   }
   offset = d
   d2 = d * d
   vis = make(map[int]bool)
   winMap = make(map[int]bool)
   // Anton and Dasha both have reflection available (1)
   if dfs(x, y, 1, 1, 0) {
       fmt.Println("Anton")
   } else {
       fmt.Println("Dasha")
   }
}
