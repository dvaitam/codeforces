package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var h0, t0, R int
   if _, err := fmt.Fscan(in, &h0, &t0, &R); err != nil {
       return
   }
   var n int
   fmt.Fscan(in, &n)
   growH := make([]int, n+1)
   growT := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &growH[i], &growT[i])
   }
   var m int
   fmt.Fscan(in, &m)
   growHT := make([]int, m+1)
   growTT := make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(in, &growHT[i], &growTT[i])
   }
   // State id: h*(R+1)+t, valid if h+t <= R
   dim := R + 1
   N := dim * dim
   valid := make([]bool, N)
   for h := 0; h <= R; h++ {
       for t := 0; t <= R; t++ {
           if h+t <= R {
               valid[h*dim+t] = true
           }
       }
   }
   // Build graph
   neighbors := make([][]int, N)
   rev := make([][]int, N)
   losingEdge := make([]bool, N)
   for h := 0; h <= R; h++ {
       for t := 0; t <= R; t++ {
           v := h*dim + t
           if !valid[v] || (h == 0 && t == 0) {
               continue
           }
           // head cuts
           for i := 1; i <= n && i <= h; i++ {
               nh := h - i + growH[i]
               nt := t + growT[i]
               if nh+nt > R {
                   losingEdge[v] = true
               } else {
                   u := nh*dim + nt
                   if valid[u] {
                       neighbors[v] = append(neighbors[v], u)
                       rev[u] = append(rev[u], v)
                   }
               }
           }
           // tail cuts
           for i := 1; i <= m && i <= t; i++ {
               nh := h + growHT[i]
               nt := t - i + growTT[i]
               if nh+nt > R {
                   losingEdge[v] = true
               } else {
                   u := nh*dim + nt
                   if valid[u] {
                       neighbors[v] = append(neighbors[v], u)
                       rev[u] = append(rev[u], v)
                   }
               }
           }
       }
   }
   // BFS for win distances
   const INF = -1
   distWin := make([]int, N)
   for i := range distWin {
       distWin[i] = INF
   }
   winID := 0 // h=0,t=0
   distWin[winID] = 0
   queue := []int{winID}
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, p := range rev[u] {
           if distWin[p] == INF {
               distWin[p] = distWin[u] + 1
               queue = append(queue, p)
           }
       }
   }
   // Determine draw states
   inDraw := make([]bool, N)
   // Tarjan SCC on nodes with distWin==-1
   index := make([]int, N)
   low := make([]int, N)
   onStack := make([]bool, N)
   for i := range index {
       index[i] = -1
   }
   var stk []int
   idx := 0
   sccID := make([]int, N)
   for i := range sccID {
       sccID[i] = -1
   }
   sccCount := 0
   sccSize := []int{}
   var dfs func(v int)
   dfs = func(v int) {
       index[v] = idx; low[v] = idx; idx++
       stk = append(stk, v); onStack[v] = true
       for _, u := range neighbors[v] {
           if !valid[u] || distWin[u] != INF {
               continue
           }
           if index[u] == -1 {
               dfs(u)
               low[v] = min(low[v], low[u])
           } else if onStack[u] {
               low[v] = min(low[v], index[u])
           }
       }
       if low[v] == index[v] {
           // new SCC
           size := 0
           for {
               w := stk[len(stk)-1]
               stk = stk[:len(stk)-1]
               onStack[w] = false
               sccID[w] = sccCount
               size++
               if w == v {
                   break
               }
           }
           sccSize = append(sccSize, size)
           sccCount++
       }
   }
   for v := 0; v < N; v++ {
       if !valid[v] || distWin[v] != INF || (v == winID) {
           continue
       }
       if index[v] == -1 {
           dfs(v)
       }
   }
   // build reverse draw graph
   revDraw := make([][]int, N)
   var drawStart []int
   for v := 0; v < N; v++ {
       if !valid[v] || distWin[v] != INF {
           continue
       }
       // detect cycle nodes
       if sccSize[sccID[v]] > 1 {
           inDraw[v] = true
           drawStart = append(drawStart, v)
       }
       for _, u := range neighbors[v] {
           if !valid[u] || distWin[u] != INF {
               continue
           }
           revDraw[u] = append(revDraw[u], v)
           // self-loop
           if u == v {
               inDraw[v] = true
               drawStart = append(drawStart, v)
           }
       }
   }
   // BFS for draw reachability
   dq := drawStart
   for qi := 0; qi < len(dq); qi++ {
       u := dq[qi]
       for _, p := range revDraw[u] {
           if !inDraw[p] {
               inDraw[p] = true
               dq = append(dq, p)
           }
       }
   }
   // Lose region DP: nodes where distWin==-1 && !inDraw
   // Build outdegree and reverse edges in lose graph
   outdeg := make([]int, N)
   revLose := make([][]int, N)
   for v := 0; v < N; v++ {
       if !valid[v] || distWin[v] != INF || inDraw[v] || (v == winID) {
           continue
       }
       for _, u := range neighbors[v] {
           if valid[u] && distWin[u] == INF && !inDraw[u] {
               outdeg[v]++
               revLose[u] = append(revLose[u], v)
           }
       }
   }
   distLose := make([]int, N)
   // queue nodes with outdeg 0
   var lq []int
   for v := 0; v < N; v++ {
       if !valid[v] || distWin[v] != INF || inDraw[v] || (v == winID) {
           continue
       }
       if outdeg[v] == 0 {
           lq = append(lq, v)
       }
   }
   for qi := 0; qi < len(lq); qi++ {
       v := lq[qi]
       // compute distLose[v]
       best := 0
       if losingEdge[v] {
           best = 1
       }
       for _, u := range neighbors[v] {
           // to L
           nh := u / dim; nt := u % dim
           if !valid[u] {
               continue
           }
           if distWin[u] != INF || inDraw[u] {
               continue
           }
           best = max(best, distLose[u]+1)
       }
       distLose[v] = best
       for _, p := range revLose[v] {
           outdeg[p]--
           if outdeg[p] == 0 {
               lq = append(lq, p)
           }
       }
   }
   // Output
   start := h0*dim + t0
   if distWin[start] != INF {
       fmt.Println("Ivan")
       fmt.Println(distWin[start])
   } else if inDraw[start] {
       fmt.Println("Draw")
   } else {
       fmt.Println("Zmey")
       fmt.Println(distLose[start])
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
