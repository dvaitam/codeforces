package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   g           [][]int
   down1, down2 []int
   sonDown1, sonDown2 []int
   upArr, sonUp []int
   branchArr, lenArr []int
   sonBranch1, sonBranch2 []int
   dad []int
   branch0, length0, a1, a2, b1, b2 int
)

func Down(node, father int) {
   dad[node] = father
   for _, son := range g[node] {
       if son != father {
           Down(son, node)
           downNow := down1[son] + 1
           sonNow := sonDown1[son]
           if sonNow == 0 {
               sonNow = son
           }
           if downNow > down1[node] {
               down2[node] = down1[node]
               sonDown2[node] = sonDown1[node]
               down1[node] = downNow
               sonDown1[node] = sonNow
           } else if downNow > down2[node] {
               down2[node] = downNow
               sonDown2[node] = sonNow
           }
       }
   }
}

func Up(node, father int) {
   best1, which1, son1 := 0, 0, 0
   best2, which2, son2 := 0, 0, 0
   for _, son := range g[node] {
       if son != father {
           bestNow := down1[son] + 2
           sonNow := sonDown1[son]
           if sonNow == 0 {
               sonNow = son
           }
           if bestNow > best1 {
               best2, son2, which2 = best1, son1, which1
               best1, son1, which1 = bestNow, sonNow, son
           } else {
               best2, son2, which2 = bestNow, sonNow, son
           }
       }
   }
   for _, son := range g[node] {
       if son != father {
           upArr[son] = upArr[node] + 1
           sonUp[son] = sonUp[node]
           bestNow := best1
           sonNow := son1
           if which1 == son {
               bestNow = best2
               sonNow = son2
           }
           if bestNow > upArr[son] {
               upArr[son] = bestNow
               sonUp[son] = sonNow
           }
           Up(son, node)
       }
   }
}

func DFS(node, father int) {
   for _, son := range g[node] {
       if son != father {
           DFS(son, node)
           if branchArr[son] != 0 && (branchArr[son]+1 > branchArr[node] || (branchArr[son]+1 == branchArr[node] && lenArr[son] > lenArr[node])) {
               branchArr[node] = branchArr[son] + 1
               lenArr[node] = lenArr[son]
               sonBranch1[node] = sonBranch1[son]
               sonBranch2[node] = sonBranch2[son]
           }
       }
   }
   if down1[node] != 0 && down2[node] != 0 && (1 > branchArr[node] || (1 == branchArr[node] && down1[node]+down2[node] > lenArr[node])) {
       branchArr[node] = 1
       lenArr[node] = down1[node] + down2[node]
       sonBranch1[node] = sonDown1[node]
       sonBranch2[node] = sonDown2[node]
   }
}

func Solve(n int) {
   for node := 1; node <= n; node++ {
       branch1, branch2, length1, length2 := 0, 0, 0, 0
       son11, son12, son21, son22 := 0, 0, 0, 0
       for _, son := range g[node] {
           if son != dad[node] && branchArr[son] != 0 {
               if branchArr[son] > branch1 || (branchArr[son] == branch1 && lenArr[son] > length1) {
                   branch2, length2, son21, son22 = branch1, length1, son11, son12
                   branch1, length1, son11, son12 = branchArr[son], lenArr[son], sonBranch1[son], sonBranch2[son]
               } else if branchArr[son] > branch2 || (branchArr[son] == branch2 && lenArr[son] > length2) {
                   branch2, length2, son21, son22 = branchArr[son], lenArr[son], sonBranch1[son], sonBranch2[son]
               }
           }
       }
       if branch1 != 0 && branch2 != 0 {
           if branch1+branch2+1 > branch0 || (branch1+branch2+1 == branch0 && length1+length2 > length0) {
               branch0 = branch1 + branch2 + 1
               length0 = length1 + length2
               a1, a2, b1, b2 = son11, son12, son21, son22
           }
           continue
       }
       best1, best2, son1, son2 := 0, 0, 0, 0
       for _, son := range g[node] {
           if son != dad[node] && branchArr[son] == 0 {
               downNow := down1[son] + 1
               sonNow := sonDown1[son]
               if sonNow == 0 {
                   sonNow = son
               }
               if downNow > best1 {
                   best2, son2 = best1, son1
                   best1, son1 = downNow, sonNow
               } else if downNow > best2 {
                   best2, son2 = downNow, sonNow
               }
           }
       }
       if upArr[node] != 0 {
           if upArr[node] > best1 {
               best2, son2 = best1, son1
               best1, son1 = upArr[node], sonUp[node]
           } else if upArr[node] > best2 {
               best2, son2 = upArr[node], sonUp[node]
           }
       }
       if best1 != 0 && best2 != 0 {
           if branch1+1 > branch0 || (branch1+1 == branch0 && length1+best1+best2 > length0) {
               branch0 = branch1 + 1
               length0 = length1 + best1 + best2
               a1, a2, b1, b2 = son11, son12, son1, son2
           }
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   g = make([][]int, n+1)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       g[u] = append(g[u], v)
       g[v] = append(g[v], u)
   }
   down1 = make([]int, n+1)
   down2 = make([]int, n+1)
   sonDown1 = make([]int, n+1)
   sonDown2 = make([]int, n+1)
   upArr = make([]int, n+1)
   sonUp = make([]int, n+1)
   branchArr = make([]int, n+1)
   lenArr = make([]int, n+1)
   sonBranch1 = make([]int, n+1)
   sonBranch2 = make([]int, n+1)
   dad = make([]int, n+1)
   branch0, length0, a1, a2, b1, b2 = 0, 0, 0, 0, 0, 0

   Down(1, 0)
   Up(1, 0)
   DFS(1, 0)
   Solve(n)
   fmt.Printf("%d %d\n%d %d\n", a1, b1, a2, b2)
}
