package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

const INF = 1000000000

// Item for priority queue
type Item struct {
   id   int
   dist int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[0 : n-1]
   return it
}

func dijkstra(m, n int, weight [][]int, passable [][]bool, src, dst int) ([]int, []int) {
   N := m * n
   dist := make([]int, N)
   prev := make([]int, N)
   for i := 0; i < N; i++ {
       dist[i] = INF
       prev[i] = -1
   }
   dist[src] = 0
   pq := &PriorityQueue{}
   heap.Init(pq)
   heap.Push(pq, Item{src, 0})
   dr := []int{-1, 1, 0, 0}
   dc := []int{0, 0, -1, 1}
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u := it.id
       d := it.dist
       if d != dist[u] {
           continue
       }
       if u == dst {
           break
       }
       r := u / n
       c := u % n
       w := weight[r][c]
       for k := 0; k < 4; k++ {
           nr := r + dr[k]
           nc := c + dc[k]
           if nr < 0 || nr >= m || nc < 0 || nc >= n {
               continue
           }
           if !passable[nr][nc] {
               continue
           }
           v := nr*n + nc
           nd := d + w
           if nd < dist[v] {
               dist[v] = nd
               prev[v] = u
               heap.Push(pq, Item{v, nd})
           }
       }
   }
   return dist, prev
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var m, n, k int
   fmt.Fscan(reader, &m, &n, &k)
   grid := make([]string, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   weight := make([][]int, m)
   passable := make([][]bool, m)
   junctions := make(map[byte]int)
   for i := 0; i < m; i++ {
       weight[i] = make([]int, n)
       passable[i] = make([]bool, n)
       for j := 0; j < n; j++ {
           ch := grid[i][j]
           if ch == '#' {
               passable[i][j] = false
           } else {
               passable[i][j] = true
               if ch >= '1' && ch <= '9' {
                   weight[i][j] = int(ch - '0')
               } else {
                   // junction
                   weight[i][j] = 1
                   junctions[ch] = i*n + j
               }
           }
       }
   }
   var rs, cs, re, ce int
   var s string
   fmt.Fscan(reader, &rs, &cs, &s, &re, &ce)
   rs--
   cs--
   re--
   ce--
   // build sequence of points
   seq := make([]int, 0, len(s)+2)
   start := rs*n + cs
   seq = append(seq, start)
   for i := 0; i < len(s); i++ {
       id, ok := junctions[s[i]]
       if !ok {
           // should not happen
           id = -1
       }
       seq = append(seq, id)
   }
   dest := re*n + ce
   seq = append(seq, dest)
   // full path of nodes
   path := make([]int, 0)
   // for each segment
   for i := 0; i+1 < len(seq); i++ {
       u := seq[i]
       v := seq[i+1]
       _, prev := dijkstra(m, n, weight, passable, u, v)
       // reconstruct path from u to v
       tmp := make([]int, 0)
       cur := v
       for cur != -1 && cur != u {
           tmp = append(tmp, cur)
           cur = prev[cur]
       }
       // include u
       tmp = append(tmp, u)
       // reverse
       for l, r := 0, len(tmp)-1; l < r; l, r = l+1, r-1 {
           tmp[l], tmp[r] = tmp[r], tmp[l]
       }
       if i > 0 {
           // skip duplicate start
           tmp = tmp[1:]
       }
       path = append(path, tmp...)
   }
   // simulate time
   time := 0
   pos := path[0]
   for i := 0; i+1 < len(path); i++ {
       u := path[i]
       r := u / n
       c := u % n
       w := weight[r][c]
       if time+w > k {
           break
       }
       time += w
       pos = path[i+1]
   }
   // output pos
   rf := pos/n + 1
   cf := pos % n + 1
   fmt.Fprintf(writer, "%d %d", rf, cf)
}
