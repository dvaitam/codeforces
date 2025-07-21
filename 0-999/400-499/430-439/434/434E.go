package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n int
   modY, k, targetX int64
   vals []int64
   adj [][]int
   used []bool
   sz []int
   kPows, invKpows []int64
   invK int64
   ansV1, ansV2 int64
)

// fast exponentiation
func modPow(a, e, m int64) int64 {
   res := int64(1)
   a %= m
   for e > 0 {
       if e&1 == 1 {
           res = (res * a) % m
       }
       a = (a * a) % m
       e >>= 1
   }
   return res
}

// compute subtree sizes
func dfsSize(u, p int) {
   sz[u] = 1
   for _, v := range adj[u] {
       if v != p && !used[v] {
           dfsSize(v, u)
           sz[u] += sz[v]
       }
   }
}

// find centroid
func findCentroid(u, p, total int) int {
   for _, v := range adj[u] {
       if v != p && !used[v] && sz[v] > total/2 {
           return findCentroid(v, u, total)
       }
   }
   return u
}

// data for nodes
type nodeData struct { hF, hB int64; d int }

// collect data from subtree
func dfsCollect(u, p, depth int, hF, hB int64, data *[]nodeData) {
   hF = (hF*k + vals[u]) % modY
   hB = (hB + kPows[depth]*vals[u]) % modY
   *data = append(*data, nodeData{hF, hB, depth})
   for _, v := range adj[u] {
       if v != p && !used[v] {
           dfsCollect(v, u, depth+1, hF, hB, data)
       }
   }
}

// process at centroid c
func processCentroid(c int) {
   // counters for V1
   var prevFx, prevBx int64
   // maps for V2
   prevF := make(map[int64]int)
   // list for prev nodes for V2 case1
   var prevList []nodeData

   // include the centroid itself as a trivial path length 0
   // path c->c: hF=vals[c], hB=vals[c], d=0
   // but for u or w equal c, handled automatically by including here
   prevList = append(prevList, nodeData{vals[c] % modY, vals[c] % modY, 0})
   if vals[c]%modY != targetX {
       // f(c,c)!=x, so for V2, u=c can be used
   }
   // track prevF for u=c if hF!=x
   if vals[c]%modY != targetX {
       prevF[vals[c]%modY] = 1
   }
   // for V1, update prev counts
   if vals[c]%modY == targetX {
       prevFx = 1; prevBx = 1
   } else {
       prevFx = 0; prevBx = 0
   }

   for _, v := range adj[c] {
       if used[v] {
           continue
       }
       var curr []nodeData
       dfsCollect(v, c, 1, vals[c]%modY, vals[c]%modY, &curr)

       // V1: count f(u,c)==x and f(c,w)==x
       var currFx, currBx int64
       for _, nd := range curr {
           if nd.hF == targetX {
               currFx++
           }
           if nd.hB == targetX {
               currBx++
           }
       }
       if vals[c]%modY != targetX {
           ansV1 += prevBx * currFx
           ansV1 += prevFx * currBx
       }
       prevFx += currFx
       prevBx += currBx

       // prepare map of current hF counts for V2
       currF := make(map[int64]int)
       for _, nd := range curr {
           if nd.hF != targetX {
               currF[nd.hF]++
           }
       }

       // V2 case: u in prev, w in curr
       for _, udata := range prevList {
           if udata.hB == targetX {
               continue
           }
           t := (targetX - udata.hB) % modY
           if t < 0 {
               t += modY
           }
           t = (t * invKpows[udata.d]) % modY
           t = (t + vals[c]) % modY
           if t != targetX {
               if cnt, ok := currF[t]; ok {
                   ansV2 += int64(cnt)
               }
           }
       }
       // V2 case: w in prev (as u), u in curr
       for _, nd := range curr {
           if nd.hB == targetX {
               continue
           }
           t := (targetX - nd.hB) % modY
           if t < 0 {
               t += modY
           }
           t = (t * invKpows[nd.d]) % modY
           t = (t + vals[c]) % modY
           if t != targetX {
               if cnt, ok := prevF[t]; ok {
                   ansV2 += int64(cnt)
               }
           }
       }

       // merge curr into prev
       for _, nd := range curr {
           prevList = append(prevList, nd)
           if nd.hF != targetX {
               prevF[nd.hF]++
           }
       }
   }
}

// centroid decomposition
func decompose(u int) {
   dfsSize(u, -1)
   c := findCentroid(u, -1, sz[u])
   used[c] = true
   processCentroid(c)
   for _, v := range adj[c] {
       if !used[v] {
           decompose(v)
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &modY, &k, &targetX)
   vals = make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &vals[i])
   }
   adj = make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--; v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   used = make([]bool, n)
   sz = make([]int, n)
   // precompute powers
   kPows = make([]int64, n+2)
   invKpows = make([]int64, n+2)
   kPows[0] = 1
   for i := 1; i <= n+1; i++ {
       kPows[i] = (kPows[i-1] * k) % modY
   }
   invK = modPow(k, modY-2, modY)
   invKpows[0] = 1
   for i := 1; i <= n+1; i++ {
       invKpows[i] = (invKpows[i-1] * invK) % modY
   }
   // decompose
   decompose(0)
   total := int64(n)
   totTrip := total * total * total
   ans := totTrip - ansV1 - ansV2
   fmt.Fprint(writer, ans)
}
