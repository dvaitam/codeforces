package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007
const KMAX = 200

func main() {
   in := bufio.NewReader(os.Stdin)
   var k int64
   var n int
   fmt.Fscan(in, &k, &n)
   blocks := make([][2]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &blocks[i][0], &blocks[i][1])
   }
   // BFS up to KMAX to compute T[k] and distances
   // grid offset
   off := KMAX*2 + 5
   size := off*2 + 1
   dist := make([][]int, size)
   for i := range dist {
       dist[i] = make([]int, size)
       for j := range dist[i] {
           dist[i][j] = -1
       }
   }
   // moves
   moves := [8][2]int{{1,2},{2,1},{2,-1},{1,-2},{-1,-2},{-2,-1},{-2,1},{-1,2}}
   // queue
   type P struct{ x, y int }
   q := make([]P, 0, 500000)
   // start
   sx, sy := off, off
   dist[sx][sy] = 0
   q = append(q, P{sx, sy})
   cnt := make([]int64, KMAX+1)
   cnt[0] = 1
   // BFS
   for head := 0; head < len(q); head++ {
       p := q[head]
       d := dist[p.x][p.y]
       if d >= KMAX {
           continue
       }
       nd := d + 1
       for _, mv := range moves {
           nx, ny := p.x+mv[0], p.y+mv[1]
           if nx < 0 || ny < 0 || nx >= size || ny >= size {
               continue
           }
           if dist[nx][ny] != -1 {
               continue
           }
           dist[nx][ny] = nd
           cnt[nd]++
           q = append(q, P{nx, ny})
       }
   }
   // prefix sum T
   T := make([]int64, KMAX+1)
   T[0] = cnt[0]
   for i := 1; i <= KMAX; i++ {
       T[i] = T[i-1] + cnt[i]
   }
   // first and second differences
   D1 := make([]int64, KMAX+1)
   for i := 1; i <= KMAX; i++ {
       D1[i] = T[i] - T[i-1]
   }
   D2 := make([]int64, KMAX+1)
   for i := 2; i <= KMAX; i++ {
       D2[i] = D1[i] - D1[i-1]
   }
   // determine threshold where third difference D3 stable
   thr := 50
   // compute third differences D3[k] = D2[k] - D2[k-1]
   D3 := make([]int64, KMAX+1)
   for i := 3; i <= KMAX; i++ {
       D3[i] = D2[i] - D2[i-1]
   }
   baseD3 := D3[thr]
   // cubic polynomial f(k) = A*k^3 + B*k^2 + C*k + D, D = T[0]
   // baseD3 = 6*A
   A3 := baseD3 / 6
   // from D2[k] = 6A*(k-1) + 2B => B = (D2[thr] - 6A*(thr-1)) / 2
   B3 := (D2[thr] - 6*A3*int64(thr-1)) / 2
   // from D1[k] = A*(3k^2-3k+1) + B*(2k-1) + C => C = D1[thr] - A*(3thr^2-3thr+1) - B*(2*thr-1)
   C3 := D1[thr] - A3*(3*int64(thr)*int64(thr) - 3*int64(thr) + 1) - B3*(2*int64(thr)-1)
   D0 := T[0]
   // compute total reachable
   var total int64
   if k <= KMAX {
       total = T[k]
   } else {
       kk := k % MOD
       // f(k) = ((A*k + B)*k + C)*k + D
       total = A3%MOD * kk %MOD * kk %MOD * kk %MOD
       total = (total + B3%MOD * kk %MOD * kk %MOD) %MOD
       total = (total + C3%MOD * kk %MOD) %MOD
       total = (total + D0%MOD) %MOD
       if total < 0 {
           total += MOD
       }
   }
   // subtract blocked reachable
   var sub int64
   for _, b := range blocks {
       x := b[0] + off
       y := b[1] + off
       if x >= 0 && x < size && y >= 0 && y < size {
           d := int64(dist[x][y])
           if d >= 0 && d <= k {
               sub++
           }
       }
   }
   ans := (total - sub) % MOD
   if ans < 0 {
       ans += MOD
   }
   fmt.Println(ans)
}
