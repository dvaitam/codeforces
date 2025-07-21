package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       line := make([]byte, m)
       if _, err := fmt.Fscan(reader, &line); err != nil {
           return
       }
       grid[i] = line
   }
   // prefix sums for rows and columns
   rowPS := make([][]int, n)
   for i := 0; i < n; i++ {
       rowPS[i] = make([]int, m+1)
       for j := 0; j < m; j++ {
           rowPS[i][j+1] = rowPS[i][j]
           if grid[i][j] == 'B' {
               rowPS[i][j+1]++
           }
       }
   }
   colPS := make([][]int, m)
   for j := 0; j < m; j++ {
       colPS[j] = make([]int, n+1)
       for i := 0; i < n; i++ {
           colPS[j][i+1] = colPS[j][i]
           if grid[i][j] == 'B' {
               colPS[j][i+1]++
           }
       }
   }
   // collect black cells
   type cell struct{ r, c int }
   blacks := make([]cell, 0, n*m)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == 'B' {
               blacks = append(blacks, cell{i, j})
           }
       }
   }
   // check each pair
   k := len(blacks)
   for a := 0; a < k; a++ {
       for b := a + 1; b < k; b++ {
           r1, c1 := blacks[a].r, blacks[a].c
           r2, c2 := blacks[b].r, blacks[b].c
           ok := false
           // direct horizontal
           if r1 == r2 {
               cmin, cmax := c1, c2
               if cmin > cmax {
                   cmin, cmax = cmax, cmin
               }
               if rowPS[r1][cmax+1]-rowPS[r1][cmin] == cmax-cmin+1 {
                   ok = true
               }
           }
           // direct vertical
           if !ok && c1 == c2 {
               rmin, rmax := r1, r2
               if rmin > rmax {
                   rmin, rmax = rmax, rmin
               }
               if colPS[c1][rmax+1]-colPS[c1][rmin] == rmax-rmin+1 {
                   ok = true
               }
           }
           // L-shape via (r1,c2)
           if !ok {
               cmin, cmax := c1, c2
               if cmin > cmax {
                   cmin, cmax = cmax, cmin
               }
               rmin, rmax := r1, r2
               if rmin > rmax {
                   rmin, rmax = rmax, rmin
               }
               seg1 := rowPS[r1][cmax+1]-rowPS[r1][cmin] == cmax-cmin+1
               seg2 := colPS[c2][rmax+1]-colPS[c2][rmin] == rmax-rmin+1
               if seg1 && seg2 {
                   ok = true
               }
           }
           // L-shape via (r2,c1)
           if !ok {
               cmin, cmax := c1, c2
               if cmin > cmax {
                   cmin, cmax = cmax, cmin
               }
               rmin, rmax := r1, r2
               if rmin > rmax {
                   rmin, rmax = rmax, rmin
               }
               seg1 := colPS[c1][rmax+1]-colPS[c1][rmin] == rmax-rmin+1
               seg2 := rowPS[r2][cmax+1]-rowPS[r2][cmin] == cmax-cmin+1
               if seg1 && seg2 {
                   ok = true
               }
           }
           if !ok {
               fmt.Println("NO")
               return
           }
       }
   }
   fmt.Println("YES")
}
