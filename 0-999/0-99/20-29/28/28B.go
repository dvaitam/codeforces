package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   d := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &d[i])
   }
   // Build undirected graph of allowed swaps
   adj := make([][]int, n)
   for i := 0; i < n; i++ {
       di := d[i]
       for _, j := range []int{i - di, i + di} {
           if j >= 0 && j < n {
               adj[i] = append(adj[i], j)
               adj[j] = append(adj[j], i)
           }
       }
   }
   visited := make([]bool, n)
   // Check each connected component
   for i := 0; i < n; i++ {
       if visited[i] {
           continue
       }
       // BFS
       queue := []int{i}
       visited[i] = true
       compIdx := []int{i}
       for q := 0; q < len(queue); q++ {
           u := queue[q]
           for _, v := range adj[u] {
               if !visited[v] {
                   visited[v] = true
                   queue = append(queue, v)
                   compIdx = append(compIdx, v)
               }
           }
       }
       // Compare initial values and target values in component
       // initial values are indices+1
       initVals := make([]int, len(compIdx))
       target := make([]int, len(compIdx))
       for k, idx := range compIdx {
           initVals[k] = idx + 1
           target[k] = p[idx]
       }
       // sort both slices
       sortInts(initVals)
       sortInts(target)
       for k := range initVals {
           if initVals[k] != target[k] {
               fmt.Println("NO")
               return
           }
       }
   }
   fmt.Println("YES")
}

// sortInts sorts a slice of ints in ascending order (in-place)
func sortInts(a []int) {
   // simple insertion sort for small n
   for i := 1; i < len(a); i++ {
       key := a[i]
       j := i - 1
       for j >= 0 && a[j] > key {
           a[j+1] = a[j]
           j--
       }
       a[j+1] = key
   }
}
