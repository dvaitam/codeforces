package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   MOD   = 998244353
   MAXC  = 510
   OFFSET = 5
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // grid indices [0..MAXC)
   var infected [MAXC][MAXC]bool
   // read initial infected
   coords := make([][2]int, n)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       x += OFFSET; y += OFFSET
       infected[x][y] = true
       coords[i][0], coords[i][1] = x, y
   }
   // neighbor count for healthy
   var cnt [MAXC][MAXC]int
   type P struct{ x, y int }
   var q []P
   // directions depend on parity
   for _, p := range coords {
       x, y := p[0], p[1]
       // for each neighbor
       // common two
       for _, d := range [][2]int{{1, 0}, {-1, 0}} {
           nx, ny := x + d[0], y + d[1]
           if nx >= 0 && nx < MAXC && ny >= 0 && ny < MAXC && !infected[nx][ny] {
               cnt[nx][ny]++
               if cnt[nx][ny] == 2 {
                   q = append(q, P{nx, ny})
               }
           }
       }
       // third neighbor
       if x&1 == 0 {
           nx, ny := x+1, y-1
           if nx >= 0 && nx < MAXC && ny >= 0 && ny < MAXC && !infected[nx][ny] {
               cnt[nx][ny]++
               if cnt[nx][ny] == 2 {
                   q = append(q, P{nx, ny})
               }
           }
       } else {
           nx, ny := x-1, y+1
           if nx >= 0 && nx < MAXC && ny >= 0 && ny < MAXC && !infected[nx][ny] {
               cnt[nx][ny]++
               if cnt[nx][ny] == 2 {
                   q = append(q, P{nx, ny})
               }
           }
       }
   }
   // closure: BFS for infection spread
   for head := 0; head < len(q); head++ {
       p := q[head]
       x, y := p.x, p.y
       if infected[x][y] {
           continue
       }
       // infect
       infected[x][y] = true
       // update neighbors of this new infection
       for _, d := range [][2]int{{1, 0}, {-1, 0}} {
           nx, ny := x + d[0], y + d[1]
           if nx >= 0 && nx < MAXC && ny >= 0 && ny < MAXC && !infected[nx][ny] {
               cnt[nx][ny]++
               if cnt[nx][ny] == 2 {
                   q = append(q, P{nx, ny})
               }
           }
       }
       if x&1 == 0 {
           nx, ny := x+1, y-1
           if nx >= 0 && nx < MAXC && ny >= 0 && ny < MAXC && !infected[nx][ny] {
               cnt[nx][ny]++
               if cnt[nx][ny] == 2 {
                   q = append(q, P{nx, ny})
               }
           }
       } else {
           nx, ny := x-1, y+1
           if nx >= 0 && nx < MAXC && ny >= 0 && ny < MAXC && !infected[nx][ny] {
               cnt[nx][ny]++
               if cnt[nx][ny] == 2 {
                   q = append(q, P{nx, ny})
               }
           }
       }
   }
   // collect closure nodes and assign ids
   id := make([][]int, MAXC)
   for i := range id {
       id[i] = make([]int, MAXC)
       for j := range id[i] {
           id[i][j] = -1
       }
   }
   cntNodes := 0
   for x := 0; x < MAXC; x++ {
       for y := 0; y < MAXC; y++ {
           if infected[x][y] {
               id[x][y] = cntNodes
               cntNodes++
           }
       }
   }
   // DSU
   parent := make([]int, cntNodes)
   for i := 0; i < cntNodes; i++ {
       parent[i] = -1
   }
   var find func(int) int
   find = func(a int) int {
       if parent[a] < 0 {
           return a
       }
       parent[a] = find(parent[a])
       return parent[a]
   }
   union := func(a, b int) bool {
       a = find(a); b = find(b)
       if a == b {
           return false
       }
       if parent[a] > parent[b] {
           a, b = b, a
       }
       parent[a] += parent[b]
       parent[b] = a
       return true
   }
   // check cycles
   sick := false
   for x := 0; x < MAXC && !sick; x++ {
       for y := 0; y < MAXC; y++ {
           if !infected[x][y] {
               continue
           }
           u := id[x][y]
           // neighbors
           // two common
           for _, d := range [][2]int{{1, 0}, {-1, 0}} {
               nx, ny := x + d[0], y + d[1]
               if nx >= 0 && nx < MAXC && ny >= 0 && ny < MAXC && infected[nx][ny] {
                   v := id[nx][ny]
                   // to avoid double, only consider u < v
                   if u < v {
                       if !union(u, v) {
                           sick = true
                           break
                       }
                   }
               }
           }
           if sick {
               break
           }
           if x&1 == 0 {
               nx, ny := x+1, y-1
               if nx >= 0 && nx < MAXC && ny >= 0 && ny < MAXC && infected[nx][ny] {
                   v := id[nx][ny]
                   if u < v {
                       if !union(u, v) {
                           sick = true
                       }
                   }
               }
           } else {
               nx, ny := x-1, y+1
               if nx >= 0 && nx < MAXC && ny >= 0 && ny < MAXC && infected[nx][ny] {
                   v := id[nx][ny]
                   if u < v {
                       if !union(u, v) {
                           sick = true
                       }
                   }
               }
           }
           if sick {
               break
           }
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   if sick {
       fmt.Fprintln(out, "SICK")
   } else {
       // maximum moves = 2*|S_inf| - n
       moves := (2*cntNodes - n) % MOD
       if moves < 0 {
           moves += MOD
       }
       fmt.Fprintln(out, "RECOVERED")
       fmt.Fprintln(out, moves)
   }
}
