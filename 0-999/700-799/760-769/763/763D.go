package main

import "fmt"

var (
   n    int
   E    [][]int
   Ch   [][]int
   C    []int
   D    []int
   Num  []int
   S    []int
   L    [][]int
   cnt  int
   R    int
   res  int
   Ans  int
)

func dfs1(u, p int) {
   C[u] = 1
   for _, v := range E[u] {
       if v == p {
           continue
       }
       dfs1(v, u)
       C[u] += C[v]
   }
}

func dfs2(u, p int) {
   D[u] = 0
   for _, v := range E[u] {
       if v == p {
           continue
       }
       Ch[u] = append(Ch[u], v)
       dfs2(v, u)
       if D[v]+1 > D[u] {
           D[u] = D[v] + 1
       }
   }
}

func ins(u int, shapeMap map[string]int) {
   tp := make([]int, len(Ch[u]))
   for i, v := range Ch[u] {
       tp[i] = Num[v]
   }
   if len(tp) > 1 {
       sortInts(tp)
   }
   key := fmt.Sprint(tp)
   id, ok := shapeMap[key]
   if !ok {
       cnt++
       shapeMap[key] = cnt
       Num[u] = cnt
   } else {
       Num[u] = id
   }
}

func doDFS(u, depth int) {
   // remove current
   S[Num[u]]--
   if S[Num[u]] == 0 {
       R--
   }
   // update answer
   if res < R+depth {
       res = R + depth
       Ans = u
   }
   for _, v := range Ch[u] {
       doDFS(v, depth+1)
   }
   // restore
   if S[Num[u]] == 0 {
       R++
   }
   S[Num[u]]++
}
// sortInts sorts slice of ints in-place using quicksort with mid pivot
func sortInts(a []int) {
   if len(a) < 2 {
       return
   }
   quickSort(a, 0, len(a)-1)
}

func quickSort(a []int, lo, hi int) {
   if lo < hi {
       p := partition(a, lo, hi)
       quickSort(a, lo, p-1)
       quickSort(a, p+1, hi)
   }
}

func partition(a []int, lo, hi int) int {
   mid := lo + (hi-lo)/2
   a[mid], a[hi] = a[hi], a[mid]
   pivot := a[hi]
   i := lo
   for j := lo; j < hi; j++ {
       if a[j] < pivot {
           a[i], a[j] = a[j], a[i]
           i++
       }
   }
   a[i], a[hi] = a[hi], a[i]
   return i
}

func main() {
   fmt.Scan(&n)
   E = make([][]int, n+1)
   for i := 1; i < n; i++ {
       var a, b int
       fmt.Scan(&a, &b)
       E[a] = append(E[a], b)
       E[b] = append(E[b], a)
   }
   C = make([]int, n+1)
   dfs1(1, 0)
   // find centroid
   mn := n + 1
   x := 1
   for i := 1; i <= n; i++ {
       if C[i]*2 > n && C[i] < mn {
           mn = C[i]
           x = i
       }
   }
   Ch = make([][]int, n+1)
   D = make([]int, n+1)
   dfs2(x, 0)
   // group by depth
   L = make([][]int, n)
   for i := 1; i <= n; i++ {
       d := D[i]
       if d < len(L) {
           L[d] = append(L[d], i)
       }
   }
   Num = make([]int, n+1)
   // assign canonical numbers depth by depth
   for d := 0; d < len(L); d++ {
       shapeMap := make(map[string]int)
       for _, u := range L[d] {
           ins(u, shapeMap)
       }
   }
   // count occurrences
   S = make([]int, cnt+2)
   for i := 1; i <= n; i++ {
       S[Num[i]]++
   }
   R = cnt
   res = 0
   Ans = x
   doDFS(x, 1)
   fmt.Println(Ans)
}
