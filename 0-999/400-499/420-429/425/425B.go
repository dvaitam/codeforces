package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   grid := make([][]int, n)
   for i := 0; i < n; i++ {
       grid[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &grid[i][j])
       }
   }
   // search flips by iterative deepening
   type Cell struct{ r, c int }
   dirs := []Cell{{1,0},{-1,0},{0,1},{0,-1}}
   // prepare initial grid
   initGrid := grid
   // helper: apply flips to grid
   makeGrid := func(flips []Cell) [][]int {
       g := make([][]int, n)
       for i := 0; i < n; i++ {
           g[i] = make([]int, m)
           copy(g[i], initGrid[i])
       }
       for _, u := range flips {
           g[u.r][u.c] ^= 1
       }
       return g
   }
   // helper: find first violation component and collect all violations
   type CompInfo struct{
       cells []Cell; v int; minR, maxR, minC, maxC int; holes []Cell
   }
   findViolations := func(g [][]int) ([]CompInfo, bool) {
       vis := make([][]bool, n)
       for i := range vis { vis[i] = make([]bool, m) }
       var infos []CompInfo
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if vis[i][j] { continue }
               v := g[i][j]
               queue := []Cell{{i, j}}
               vis[i][j] = true
               comp := []Cell{{i, j}}
               for qi := 0; qi < len(queue); qi++ {
                   u := queue[qi]
                   for _, d := range dirs {
                       ni, nj := u.r+d.r, u.c+d.c
                       if ni>=0 && ni<n && nj>=0 && nj<m && !vis[ni][nj] && g[ni][nj]==v {
                           vis[ni][nj] = true
                           queue = append(queue, Cell{ni, nj})
                           comp = append(comp, Cell{ni, nj})
                       }
                   }
               }
               minR, maxR, minC, maxC := n, -1, m, -1
               for _, u := range comp {
                   if u.r < minR { minR = u.r }
                   if u.r > maxR { maxR = u.r }
                   if u.c < minC { minC = u.c }
                   if u.c > maxC { maxC = u.c }
               }
               h := maxR - minR + 1
               w := maxC - minC + 1
               area := h * w
               if area != len(comp) {
                   // compute holes
                   holes := []Cell{}
                   for rr := minR; rr <= maxR; rr++ {
                       for cc := minC; cc <= maxC; cc++ {
                           if g[rr][cc] != v {
                               holes = append(holes, Cell{rr, cc})
                           }
                       }
                   }
                   infos = append(infos, CompInfo{comp, v, minR, maxR, minC, maxC, holes})
               }
           }
       }
       return infos, len(infos)==0
   }
   // DFS search
   var dfs func(depth, D int, flips []Cell) bool
   dfs = func(depth, D int, flips []Cell) bool {
       g := makeGrid(flips)
       infos, ok := findViolations(g)
       if ok {
           return true
       }
       // heuristic: at least one flip per violating component
       lb := len(infos)
       if depth+lb > D {
           return false
       }
       // pick first violation comp
       info := infos[0]
       // candidate flips: fill holes
       // also candidate removal: comp cell nearest center
       // find center
       midR := (info.minR + info.maxR) / 2
       midC := (info.minC + info.maxC) / 2
       bestIdx := 0
       bestDist := n + m
       for idx, u := range info.cells {
           d := abs(u.r-midR) + abs(u.c-midC)
           if d < bestDist {
               bestDist = d
               bestIdx = idx
           }
       }
       rem := info.cells[bestIdx]
       // generate candidates
       var cands []Cell
       cands = append(cands, info.holes...)
       cands = append(cands, rem)
       // try candidates
       used := map[int]bool{}
       for _, u := range flips {
           used[u.r*m+u.c] = true
       }
       for _, u := range cands {
           key := u.r*m + u.c
           if used[key] {
               continue
           }
           // flip
           if dfs(depth+1, D, append(flips, u)) {
               return true
           }
       }
       return false
   }
   ans := -1
   for D := 0; D <= k; D++ {
       if dfs(0, D, nil) {
           ans = D
           break
       }
   }
   fmt.Println(ans)
}

func abs(x int) int { if x<0 { return -x }; return x }
