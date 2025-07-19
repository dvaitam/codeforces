package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct { x, y int }

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

// calc returns sum of (a[i]+1) for i from 0 to p
func calc(a []int, p int) int {
   sum := 0
   for i := 0; i <= p; i++ {
       sum += a[i] + 1
   }
   return sum
}

func trans0(a *[]int, k int, ord *[]pair) {
   arr := *a
   if arr[k] != 0 {
       L := calc(arr, 0) - 1
       R := calc(arr, k) - 2
       *ord = append(*ord, pair{L, R})
       arr[0] += arr[k]
       arr[k] = 0
       for p := 1; p < k-p; p++ {
           arr[p], arr[k-p] = arr[k-p], arr[p]
       }
   }
   for p := 1; p < k-p; p++ {
       if arr[p] > arr[k-p] {
           break
       }
       if arr[p] < arr[k-p] {
           L := calc(arr, 0) - 1
           R := calc(arr, k) - 2
           *ord = append(*ord, pair{L, R})
           for j := 1; j < k-j; j++ {
               arr[j], arr[k-j] = arr[k-j], arr[j]
           }
           break
       }
   }
   *a = arr
}

func trans1(a *[]int, k int, ord *[]pair) {
   arr := *a
   for {
       F := false
       for i := 1; i <= k; i++ {
           if arr[i] != 0 {
               F = true
               break
           }
       }
       if !F {
           break
       }
       L := calc(arr, 0) - 1
       R := calc(arr, k) - 2
       *ord = append(*ord, pair{L, R})
       arr[0] += arr[k]
       arr[k] = 0
       for i := 1; i < k-i; i++ {
           arr[i], arr[k-i] = arr[k-i], arr[i]
       }
       L2 := calc(arr, 1) - 1
       R2 := calc(arr, k) - 1
       *ord = append(*ord, pair{L2, R2})
       for i := 1; i < k-i; i++ {
           arr[i+1], arr[k-i+1] = arr[k-i+1], arr[i+1]
       }
   }
   *a = arr
}

func trans2(a *[]int, k int, ord *[]pair) {
   arr := *a
   // even indices
   for {
       ok := false
       for i := 2; i <= k; i += 2 {
           if arr[i] != 0 {
               ok = true
               break
           }
       }
       if !ok {
           break
       }
       L := calc(arr, 0) - 1
       R := calc(arr, k) - 2
       *ord = append(*ord, pair{L, R})
       arr[0] += arr[k]
       arr[k] = 0
       for i := 1; i < k-i; i++ {
           arr[i], arr[k-i] = arr[k-i], arr[i]
       }
       L2 := calc(arr, 1) - 1
       R2 := calc(arr, k) - 1
       *ord = append(*ord, pair{L2, R2})
       for i := 1; i < k-i; i++ {
           arr[i+1], arr[k-i+1] = arr[k-i+1], arr[i+1]
       }
   }
   // odd indices
   for {
       ok := false
       for i := 1; i <= k; i += 2 {
           if arr[i] != 0 {
               ok = true
               break
           }
       }
       if !ok {
           break
       }
       L := calc(arr, 0)
       R := calc(arr, k) - 1
       *ord = append(*ord, pair{L, R})
       // ensure arr has index k+1
       if len(arr) > k+1 {
           arr[k+1] += arr[1]
       } else {
           arr = append(arr, arr[1])
       }
       arr[1] = 0
       for i := 1; i < k-i; i++ {
           arr[i+1], arr[k-i+1] = arr[k-i+1], arr[i+1]
       }
       L2 := calc(arr, 0) - 1
       R2 := calc(arr, k) - 2
       *ord = append(*ord, pair{L2, R2})
       for i := 1; i < k-i; i++ {
           arr[i], arr[k-i] = arr[k-i], arr[i]
       }
   }
   *a = arr
}

func trans(a *[]int, k int, ord *[]pair) {
   arr := *a
   if len(arr) == k+1 {
       trans0(&arr, k, ord)
       *a = arr
       return
   }
   for len(arr) > k+2 {
       if arr[len(arr)-1] == 0 {
           arr = arr[:len(arr)-1]
           continue
       }
       p := len(arr) - 1
       L := calc(arr, p-k) - 1
       R := calc(arr, p) - 2
       *ord = append(*ord, pair{L, R})
       arr[p-k] += arr[p]
       arr = arr[:p]
       for i := 1; i < k-i; i++ {
           idx1 := p-k + i
           idx2 := p - i
           arr[idx1], arr[idx2] = arr[idx2], arr[idx1]
       }
   }
   if arr[len(arr)-1] != 0 {
       p := len(arr) - 1
       L := calc(arr, p-k) - 1
       R := calc(arr, p) - 2
       *ord = append(*ord, pair{L, R})
       arr[p-k] += arr[p]
       arr[p] = 0
       for i := 1; i < k-i; i++ {
           idx1 := p-k + i
           idx2 := p - i
           arr[idx1], arr[idx2] = arr[idx2], arr[idx1]
       }
   }
   if k%2 == 1 {
       trans1(&arr, k, ord)
   } else {
       trans2(&arr, k, ord)
   }
   *a = arr
}

func solve() {
   var n, k int
   fmt.Fscan(reader, &n, &k)
   var s, t string
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &t)
   a := []int{0}
   b := []int{0}
   for _, ch := range s {
       if ch == '1' {
           a = append(a, 0)
       } else {
           a[len(a)-1]++
       }
   }
   for _, ch := range t {
       if ch == '1' {
           b = append(b, 0)
       } else {
           b[len(b)-1]++
       }
   }
   if len(a) != len(b) {
       fmt.Fprintln(writer, -1)
       return
   }
   same := true
   for i := range a {
       if a[i] != b[i] {
           same = false
           break
       }
   }
   if same {
       fmt.Fprintln(writer, 0)
       return
   }
   if k == 0 || len(a) <= k {
       fmt.Fprintln(writer, -1)
       return
   }
   ord := make([]pair, 0)
   trans(&b, k, &ord)
   ord2 := make([]pair, len(ord))
   copy(ord2, ord)
   ord = ord[:0]
   trans(&a, k, &ord)
   // check equal
   for i := range a {
       if a[i] != b[i] {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   // append reversed ord2
   for i := len(ord2) - 1; i >= 0; i-- {
       ord = append(ord, ord2[i])
   }
   fmt.Fprintln(writer, len(ord))
   for _, p := range ord {
       fmt.Fprintln(writer, p.x+1, p.y+1)
   }
}

func main() {
   defer writer.Flush()
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       solve()
       T--
   }
}
