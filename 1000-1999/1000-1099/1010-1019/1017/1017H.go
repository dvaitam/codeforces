package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 998244353

type Query struct {
   l, r, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   shelf := make([]int, n+1)
   // count total occurrences per film type
   // track max shelf value
   maxShelf := 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &shelf[i])
       if shelf[i] > maxShelf {
           maxShelf = shelf[i]
       }
   }
   // read queries and determine maxK
   type QInput struct{ l, r, k, id int }
   qInputs := make([]QInput, q)
   maxK := 0
   for i := 0; i < q; i++ {
       fmt.Fscan(reader, &qInputs[i].l, &qInputs[i].r, &qInputs[i].k)
       qInputs[i].id = i + 1
       if qInputs[i].k > maxK {
           maxK = qInputs[i].k
       }
   }
   // group queries by k
   queriesByK := make([][]Query, maxK+1)
   for _, qi := range qInputs {
       k := qi.k
       queriesByK[k] = append(queriesByK[k], Query{qi.l, qi.r, qi.id})
   }
   // totalCount computes frequency of each shelf value
   totalCount := make([]int, maxShelf+1)
   for i := 1; i <= n; i++ {
       totalCount[shelf[i]]++
   }
   // compute inverses up to n + maxK + 5
   invSize := n + maxK + 5
   invArr := make([]int64, invSize)
   invArr[1] = 1
   for i := 2; i < invSize; i++ {
       invArr[i] = invArr[int(mod)%i] * (mod - int64(mod)/int64(i)) % mod
   }
   // block size for Mo's algorithm
   const blockSize = 1000
   blockIdx := make([]int, n+1)
   for i := 1; i <= n; i++ {
       blockIdx[i] = (i-1)/blockSize + 1
   }
   // answers
   answer := make([]int64, q+1)

   // process queries for each k
   for k := 0; k <= maxK; k++ {
       qs := queriesByK[k]
       if len(qs) == 0 {
           continue
       }
       // cnt per film type in current window
       cnt := make([]int, maxShelf+1)
       currentK := k
       currentProd := int64(1)
       zz := int64(m) * int64(currentK) % mod
       // precompute pArr
       pArr := make([]int64, n+2)
       pArr[n] = 1
       for i := n - 1; i >= 1; i-- {
           pArr[i] = pArr[i+1] * (zz + int64(n-i)) % mod
       }
       // sort queries by Mo's order
       sort.Slice(qs, func(i, j int) bool {
           bi := blockIdx[qs[i].l]
           bj := blockIdx[qs[j].l]
           if bi != bj {
               return bi < bj
           }
           return qs[i].r < qs[j].r
       })
       currentL, currentR := 1, 0
       // add and remove functions
       add := func(x int) {
           v := shelf[x]
           currentProd = currentProd * int64(totalCount[v]+currentK-cnt[v]) % mod
           cnt[v]++
       }
       remove := func(x int) {
           v := shelf[x]
           cnt[v]--
           currentProd = currentProd * invArr[totalCount[v]+currentK-cnt[v]] % mod
       }
       // process each query
       for _, qu := range qs {
           L, R := qu.l, qu.r
           for currentL > L {
               currentL--
               add(currentL)
           }
           for currentR < R {
               currentR++
               add(currentR)
           }
           for currentL < L {
               remove(currentL)
               currentL++
           }
           for currentR > R {
               remove(currentR)
               currentR--
           }
           length := currentR - currentL + 1
           ans := currentProd * pArr[length] % mod
           answer[qu.id] = ans
       }
   }
   // output answers
   for i := 1; i <= q; i++ {
       fmt.Fprintln(writer, answer[i])
   }
}
