package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU for land cells, tracking if a component touches left or right boundary
type DSU struct {
   parent []int32
   left   []bool
   right  []bool
}

func NewDSU(n int) *DSU {
   return &DSU{
       parent: make([]int32, n),
       left:   make([]bool, n),
       right:  make([]bool, n),
   }
}

func (d *DSU) Find(x int32) int32 {
   for d.parent[x] != x {
       d.parent[x] = d.parent[d.parent[x]]
       x = d.parent[x]
   }
   return x
}

// Union components of x and y, merging flags
func (d *DSU) Union(x, y int32) {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return
   }
   // attach ry under rx
   d.parent[ry] = rx
   d.left[rx] = d.left[rx] || d.left[ry]
   d.right[rx] = d.right[rx] || d.right[ry]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var r, c, n int
   fmt.Fscan(reader, &r, &c, &n)
   ops := make([][2]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &ops[i][0], &ops[i][1])
   }

   // idx maps cell to DSU id (1-based), zero means sea
   rc := r * c
   idx := make([]int32, rc)
   dsu := NewDSU(n + 5)
   var curID int32 = 0
   ans := 0

   // direction offsets: up, down, left, right (with wrap for left/right)
   for _, op := range ops {
       i := op[0] - 1
       j := op[1] - 1
       pos := i*c + j
       // compute boundary flags for this cell
       newLeft := (j == 0)
       newRight := (j == c-1)
       // collect unique neighbor roots and aggregate flags
       var orLeft, orRight bool
       neigh := [4]int32{-1, -1, -1, -1}
       cnt := 0
       // up, down
       for di := -1; di <= 1; di += 2 {
           ni := i + di
           if ni < 0 || ni >= r {
               continue
           }
           nj := j
           p2 := ni*c + nj
           id2 := idx[p2]
           if id2 == 0 {
               continue
           }
           root := dsu.Find(id2)
           // check duplicate
           dup := false
           for k := 0; k < cnt; k++ {
               if neigh[k] == root {
                   dup = true
                   break
               }
           }
           if dup {
               continue
           }
           neigh[cnt] = root
           cnt++
           if dsu.left[root] {
               orLeft = true
           }
           if dsu.right[root] {
               orRight = true
           }
       }
       // left, right with wrap
       for _, dj := range []int{-1, 1} {
           ni := i
           nj := j + dj
           if nj < 0 {
               nj = c - 1
           } else if nj >= c {
               nj = 0
           }
           p2 := ni*c + nj
           id2 := idx[p2]
           if id2 == 0 {
               continue
           }
           root := dsu.Find(id2)
           dup := false
           for k := 0; k < cnt; k++ {
               if neigh[k] == root {
                   dup = true
                   break
               }
           }
           if dup {
               continue
           }
           neigh[cnt] = root
           cnt++
           if dsu.left[root] {
               orLeft = true
           }
           if dsu.right[root] {
               orRight = true
           }
       }
       totalLeft := newLeft || orLeft
       totalRight := newRight || orRight
       // if adding this land creates barrier, skip
       if totalLeft && totalRight {
           continue
       }
       // accept removal
       curID++
       ans++
       id := curID
       idx[pos] = id
       dsu.parent[id] = id
       dsu.left[id] = newLeft
       dsu.right[id] = newRight
       // union with neighbors
       for k := 0; k < cnt; k++ {
           dsu.Union(id, neigh[k])
       }
   }
   fmt.Fprintln(writer, ans)
}
