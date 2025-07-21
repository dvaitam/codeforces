package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Gripper struct {
   x, y     int64
   mi, pi, ri int64
}

// segment tree node
type STNode struct {
   l, r       int
   xmin, xmax int64
   ymin, ymax int64
}

var (
   grippers []Gripper
   xs, ys   []int64
   pis, ris []int64
   mis      []int64
   visited  []bool
   st       []STNode
) 

func build(node, l, r int) {
   st[node].l = l
   st[node].r = r
   if l == r {
       // leaf
       x, y := xs[l], ys[l]
       st[node].xmin, st[node].xmax = x, x
       st[node].ymin, st[node].ymax = y, y
   } else {
       mid := (l + r) >> 1
       build(node<<1, l, mid)
       build(node<<1|1, mid+1, r)
       left, right := &st[node<<1], &st[node<<1|1]
       // combine
       st[node].xmin = min64(left.xmin, right.xmin)
       st[node].xmax = max64(left.xmax, right.xmax)
       st[node].ymin = min64(left.ymin, right.ymin)
       st[node].ymax = max64(left.ymax, right.ymax)
   }
}

func update(node, idx int) {
   l, r := st[node].l, st[node].r
   if l == r {
       // remove leaf
       st[node].xmin = 1<<62
       st[node].xmax = -1<<62
       st[node].ymin = 1<<62
       st[node].ymax = -1<<62
       return
   }
   mid := (l + r) >> 1
   if idx <= mid {
       update(node<<1, idx)
   } else {
       update(node<<1|1, idx)
   }
   left, right := &st[node<<1], &st[node<<1|1]
   st[node].xmin = min64(left.xmin, right.xmin)
   st[node].xmax = max64(left.xmax, right.xmax)
   st[node].ymin = min64(left.ymin, right.ymin)
   st[node].ymax = max64(left.ymax, right.ymax)
}

// query to collect reachable grippers
func query(node int, cx, cy, cr int64, K int, queue *[]int) {
   nd := &st[node]
   // range index check
   if nd.l > K {
       return
   }
   // spatial prune by bounding box outside circle bounding square
   if nd.xmin > cx+cr || nd.xmax < cx-cr || nd.ymin > cy+cr || nd.ymax < cy-cr {
       return
   }
   if nd.l == nd.r {
       i := nd.l
       if visited[i] {
           return
       }
       dx := xs[i] - cx
       if dx < 0 {
           dx = -dx
       }
       dy := ys[i] - cy
       if dy < 0 {
           dy = -dy
       }
       if dx*dx+dy*dy <= cr*cr {
           // visit
           visited[i] = true
           update(1, i)
           *queue = append(*queue, i)
       }
       return
   }
   // internal node, go deeper
   query(node<<1, cx, cy, cr, K, queue)
   query(node<<1|1, cx, cy, cr, K, queue)
}

func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func max64(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var x0, y0, p0, r0 int64
   var n int
   fmt.Fscan(in, &x0, &y0, &p0, &r0, &n)
   grippers = make([]Gripper, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &grippers[i].x, &grippers[i].y,
           &grippers[i].mi, &grippers[i].pi, &grippers[i].ri)
   }
   // sort by mi
   sort.Slice(grippers, func(i, j int) bool {
       return grippers[i].mi < grippers[j].mi
   })
   xs = make([]int64, n)
   ys = make([]int64, n)
   pis = make([]int64, n)
   ris = make([]int64, n)
   mis = make([]int64, n)
   for i, g := range grippers {
       xs[i], ys[i], mis[i] = g.x, g.y, g.mi
       pis[i], ris[i] = g.pi, g.ri
   }
   visited = make([]bool, n)
   // build segment tree
   st = make([]STNode, 4*n+4)
   build(1, 0, n-1)
   // BFS
   var queue []int
   var cnt int
   // initial collect
   // find K for initial power
   K := sort.Search(n, func(i int) bool {
       return mis[i] > p0
   }) - 1
   if K >= 0 {
       query(1, x0, y0, r0, K, &queue)
   }
   // process queue
   for qi := 0; qi < len(queue); qi++ {
       idx := queue[qi]
       cnt++
       // for this gripper, collect new
       cx, cy := xs[idx], ys[idx]
       cp, cr := pis[idx], ris[idx]
       K = sort.Search(n, func(i int) bool {
           return mis[i] > cp
       }) - 1
       if K >= 0 {
           query(1, cx, cy, cr, K, &queue)
       }
   }
   // output
   fmt.Println(cnt)
}
