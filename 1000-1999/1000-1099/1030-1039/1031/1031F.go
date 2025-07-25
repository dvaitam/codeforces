package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

func main() {
   const N = 1000000
   // smallest prime factor
   spf := make([]int, N+1)
   for i := 2; i <= N; i++ {
       if spf[i] == 0 {
           for j := i; j <= N; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   // compute divisor counts up to N
   d := make([]int, N+1)
   d[1] = 1
   for i := 2; i <= N; i++ {
       p := spf[i]
       cnt := 0
       x := i
       for x%p == 0 {
           x /= p
           cnt++
       }
       // now x = i / p^cnt
       d[i] = d[x] * (cnt + 1)
   }
   // max divisor count
   maxD := 0
   for i := 1; i <= N; i++ {
       if d[i] > maxD {
           maxD = d[i]
       }
   }
   M := maxD
   // generate patterns of exponents for D <= M
   type pat struct{ exps []int; D int }
   patterns := make([]pat, 0)
   patMap := make(map[string]int)
   // add empty pattern for D=1
   patterns = append(patterns, pat{exps: []int{}, D: 1})
   patMap[""] = 0
   // helper to record pattern
   var record = func(exps []int, D int) int {
       var key string
       // build key
       for i, e := range exps {
           if i > 0 {
               key += ","
           }
           key += strconv.Itoa(e)
       }
       if id, ok := patMap[key]; ok {
           return id
       }
       id := len(patterns)
       // copy exps
       c := make([]int, len(exps))
       copy(c, exps)
       patterns = append(patterns, pat{exps: c, D: D})
       patMap[key] = id
       return id
   }
   // dfs generate
   var dfs func(cur []int, lastE, curP int)
   dfs = func(cur []int, lastE, curP int) {
       for e := lastE; e >= 1; e-- {
           np := curP * (e + 1)
           if np > M {
               continue
           }
           // new exps sorted descending: cur is sorted desc and e <= lastE
           nxt := append(cur, e)
           // no need to sort, condition ensures desc order
           id := record(nxt, np)
           dfs(nxt, e, np)
       }
   }
   dfs([]int{}, M-1, 1)
   P := len(patterns)
   // build neighbors
   nbr := make([][]int, P)
   // temp key builder
   for i, p := range patterns {
       exps := p.exps
       // decrement operations
       for j, e := range exps {
           var nxt []int
           if e > 1 {
               nxt = make([]int, len(exps))
               copy(nxt, exps)
               nxt[j] = e - 1
           } else {
               // remove element
               nxt = append([]int{}, exps[:j]...)
               nxt = append(nxt, exps[j+1:]...)
           }
           // ensure descending order
           sort.Slice(nxt, func(a, b int) bool { return nxt[a] > nxt[b] })
           // compute key and find id
           var key string
           for k, ee := range nxt {
               if k > 0 { key += "," }
               key += strconv.Itoa(ee)
           }
           if id, ok := patMap[key]; ok {
               nbr[i] = append(nbr[i], id)
           }
       }
       // increment operations
       for j, e := range exps {
           var nxt = make([]int, len(exps))
           copy(nxt, exps)
           nxt[j] = e + 1
           // check product
           prod := 1
           for _, ee := range nxt {
               prod *= (ee + 1)
               if prod > M { break }
           }
           if prod <= M {
               // sort desc
               sort.Slice(nxt, func(i, j int) bool { return nxt[i] > nxt[j] })
               var key string
               for k, ee := range nxt {
                   if k > 0 { key += "," }
                   key += strconv.Itoa(ee)
               }
               if id, ok := patMap[key]; ok {
                   nbr[i] = append(nbr[i], id)
               }
           }
       }
       // add new prime (exp=1)
       // only if product*2<=M
       if p.D*2 <= M {
           nxt := append([]int{}, exps...)
           nxt = append(nxt, 1)
           sort.Slice(nxt, func(i, j int) bool { return nxt[i] > nxt[j] })
           var key string
           for k, ee := range nxt {
               if k > 0 { key += "," }
               key += strconv.Itoa(ee)
           }
           if id, ok := patMap[key]; ok {
               nbr[i] = append(nbr[i], id)
           }
       }
       // dedupe neighbors
       if len(nbr[i]) > 1 {
           sort.Ints(nbr[i])
           u := nbr[i][:1]
           for _, v := range nbr[i][1:] {
               if v != u[len(u)-1] {
                   u = append(u, v)
               }
           }
           nbr[i] = u
       }
   }
   // BFS per D
   INF := 1<<30
   dist := make([][]int, M+1)
   for D := 1; D <= M; D++ {
       dist[D] = make([]int, P)
       for i := 0; i < P; i++ {
           dist[D][i] = INF
       }
       // multi-source
       var q []int
       for i, p := range patterns {
           if p.D == D {
               dist[D][i] = 0
               q = append(q, i)
           }
       }
       // BFS
       for qi := 0; qi < len(q); qi++ {
           u := q[qi]
           du := dist[D][u]
           for _, v := range nbr[u] {
               if dist[D][v] > du+1 {
                   dist[D][v] = du + 1
                   q = append(q, v)
               }
           }
       }
   }
   // read input and answer
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var tc int
   fmt.Fscan(in, &tc)
   for ; tc > 0; tc-- {
       var a, b int
       fmt.Fscan(in, &a, &b)
       patA := sigToPat(a, spf, patMap)
       patB := sigToPat(b, spf, patMap)
       ans := INF
       for D := 1; D <= M; D++ {
           da := dist[D][patA]
           db := dist[D][patB]
           if da+db < ans {
               ans = da + db
           }
       }
       fmt.Fprintln(out, ans)
   }
}

// sigToPat computes signature pattern id for n
func sigToPat(n int, spf []int, patMap map[string]int) int {
   var exps []int
   for n > 1 {
       p := spf[n]
       cnt := 0
       for n%p == 0 {
           n /= p
           cnt++
       }
       exps = append(exps, cnt)
   }
   if len(exps) == 0 {
       return patMap[""]
   }
   sort.Slice(exps, func(i, j int) bool { return exps[i] > exps[j] })
   var key string
   for i, e := range exps {
       if i > 0 { key += "," }
       key += strconv.Itoa(e)
   }
   return patMap[key]
}
