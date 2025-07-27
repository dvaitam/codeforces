package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// Treap for int keys
type Node struct {
   key      int
   priority int
   left     *Node
   right    *Node
}

// rotateRight rotates right
func rotateRight(y *Node) *Node {
   x := y.left
   y.left = x.right
   x.right = y
   return x
}

// rotateLeft rotates left
func rotateLeft(x *Node) *Node {
   y := x.right
   x.right = y.left
   y.left = x
   return y
}

// insert key into treap
func insert(root *Node, key int) *Node {
   if root == nil {
       return &Node{key: key, priority: rand.Int()}
   }
   if key < root.key {
       root.left = insert(root.left, key)
       if root.left.priority > root.priority {
           root = rotateRight(root)
       }
   } else {
       root.right = insert(root.right, key)
       if root.right.priority > root.priority {
           root = rotateLeft(root)
       }
   }
   return root
}

// delete key from treap
func deleteNode(root *Node, key int) *Node {
   if root == nil {
       return nil
   }
   if key < root.key {
       root.left = deleteNode(root.left, key)
   } else if key > root.key {
       root.right = deleteNode(root.right, key)
   } else {
       // root.key == key
       if root.left == nil {
           return root.right
       }
       if root.right == nil {
           return root.left
       }
       // both children exist: rotate with higher priority
       if root.left.priority > root.right.priority {
           root = rotateRight(root)
           root.right = deleteNode(root.right, key)
       } else {
           root = rotateLeft(root)
           root.left = deleteNode(root.left, key)
       }
   }
   return root
}

// find predecessor (max key < x)
func findPrev(root *Node, x int) *Node {
   var res *Node
   cur := root
   for cur != nil {
       if cur.key < x {
           res = cur
           cur = cur.right
       } else {
           cur = cur.left
       }
   }
   return res
}

// find successor (min key > x)
func findNext(root *Node, x int) *Node {
   var res *Node
   cur := root
   for cur != nil {
       if cur.key > x {
           res = cur
           cur = cur.left
       } else {
           cur = cur.right
       }
   }
   return res
}

// find minimum key
func findMin(root *Node) int {
   cur := root
   for cur.left != nil {
       cur = cur.left
   }
   return cur.key
}

// find maximum key
func findMax(root *Node) int {
   cur := root
   for cur.right != nil {
       cur = cur.right
   }
   return cur.key
}

// MaxHeap for ints
type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   rand.Seed(time.Now().UnixNano())
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   fmt.Fscan(reader, &n, &q)
   var x int
   var root *Node
   // gap counts and heap
   gapCount := make(map[int]int)
   h := &MaxHeap{}
   heap.Init(h)
   coords := make([]int, 0, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x)
       coords = append(coords, x)
       root = insert(root, x)
   }
   if n > 1 {
       // build gaps
       // sort coords
       // simple quick sort using Go
       sortInts(coords)
       for i := 1; i < len(coords); i++ {
           d := coords[i] - coords[i-1]
           gapCount[d]++
           heap.Push(h, d)
       }
   }
   sz := n
   // function to get answer
   getAns := func() int {
       if sz <= 1 {
           return 0
       }
       mn := findMin(root)
       mx := findMax(root)
       span := mx - mn
       // get max gap
       var maxGap int
       for h.Len() > 0 {
           top := (*h)[0]
           if gapCount[top] > 0 {
               maxGap = top
               break
           }
           heap.Pop(h)
       }
       return span - maxGap
   }
   // print initial
   fmt.Fprintln(writer, getAns())
   // process queries
   for i := 0; i < q; i++ {
       var t int
       fmt.Fscan(reader, &t, &x)
       if t == 1 {
           // add x
           // find neighbors
           prev := findPrev(root, x)
           next := findNext(root, x)
           if prev != nil && next != nil {
               d := next.key - prev.key
               gapCount[d]--
           }
           if prev != nil {
               d := x - prev.key
               gapCount[d]++
               heap.Push(h, d)
           }
           if next != nil {
               d := next.key - x
               gapCount[d]++
               heap.Push(h, d)
           }
           root = insert(root, x)
           sz++
       } else {
           // remove x
           prev := findPrev(root, x)
           next := findNext(root, x)
           if prev != nil {
               d := x - prev.key
               gapCount[d]--
           }
           if next != nil {
               d := next.key - x
               gapCount[d]--
           }
           if prev != nil && next != nil {
               d := next.key - prev.key
               gapCount[d]++
               heap.Push(h, d)
           }
           root = deleteNode(root, x)
           sz--
       }
       fmt.Fprintln(writer, getAns())
   }
}

// simple sort for ints (QuickSort)
func sortInts(a []int) {
   if len(a) < 2 {
       return
   }
   quickSort(a, 0, len(a)-1)
}

func quickSort(a []int, lo, hi int) {
   if lo >= hi {
       return
   }
   p := a[(lo+hi)/2]
   i, j := lo, hi
   for i <= j {
       for a[i] < p {
           i++
       }
       for a[j] > p {
           j--
       }
       if i <= j {
           a[i], a[j] = a[j], a[i]
           i++
           j--
       }
   }
   if lo < j {
       quickSort(a, lo, j)
   }
   if i < hi {
       quickSort(a, i, hi)
   }
}
