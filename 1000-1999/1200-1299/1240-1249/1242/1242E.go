package main

import (
   "bufio"
   "container/list"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   t := make([]int, n)
   pocz := make([]int, n)
   total := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &t[i])
       pocz[i] = total
       total += t[i]
   }
   // parent for DSU
   f := make([]int, total)
   for i := 0; i < total; i++ {
       f[i] = i
   }
   // number of components
   comp := total
   // result labels per component root
   res := make([]int, total)
   // order groups by descending t
   p := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
   }
   sort.Slice(p, func(i, j int) bool { return t[p[i]] > t[p[j]] })
   // helper DSU find
   var find func(int) int
   find = func(a int) int {
       if f[a] != a {
           f[a] = find(f[a])
       }
       return f[a]
   }
   // union a->b
   uni := func(a, b int) {
       a = find(a)
       b = find(b)
       if a != b {
           f[a] = b
           comp--
       }
   }
   // function daj(x) returns slice of indices for group x
   daj := func(x int) []int {
       arr := make([]int, t[x])
       for i := 0; i < t[x]; i++ {
           arr[i] = pocz[x] + i
       }
       return arr
   }
   // initialize deque with largest group
   ak := list.New()
   for _, v := range daj(p[0]) {
       ak.PushBack(v)
   }
   // process other groups
   for idx := 1; idx < n; idx++ {
       i := p[idx]
       nextSize := 3
       if idx < n-1 {
           nextSize = t[p[idx+1]]
       }
       pom := daj(i)
       // save last of ak
       back := ak.Back().Value.(int)
       // first two unions
       uni(pom[len(pom)-1], ak.Back().Value.(int))
       ak.Remove(ak.Back())
       pom = pom[:len(pom)-1]
       uni(pom[len(pom)-1], ak.Back().Value.(int))
       ak.Remove(ak.Back())
       // additional merges while possible
       for len(pom) > 1 && len(pom)+ak.Len()-2 >= nextSize {
           pom = pom[:len(pom)-1]
           ak.Remove(ak.Back())
           uni(pom[len(pom)-1], ak.Back().Value.(int))
       }
       // pop one more from pom
       pom = pom[:len(pom)-1]
       // push remaining pom reversed onto back of ak
       for j := len(pom) - 1; j >= 0; j-- {
           ak.PushBack(pom[j])
       }
       // push saved back to front
       ak.PushFront(back)
   }
   // output number of components
   fmt.Fprintln(writer, comp)
   // label and print assignments
   ver := 1
   for i := 0; i < n; i++ {
       for j := 0; j < t[i]; j++ {
           idx := pocz[i] + j
           root := find(idx)
           if res[root] == 0 {
               res[root] = ver
               ver++
           }
           fmt.Fprint(writer, res[root], " ")
       }
       fmt.Fprintln(writer)
   }
}
