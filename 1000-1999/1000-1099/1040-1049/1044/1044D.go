package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var q int
   fmt.Fscan(in, &q)
   // DSU with xor potentials on prefix indices
   idMap := make(map[int]int)
   parent := make([]int, 0, q*2)
   xorToParent := make([]int, 0, q*2)
   rank := make([]int, 0, q*2)
   nextID := 0
   // initialize node for prefix index -1 (empty prefix)
   idMap[-1] = nextID
   parent = append(parent, nextID)
   xorToParent = append(xorToParent, 0)
   rank = append(rank, 0)
   nextID++

   var last int
   // closure to get or create DSU node for a prefix
   getID := func(pos int) int {
       if id, ok := idMap[pos]; ok {
           return id
       }
       id := nextID
       idMap[pos] = id
       parent = append(parent, id)
       xorToParent = append(xorToParent, 0)
       rank = append(rank, 0)
       nextID++
       return id
   }
   // recursive find with path compression
   var find func(int) (int, int)
   find = func(u int) (int, int) {
       if parent[u] != u {
           p := parent[u]
           r, xr := find(p)
           xorToParent[u] ^= xr
           parent[u] = r
       }
       return parent[u], xorToParent[u]
   }
   // union u and v with xor(u,v) == w; returns false if conflict
   union := func(u, v, w int) bool {
       ru, xu := find(u)
       rv, xv := find(v)
       if ru == rv {
           return (xu ^ xv) == w
       }
       // union by rank
       if rank[ru] < rank[rv] {
           ru, rv = rv, ru
           xu, xv = xv, xu
       }
       parent[rv] = ru
       xorToParent[rv] = xu ^ xv ^ w
       if rank[ru] == rank[rv] {
           rank[ru]++
       }
       return true
   }

   for i := 0; i < q; i++ {
       var t int
       fmt.Fscan(in, &t)
       if t == 1 {
           var lp, rp, xp int
           fmt.Fscan(in, &lp, &rp, &xp)
           l := lp ^ last
           r := rp ^ last
           x := xp ^ last
           if l > r {
               l, r = r, l
           }
           u := getID(l - 1)
           v := getID(r)
           // ignore conflicting updates
           union(u, v, x)
       } else {
           var lp, rp int
           fmt.Fscan(in, &lp, &rp)
           l := lp ^ last
           r := rp ^ last
           if l > r {
               l, r = r, l
           }
           u := getID(l - 1)
           v := getID(r)
           ru, xu := find(u)
           rv, xv := find(v)
           var ans int
           if ru != rv {
               ans = -1
           } else {
               ans = xu ^ xv
           }
           // output and update last
           if ans < 0 {
               fmt.Fprintln(out, -1)
               last = 1
           } else {
               fmt.Fprintln(out, ans)
               last = ans
           }
       }
   }
}
