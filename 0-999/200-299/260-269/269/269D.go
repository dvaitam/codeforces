package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Edge represents a directed edge with weight
type Edge struct {
   to int
   w  int64
}

var (
   n        int
   t        int64
   panels   []Panel
   coords   []int64
   segLen   []int64
   prefix   []int64
   treeVal  []int
   edges    [][]Edge
)

// Panel represents a horizontal panel
type Panel struct {
   h  int64
   l  int64
   r  int64
   id int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   if _, err := fmt.Fscan(in, &n, &t); err != nil {
       return
   }
   panels = make([]Panel, n+2)
   // top panel id=0
   panels[0] = Panel{h: t, l: -1000000000, r: 1000000000, id: 0}
   // input panels id=1..n
   for i := 1; i <= n; i++ {
       var hi, li, ri int64
       fmt.Fscan(in, &hi, &li, &ri)
       panels[i] = Panel{h: hi, l: li, r: ri, id: i}
   }
   // bottom panel id=n+1
   bottomID := n + 1
   panels = append(panels, Panel{h: 0, l: -1000000000, r: 1000000000, id: bottomID})
   // collect coords
   coords = make([]int64, 0, 2*(n+2))
   for _, p := range panels {
       coords = append(coords, p.l, p.r)
   }
   sort.Slice(coords, func(i, j int) bool { return coords[i] < coords[j] })
   uniq := coords[:0]
   for i, x := range coords {
       if i == 0 || x != coords[i-1] {
           uniq = append(uniq, x)
       }
   }
   coords = uniq
   // build segment lengths
   m := len(coords) - 1
   segLen = make([]int64, m)
   for i := 0; i < m; i++ {
       segLen[i] = coords[i+1] - coords[i]
   }
   prefix = make([]int64, m+1)
   for i := 0; i < m; i++ {
       prefix[i+1] = prefix[i] + segLen[i]
   }
   // prepare edges
   totalNodes := n + 2
   edges = make([][]Edge, totalNodes)
   // segment tree for cover values
   treeVal = make([]int, 4*m)
   // initial cover = top (0)
   if m > 0 {
       treeVal[1] = 0
   }
   // sort panels 1..n by descending height
   order := make([]Panel, 0, n)
   for i := 1; i <= n; i++ {
       order = append(order, panels[i])
   }
   sort.Slice(order, func(i, j int) bool { return order[i].h > order[j].h })
   // process each panel
   for _, p := range order {
       l := sort.Search(len(coords), func(i int) bool { return coords[i] >= p.l })
       r := sort.Search(len(coords), func(i int) bool { return coords[i] >= p.r })
       if l < r {
           processRange(1, 0, m, l, r, p.id)
       }
   }
   // connect bottom
   if m > 0 {
       processRange(1, 0, m, 0, m, bottomID)
   }
   // DP over nodes sorted by descending height
   nodes := make([]Panel, 0, totalNodes)
   for i := range panels {
       nodes = append(nodes, panels[i])
   }
   sort.Slice(nodes, func(i, j int) bool { return nodes[i].h > nodes[j].h })
   dp := make([]int64, totalNodes)
   const INF = int64(4e18)
   dp[0] = INF
   for _, u := range nodes {
       uid := u.id
       for _, e := range edges[uid] {
           flow := dp[uid]
           if flow > e.w {
               flow = e.w
           }
           if flow > dp[e.to] {
               dp[e.to] = flow
           }
       }
   }
   fmt.Println(dp[bottomID])
}

// processRange splits segments in [ql,qr) and connects previous covers to curID
func processRange(idx, l, r, ql, qr, curID int) {
   if r <= ql || l >= qr {
       return
   }
   if ql <= l && r <= qr && treeVal[idx] != -1 {
       parent := treeVal[idx]
       length := prefix[r] - prefix[l]
       edges[parent] = append(edges[parent], Edge{to: curID, w: length})
       treeVal[idx] = curID
       return
   }
   if r-l == 1 {
       // leaf but mixed; should not happen
       if treeVal[idx] != curID {
           treeVal[idx] = curID
       }
       return
   }
   // push down
   if treeVal[idx] != -1 {
       treeVal[idx*2] = treeVal[idx]
       treeVal[idx*2+1] = treeVal[idx]
       treeVal[idx] = -1
   }
   mid := (l + r) >> 1
   processRange(idx*2, l, mid, ql, qr, curID)
   processRange(idx*2+1, mid, r, ql, qr, curID)
   leftv := treeVal[idx*2]
   rightv := treeVal[idx*2+1]
   if leftv == rightv {
       treeVal[idx] = leftv
   } else {
       treeVal[idx] = -1
   }
}
