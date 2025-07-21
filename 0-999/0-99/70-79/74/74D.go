package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math/rand"
   "os"
)

// Segment represents a free interval [start,end]
type segment struct { start, end int }

// segHeap is a max-heap of segments by length, then by start (rightmost)
type segHeap []segment
func (h segHeap) Len() int { return len(h) }
func (h segHeap) Less(i, j int) bool {
   li := h[i].end - h[i].start + 1
   lj := h[j].end - h[j].start + 1
   if li != lj {
       return li > lj
   }
   return h[i].start > h[j].start
}
func (h segHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *segHeap) Push(x interface{}) { *h = append(*h, x.(segment)) }
func (h *segHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

// Treap for order statistics
type Node struct {
   key, pri, size int
   left, right    *Node
}
func (n *Node) upd() {
   n.size = 1
   if n.left != nil {
       n.size += n.left.size
   }
   if n.right != nil {
       n.size += n.right.size
   }
}
func rotateRight(t *Node) *Node {
   l := t.left
   t.left = l.right
   l.right = t
   t.upd()
   l.upd()
   return l
}
func rotateLeft(t *Node) *Node {
   r := t.right
   t.right = r.left
   r.left = t
   t.upd()
   r.upd()
   return r
}
func insert(t *Node, key int) *Node {
   if t == nil {
       return &Node{key: key, pri: rand.Int(), size: 1}
   }
   if key < t.key {
       t.left = insert(t.left, key)
       if t.left.pri > t.pri {
           t = rotateRight(t)
       }
   } else {
       t.right = insert(t.right, key)
       if t.right.pri > t.pri {
           t = rotateLeft(t)
       }
   }
   t.upd()
   return t
}
func deleteKey(t *Node, key int) *Node {
   if t == nil {
       return nil
   }
   if key < t.key {
       t.left = deleteKey(t.left, key)
   } else if key > t.key {
       t.right = deleteKey(t.right, key)
   } else {
       if t.left == nil && t.right == nil {
           return nil
       }
       if t.left == nil {
           t = rotateLeft(t)
           t.left = deleteKey(t.left, key)
       } else if t.right == nil {
           t = rotateRight(t)
           t.right = deleteKey(t.right, key)
       } else if t.left.pri > t.right.pri {
           t = rotateRight(t)
           t.right = deleteKey(t.right, key)
       } else {
           t = rotateLeft(t)
           t.left = deleteKey(t.left, key)
       }
   }
   t.upd()
   return t
}
// order returns count of keys < key
func order(t *Node, key int) int {
   if t == nil {
       return 0
   }
   if key <= t.key {
       return order(t.left, key)
   }
   leftSize := 0
   if t.left != nil {
       leftSize = t.left.size
   }
   return leftSize + 1 + order(t.right, key)
}

func main() {
   rand.Seed(1)
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, q int
   fmt.Fscan(in, &n, &q)
   // free segments maps
   s2e := make(map[int]int)
   e2s := make(map[int]int)
   h := &segHeap{}
   heap.Init(h)
   // initial segment
   s2e[1] = n
   e2s[n] = 1
   heap.Push(h, segment{1, n})
   // employee positions
   id2pos := make(map[int]int)
   var root *Node
   for i := 0; i < q; i++ {
       var t int
       fmt.Fscan(in, &t)
       if t == 0 {
           var l, r int
           fmt.Fscan(in, &l, &r)
           cnt := order(root, r+1) - order(root, l)
           fmt.Fprintln(out, cnt)
       } else {
           id := t
           if pos, ok := id2pos[id]; !ok {
               // arrival
               var seg segment
               // get valid largest segment
               for {
                   if h.Len() == 0 {
                       break
                   }
                   seg = heap.Pop(h).(segment)
                   if end, ex := s2e[seg.start]; ex && end == seg.end {
                       break
                   }
               }
               start, end := seg.start, seg.end
               delete(s2e, start)
               delete(e2s, end)
               // choose position
               pos = (start + end + 1) / 2
               id2pos[id] = pos
               root = insert(root, pos)
               // left split
               if start <= pos-1 {
                   s2e[start] = pos - 1
                   e2s[pos-1] = start
                   heap.Push(h, segment{start, pos - 1})
               }
               // right split
               if pos+1 <= end {
                   s2e[pos+1] = end
                   e2s[end] = pos + 1
                   heap.Push(h, segment{pos + 1, end})
               }
           } else {
               // departure
               root = deleteKey(root, pos)
               delete(id2pos, id)
               leftStart, leftEnd := 0, 0
               rightStart, rightEnd := 0, 0
               // check left adjacent
               if lstart, ex := e2s[pos-1]; ex {
                   leftStart, leftEnd = lstart, pos-1
                   delete(s2e, leftStart)
                   delete(e2s, leftEnd)
               }
               // check right adjacent
               if rend, ex := s2e[pos+1]; ex {
                   rightStart, rightEnd = pos+1, rend
                   delete(s2e, rightStart)
                   delete(e2s, rightEnd)
               }
               // new merged segment
               ns := pos
               ne := pos
               if leftEnd > 0 {
                   ns = leftStart
               }
               if rightEnd > 0 {
                   ne = rightEnd
               }
               s2e[ns] = ne
               e2s[ne] = ns
               heap.Push(h, segment{ns, ne})
           }
       }
   }
}
