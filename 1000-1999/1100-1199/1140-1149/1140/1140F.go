package main

import (
   "bufio"
   "bytes"
   "fmt"
   "os"
)

const K = 300000

type pair struct{ x, y int }
type op struct{ u, v, height int }

var (
   n int
   seg [][]pair
   parent, heightArr, Xcnt, Ycnt []int
   rollback []op
   cur int64
   ans []int64
)

// add edge (x,y) active in interval [s,t] into segment tree
func add(node, l, r, s, t, x, y int) {
   if s > t {
       return
   }
   if s <= l && r <= t {
       seg[node] = append(seg[node], pair{x, y})
       return
   }
   mid := (l + r) >> 1
   if t <= mid {
       add(node*2, l, mid, s, t, x, y)
   } else if s > mid {
       add(node*2+1, mid+1, r, s, t, x, y)
   } else {
       add(node*2, l, mid, s, mid, x, y)
       add(node*2+1, mid+1, r, mid+1, t, x, y)
   }
}

func find(u int) int {
   for parent[u] != u {
       u = parent[u]
   }
   return u
}

func dfs(node, l, r int) {
   before := len(rollback)
   // process edges
   for _, e := range seg[node] {
       a, b := e.x, e.y
       u := find(a)
       v := find(K + b)
       if u == v {
           continue
       }
       if heightArr[u] > heightArr[v] {
           u, v = v, u
       }
       // record rollback info
       rollback = append(rollback, op{u, v, heightArr[v]})
       // update cur
       cur -= int64(Xcnt[u]) * int64(Ycnt[u])
       cur -= int64(Xcnt[v]) * int64(Ycnt[v])
       // merge u into v
       parent[u] = v
       Xcnt[v] += Xcnt[u]
       Ycnt[v] += Ycnt[u]
       if heightArr[v] < heightArr[u]+1 {
           heightArr[v] = heightArr[u] + 1
       }
       // new contribution
       cur += int64(Xcnt[v]) * int64(Ycnt[v])
   }
   if l == r {
       ans[l] = cur
   } else {
       mid := (l + r) >> 1
       dfs(node*2, l, mid)
       dfs(node*2+1, mid+1, r)
   }
   // rollback
   for len(rollback) > before {
       opv := rollback[len(rollback)-1]
       rollback = rollback[:len(rollback)-1]
       u, v := opv.u, opv.v
       // remove merged contribution
       cur -= int64(Xcnt[v]) * int64(Ycnt[v])
       // rollback u
       parent[u] = u
       Xcnt[v] -= Xcnt[u]
       Ycnt[v] -= Ycnt[u]
       heightArr[v] = opv.height
       // restore contributions
       cur += int64(Xcnt[u]) * int64(Ycnt[u])
       cur += int64(Xcnt[v]) * int64(Ycnt[v])
   }
}

// fast IO
var rdr = bufio.NewReader(os.Stdin)
var wrtr = bufio.NewWriter(os.Stdout)

func readInt() int {
   var c byte
   var x int
   c, _ = rdr.ReadByte()
   for c < '0' || c > '9' {
       c, _ = rdr.ReadByte()
   }
   for c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, _ = rdr.ReadByte()
   }
   return x
}

func main() {
   defer wrtr.Flush()
   n = readInt()
   seg = make([][]pair, 4*n+4)
   ans = make([]int64, n+2)
   parent = make([]int, 2*K+2)
   heightArr = make([]int, 2*K+2)
   Xcnt = make([]int, 2*K+2)
   Ycnt = make([]int, 2*K+2)
   // map for active intervals
   m := make(map[[2]int]int)
   // read inputs
   for i := 1; i <= n; i++ {
       x := readInt()
       y := readInt()
       key := [2]int{x, y}
       if t, ok := m[key]; ok {
           add(1, 1, n, t, i-1, x, y)
           delete(m, key)
       } else {
           m[key] = i
       }
   }
   // remaining active till n
   for key, t := range m {
       add(1, 1, n, t, n, key[0], key[1])
   }
   // init dsu
   for i := 1; i <= 2*K; i++ {
       parent[i] = i
       heightArr[i] = 1
       if i <= K {
           Xcnt[i] = 1
       } else {
           Ycnt[i] = 1
       }
   }
   cur = 0
   // run
   dfs(1, 1, n)
   // output
   var buf bytes.Buffer
   for i := 1; i <= n; i++ {
       buf.WriteString(fmt.Sprintf("%d", ans[i]))
       if i < n {
           buf.WriteByte(' ')
       }
   }
   buf.WriteByte('\n')
   wrtr.Write(buf.Bytes())
}
