package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // count factors and record zero location
   twos := make([][]int, n)
   fives := make([][]int, n)
   zeroExists := false
   zi, zj := 0, 0
   for i := 0; i < n; i++ {
       twos[i] = make([]int, n)
       fives[i] = make([]int, n)
       for j := 0; j < n; j++ {
           var x int
           fmt.Fscan(in, &x)
           if x == 0 {
               if !zeroExists {
                   zeroExists = true
                   zi, zj = i, j
               }
               // avoid zero in DP
               twos[i][j], fives[i][j] = INF, INF
           } else {
               c2, c5 := 0, 0
               for x%2 == 0 {
                   c2++
                   x /= 2
               }
               for x%5 == 0 {
                   c5++
                   x /= 5
               }
               twos[i][j], fives[i][j] = c2, c5
           }
       }
   }
   // DP minimizing twos (tie-break by fives)
   type pair struct{ t, f int }
   dp2 := make([][]pair, n)
   dir2 := make([][]byte, n)
   for i := 0; i < n; i++ {
       dp2[i] = make([]pair, n)
       dir2[i] = make([]byte, n)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if i == 0 && j == 0 {
               dp2[0][0] = pair{twos[0][0], fives[0][0]}
           } else if i == 0 {
               prev := dp2[0][j-1]
               dp2[0][j] = pair{prev.t + twos[0][j], prev.f + fives[0][j]}
               dir2[0][j] = 'R'
           } else if j == 0 {
               prev := dp2[i-1][0]
               dp2[i][0] = pair{prev.t + twos[i][0], prev.f + fives[i][0]}
               dir2[i][0] = 'D'
           } else {
               up := pair{dp2[i-1][j].t + twos[i][j], dp2[i-1][j].f + fives[i][j]}
               left := pair{dp2[i][j-1].t + twos[i][j], dp2[i][j-1].f + fives[i][j]}
               // choose smaller twos, tie-break by fives
               if up.t < left.t || (up.t == left.t && up.f < left.f) {
                   dp2[i][j], dir2[i][j] = up, 'D'
               } else {
                   dp2[i][j], dir2[i][j] = left, 'R'
               }
           }
       }
   }
   // DP minimizing fives (tie-break by twos)
   dp5 := make([][]pair, n)
   dir5 := make([][]byte, n)
   for i := 0; i < n; i++ {
       dp5[i] = make([]pair, n)
       dir5[i] = make([]byte, n)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if i == 0 && j == 0 {
               dp5[0][0] = pair{twos[0][0], fives[0][0]}
           } else if i == 0 {
               prev := dp5[0][j-1]
               dp5[0][j] = pair{prev.t + twos[0][j], prev.f + fives[0][j]}
               dir5[0][j] = 'R'
           } else if j == 0 {
               prev := dp5[i-1][0]
               dp5[i][0] = pair{prev.t + twos[i][0], prev.f + fives[i][0]}
               dir5[i][0] = 'D'
           } else {
               up := pair{dp5[i-1][j].t + twos[i][j], dp5[i-1][j].f + fives[i][j]}
               left := pair{dp5[i][j-1].t + twos[i][j], dp5[i][j-1].f + fives[i][j]}
               // choose smaller fives, tie-break by twos
               if up.f < left.f || (up.f == left.f && up.t < left.t) {
                   dp5[i][j], dir5[i][j] = up, 'D'
               } else {
                   dp5[i][j], dir5[i][j] = left, 'R'
               }
           }
       }
   }
   // compute best among two paths
   end2 := dp2[n-1][n-1]
   z2 := min(end2.t, end2.f)
   end5 := dp5[n-1][n-1]
   z5 := min(end5.t, end5.f)
   best := z2
   use5 := false
   if z5 < best {
       best = z5
       use5 = true
   }
   var path []byte
   // check zero path
   if zeroExists && best > 1 {
       best = 1
       // path through zero
       // rights to zj, downs to zi, then rights to end, downs to end
       for k := 0; k < zj; k++ {
           path = append(path, 'R')
       }
       for k := 0; k < zi; k++ {
           path = append(path, 'D')
       }
       for k := zj; k < n-1; k++ {
           path = append(path, 'R')
       }
       for k := zi; k < n-1; k++ {
           path = append(path, 'D')
       }
   } else {
       // backtrack selected DP
       i, j := n-1, n-1
       if use5 {
           for i > 0 || j > 0 {
               if dir5[i][j] == 'D' {
                   path = append(path, 'D')
                   i--
               } else {
                   path = append(path, 'R')
                   j--
               }
           }
       } else {
           for i > 0 || j > 0 {
               if dir2[i][j] == 'D' {
                   path = append(path, 'D')
                   i--
               } else {
                   path = append(path, 'R')
                   j--
               }
           }
       }
       // reverse
       for l, r := 0, len(path)-1; l < r; l, r = l+1, r-1 {
           path[l], path[r] = path[r], path[l]
       }
   }
   // output
   fmt.Println(best)
   fmt.Println(string(path))
}
