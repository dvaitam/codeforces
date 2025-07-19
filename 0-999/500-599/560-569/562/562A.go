package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const maxn = 200100

var (
   e    [maxn][2]int
   ed   [maxn]float64
   restr [maxn]bool
   g    [maxn][]int
   wv   [maxn]float64
   cnt  [maxn]int
)

func getOth(id, v int) int {
   if e[id][0] == v {
       return e[id][1]
   }
   return e[id][0]
}

func getDer(v int, dist float64, par int) float64 {
   ans := math.Pow(dist, 0.5) * wv[v]
   for _, eid := range g[v] {
       nv := getOth(eid, v)
       if nv == par {
           continue
       }
       ans += getDer(nv, dist+ed[eid], v)
   }
   return ans
}

func getSum(v int, dist float64, par int) float64 {
   ans := math.Pow(dist, 1.5) * wv[v]
   for _, eid := range g[v] {
       nv := getOth(eid, v)
       if nv == par {
           continue
       }
       ans += getSum(nv, dist+ed[eid], v)
   }
   return ans
}

func getSize(v, par int) int {
   cnt[v] = 1
   for _, eid := range g[v] {
       nv := getOth(eid, v)
       if nv == par || restr[eid] {
           continue
       }
       cnt[v] += getSize(nv, v)
   }
   return cnt[v]
}

func getCentr(v, gcnt, par int) int {
   for _, eid := range g[v] {
       nv := getOth(eid, v)
       if nv == par || restr[eid] {
           continue
       }
       if cnt[nv]*2 > gcnt {
           return getCentr(nv, gcnt, v)
       }
   }
   return v
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &wv[i])
   }
   for i := 0; i < n-1; i++ {
       var a, b int
       var l float64
       fmt.Fscan(in, &a, &b, &l)
       a--
       b--
       e[i][0], e[i][1] = a, b
       ed[i] = l
       g[a] = append(g[a], i)
       g[b] = append(g[b], i)
   }
   root := 0
   for {
       gcnt := getSize(root, -1)
       if gcnt == 1 {
           break
       }
       root = getCentr(root, gcnt, -1)
       mx := 0.0
       mxid := -1
       for _, eid := range g[root] {
           if restr[eid] {
               continue
           }
           nv := getOth(eid, root)
           cur := getDer(nv, ed[eid], root)
           if cur > mx {
               mx = cur
               mxid = eid
           }
       }
       if mxid == -1 {
           break
       }
       newRoot := getOth(mxid, root)
       if getSum(root, 0, -1) <= getSum(newRoot, 0, -1) {
           break
       }
       root = newRoot
       restr[mxid] = true
   }
   res := getSum(root, 0, -1)
   fmt.Printf("%d %.10f\n", root+1, res)
}
