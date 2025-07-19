package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   X := make([]int, n+1)
   Y := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &X[i], &Y[i])
   }
   A := make([]int, 10)
   for i := 1; i <= 9; i++ {
       fmt.Fscan(reader, &A[i])
   }
   sort.Ints(A[1:10])
   // map sums to indices
   hmap := make(map[int]int)
   cnt := 0
   for i := 1; i <= 9; i++ {
       for j := i + 1; j <= 9; j++ {
           for k := j + 1; k <= 9; k++ {
               s := A[i] + A[j] + A[k]
               if _, ok := hmap[s]; !ok {
                   cnt++
                   hmap[s] = cnt
               }
               t := n - s
               if _, ok := hmap[t]; !ok {
                   cnt++
                   hmap[t] = cnt
               }
           }
       }
   }
   // compress Y coordinates
   tmp := make([]int, n)
   for i := 1; i <= n; i++ {
       tmp[i-1] = i
   }
   sort.Slice(tmp, func(i, j int) bool { return Y[tmp[i]] < Y[tmp[j]] })
   SymY := make([]int, n+2)
   Yx := make([]int, n+2)
   Yc := make([]int, n+1)
   tY := 1
   for i := 1; i <= n; i++ {
       v := tmp[i-1]
       SymY[tY] = Y[v]
       if i < n && Y[v] == Y[tmp[i]] {
           Yx[i] = -1
           Yc[v] = tY
       } else {
           Yx[i] = tY
           Yc[v] = tY
           tY++
       }
   }
   // compress X coordinates
   for i := 1; i <= n; i++ {
       tmp[i-1] = i
   }
   sort.Slice(tmp, func(i, j int) bool { return X[tmp[i]] < X[tmp[j]] })
   SymX := make([]int, n+2)
   Xx := make([]int, n+2)
   Xc := make([]int, n+1)
   tX := 1
   for i := 1; i <= n; i++ {
       v := tmp[i-1]
       SymX[tX] = X[v]
       if i < n && X[v] == X[tmp[i]] {
           Xx[i] = -1
           Xc[v] = tX
       } else {
           Xx[i] = tX
           Xc[v] = tX
           tX++
       }
   }
   // BIT and snapshots
   T := make([]int, n+2)
   Tx := make([][]int, cnt+1)
   for i := 1; i <= cnt; i++ {
       Tx[i] = make([]int, n+2)
   }
   // build BIT ordering by X
   for i := 1; i <= n; i++ {
       v := tmp[i-1]
       yv := Yc[v]
       for j := yv; j <= n; j += j & -j {
           T[j]++
       }
       if s, ok := hmap[i]; ok && Xx[i] != -1 {
           copy(Tx[s], T)
       }
   }
   // check function
   check := func(x, y, goal int) bool {
       if x < 1 || x > n || y < 1 || y > n {
           return false
       }
       if Xx[x] == -1 || Yx[y] == -1 {
           return false
       }
       s, ok := hmap[x]
       if !ok {
           return false
       }
       sum := goal
       yy := Yx[y]
       for j := yy; j > 0; j -= j & -j {
           sum -= Tx[s][j]
       }
       return sum == 0
   }
   // next permutation over A[1..9]
   nextPerm := func(a []int) bool {
       m := len(a)
       i := m - 2
       for i >= 0 && a[i] >= a[i+1] {
           i--
       }
       if i < 0 {
           return false
       }
       j := m - 1
       for a[j] <= a[i] {
           j--
       }
       a[i], a[j] = a[j], a[i]
       for l, r := i+1, m-1; l < r; l, r = l+1, r-1 {
           a[l], a[r] = a[r], a[l]
       }
       return true
   }
   // search
   for {
       a1 := A[1] + A[2] + A[3]
       a2 := a1 + A[4] + A[5] + A[6]
       a3 := A[1] + A[4] + A[7]
       a4 := a3 + A[2] + A[5] + A[8]
       if check(a3, a1, A[1]) && check(a4, a1, A[1]+A[2]) &&
          check(a3, a2, A[1]+A[4]) && check(a4, a2, A[1]+A[2]+A[4]+A[5]) {
           x1 := float64(SymX[Xx[a3]]) + 0.5
           x2 := float64(SymX[Xx[a4]]) + 0.5
           y1 := float64(SymY[Yx[a1]]) + 0.5
           y2 := float64(SymY[Yx[a2]]) + 0.5
           fmt.Printf("%.12f %.12f\n%.12f %.12f\n", x1, x2, y1, y2)
           return
       }
       if !nextPerm(A[1:10]) {
           fmt.Println(-1)
           return
       }
   }
}
