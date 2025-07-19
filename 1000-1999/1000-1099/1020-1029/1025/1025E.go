package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Move represents a token move from (r1,c1) to (r2,c2)
type Move struct{ r1, c1, r2, c2 int }

func solve(n, m int, R, C, TR, TC []int, A [][]int) []Move {
   res := make([]Move, 0, 4*m*n)
   // directions: up, down, left, right
   dr := []int{-1, 1, 0, 0}
   dc := []int{0, 0, -1, 1}
   // move token id to (r,c) in straight line
   var move func(id, r, c int)
   move = func(id, r, c int) {
       if R[id] != r && C[id] != c {
           panic("invalid move")
       }
       dir := -1
       switch {
       case R[id] > r:
           dir = 0
       case R[id] < r:
           dir = 1
       case C[id] > c:
           dir = 2
       case C[id] < c:
           dir = 3
       }
       if dir < 0 {
           return
       }
       for R[id] != r || C[id] != c {
           pr, pc := R[id], C[id]
           A[pr][pc] = -1
           R[id] += dr[dir]
           C[id] += dc[dir]
           A[R[id]][C[id]] = id
           res = append(res, Move{pr, pc, R[id], C[id]})
       }
   }
   // initial sort phase
   type item struct{ key1, key2, id int }
   v := make([]item, m)
   for i := 0; i < m; i++ {
       k1 := R[i]
       k2 := R[i]
       if R[i] == 0 {
           k2 = C[i]
       } else {
           k2 = -C[i]
       }
       v[i] = item{k1, k2, i}
   }
   sort.Slice(v, func(i, j int) bool {
       if v[i].key1 != v[j].key1 {
           return v[i].key1 < v[j].key1
       }
       return v[i].key2 < v[j].key2
   })
   for i := 0; i < m; i++ {
       id := v[i].id
       tr, tc := 0, i
       if C[id] < tc {
           move(id, R[id], tc)
       }
       move(id, tr, C[id])
       move(id, tr, tc)
   }
   // middle sort phase
   v2 := make([]item, m)
   for i := 0; i < m; i++ {
       v2[i] = item{TR[i], TC[i], i}
   }
   sort.Slice(v2, func(i, j int) bool {
       if v2[i].key1 != v2[j].key1 {
           return v2[i].key1 < v2[j].key1
       }
       return v2[i].key2 < v2[j].key2
   })
   for idx := m - 1; idx >= 0; idx-- {
       if v2[idx].key1 > 1 {
           id := v2[idx].id
           move(id, 1, C[id])
           move(id, R[id], TC[id])
           move(id, TR[id], TC[id])
       }
   }
   // final adjustments on row 0
   var v3 []item
   for i := 0; i < m; i++ {
       if TR[i] <= 1 {
           k2 := TC[i]
           if TR[i] != 0 {
               k2 = -TC[i]
           }
           v3 = append(v3, item{TR[i], k2, i})
       }
   }
   sort.Slice(v3, func(i, j int) bool {
       if v3[i].key1 != v3[j].key1 {
           return v3[i].key1 < v3[j].key1
       }
       return v3[i].key2 < v3[j].key2
   })
   am := len(v3)
   for i := 0; i < am; i++ {
       id := v3[i].id
       for C[id] != i {
           if A[0][C[id]-1] == -1 {
               move(id, 0, C[id]-1)
           } else {
               id2 := A[0][C[id]-1]
               move(id2, 1, C[id2])
               move(id, 0, C[id]-1)
               move(id2, 1, C[id2]+1)
               move(id2, 0, C[id2])
           }
       }
   }
   for idx := am - 1; idx >= 0; idx-- {
       id := v3[idx].id
       if TR[id] != 0 {
           move(id, 0, n-1)
           move(id, 1, n-1)
       }
       move(id, TR[id], TC[id])
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   R := make([]int, m)
   C := make([]int, m)
   TR := make([]int, m)
   TC := make([]int, m)
   A := make([][]int, n)
   for i := 0; i < n; i++ {
       row := make([]int, n)
       for j := 0; j < n; j++ {
           row[j] = -1
       }
       A[i] = row
   }
   for i := 0; i < m; i++ {
       var r, c int
       fmt.Fscan(in, &r, &c)
       r--
       c--
       R[i], C[i] = r, c
       A[r][c] = i
   }
   for i := 0; i < m; i++ {
       var r, c int
       fmt.Fscan(in, &r, &c)
       TR[i], TC[i] = r-1, c-1
   }
   moves := solve(n, m, R, C, TR, TC, A)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(moves))
   for _, mv := range moves {
       fmt.Fprintf(w, "%d %d %d %d\n", mv.r1+1, mv.c1+1, mv.r2+1, mv.c2+1)
   }
}
