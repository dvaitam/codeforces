package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU for rooms
type DSU struct {
   p []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
   }
   return &DSU{p: p}
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

// Union x and y; returns new root
func (d *DSU) Union(x, y int) int {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return rx
   }
   // attach ry under rx
   d.p[ry] = rx
   return rx
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   u := make([]int, m+1)
   v := make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(in, &u[i], &v[i])
   }
   // initial owners
   initOwner := make([]int, m+1)
   // final owners
   finalOwner := make([]int, m+1)
   // person rooms
   initRoom := make(map[string]int, k)
   finalRoom := make(map[string]int, k)
   // read initial positions
   for i := 0; i < k; i++ {
       var name string
       var room, cnt int
       fmt.Fscan(in, &name, &room, &cnt)
       initRoom[name] = room
       for j := 0; j < cnt; j++ {
           var key int
           fmt.Fscan(in, &key)
           initOwner[key] = room
       }
   }
   // read final positions
   for i := 0; i < k; i++ {
       var name string
       var room, cnt int
       fmt.Fscan(in, &name, &room, &cnt)
       finalRoom[name] = room
       for j := 0; j < cnt; j++ {
           var key int
           fmt.Fscan(in, &key)
           finalOwner[key] = room
       }
   }
   // build DSU closure of openable edges
   dsu := NewDSU(n)
   // keys present in comp
   keysIn := make([][]int, n+1)
   for key := 1; key <= m; key++ {
       r := initOwner[key]
       if r >= 1 && r <= n {
           keysIn[r] = append(keysIn[r], key)
       }
   }
   used := make([]bool, m+1)
   // queue of components to process
   queue := make([]int, 0, n)
   inq := make([]bool, n+1)
   // init queue with rooms having keys
   for room := 1; room <= n; room++ {
       if len(keysIn[room]) > 0 {
           queue = append(queue, room)
           inq[room] = true
       }
   }
   head := 0
   for head < len(queue) {
       r0 := queue[head]
       head++
       inq[r0] = false
       r := dsu.Find(r0)
       // process all keys in this comp
       for _, key := range keysIn[r] {
           if used[key] {
               continue
           }
           // if key-holder in comp can reach an endpoint
           u0, v0 := u[key], v[key]
           if dsu.Find(u0) != r && dsu.Find(v0) != r {
               continue
           }
           used[key] = true
           // unlock and merge endpoints
           oldU := dsu.Find(u0)
           oldV := dsu.Find(v0)
           newR := dsu.Union(u0, v0)
           other := oldU
           if newR == oldU {
               other = oldV
           }
           // merge keysIn
           if other != newR {
               keysIn[newR] = append(keysIn[newR], keysIn[other]...)
               keysIn[other] = nil
           }
           // enqueue new comp
           if !inq[newR] {
               queue = append(queue, newR)
               inq[newR] = true
           }
       }
   }
   // validate final positions
   for name, r1 := range initRoom {
       r2, ok := finalRoom[name]
       if !ok || dsu.Find(r1) != dsu.Find(r2) {
           fmt.Println("NO")
           return
       }
   }
   for key := 1; key <= m; key++ {
       r1 := initOwner[key]
       r2 := finalOwner[key]
       if dsu.Find(r1) != dsu.Find(r2) {
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}
