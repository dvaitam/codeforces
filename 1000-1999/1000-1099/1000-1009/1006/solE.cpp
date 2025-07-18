#include <bits/stdc++.h>
using namespace std;
#define REP(i, a, b) for (int i = (a), i##_end_ = (b); i <= i##_end_; ++i)
#define RPD(i, b, a) for (int i = (b), i##_end_ = (a); i >= i##_end_; --i)
#define pii pair<int, int>
#define PB push_back
#define SZ(x) (int)((x).size())

typedef long long LL;
const int oo = 0x3f3f3f3f;
const int MAXN = 200010;

inline int read()
{
    char c = getchar(); int x = 0, f = 1;
    while(c < '0' || c > '9'){if(c == '-') f = -1; c = getchar();}
    while(c >= '0' && c <= '9'){x = x * 10 + (c - '0'); c = getchar();}
    return x * f;
}

vector<int> g[MAXN];
int st[MAXN], ed[MAXN], id[MAXN], num;

void dfs(int x)
{
    st[x] = ++ num;
    id[num] = x;
    for(auto y : g[x]) dfs(y);
    ed[x] = num;
}

int main()
{
    int n = read(), q = read();
    REP(i, 2, n){
        int fa = read();
        g[fa].PB(i);
    }
    REP(i, 1, n) sort(g[i].begin(), g[i].end());
    dfs(1);
    //REP(i, 1, n) cout << st[i] << ' '; cout << endl;
    //REP(i, 1, n) cout << id[i] << ' '; cout << endl;
    REP(i, 1, q){
        int u = read(), v = read();
        printf("%d\n", st[u] + v - 1 <= ed[u] ? id[st[u] + v - 1] : -1);
    }
    return 0;
}