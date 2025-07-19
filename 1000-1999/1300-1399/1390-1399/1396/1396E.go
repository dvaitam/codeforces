package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const N = 100005
const M = N * 2

var (
   n    int
   W    int64
   E    int
   to   [M]int
   nxt  [M]int
   first [N]int
   p    [N]int
   sizeArr [N]int
   w    [N]int
   minSum, maxSum int64
   buc  [N][]int
   ans  []pr
   send [N][]int
)

type pr struct{first, second int}
type tup struct{u, v, wi int}

func addedge(u, v int) {
   E++
   to[E] = v; nxt[E] = first[u]; first[u] = E
   E++
   to[E] = u; nxt[E] = first[v]; first[v] = E
}

func dfs1(x int) {
   sizeArr[x] = 1
   for e := first[x]; e != 0; e = nxt[e] {
       y := to[e]
       if y == p[x] {
           continue
       }
       p[y] = x
       dfs1(y)
       sizeArr[x] += sizeArr[y]
       minSum += int64(sizeArr[y] & 1)
       t := sizeArr[y]
       if n - sizeArr[y] < t {
           t = n - sizeArr[y]
       }
       w[y] = t
       maxSum += int64(t)
   }
}

// match pairs for construction
func match(cnt int, cSlice []pr, ret []tup) int {
   // copy local
   c := make([]pr, cnt)
   copy(c, cSlice[:cnt])
   sort.Slice(c, func(i, j int) bool {
       if c[i].first != c[j].first {
           return c[i].first < c[j].first
       }
       return c[i].second < c[j].second
   })
   s := 0
   for i := 0; i < cnt; i++ {
       s += c[i].first
   }
   // precondition: odd sum and largest <= (s+1)/2
   if c[cnt-1].first*2 > s+1 || s%2 == 0 {
       panic("match precondition")
   }
   L := 0
   ret[L] = tup{0, c[cnt-1].second, 1}
   L++
   s--
   c[cnt-1].first--
   cnt--
   if cnt <= 0 {
       return L
   }
   sort.Slice(c[:cnt], func(i, j int) bool {
       if c[i].first != c[j].first {
           return c[i].first < c[j].first
       }
       return c[i].second < c[j].second
   })
   ss := 0
   var i int
   var u int
   for i = 0; i < cnt; i++ {
       ss += c[i].first
       if ss >= s/2 {
           u = s/2 - (ss - c[i].first)
           // insert one slot
           c = append(c, pr{})
           copy(c[i+2:cnt+1], c[i+1:cnt])
           c[i].first = u
           c[i+1].first -= u
           cnt++
           break
       }
   }
   if i == cnt {
       panic("match split failed")
   }
   j := i + 1
   i = 0
   for s > 0 {
       for i < cnt && c[i].first == 0 {
           i++
       }
       for j < cnt && c[j].first == 0 {
           j++
       }
       if i >= cnt || j >= cnt {
           break
       }
       if c[i].first < c[j].first {
           u = c[i].first
       } else {
           u = c[j].first
       }
       c[i].first -= u
       c[j].first -= u
       s -= 2 * u
       ret[L] = tup{c[i].second, c[j].second, u}
       L++
   }
   return L
}

func dfs2(x int) {
   for e := first[x]; e != 0; e = nxt[e] {
       y := to[e]
       if p[y] == x {
           dfs2(y)
       }
   }
   // collect candidates
   var cList []pr
   for e := first[x]; e != 0; e = nxt[e] {
       y := to[e]
       if p[y] == x {
           cList = append(cList, pr{w[y], y})
       }
   }
   if x != 1 {
       cList = append(cList, pr{w[x], p[x]})
   }
   cnt := len(cList)
   ret := make([]tup, cnt+5)
   L := match(cnt, cList, ret)
   // first pairing
   for i := 0; i < L; i++ {
       u, v, wi := ret[i].u, ret[i].v, ret[i].wi
       if u == 0 {
           if v != p[x] {
               sz := len(send[v])
               y := send[v][sz-1]
               send[v] = send[v][:sz-1]
               ans = append(ans, pr{y, x})
           } else {
               send[x] = append(send[x], x)
           }
       } else {
           if u != p[x] && v != p[x] {
               for j := 0; j < wi; j++ {
                   su := len(send[u]); sv := len(send[v])
                   uu := send[u][su-1]; vv := send[v][sv-1]
                   send[u] = send[u][:su-1]
                   send[v] = send[v][:sv-1]
                   ans = append(ans, pr{uu, vv})
               }
           }
       }
   }
   // second transfer
   for i := 0; i < L; i++ {
       u, v, _ := ret[i].u, ret[i].v, ret[i].wi
       if u != 0 {
           if u != p[x] && v != p[x] {
               continue
           }
           if u == p[x] {
               u, v = v, u
           }
           // move send[u] to send[x]
           send[x] = append(send[x], send[u]...)
           send[u] = send[u][:0]
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &W)
   if n&1 == 1 {
       fmt.Fprintln(writer, "NO")
       return
   }
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       addedge(u, v)
   }
   dfs1(1)
   if W < minSum || maxSum < W || ((minSum^W)&1) == 1 {
       fmt.Fprintln(writer, "NO")
       return
   }
   for i := 2; i <= n; i++ {
       buc[w[i]] = append(buc[w[i]], i)
   }
   adj := (maxSum - W) / 2
   var S, s0, s1 int64
   var idx int
   for i := n/2; i >= 1; i-- {
       if i&1 == 1 {
           s1 += int64(len(buc[i]))
           S += s1
       } else {
           s0 += int64(len(buc[i]))
           S += s0
       }
       if S >= adj {
           idx = i
           // compute used sum for current parity
           var used int64
           if i&1 == 1 {
               used = s1
           } else {
               used = s0
           }
           uVal := adj - (S - used)
           // adjust w for large buckets
           for j := n/2; j > idx; j-- {
               for _, x := range buc[j] {
                   if ((idx ^ j) & 1) == 1 {
                       w[x] = idx - 1
                   } else {
                       w[x] = idx
                   }
               }
           }
           // fill remaining
           for j := idx; uVal > 0; j += 2 {
               for _, x := range buc[j] {
                   if uVal <= 0 {
                       break
                   }
                   w[x] = idx - 2
                   uVal--
               }
           }
           break
       }
   }
   // sum check
   var sumW int64
   for i := 2; i <= n; i++ {
       sumW += int64(w[i])
   }
   if sumW != W {
       panic("sum mismatch")
   }
   ans = make([]pr, 0, n)
   dfs2(1)
   fmt.Fprintln(writer, "YES")
   for _, e := range ans {
       fmt.Fprintln(writer, e.first, e.second)
   }
}
