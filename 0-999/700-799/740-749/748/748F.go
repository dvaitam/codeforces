package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// groupHeap is a max-heap of group indices by remaining size
type groupHeap []int

func (h groupHeap) Len() int { return len(h) }
func (h groupHeap) Less(i, j int) bool {
   // max-heap: larger group first
   return len(groups[h[i]]) > len(groups[h[j]])
}
func (h groupHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *groupHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *groupHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

var (
   n, k int
   adj [][]int
   special []bool
   totalSpecial int
   parent []int
   order []int
   subCnt []int
   groups [][]int
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &k)
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   special = make([]bool, n+1)
   totalSpecial = 2 * k
   spe := make([]int, totalSpecial)
   for i := 0; i < totalSpecial; i++ {
       fmt.Fscan(in, &spe[i])
       special[spe[i]] = true
   }
   // DFS order from 1 to compute parent and subtree counts
   parent = make([]int, n+1)
   order = make([]int, 0, n)
   stack := []int{1}
   parent[1] = 0
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       order = append(order, u)
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           stack = append(stack, v)
       }
   }
   subCnt = make([]int, n+1)
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       cnt := 0
       if special[u] {
           cnt = 1
       }
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           cnt += subCnt[v]
       }
       subCnt[u] = cnt
   }
   // find centroid: node with max component <= k
   centroid := 1
   for u := 1; u <= n; u++ {
       maxc := totalSpecial - subCnt[u]
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           if subCnt[v] > maxc {
               maxc = subCnt[v]
           }
       }
       if maxc <= k {
           centroid = u
           break
       }
   }
   // collect special nodes in groups by subtree of centroid
   groups = [][]int{}
   // if centroid is special, its own group
   if special[centroid] {
       groups = append(groups, []int{centroid})
   }
   for _, v := range adj[centroid] {
       // collect special in subtree v, avoid going back to centroid
       temp := make([]int, 0)
       stU := []int{v}
       stP := []int{centroid}
       for len(stU) > 0 {
           u := stU[len(stU)-1]
           stU = stU[:len(stU)-1]
           p := stP[len(stP)-1]
           stP = stP[:len(stP)-1]
           if special[u] {
               temp = append(temp, u)
           }
           for _, w := range adj[u] {
               if w == p {
                   continue
               }
               stU = append(stU, w)
               stP = append(stP, u)
           }
       }
       if len(temp) > 0 {
           groups = append(groups, temp)
       }
   }
   // priority queue of groups
   h := &groupHeap{}
   heap.Init(h)
   for i := range groups {
       if len(groups[i]) > 0 {
           heap.Push(h, i)
       }
   }
   pairsU := make([]int, 0, k)
   pairsV := make([]int, 0, k)
   // pair across groups
   for h.Len() > 1 {
       i := heap.Pop(h).(int)
       j := heap.Pop(h).(int)
       u := groups[i][len(groups[i])-1]
       v := groups[j][len(groups[j])-1]
       groups[i] = groups[i][:len(groups[i])-1]
       groups[j] = groups[j][:len(groups[j])-1]
       pairsU = append(pairsU, u)
       pairsV = append(pairsV, v)
       if len(groups[i]) > 0 {
           heap.Push(h, i)
       }
       if len(groups[j]) > 0 {
           heap.Push(h, j)
       }
   }
   // output
   fmt.Fprintln(out, 1)
   fmt.Fprintln(out, centroid)
   for idx := 0; idx < len(pairsU); idx++ {
       fmt.Fprintf(out, "%d %d %d\n", pairsU[idx], pairsV[idx], centroid)
   }
}
