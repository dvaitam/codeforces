package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct{ r, c int }

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var n, m, k, q int
       fmt.Fscan(in, &n, &m, &k, &q)
       // initialize grid 1-indexed
       grid := make([][]bool, n+2)
       for i := range grid {
           grid[i] = make([]bool, m+2)
       }
       ops := make([]pair, 0, q)
       // helper functions
       placeH := func(r, c int) {
           for p := c; p <= c+k-1; p++ {
               grid[r][p] = true
           }
       }
       placeV := func(r, c int) {
           for p := r; p <= r+k-1; p++ {
               grid[p][c] = true
           }
       }
       place := func(r, c int, tp byte) {
           ops = append(ops, pair{r, c})
           if tp == 'H' {
               placeH(r, c)
           } else {
               placeV(r, c)
           }
       }
       clearRow := func(r int) {
           for c := 1; c <= m; c++ {
               grid[r][c] = false
           }
       }
       clearCol := func(c int) {
           for r := 1; r <= n; r++ {
               grid[r][c] = false
           }
       }
       // process queries
       var typestr string
       if k == 1 {
           ptr := 1
           for i := 0; i < q; i++ {
               fmt.Fscan(in, &typestr)
               tp := typestr[0]
               place(1, ptr, tp)
               if ptr == m {
                   clearRow(1)
                   ptr = 1
               } else {
                   ptr++
               }
           }
       } else if k == min(n, m) {
           if k == max(n, m) {
               for i := 0; i < q; i++ {
                   fmt.Fscan(in, &typestr)
                   tp := typestr[0]
                   place(1, 1, tp)
                   if tp == 'H' {
                       clearRow(1)
                   } else {
                       clearCol(1)
                   }
               }
           } else if k == n {
               ptr := 1
               for i := 0; i < q; i++ {
                   fmt.Fscan(in, &typestr)
                   tp := typestr[0]
                   if tp == 'H' {
                       place(ptr, 2, tp)
                       if ptr == n {
                           for j := 2; j <= 2+k-1; j++ {
                               clearCol(j)
                           }
                           ptr = 1
                       } else {
                           ptr++
                       }
                   } else {
                       place(1, 1, tp)
                       clearCol(1)
                       for r := 1; r <= n; r++ {
                           if grid[r][2] && m == k+1 {
                               for c := 2; c <= m; c++ {
                                   grid[r][c] = false
                               }
                               ptr = 1
                           }
                       }
                   }
               }
           } else {
               // k == m
               ptr := 1
               for i := 0; i < q; i++ {
                   fmt.Fscan(in, &typestr)
                   tp := typestr[0]
                   if tp == 'V' {
                       place(2, ptr, tp)
                       if ptr == m {
                           for j := 2; j <= 2+k-1; j++ {
                               clearRow(j)
                           }
                           ptr = 1
                       } else {
                           ptr++
                       }
                   } else {
                       place(1, 1, tp)
                       clearRow(1)
                       for c := 1; c <= m; c++ {
                           if grid[2][c] && n == k+1 {
                               for r := 2; r <= n; r++ {
                                   grid[r][c] = false
                               }
                               ptr = 1
                           }
                       }
                   }
               }
           }
       } else {
           for i := 0; i < q; i++ {
               fmt.Fscan(in, &typestr)
               tp := typestr[0]
               h_clean := !((!grid[1][m-k+1]) && grid[1][m])
               v_clean := !((!grid[n-k+1][1]) && grid[n][1])
               if tp == 'H' {
                   if h_clean {
                       to_place := -1
                       for ii := 1; ii <= n; ii++ {
                           if !grid[ii][m-k+1] {
                               to_place = ii
                               break
                           }
                       }
                       if to_place < 0 {
                           to_place = -1 // should not happen
                       }
                       if to_place < n-k {
                           place(to_place, m-k+1, tp)
                       } else if to_place == n-k {
                           place(to_place, m-k+1, tp)
                           for j := m-k+1; j <= m; j++ {
                               if grid[n][j] {
                                   clearCol(j)
                               }
                           }
                       } else if v_clean {
                           place(to_place, m-k+1, tp)
                           if to_place == n-k+1 {
                               if grid[n-k+1][m-k] {
                                   clearRow(n-k+1)
                               }
                           }
                           if to_place == n {
                               for j := m-k+1; j <= m; j++ {
                                   clearCol(j)
                               }
                           }
                       } else {
                           // goto hell
                           // find to_del
                           to_del := -1
                           for r := n-k+1; r <= n; r++ {
                               if grid[r][m-k] {
                                   to_del = r
                                   break
                               }
                           }
                           if to_del < 0 {
                               // find alternate to_place
                               to_place2 := -1
                               for r := n-k+1; r <= n; r++ {
                                   if !grid[r][m-k+1] {
                                       to_place2 = r
                                       break
                                   }
                               }
                               place(to_place2, m-k+1, tp)
                               if to_place2 == n {
                                   for c := m-k+1; c <= m; c++ {
                                       if grid[1][c] {
                                           clearCol(c)
                                       }
                                   }
                               }
                           } else {
                               place(to_del, m-k+1, tp)
                               clearRow(to_del)
                           }
                       }
                   } else {
                       // h_clean failed or v_clean branch: remove existing block
                       to_del := -1
                       for r := n-k+1; r <= n; r++ {
                           if grid[r][m-k] {
                               to_del = r
                               break
                           }
                       }
                       if to_del < 0 {
                           to_place2 := -1
                           for r := n-k+1; r <= n; r++ {
                               if !grid[r][m-k+1] {
                                   to_place2 = r
                                   break
                               }
                           }
                           place(to_place2, m-k+1, tp)
                           if to_place2 == n {
                               for c := m-k+1; c <= m; c++ {
                                   if grid[1][c] {
                                       clearCol(c)
                                   }
                               }
                           }
                       } else {
                           place(to_del, m-k+1, tp)
                           clearRow(to_del)
                       }
                   }
               } else {
                   if v_clean {
                       to_place := -1
                       for ii := 1; ii <= m; ii++ {
                           if !grid[n-k+1][ii] {
                               to_place = ii
                               break
                           }
                       }
                       if to_place < m-k {
                           place(n-k+1, to_place, tp)
                       } else if to_place == m-k {
                           place(n-k+1, to_place, tp)
                           for j := n-k+1; j <= n; j++ {
                               if grid[j][m] {
                                   clearRow(j)
                               }
                           }
                       } else if h_clean {
                           place(n-k+1, to_place, tp)
                           if to_place == m-k+1 {
                               if grid[n-k][m-k+1] {
                                   clearCol(m-k+1)
                               }
                           }
                           if to_place == m {
                               for j := n-k+1; j <= n; j++ {
                                   clearRow(j)
                               }
                           }
                       } else {
                           // handle failure similar to hell2
                           to_del := -1
                           for c := m-k+1; c <= m; c++ {
                               if grid[n-k][c] {
                                   to_del = c
                                   break
                               }
                           }
                           if to_del < 0 {
                               to_place2 := -1
                               for c := m-k+1; c <= m; c++ {
                                   if !grid[n-k+1][c] {
                                       to_place2 = c
                                       break
                                   }
                               }
                               place(n-k+1, to_place2, tp)
                               if to_place2 == m {
                                   for r := n-k+1; r <= n; r++ {
                                       if grid[r][1] {
                                           clearRow(r)
                                       }
                                   }
                               }
                           } else {
                               place(n-k+1, to_del, tp)
                               clearCol(to_del)
                           }
                       }
                   }
               }
           }
       }
       // output ops
       for _, p := range ops {
           fmt.Fprintln(out, p.r, p.c)
       }
   }
}
