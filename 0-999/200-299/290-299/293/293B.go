package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

var (
   n, m, k int
   board [][]int
   coords [][2]int
   prevList [][]int
   color []int
   total int
)

// DFS over cells in row-major order
func dfs(pos int) {
   if pos == len(coords) {
       total++
       if total >= MOD {
           total -= MOD
       }
       return
   }
   // current cell
   if board[coords[pos][0]][coords[pos][1]] != 0 {
       // precolored
       c := board[coords[pos][0]][coords[pos][1]]
       // check conflict with prev
       for _, v := range prevList[pos] {
           if color[v] == c {
               return
           }
       }
       color[pos] = c
       dfs(pos + 1)
       color[pos] = 0
       return
   }
   // try all colors
   for c := 1; c <= k; c++ {
       ok := true
       for _, v := range prevList[pos] {
           if color[v] == c {
               ok = false
               break
           }
       }
       if !ok {
           continue
       }
       color[pos] = c
       dfs(pos + 1)
       color[pos] = 0
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m, &k)
   board = make([][]int, n)
   for i := 0; i < n; i++ {
       board[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &board[i][j])
       }
   }
   // Quick impossibility
   if n > k || m > k || n+m-1 > k {
       fmt.Fprintln(out, 0)
       return
   }
   // prepare coords
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           coords = append(coords, [2]int{i, j})
       }
   }
   N := len(coords)
   prevList = make([][]int, N)
   // precompute comparability: for each pos, list previous comparable
   for i := 0; i < N; i++ {
       x1, y1 := coords[i][0], coords[i][1]
       for j := 0; j < i; j++ {
           x2, y2 := coords[j][0], coords[j][1]
           // if comparable
           if (x1 <= x2 && y1 <= y2) || (x2 <= x1 && y2 <= y1) {
               // since j < i in order, j is previous
               prevList[i] = append(prevList[i], j)
           }
       }
   }
   color = make([]int, N)
   total = 0
   dfs(0)
   fmt.Fprintln(out, total)
}
