#include <iostream>
#include <cstdio>
#include <cstring>
#include <algorithm>
#include <queue>
#include <vector>
#include <map>

using namespace std;

#define mp make_pair
const int MAXN = int(2e5);
typedef pair<int, int> pii;
bool vis[MAXN];

int N, G[MAXN], adj[MAXN], ind[MAXN], q[MAXN], fr, re, edge_cnt;
struct E{int nxt, des, flg; E(){} E(int a, int b, int c): nxt(a), des(b), flg(c) {} } e[int(2e6)];
void adde(int a, int b, int c) { e[edge_cnt] = E(G[a], b, c); G[a] = edge_cnt++; }
void adde2(int a, int b, int c) { e[edge_cnt] = E(adj[a], b, c); adj[a] = edge_cnt++; }
void dfs_tree(int x, int f, int t = 0){
    if (~f) adde2(x, f, t), ++ind[f];
    for (int i = G[x]; ~i; i = e[i].nxt){
        if (e[i].des != f) dfs_tree(e[i].des, x, e[i].flg);
    }
}
void dfs(int x){
    if (vis[x]) return;
    for (int i = adj[x]; ~i; i = e[i].nxt){
        vis[x] = true;
        dfs(e[i].des);
    }
}
void print(int v){
    int qc = 0;
    while (v) q[qc++] = v%10, v /= 10;
    while ((--qc) >= 0) putchar(q[qc] + '0');
}
void scanf_(int &v){
    int ch; while ((ch = getchar()) > '9' || ch < '0');
    v = ch-'0'; while ((ch = getchar()) >= '0' && ch <= '9') v = v*10 +ch-'0';
}
vector<int> ans;
int main(){
    //freopen("h:\\data.in", "r", stdin);
    //freopen("h:\\data.out", "w", stdout);
    while (~scanf("%d", &N)){
        memset(G, 0xff, sizeof(G));
        memset(ind, 0x00, sizeof(ind));
        memset(adj, 0xff, sizeof(adj));
        memset(vis, 0x00, sizeof(vis));
        edge_cnt = 0;
        
        
        for (int i = 1; i < N; ++i){
            int a, b, c; scanf_(a), scanf_(b), scanf_(c);
            adde(a, b, c); adde(b, a, c);
        }
        dfs_tree(1, -1);
        
        fr = re = 0;
        for (int i = 1; i <= N; ++i) if (ind[i] == 0) q[re++] = i;

        ans.clear();
        while (fr < re){
            int u = q[fr++];
            for (int i = adj[u]; ~i; i = e[i].nxt){
                int v = e[i].des;
                if ((--ind[v]) == 0) q[re++] = v;
                if (e[i].flg == 1 || vis[u]) continue;
                ans.push_back(u); dfs(u);
            }
        }
        
        printf("%d\n", ans.size());
        for (int i = 0, j = ans.size(); i < j; ++i) { print(ans[i]); putchar(' '); }
        puts("");
        
    }
    return 0;
}