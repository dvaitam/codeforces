/**
  * @file pf.cpp
  * @brief codeforces

  *

  * @author yao
  * @date 2018-10-11
  */
#include <cstdlib>
#include <cstdio>
#include <cctype>
#include <cstring>
#include <utility>
#include <algorithm>
#include <functional>
#define N 1048576
#define ft first
#define sd second
#define INF ((int)1e9)
#ifdef DBG
#   define dbg_pri(x...) fprintf(stderr,x)
#else
#   define dbg_pri(x...) 0
#endif //DBG

typedef unsigned int uint;
typedef long long int lli;
typedef unsigned long long int ulli;

int sib[N], chi[N];
int a[N], b[N];
int h[N];
int n,k;

void dfs(int u)
{
    h[u] = INF;
    for(int v=chi[u]; v; v = sib[v])
    {
        dfs(v);
        h[u] = std::min(h[u],h[v]+1);
        a[u] += a[v];
        b[u] = std::max(b[u], b[v] - a[v]);
    }
    b[u] += a[u];
    if(h[u] == INF) h[u] = 0, b[u] = a[u] = 1;
    if(h[u] >= k) a[u] = 0;
    dbg_pri("poi 1: %d %d %d %d\n",u,a[u],b[u],h[u]);
}

int main()
{
    scanf("%d%d",&n,&k);
    for(int i=2;i<=n;++i)
    {
        int p;
        scanf("%d",&p);
        sib[i] = chi[p], chi[p] = i;
    }
    dfs(1);
    printf("%d\n",b[1]);
    return 0;
}