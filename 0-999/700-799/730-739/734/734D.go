package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var x0, y0 int64
   fmt.Fscan(reader, &n)
   fmt.Fscan(reader, &x0, &y0)

   const INF = int64(4e18)
   // directions: 0 N, 1 NE, 2 E, 3 SE, 4 S, 5 SW, 6 W, 7 NW
   dist := [8]int64{INF, INF, INF, INF, INF, INF, INF, INF}
   piece := [8]byte{}

   for i := 0; i < n; i++ {
       var typ string
       var x, y int64
       fmt.Fscan(reader, &typ, &x, &y)
       t := typ[0]
       dx := x - x0
       dy := y - y0
       dir := -1
       var d int64
       if dx == 0 {
           if dy > 0 {
               dir, d = 0, dy
           } else if dy < 0 {
               dir, d = 4, -dy
           }
       } else if dy == 0 {
           if dx > 0 {
               dir, d = 2, dx
           } else if dx < 0 {
               dir, d = 6, -dx
           }
       } else if abs(dx) == abs(dy) {
           switch {
           case dx > 0 && dy > 0:
               dir, d = 1, dx
           case dx > 0 && dy < 0:
               dir, d = 3, dx
           case dx < 0 && dy < 0:
               dir, d = 5, -dx
           case dx < 0 && dy > 0:
               dir, d = 7, -dx
           }
       }
       if dir >= 0 && d < dist[dir] {
           dist[dir] = d
           piece[dir] = t
       }
   }

   for dir, t := range piece {
       if t == 0 {
           continue
       }
       // even dirs: orthogonal; odd: diagonal
       if dir%2 == 0 {
           // orthogonal directions
           if t == 'R' || t == 'Q' {
               fmt.Println("YES")
               return
           }
       } else {
           // diagonal directions
           if t == 'B' || t == 'Q' {
               fmt.Println("YES")
               return
           }
       }
   }
   fmt.Println("NO")
}

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}
