package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m, q int
   fmt.Fscan(in, &n, &m, &q)
   // 1-based arrays
   grid := make([][]byte, n+2)
   for i := range grid {
       grid[i] = make([]byte, m+2)
   }
   up := make([][]int, n+2)
   down := make([][]int, n+2)
   leftA := make([][]int, n+2)
   rightA := make([][]int, n+2)
   for i := 0; i <= n+1; i++ {
       up[i] = make([]int, m+2)
       down[i] = make([]int, m+2)
       leftA[i] = make([]int, m+2)
       rightA[i] = make([]int, m+2)
   }
   // read initial grid
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           var v int
           fmt.Fscan(in, &v)
           if v != 0 {
               grid[i][j] = 1
           }
       }
   }
   // init up and left
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if grid[i][j] == 1 {
               up[i][j] = up[i-1][j] + 1
               leftA[i][j] = leftA[i][j-1] + 1
           }
       }
       for j := m; j >= 1; j-- {
           if grid[i][j] == 1 {
               rightA[i][j] = rightA[i][j+1] + 1
           }
       }
   }
   // init down
   for j := 1; j <= m; j++ {
       for i := n; i >= 1; i-- {
           if grid[i][j] == 1 {
               down[i][j] = down[i+1][j] + 1
           }
       }
   }
   // temp arrays for histogram
   h := make([]int, m+2)
   Lm := make([]int, m+2)
   Rm := make([]int, m+2)
   hr := make([]int, n+2)
   Lr := make([]int, n+2)
   Rr := make([]int, n+2)
   st := make([]int, max(n, m)+2)
   for qq := 0; qq < q; qq++ {
       var op, x, y int
       fmt.Fscan(in, &op, &x, &y)
       if op == 1 {
           // flip
           if grid[x][y] == 1 {
               grid[x][y] = 0
           } else {
               grid[x][y] = 1
           }
           // update column y up and down
           for i := 1; i <= n; i++ {
               if grid[i][y] == 1 {
                   up[i][y] = up[i-1][y] + 1
               } else {
                   up[i][y] = 0
               }
           }
           for i := n; i >= 1; i-- {
               if grid[i][y] == 1 {
                   down[i][y] = down[i+1][y] + 1
               } else {
                   down[i][y] = 0
               }
           }
           // update row x left and right
           for j := 1; j <= m; j++ {
               if grid[x][j] == 1 {
                   leftA[x][j] = leftA[x][j-1] + 1
               } else {
                   leftA[x][j] = 0
               }
           }
           for j := m; j >= 1; j-- {
               if grid[x][j] == 1 {
                   rightA[x][j] = rightA[x][j+1] + 1
               } else {
                   rightA[x][j] = 0
               }
           }
       } else {
           // query
           if grid[x][y] == 0 {
               fmt.Fprintln(out, 0)
               continue
           }
           ans := 0
           // top edge at x, histogram down[x][1..m]
           // copy
           for j := 1; j <= m; j++ {
               h[j] = down[x][j]
           }
           // compute Lm
           top := 0
           for j := 1; j <= m; j++ {
               for top > 0 && h[st[top-1]] >= h[j] {
                   top--
               }
               if top == 0 {
                   Lm[j] = 0
               } else {
                   Lm[j] = st[top-1]
               }
               st[top] = j; top++
           }
           // compute Rm
           top = 0
           for j := m; j >= 1; j-- {
               for top > 0 && h[st[top-1]] >= h[j] {
                   top--
               }
               if top == 0 {
                   Rm[j] = m + 1
               } else {
                   Rm[j] = st[top-1]
               }
               st[top] = j; top++
           }
           for j := 1; j <= m; j++ {
               l := Lm[j] + 1
               r := Rm[j] - 1
               if l <= y && y <= r {
                   area := h[j] * (r - l + 1)
                   if area > ans {
                       ans = area
                   }
               }
           }
           // bottom edge at x, histogram up[x]
           for j := 1; j <= m; j++ {
               h[j] = up[x][j]
           }
           // Lm
           top = 0
           for j := 1; j <= m; j++ {
               for top > 0 && h[st[top-1]] >= h[j] {
                   top--
               }
               if top == 0 {
                   Lm[j] = 0
               } else {
                   Lm[j] = st[top-1]
               }
               st[top] = j; top++
           }
           // Rm
           top = 0
           for j := m; j >= 1; j-- {
               for top > 0 && h[st[top-1]] >= h[j] {
                   top--
               }
               if top == 0 {
                   Rm[j] = m + 1
               } else {
                   Rm[j] = st[top-1]
               }
               st[top] = j; top++
           }
           for j := 1; j <= m; j++ {
               l := Lm[j] + 1
               r := Rm[j] - 1
               if l <= y && y <= r {
                   area := h[j] * (r - l + 1)
                   if area > ans {
                       ans = area
                   }
               }
           }
           // left edge at y, histogram rightA[1..n][y]
           for i := 1; i <= n; i++ {
               hr[i] = rightA[i][y]
           }
           // Lr
           top = 0
           for i := 1; i <= n; i++ {
               for top > 0 && hr[st[top-1]] >= hr[i] {
                   top--
               }
               if top == 0 {
                   Lr[i] = 0
               } else {
                   Lr[i] = st[top-1]
               }
               st[top] = i; top++
           }
           // Rr
           top = 0
           for i := n; i >= 1; i-- {
               for top > 0 && hr[st[top-1]] >= hr[i] {
                   top--
               }
               if top == 0 {
                   Rr[i] = n + 1
               } else {
                   Rr[i] = st[top-1]
               }
               st[top] = i; top++
           }
           for i := 1; i <= n; i++ {
               l := Lr[i] + 1
               r := Rr[i] - 1
               if l <= x && x <= r {
                   area := hr[i] * (r - l + 1)
                   if area > ans {
                       ans = area
                   }
               }
           }
           // right edge at y, histogram leftA[][y]
           for i := 1; i <= n; i++ {
               hr[i] = leftA[i][y]
           }
           // Lr
           top = 0
           for i := 1; i <= n; i++ {
               for top > 0 && hr[st[top-1]] >= hr[i] {
                   top--
               }
               if top == 0 {
                   Lr[i] = 0
               } else {
                   Lr[i] = st[top-1]
               }
               st[top] = i; top++
           }
           // Rr
           top = 0
           for i := n; i >= 1; i-- {
               for top > 0 && hr[st[top-1]] >= hr[i] {
                   top--
               }
               if top == 0 {
                   Rr[i] = n + 1
               } else {
                   Rr[i] = st[top-1]
               }
               st[top] = i; top++
           }
           for i := 1; i <= n; i++ {
               l := Lr[i] + 1
               r := Rr[i] - 1
               if l <= x && x <= r {
                   area := hr[i] * (r - l + 1)
                   if area > ans {
                       ans = area
                   }
               }
           }
           fmt.Fprintln(out, ans)
       }
   }
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
