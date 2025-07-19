package main
import (
   "bufio"
   "fmt"
   "os"
)
func main() {
   reader := bufio.NewReader(os.Stdin)
   var N, K int
   var V int64
   if _, err := fmt.Fscan(reader, &N, &K, &V); err != nil {
       return
   }
   A := make([]int64, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &A[i])
   }
   type op struct{ x int64; from, to int }
   ops := make([]op, 0, 3*N)
   // initial move if V % K == 0
   if V%int64(K) == 0 {
       ops = append(ops, op{A[0], 0, 1})
       A[1] += A[0]
       A[0] = 0
   }
   // build L
   L := make([]int, K)
   for i := 0; i < K; i++ {
       L[i] = -1
   }
   mod0 := int(A[0] % int64(K))
   if mod0 < 0 {
       mod0 += K
   }
   L[mod0] = 0
   aim := int(V % int64(K))
   if aim < 0 {
       aim += K
   }
   for n := 0; n < N && L[aim] == -1; n++ {
       v := int(A[n] % int64(K))
       if v < 0 {
           v += K
       }
       if L[v] == -1 {
           L[v] = n
       }
       for k := 0; k < K; k++ {
           if L[k] != -1 && L[k] != n {
               kk := k + v
               if kk >= K {
                   kk -= K
               }
               if L[kk] == -1 {
                   L[kk] = n
               }
           }
       }
   }
   if L[aim] == -1 {
       fmt.Println("NO")
       return
   }
   target := L[aim]
   V -= A[target]
   A[target] = 0
   // adjust to make V divisible by K
   for V%int64(K) != 0 {
       v := V % int64(K)
       if v < 0 {
           v += int64(K)
       }
       idx := L[int(v)]
       // ensure A[idx] != 0
       ops = append(ops, op{A[idx], idx, target})
       V -= A[idx]
       A[idx] = 0
   }
   other := 0
   if target == 0 {
       other = 1
   }
   if V < 0 {
       ops = append(ops, op{(-V) / int64(K), target, other})
   } else if V > 0 {
       for i := 0; i < N; i++ {
           if A[i] > 0 && i != other && i != target {
               ops = append(ops, op{A[i], i, other})
               A[other] += A[i]
               A[i] = 0
           }
       }
       if A[other] < V {
           fmt.Println("NO")
           return
       }
       ops = append(ops, op{V / int64(K), other, target})
   }
   // output
   fmt.Println("YES")
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for _, oc := range ops {
       x := oc.x
       if x < 1 {
           x = 1
       }
       fmt.Fprintf(writer, "%d %d %d\n", x, oc.from+1, oc.to+1)
   }
}
