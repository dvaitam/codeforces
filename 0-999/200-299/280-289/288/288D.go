package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   parent := make([]int, n+1)
   order := make([]int, 0, n)
   stack := make([]int, 0, n)
   stack = append(stack, 1)
   parent[1] = 0
   // build parent and order (pre-order)
   for i := 0; i < len(stack); i++ {
       v := stack[i]
       order = append(order, v)
       for _, u := range adj[v] {
           if u == parent[v] {
               continue
           }
           parent[u] = v
           stack = append(stack, u)
       }
   }
   // compute subtree sizes
   subtree := make([]int, n+1)
   for i := n - 1; i >= 0; i-- {
       v := order[i]
       sz := 1
       for _, u := range adj[v] {
           if u == parent[v] {
               continue
           }
           sz += subtree[u]
       }
       subtree[v] = sz
   }
   // total pairs C(n,2)
   totalPairs := new(big.Int)
   totalPairs.SetUint64(uint64(n))
   temp := new(big.Int).SetUint64(uint64(n - 1))
   totalPairs.Mul(totalPairs, temp)
   totalPairs.Div(totalPairs, big.NewInt(2))

   sum := big.NewInt(0)
   // iterate v to accumulate sum(f(v)^2 - e(v)^2)
   for v := 1; v <= n; v++ {
       // compute sum of C(size_i,2) over components after removing v
       var sumCi uint64 = 0
       for _, u := range adj[v] {
           var comp int
           if u == parent[v] {
               comp = n - subtree[v]
           } else {
               comp = subtree[u]
           }
           if comp >= 2 {
               sumCi += uint64(comp) * uint64(comp-1) / 2
           }
       }
       // f(v) = totalPairs - sumCi
       fv := new(big.Int).Set(totalPairs)
       fv.Sub(fv, new(big.Int).SetUint64(sumCi))
       // f(v)^2
       f2 := new(big.Int).Mul(fv, fv)
       // e(v): number of paths through edge parent-v
       var e2 *big.Int
       if parent[v] != 0 {
           sz := uint64(subtree[v])
           other := uint64(n) - sz
           e := sz * other
           be := new(big.Int).SetUint64(e)
           e2 = new(big.Int).Mul(be, be)
           // subtract e2
           f2.Sub(f2, e2)
       }
       sum.Add(sum, f2)
   }
   // compute totalPairs^2
   totalSq := new(big.Int).Mul(totalPairs, totalPairs)
   // answer = totalSq - sum
   ans := new(big.Int).Sub(totalSq, sum)
   // output
   writer := bufio.NewWriter(os.Stdout)
   fmt.Fprint(writer, ans.String())
   writer.Flush()
}
