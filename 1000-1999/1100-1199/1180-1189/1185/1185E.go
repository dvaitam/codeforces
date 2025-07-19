package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   const INF = 1000000000
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       grid := make([][]int, n)
       for i := 0; i < n; i++ {
           var s string
           fmt.Fscan(reader, &s)
           grid[i] = make([]int, m)
           for j := 0; j < m; j++ {
               grid[i][j] = int(s[j] - 'a')
           }
       }
       // find max letter and last position
       mx := -1
       lastR, lastC := 0, 0
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               v := grid[i][j]
               if v > mx {
                   mx = v
                   lastR, lastC = i, j
               }
           }
       }
       // initialize arrays
       row := make([]int, mx+1)
       col := make([]int, mx+1)
       mnr := make([]int, mx+1)
       mxr := make([]int, mx+1)
       mnc := make([]int, mx+1)
       mxc := make([]int, mx+1)
       for c := 0; c <= mx; c++ {
           row[c] = -1
           col[c] = -1
           mnr[c] = n
           mnc[c] = m
           mxr[c] = -1
           mxc[c] = -1
       }
       // record positions
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               c := grid[i][j]
               // row
               if row[c] == -1 {
                   row[c] = i
               } else if row[c] != i {
                   row[c] = INF
               }
               // col
               if col[c] == -1 {
                   col[c] = j
               } else if col[c] != j {
                   col[c] = INF
               }
               // bounds
               if i < mnr[c] {
                   mnr[c] = i
               }
               if i > mxr[c] {
                   mxr[c] = i
               }
               if j < mnc[c] {
                   mnc[c] = j
               }
               if j > mxc[c] {
                   mxc[c] = j
               }
           }
       }
       // build result and ops
       b := make([][]int, n)
       for i := range b {
           b[i] = make([]int, m)
           for j := range b[i] {
               b[i][j] = -1
           }
       }
       ops := make([][4]int, mx+1)
       ok := true
       for c := 0; c <= mx; c++ {
           if row[c] == -1 {
               // absent, use last
               ops[c] = [4]int{lastR, lastC, lastR, lastC}
           } else if row[c] == INF && col[c] == INF {
               ok = false
               break
           } else if row[c] != INF {
               r := row[c]
               ops[c] = [4]int{r, mnc[c], r, mxc[c]}
               for x := mnc[c]; x <= mxc[c]; x++ {
                   b[r][x] = c
               }
           } else {
               // vertical
               cc := col[c]
               ops[c] = [4]int{mnr[c], cc, mxr[c], cc}
               for i := mnr[c]; i <= mxr[c]; i++ {
                   b[i][cc] = c
               }
           }
       }
       // validate
       if ok {
           for i := 0; i < n && ok; i++ {
               for j := 0; j < m; j++ {
                   if b[i][j] != grid[i][j] {
                       ok = false
                       break
                   }
               }
           }
       }
       if !ok {
           fmt.Fprintln(writer, "NO")
           continue
       }
       // output
       fmt.Fprintln(writer, "YES")
       fmt.Fprintln(writer, mx+1)
       for c := 0; c <= mx; c++ {
           // convert to 1-based
           r1, c1, r2, c2 := ops[c][0]+1, ops[c][1]+1, ops[c][2]+1, ops[c][3]+1
           fmt.Fprintf(writer, "%d %d %d %d\n", r1, c1, r2, c2)
       }
   }
}
