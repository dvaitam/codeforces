#include<bits/stdc++.h>
using namespace std;

const int MAXN = 200010;

struct Edge{
    int to, next;
}edges[MAXN << 1];

int head[MAXN], tot;

void addEdge(int u, int v){
    edges[++ tot].to = v;
    edges[tot].next = head[u];
    head[u] = tot;
}

int n, a[MAXN], d[MAXN];
long long val[MAXN], dp[MAXN];

void dfs(int x, int fa){
    val[x] = a[x];
    for(int i = head[x]; i; i = edges[i].next){
        int y = edges[i].to;
        if(y == fa) continue;
        d[y] = d[x] + 1;
        dfs(y, x);
        val[x] += val[y];
    }
}

void init(){
    for(int i = 2; i <= n; i ++)
        dp[1] += 1LL * d[i] * a[i];
}

void tree_dp(int x, int fa){
    for(int i = head[x]; i; i = edges[i].next){
        int y = edges[i].to;
        if(y == fa) continue;
        dp[y] = dp[x] + val[1] - 2LL * val[y];
        tree_dp(y, x);
    }
}

int main(){
    cin >> n;
    for(int i = 1; i <= n; i ++)
        scanf("%d", a + i);
    for(int i = 2; i <= n; i ++){
        int u, v; scanf("%d%d", &u, &v);
        addEdge(u, v);
        addEdge(v, u);
    }
    dfs(1, 0);
    init();
    tree_dp(1, 0);
    cout << *max_element(dp + 1, dp + n + 1);
    return 0;
}