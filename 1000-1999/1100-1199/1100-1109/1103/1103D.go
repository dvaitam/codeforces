package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   defer writer.Flush()
   // read n, k
   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // read a and e
   a := make([]int64, n)
   e := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &e[i])
   }
   // compute gcd of all a
   g := a[0]
   for i := int64(1); i < n; i++ {
       g = gcd(g, a[i])
   }
   // prime factors of g
   var p []int64
   tmp := g
   for i := int64(2); i*i <= tmp; i++ {
       if tmp%i == 0 {
           p = append(p, i)
           for tmp%i == 0 {
               tmp /= i
           }
       }
   }
   if tmp > 1 {
       p = append(p, tmp)
   }
   m := len(p)
   // group by prime-power vector
   type group struct {
       cur   []int64
       costs []int64
   }
   groups := []group{}
   index := make(map[string]int)
   for i := int64(0); i < n; i++ {
       cur := make([]int64, m)
       for j := 0; j < m; j++ {
           cur[j] = 1
           for a[i]%p[j] == 0 {
               a[i] /= p[j]
               cur[j] *= p[j]
           }
       }
       // key
       key := ""
       for j, v := range cur {
           if j > 0 {
               key += ","
           }
           key += strconv.FormatInt(v, 10)
       }
       idx, ok := index[key]
       if !ok {
           idx = len(groups)
           index[key] = idx
           // copy cur
           vec := make([]int64, m)
           copy(vec, cur)
           groups = append(groups, group{cur: vec, costs: nil})
       }
       groups[idx].costs = append(groups[idx].costs, e[i])
   }
   // dp dimensions
   maskSize := 1 << m
   const INF = int64(9e18)
   // initialize dp
   dp := make([][]int64, maskSize)
   for i := 0; i < maskSize; i++ {
       dp[i] = make([]int64, m+1)
       for j := 0; j <= m; j++ {
           dp[i][j] = INF
       }
   }
   dp[0][0] = 0
   // process groups
   for _, gr := range groups {
       // sort costs
       sort.Slice(gr.costs, func(i, j int) bool { return gr.costs[i] < gr.costs[j] })
       // build f array
       f := make([]int, maskSize)
       // dfs to mark valid subsets
       var dfs func(cur []int64, num int64, subset, pos int)
       dfs = func(cur []int64, num int64, subset, pos int) {
           if pos == len(cur) {
               f[subset] = 1
               return
           }
           dfs(cur, num, subset, pos+1)
           if num <= k/cur[pos] {
               dfs(cur, num*cur[pos], subset|(1<<pos), pos+1)
           }
       }
       dfs(gr.cur, 1, 0, 0)
       // accumulate supersets
       for bit := 0; bit < m; bit++ {
           for mask := 0; mask < maskSize; mask++ {
               if mask&(1<<bit) == 0 {
                   f[mask] += f[mask|(1<<bit)]
               }
           }
       }
       // collect op masks
       var op []int
       for mask := 1; mask < maskSize; mask++ {
           if f[mask] == 1 {
               op = append(op, mask)
           }
       }
       // number of items to consider
       l := len(gr.costs)
       if l > m {
           l = m
       }
       // iterate costs
       for s := 0; s < l; s++ {
           improved := false
           // copy dp to cp
           cp := make([][]int64, maskSize)
           for i := 0; i < maskSize; i++ {
               cp[i] = make([]int64, m+1)
               copy(cp[i], dp[i])
           }
           cost := gr.costs[s]
           for _, t := range op {
               for mask := 0; mask < maskSize; mask++ {
                   for j := 0; j < m; j++ {
                       if dp[mask][j] == INF {
                           continue
                       }
                       newMask := mask | t
                       newCnt := j + 1
                       newCost := dp[mask][j] + cost
                       if newCost < cp[newMask][newCnt] {
                           cp[newMask][newCnt] = newCost
                           improved = true
                       }
                   }
               }
           }
           dp = cp
           if !improved {
               break
           }
       }
   }
   // compute answer
   full := maskSize - 1
   var ans int64 = INF
   for j := 0; j <= m; j++ {
       v := dp[full][j]
       if v < INF {
           cand := v * int64(j)
           if cand < ans {
               ans = cand
           }
       }
   }
   if ans >= INF {
       fmt.Fprintln(writer, -1)
   } else {
       fmt.Fprintln(writer, ans)
   }
}
