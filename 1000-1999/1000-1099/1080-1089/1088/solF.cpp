#include <bits/stdc++.h>
using namespace std;
using LL = long long;
#define FOR(i, x, y) for (decay<decltype(y)>::type i = (x), _##i = (y); i < _##i; ++i)
#define FORD(i, x, y) for (decay<decltype(x)>::type i = (x), _##i = (y); i > _##i; --i)
#ifdef zerol
#define dbg(args...) do { cout << "DBG: " << #args << " -> "; err(args); } while (0)
#else
#define dbg(...)
#endif
void err() { cout << "" << endl; }
template<template<typename...> class T, typename t, typename... Args>
void err(T<t> a, Args... args) { for (auto x: a) cout << x << ' '; err(args...); }
template<typename T, typename... Args>
void err(T a, Args... args) { cout << a << ' '; err(args...); }
// -----------------------------------------------------------------------------

const int N = 1e6+10;
const int M = 21;

int a[N][M];

int w[N];

LL ans = 0;

vector<int> E[N];

void dfs(const int p, const int fa, const int dep) {
    a[dep][0] = w[p];
    FOR(i, 1, M) {
        a[dep][i] = a[dep][i-1];
        int step = 1<<(i-1);
        if(dep>=step) a[dep][i] = min(a[dep][i], a[dep-step][i-1]);
    }
    if(dep!=0) {
        ans+=w[p];
        LL nxt = a[dep-1][0];
        if(dep > 1) {
            int now = dep-2;
            FOR(i, 0, M) {
                dbg(now, 1<<i);
                nxt = min(nxt, a[now][i]*(i+2LL));
                now -= (1<<i);
                if(now<0) break;
            }
        }
        dbg(p, nxt);
        ans += nxt;
    }
    for(auto q: E[p]) if(q!=fa) {
        assert(w[q]>w[p]);
        dfs(q, p, dep+1);
    }
}

int main() {
    int n; scanf("%d", &n); int p = 1;
    FOR(i, 1, n+1) {
        scanf("%d", w+i);
        if(w[i]<w[p]) p = i;
    }
    FOR(i, 1, n) {
        int u, v;
        scanf("%d%d", &u, &v);
        E[u].push_back(v);
        E[v].push_back(u);
    }
    dfs(p, -1, 0);
    printf("%lld\n", ans);
    return 0;
}