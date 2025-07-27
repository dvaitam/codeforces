package main

import (
   "bufio"
   "fmt"
   "os"
)

type point struct{ y, x int }
type state struct{ y, x, f, day int }

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, M int
   var K, T, W int
   if _, err := fmt.Fscan(in, &N, &M); err != nil {
       return
   }
   fmt.Fscan(in, &K, &T, &W)
   grid := make([][]byte, N)
   var sy, sx, py, px int
   for i := 0; i < N; i++ {
       row := make([]byte, M)
       fmt.Fscan(in, &row)
       grid[i] = row
       for j := 0; j < M; j++ {
           if row[j] == 'M' {
               sy, sx = i, j
           } else if row[j] == 'P' {
               py, px = i, j
           }
       }
   }
   wind := make([]byte, W)
   if W > 0 {
       fmt.Fscan(in, &wind)
   }
   supplies := make([]struct{ y, x, d int }, T)
   for i := 0; i < T; i++ {
       fmt.Fscan(in, &supplies[i].y, &supplies[i].x, &supplies[i].d)
   }
   // compute sea component from start
   comp := make([][]bool, N)
   for i := range comp {
       comp[i] = make([]bool, M)
   }
   q := []point{{sy, sx}}
   comp[sy][sx] = true
   dirs := []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
   validSea := func(c byte) bool {
       return c == 'S' || c == 'M' || c == 'P'
   }
   for qi := 0; qi < len(q); qi++ {
       p := q[qi]
       for _, d := range dirs {
           ny, nx := p.y+d.y, p.x+d.x
           if ny >= 0 && ny < N && nx >= 0 && nx < M && !comp[ny][nx] && validSea(grid[ny][nx]) {
               comp[ny][nx] = true
               q = append(q, point{ny, nx})
           }
       }
   }
   if !comp[py][px] {
       fmt.Println(-1)
       return
   }
   // visited min day for y,x,f
   const INF = 1e9
   visited := make([][][]int, N)
   for i := 0; i < N; i++ {
       visited[i] = make([][]int, M)
       for j := 0; j < M; j++ {
           visited[i][j] = make([]int, K+1)
           for f := 0; f <= K; f++ {
               visited[i][j][f] = INF
           }
       }
   }
   // BFS queue
   dq := make([]state, 0, 1024)
   dq = append(dq, state{sy, sx, K, 0})
   visited[sy][sx][K] = 0
   // wind mapping
   wmap := map[byte]point{'N': {-1, 0}, 'S': {1, 0}, 'W': {0, -1}, 'E': {0, 1}, 'C': {0, 0}}
   for qi := 0; qi < len(dq); qi++ {
       st := dq[qi]
       if st.y == py && st.x == px {
           fmt.Println(st.day)
           return
       }
       if st.day >= W {
           continue
       }
       w := wmap[wind[st.day]]
       // try moves
       for _, mv := range append([]point{{0, 0}}, dirs...) {
           ny := st.y + mv.y + w.y
           nx := st.x + mv.x + w.x
           if ny < 0 || ny >= N || nx < 0 || nx >= M {
               continue
           }
           if !comp[ny][nx] {
               continue
           }
           nf := st.f - 1
           nd := st.day + 1
           if nf < 0 {
               continue
           }
           // replenish
           for _, sp := range supplies {
               if sp.y == ny && sp.x == nx && sp.d == nd {
                   nf = K
                   break
               }
           }
           if visited[ny][nx][nf] <= nd {
               continue
           }
           visited[ny][nx][nf] = nd
           dq = append(dq, state{ny, nx, nf, nd})
       }
   }
   fmt.Println(-1)
}
