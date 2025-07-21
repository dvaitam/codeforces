package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// rotate90 returns the point (x,y) rotated 90 degrees CCW around (a,b).
func rotate90(x, y, a, b int) (int, int) {
   dx := x - a
   dy := y - b
   // (-dy, dx)
   return a - dy, b + dx
}

// isSquare checks whether 4 points form a non-degenerate square.
func isSquare(px [4]int, py [4]int) bool {
   var dists []int
   for i := 0; i < 4; i++ {
       for j := i + 1; j < 4; j++ {
           dx := px[i] - px[j]
           dy := py[i] - py[j]
           dists = append(dists, dx*dx+dy*dy)
       }
   }
   sort.Ints(dists)
   // 4 sides equal and >0, 2 diagonals equal, diag = 2*side
   if dists[0] == 0 {
       return false
   }
   side := dists[0]
   // sides: dists[0..3], diagonals: dists[4], dists[5]
   for i := 0; i < 4; i++ {
       if dists[i] != side {
           return false
       }
   }
   diag := dists[4]
   if dists[4] != dists[5] {
       return false
   }
   if diag != 2*side {
       return false
   }
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // For each regiment
   for i := 0; i < n; i++ {
       // Read 4 moles
       x := [4]int{}
       y := [4]int{}
       a := [4]int{}
       b := [4]int{}
       for j := 0; j < 4; j++ {
           fmt.Fscan(in, &x[j], &y[j], &a[j], &b[j])
       }
       // Precompute rotations
       var rx [4][4]int
       var ry [4][4]int
       for j := 0; j < 4; j++ {
           cx, cy := x[j], y[j]
           for k := 0; k < 4; k++ {
               rx[j][k] = cx
               ry[j][k] = cy
               // rotate for next k
               cx, cy = rotate90(cx, cy, a[j], b[j])
           }
       }
       best := -1
       // try all combinations
       for k0 := 0; k0 < 4; k0++ {
           for k1 := 0; k1 < 4; k1++ {
               for k2 := 0; k2 < 4; k2++ {
                   for k3 := 0; k3 < 4; k3++ {
                       cost := k0 + k1 + k2 + k3
                       if best != -1 && cost >= best {
                           continue
                       }
                       px := [4]int{rx[0][k0], rx[1][k1], rx[2][k2], rx[3][k3]}
                       py := [4]int{ry[0][k0], ry[1][k1], ry[2][k2], ry[3][k3]}
                       if isSquare(px, py) {
                           best = cost
                       }
                   }
               }
           }
       }
       fmt.Println(best)
   }
}
