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
   grid := make([]string, n+1)
   for i := 1; i <= n; i++ {
       var row string
       fmt.Fscan(reader, &row)
       grid[i] = row
   }
   // build 2D prefix sum of obstacles ('#')
   sum := make([][]int, n+1)
   for i := 0; i <= n; i++ {
       sum[i] = make([]int, m+1)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           add := 0
           if grid[i][j-1] == '#' {
               add = 1
           }
           sum[i][j] = sum[i-1][j] + sum[i][j-1] - sum[i-1][j-1] + add
       }
   }
   // track positions and active status of letters A-Z
   pos := make([][2]int, 26)
   ok := make([]bool, 26)
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           c := grid[i][j-1]
           if c >= 'A' && c <= 'Z' {
               idx := int(c - 'A')
               pos[idx][0] = i
               pos[idx][1] = j
               ok[idx] = true
           }
       }
   }
   var k int
   fmt.Fscan(reader, &k)
   // direction deltas for N, S, W, E
   dmap := map[byte][2]int{
       'N': {-1, 0},
       'S': {+1, 0},
       'W': {0, -1},
       'E': {0, +1},
   }
   for q := 0; q < k; q++ {
       var cmd string
       var dist int
       fmt.Fscan(reader, &cmd, &dist)
       delta := dmap[cmd[0]]
       dx, dy := delta[0], delta[1]
       for idx := 0; idx < 26; idx++ {
           if !ok[idx] {
               continue
           }
           x0, y0 := pos[idx][0], pos[idx][1]
           x1 := x0 + dx*dist
           y1 := y0 + dy*dist
           // check bounds
           if x1 < 1 || x1 > n || y1 < 1 || y1 > m {
               ok[idx] = false
           } else {
               xa, xb := x0, x1
               if xa > xb {
                   xa, xb = xb, xa
               }
               ya, yb := y0, y1
               if ya > yb {
                   ya, yb = yb, ya
               }
               // if any obstacle in rectangle [xa..xb][ya..yb]
               if sum[xb][yb] - sum[xa-1][yb] - sum[xb][ya-1] + sum[xa-1][ya-1] > 0 {
                   ok[idx] = false
               }
           }
           pos[idx][0], pos[idx][1] = x1, y1
       }
   }
   // collect remaining active letters
   var result []byte
   for idx := 0; idx < 26; idx++ {
       if ok[idx] {
           result = append(result, byte('A'+idx))
       }
   }
   if len(result) == 0 {
       fmt.Println("no solution")
   } else {
       fmt.Println(string(result))
   }
}
