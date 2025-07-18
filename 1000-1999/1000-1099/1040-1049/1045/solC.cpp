# include <bits/stdc++.h>
# define IL inline
# define RG register
# define Fill(a, b) memset(a, b, sizeof(a))
using namespace std;
typedef long long ll;

IL int Input(){
    RG int x = 0, z = 1; RG char c = getchar();
    for(; c < '0' || c > '9'; c = getchar()) z = c == '-' ? -1 : 1;
    for(; c >= '0' && c <= '9'; c = getchar()) x = (x << 1) + (x << 3) + (c ^ 48);
    return x * z;
}

const int maxn(2e5 + 5);
const int inf(1e9);

int n, m, dfn[maxn], low[maxn], idx, cnt[maxn], fa[maxn];
int sta[maxn], top, tot, vis[maxn], deep[maxn], val[maxn];
int Top[maxn], son[maxn], size[maxn];

struct Edge{
    int first[maxn], cnt, nxt[maxn * 10], to[maxn * 10];

    IL void Init(){
        cnt = 0, Fill(first, -1);
    }

    IL void Add(RG int u, RG int v){
        nxt[cnt] = first[u], to[cnt] = v, first[u] = cnt++;
    }
} e1, e2;

void Tarjan(RG int u){
    dfn[u] = low[u] = ++idx, sta[++top] = u;
    for(RG int e = e1.first[u]; ~e; e = e1.nxt[e]){
        RG int v = e1.to[e];
        if(!dfn[v]){
            Tarjan(v), low[u] = min(low[u], low[v]);
            if(low[v] >= dfn[u]){
                RG int x = sta[top]; ++tot;
                do{
                    x = sta[top--], val[tot] = 1;
                    e2.Add(tot, x), e2.Add(x, tot);
                } while(x != v);
                e2.Add(tot, u), e2.Add(u, tot);
			}
        }
        else low[u] = min(low[u], dfn[v]);
    }
}

void Dfs1(RG int u, RG int ff) {
	size[u] = 1;
    for(RG int e = e2.first[u]; ~e; e = e2.nxt[e]){
        RG int v = e2.to[e];
		if (v == ff) continue;
		deep[v] = deep[u] + 1, cnt[v] = cnt[u] + val[v];
		fa[v] = u, Dfs1(v, u);
		size[u] += size[v];
		if (size[v] > size[son[u]]) son[u] = v;
	}
}

void Dfs2(RG int u, RG int tp, int ff) {
	Top[u] = tp;
	if (son[u]) Dfs2(son[u], tp, u);
    for(RG int e = e2.first[u]; ~e; e = e2.nxt[e]){
        RG int v = e2.to[e];
		if (v == ff || v == son[u]) continue;
		Dfs2(v, v, u);
	}
}

int LCA(int u, int v) {
	while (Top[u] ^ Top[v]) {
		if (deep[Top[u]] > deep[Top[v]]) swap(u, v);
		v = fa[Top[v]];
	}
	return deep[u] > deep[v] ? v : u;
}

int Query(int u, int v) {
	int ans = deep[u] + deep[v] - cnt[u] - cnt[v];
	u = LCA(u, v);
	return ans - 2 * deep[u] + cnt[u] + cnt[fa[u]];
}

int main(){
    e1.Init(), e2.Init();
    tot = n = Input(), m = Input();
    RG int t = Input();
    for(RG int i = 1; i <= m; ++i){
        RG int u = Input(), v = Input();
        e1.Add(u, v), e1.Add(v, u);
    }
	for (int i = 1; i <= n; ++i)
		if (!dfn[i]) Tarjan(i), cnt[i] = val[i], Dfs1(i, 0), Dfs2(i, i, 0);
    for(RG int i = 1; i <= t; ++i){
        RG int u = Input(), v = Input();
        printf("%d\n", Query(u, v));
	}
    return 0;
}