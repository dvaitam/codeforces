package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func nextPerm(a []int) bool {
   // lexicographical next permutation
   n := len(a)
   i := n - 2
   for i >= 0 && a[i] >= a[i+1] {
       i--
   }
   if i < 0 {
       return false
   }
   j := n - 1
   for a[j] <= a[i] {
       j--
   }
   a[i], a[j] = a[j], a[i]
   // reverse a[i+1:]
   for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
       a[l], a[r] = a[r], a[l]
   }
   return true
}

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   // corner definitions: xflag, yflag
   cornersX := [4]int{0, 0, 1, 1}
   cornersY := [4]int{0, 1, 0, 1}
   for ; t > 0; t-- {
       pts := [4][2]int64{}
       for i := 0; i < 4; i++ {
           fmt.Fscan(in, &pts[i][0], &pts[i][1])
       }
       perm := []int{0, 1, 2, 3}
       ans := int64(1<<62)
       for {
           // group coordinates
           ax := make([]int64, 0, 2)
           bx := make([]int64, 0, 2)
           ay := make([]int64, 0, 2)
           by := make([]int64, 0, 2)
           for k := 0; k < 4; k++ {
               p := pts[perm[k]]
               if cornersX[k] == 0 {
                   ax = append(ax, p[0])
               } else {
                   bx = append(bx, p[0])
               }
               if cornersY[k] == 0 {
                   ay = append(ay, p[1])
               } else {
                   by = append(by, p[1])
               }
           }
           // collect candidate s
           cand := map[int64]struct{}{}
           cand[0] = struct{}{}
           for _, xi := range bx {
               for _, xj := range ax {
                   s := xi - xj
                   if s >= 0 {
                       cand[s] = struct{}{}
                   }
               }
           }
           for _, yi := range by {
               for _, yj := range ay {
                   s := yi - yj
                   if s >= 0 {
                       cand[s] = struct{}{}
                   }
               }
           }
           // evaluate
           for s := range cand {
               // x dimension
               X := [4]int64{}
               idx := 0
               for _, v := range ax {
                   X[idx] = v
                   idx++
               }
               for _, v := range bx {
                   X[idx] = v - s
                   idx++
               }
               // y dimension
               Y := [4]int64{}
               idx = 0
               for _, v := range ay {
                   Y[idx] = v
                   idx++
               }
               for _, v := range by {
                   Y[idx] = v - s
                   idx++
               }
               // sort
               sort.Slice(X[:], func(i, j int) bool { return X[i] < X[j] })
               sort.Slice(Y[:], func(i, j int) bool { return Y[i] < Y[j] })
               // cost
               cost := (X[2] + X[3] - X[0] - X[1]) + (Y[2] + Y[3] - Y[0] - Y[1])
               if cost < ans {
                   ans = cost
               }
           }
           if !nextPerm(perm) {
               break
           }
       }
       fmt.Fprintln(out, ans)
   }
}
