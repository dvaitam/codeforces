package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Cost represents (moves, piece changes)
type Cost struct{ moves, changes int }

func (c Cost) less(o Cost) bool {
   if c.moves != o.moves {
       return c.moves < o.moves
   }
   return c.changes < o.changes
}

// Item is a state in priority queue
type Item struct {
   x, y, p    int
   cost       Cost
   index      int
}

// PriorityQueue implements heap.Interface
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
   return pq[i].cost.less(pq[j].cost)
}
func (pq PriorityQueue) Swap(i, j int) {
   pq[i], pq[j] = pq[j], pq[i]
   pq[i].index = i
   pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
   n := len(*pq)
   item := x.(*Item)
   item.index = n
   *pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   item := old[n-1]
   item.index = -1
   *pq = old[0 : n-1]
   return item
}

var (
   dx = []int{2, -2, 2, -2, 1, -1, 1, -1}
   dy = []int{1, 1, -1, -1, 2, 2, -2, -2}
)

func dijkstra(n, sx, sy, tx, ty int, disPrev [3]Cost) [3]Cost {
   const INF = 1000000000
   // initialize distances
   dist := make([][][]Cost, n)
   for i := 0; i < n; i++ {
       dist[i] = make([][]Cost, n)
       for j := 0; j < n; j++ {
           dist[i][j] = make([]Cost, 3)
           for p := 0; p < 3; p++ {
               dist[i][j][p] = Cost{INF, INF}
           }
       }
   }
   var pq PriorityQueue
   heap.Init(&pq)
   // starting states
   for p := 0; p < 3; p++ {
       dist[sx][sy][p] = disPrev[p]
       heap.Push(&pq, &Item{x: sx, y: sy, p: p, cost: disPrev[p]})
   }
   cnt := 0
   for pq.Len() > 0 {
       it := heap.Pop(&pq).(*Item)
       x, y, p, c := it.x, it.y, it.p, it.cost
       // skip outdated
       if c.moves != dist[x][y][p].moves || c.changes != dist[x][y][p].changes {
           continue
       }
       // reached target
       if x == tx && y == ty {
           cnt++
           if cnt == 3 {
               break
           }
       }
       // change piece
       for np := 0; np < 3; np++ {
           if np == p {
               continue
           }
           nc := Cost{c.moves + 1, c.changes + 1}
           if nc.less(dist[x][y][np]) {
               dist[x][y][np] = nc
               heap.Push(&pq, &Item{x: x, y: y, p: np, cost: nc})
           }
       }
       // moves by piece
       switch p {
       case 0: // knight
           for d := 0; d < 8; d++ {
               nx, ny := x+dx[d], y+dy[d]
               if nx < 0 || nx >= n || ny < 0 || ny >= n {
                   continue
               }
               nc := Cost{c.moves + 1, c.changes}
               if nc.less(dist[nx][ny][p]) {
                   dist[nx][ny][p] = nc
                   heap.Push(&pq, &Item{x: nx, y: ny, p: p, cost: nc})
               }
           }
       case 1: // rook
           for i := 0; i < n; i++ {
               if i == x {
                   continue
               }
               nc := Cost{c.moves + 1, c.changes}
               if nc.less(dist[i][y][p]) {
                   dist[i][y][p] = nc
                   heap.Push(&pq, &Item{x: i, y: y, p: p, cost: nc})
               }
           }
           for j := 0; j < n; j++ {
               if j == y {
                   continue
               }
               nc := Cost{c.moves + 1, c.changes}
               if nc.less(dist[x][j][p]) {
                   dist[x][j][p] = nc
                   heap.Push(&pq, &Item{x: x, y: j, p: p, cost: nc})
               }
           }
       case 2: // bishop
           sum := x + y
           diff := x - y
           for i := 0; i < n; i++ {
               if i == x {
                   continue
               }
               // diagonal sum
               j1 := sum - i
               if j1 >= 0 && j1 < n {
                   nc := Cost{c.moves + 1, c.changes}
                   if nc.less(dist[i][j1][p]) {
                       dist[i][j1][p] = nc
                       heap.Push(&pq, &Item{x: i, y: j1, p: p, cost: nc})
                   }
               }
               // diagonal diff
               j2 := i - diff
               if j2 >= 0 && j2 < n {
                   nc := Cost{c.moves + 1, c.changes}
                   if nc.less(dist[i][j2][p]) {
                       dist[i][j2][p] = nc
                       heap.Push(&pq, &Item{x: i, y: j2, p: p, cost: nc})
                   }
               }
           }
       }
   }
   return [3]Cost{dist[tx][ty][0], dist[tx][ty][1], dist[tx][ty][2]}
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   size := n * n
   posX := make([]int, size+1)
   posY := make([]int, size+1)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           var v int
           fmt.Fscan(reader, &v)
           posX[v] = i
           posY[v] = j
       }
   }
   dis := [3]Cost{}
   // start at value 1
   x1, y1 := posX[1], posY[1]
   // visit in increasing order
   for v := 2; v <= size; v++ {
       x2, y2 := posX[v], posY[v]
       dis = dijkstra(n, x1, y1, x2, y2, dis)
       x1, y1 = x2, y2
   }
   // find best among pieces
   res := dis[0]
   for p := 1; p < 3; p++ {
       if dis[p].less(res) {
           res = dis[p]
       }
   }
   fmt.Fprintln(writer, res.moves, res.changes)
}
