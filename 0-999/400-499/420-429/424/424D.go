package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

type pair struct { val, j int }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var t int
   fmt.Fscan(in, &n, &m, &t)
   var tp, tu, td int
   fmt.Fscan(in, &tp, &tu, &td)
   H := make([][]int, n)
   for i := 0; i < n; i++ {
       H[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &H[i][j])
       }
   }
   // cost arrays
   costE := make([][]int, n)
   costW := make([][]int, n)
   costS := make([][]int, n)
   costN := make([][]int, n)
   for i := 0; i < n; i++ {
       costE[i] = make([]int, m)
       costW[i] = make([]int, m)
       costS[i] = make([]int, m)
       costN[i] = make([]int, m)
       for j := 0; j < m; j++ {
           if j+1 < m {
               h1, h2 := H[i][j], H[i][j+1]
               if h1 == h2 {
                   costE[i][j] = tp
               } else if h1 < h2 {
                   costE[i][j] = tu
               } else {
                   costE[i][j] = td
               }
           }
           if j-1 >= 0 {
               h1, h2 := H[i][j], H[i][j-1]
               if h1 == h2 {
                   costW[i][j] = tp
               } else if h1 < h2 {
                   costW[i][j] = tu
               } else {
                   costW[i][j] = td
               }
           }
           if i+1 < n {
               h1, h2 := H[i][j], H[i+1][j]
               if h1 == h2 {
                   costS[i][j] = tp
               } else if h1 < h2 {
                   costS[i][j] = tu
               } else {
                   costS[i][j] = td
               }
           }
           if i-1 >= 0 {
               h1, h2 := H[i][j], H[i-1][j]
               if h1 == h2 {
                   costN[i][j] = tp
               } else if h1 < h2 {
                   costN[i][j] = tu
               } else {
                   costN[i][j] = td
               }
           }
       }
   }
   // prefix sums
   eastPS := make([][]int, n)
   westPS := make([][]int, n)
   southPS := make([][]int, n)
   northPS := make([][]int, n)
   for i := 0; i < n; i++ {
       eastPS[i] = make([]int, m)
       westPS[i] = make([]int, m)
       for j := 0; j < m; j++ {
           if j == 0 {
               eastPS[i][j] = 0
               westPS[i][j] = 0
           } else {
               eastPS[i][j] = eastPS[i][j-1] + costE[i][j-1]
               westPS[i][j] = westPS[i][j-1] + costW[i][j]
           }
       }
   }
   for j := 0; j < m; j++ {
       southPS[0] = make([]int, m)
       northPS[0] = make([]int, m)
   }
   for i := 1; i < n; i++ {
       southPS[i] = make([]int, m)
       northPS[i] = make([]int, m)
       for j := 0; j < m; j++ {
           southPS[i][j] = southPS[i-1][j] + costS[i-1][j]
           northPS[i][j] = northPS[i-1][j] + costN[i][j]
       }
   }
   bestDiff := int(1e18)
   bi1, bj1, bi2, bj2 := 0, 0, 0, 0
   // temp arrays
   X := make([]int, m)
   Y := make([]int, m)
   var sortedY []pair
   // iterate row pairs
   for i1 := 0; i1 < n; i1++ {
       for i2 := i1 + 2; i2 < n; i2++ {
           // build X,Y
           for j := 0; j < m; j++ {
               top := eastPS[i1][j]
               bot := westPS[i2][j]
               sd := southPS[i2][j] - southPS[i1][j]
               nd := northPS[i2][j] - northPS[i1][j]
               X[j] = top + bot + sd
               Y[j] = top + bot - nd
           }
           sortedY = sortedY[:0]
           // scan columns
           for j2 := 2; j2 < m; j2++ {
               // insert j1 = j2-2
               j1 := j2 - 2
               v := Y[j1]
               // binary insert
               idx := sort.Search(len(sortedY), func(i int) bool { return sortedY[i].val >= v })
               sortedY = append(sortedY, pair{})
               copy(sortedY[idx+1:], sortedY[idx:])
               sortedY[idx] = pair{v, j1}
               // query for j2
               target := X[j2] - t
               // find closest Y to target
               k := sort.Search(len(sortedY), func(i int) bool { return sortedY[i].val >= target })
               for _, kk := range []int{k - 1, k} {
                   if kk >= 0 && kk < len(sortedY) {
                       yv := sortedY[kk].val
                       tj1 := sortedY[kk].j
                       ts := X[j2] - yv
                       diff := abs(ts - t)
                       if diff < bestDiff {
                           bestDiff = diff
                           bi1, bj1, bi2, bj2 = i1, tj1, i2, j2
                           if bestDiff == 0 {
                               fmt.Printf("%d %d %d %d\n", bi1+1, bj1+1, bi2+1, bj2+1)
                               return
                           }
                       }
                   }
               }
           }
       }
   }
   // output best
   fmt.Printf("%d %d %d %d\n", bi1+1, bj1+1, bi2+1, bj2+1)
}
