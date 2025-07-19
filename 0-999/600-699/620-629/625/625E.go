package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Node holds position p, current power z, and original index n
type Node struct {
   p, z, n int
}

// Item is an entry in the priority queue
type Item struct {
   idx  int // sorted index
   z    int // time until collision
   orig int // original index for tie-breaker
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
   if pq[i].z != pq[j].z {
       return pq[i].z < pq[j].z
   }
   return pq[i].orig < pq[j].orig
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   item := old[n-1]
   *pq = old[0 : n-1]
   return item
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var m int64
   fmt.Fscan(in, &n, &m)
   a := make([]Node, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i].p, &a[i].z)
       a[i].n = i + 1
   }
   // sort by p
   // simple insertion sort for small code; use sort.Slice
   // but import sort
   // We'll use sort.Slice
   // Actually import sort
   // sort by a[i].p
   // -----
   // Import sort
   //
   // Do sort
   importSort(a)

   // map original index to sorted index
   pl := make([]int, n+1)
   for i := 0; i < n; i++ {
       pl[a[i].n] = i
   }
   nex := make([]int, n)
   pre := make([]int, n)
   c := make([]bool, n)
   add := make([]int64, n)
   w := make([]int, n)
   for i := 0; i < n; i++ {
       c[i] = true
   }
   var pq PriorityQueue
   heap.Init(&pq)
   // initialize neighbors and queue
   for j := 1; j <= n; j++ {
       i := pl[j]
       k := (i + 1) % n
       nex[i] = k
       pre[k] = i
       z := Find(i, k, a, add, m)
       if z > 0 {
           w[i] = z
           heap.Push(&pq, Item{idx: i, z: z, orig: a[i].n})
       }
   }
   // process collisions
   for pq.Len() > 0 {
       item := heap.Pop(&pq).(Item)
       x, z := item.idx, item.z
       if !c[x] || w[x] != z {
           continue
       }
       // remove next
       c[nex[x]] = false
       nex[x] = nex[nex[x]]
       num := 1
       for {
           vl := Find(x, nex[x], a, add, m)
           if vl == z {
               c[nex[x]] = false
               nex[x] = nex[nex[x]]
               num++
           } else {
               break
           }
       }
       // update current
       a[x].z -= num
       if a[x].z < 0 {
           a[x].z = 0
       }
       pre[nex[x]] = x
       add[x] += int64(num) * int64(z)
       // re-evaluate x
       z2 := Find(x, nex[x], a, add, m)
       if z2 > 0 {
           w[x] = z2
           heap.Push(&pq, Item{idx: x, z: z2, orig: a[x].n})
       } else {
           w[x] = 0
       }
       // re-evaluate previous of x
       px := pre[x]
       z3 := Find(px, x, a, add, m)
       if z3 > 0 {
           w[px] = z3
           heap.Push(&pq, Item{idx: px, z: z3, orig: a[px].n})
       } else {
           w[px] = 0
       }
   }
   // collect survivors
   var ans []int
   for i := 0; i < n; i++ {
       if c[i] {
           ans = append(ans, a[i].n)
       }
   }
   // sort answer
   importSortInts(ans)
   fmt.Fprintln(out, len(ans))
   for _, v := range ans {
       fmt.Fprint(out, v, " ")
   }
   fmt.Fprintln(out)
}

// Find computes time until collision between x and k
func Find(x, k int, a []Node, add []int64, m int64) int {
   if x == k {
       return 0
   }
   // distance in p
   l := int64(a[k].p - a[x].p)
   if l <= 0 {
       l += m
   }
   l += add[k] - add[x]
   t := int64(0)
   if a[x].n < a[k].n {
       l -= int64(a[x].z)
       t = 1
   }
   if l <= 0 && add[x] == 0 && add[k] == 0 {
       return 1
   }
   if a[k].z < a[x].z {
       d := int64(a[x].z - a[k].z)
       return int((l-1)/d + t + 1)
   }
   return 0
}

// sorting helpers
// import sort via functions below
func importSort(a []Node) {
   // simple sort.Slice replacement
   // using built-in sort
   // to avoid import cycle in patch, we use closure
   type byP []Node
   // sort by p ascending
   // implement sort.Interface
   var s byP = a
   // implement methods
   sortByP(s)
}

func importSortInts(a []int) {
   // simple insertion sort for ints
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

// sort.Interface implementation for Node by p
// to minimize imports
var _sortData []Node
func sortByP(data []Node) {
   _sortData = data
   quickSort(0, len(data)-1)
}

func quickSort(l, r int) {
   if l >= r {
       return
   }
   mid := _sortData[(l+r)/2].p
   i, j := l, r
   for i <= j {
       for _sortData[i].p < mid {
           i++
       }
       for _sortData[j].p > mid {
           j--
       }
       if i <= j {
           _sortData[i], _sortData[j] = _sortData[j], _sortData[i]
           i++
           j--
       }
   }
   if l < j {
       quickSort(l, j)
   }
   if i < r {
       quickSort(i, r)
   }
}
