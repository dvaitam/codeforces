#include <iostream>
#include <iomanip>
#include <sstream>
#include <cstdio>
#include <cmath>
#include <cstring>
#include <cctype>
#include <string>
#include <vector>
#include <list>
#include <set>
#include <map>
#include <queue>
#include <stack>
#include <algorithm>
#include <cstdio>
#include <functional>
#include <numeric>
#include <unordered_map>
#include <unordered_set>
using namespace std;
// #define int long long
#define DEBUG(x) cout << '>' << #x << ':' << x << endl;
#define REP(i,n) for(int i=0;i<(n);i++)
#define FOR(i,a,b) for(int i=(a);i<=(b);i++)
#define FORD(i,a,b) for(int i=(a);i>=(b);i--)
inline bool EQ(double a, double b) { return fabs(a-b) < 1e-9; }
// const int INF = (((1<<30)-1)<<1)+1;
const int nINF = 1<<31;
typedef long long ll;
typedef unsigned long long ull;
typedef pair<int, int> pii;
typedef pair<ll, ll> pll;
typedef pair<char, int> pci;

/////////////////////////////////////////////////////////////////////

const int MAXN = 2010;
const int INF = 1000000000;

class node {
public:
    int v, f, couple;
    pii w;
    node(int V, pii W, int F, int CC) {v = V; w = W; f = F; couple = CC;};
};

int fa[MAXN], ff[MAXN];
vector<node> g[MAXN];

// s, t, capacity, cost
void newnode(int a, int b, int c, pii d) {
    node p(b, d, c, g[b].size());
    node q(a, {-d.first, -d.second}, 0, g[a].size());
    g[a].push_back(p);
    g[b].push_back(q);
}

pii SPFA(int s, int t, int n) {
    vector<pii> d(n + 1, {INF, INF});
    vector<bool> z(n + 1, false);
    d[s] = {0, 0};
    queue<int> q;
    q.push(s);
    z[s] = true;
    while(!q.empty()) {
        int i = q.front();
        q.pop();
        for(auto p = g[i].begin(); p != g[i].end(); ++p) {
        	pii dc = {d[i].first + p->w.first, d[i].second + p->w.second};
            if(p->f > 0 && d[p->v] > dc) {
                d[p->v] = dc;
                fa[p->v] = i;
                ff[p->v] = p - g[i].begin();
                if(!z[p->v]) {
                    z[p->v] = true;
                    q.push(p->v);
                }
            }
        }
        z[i] = false;
    }
    return d[t];
}

pii min_cost(int s, int t, int n) {
    int flow = 0;
    pii cost = {0, 0};
    while(SPFA(s, t, n) < (pii){0, 0}) { //if using negative cost, replace INF with 0
        int i = t, now = INF;
        while(i != s) {
            node &p = g[fa[i]][ff[i]];
            i = fa[i];
            if(p.f < now) now = p.f;
        }
        flow += now;
        
        i = t;
        while(i != s) {
            g[fa[i]][ff[i]].f -= now;
            cost = {cost.first + now * g[fa[i]][ff[i]].w.first, cost.second + now * g[fa[i]][ff[i]].w.second};
            g[i][g[fa[i]][ff[i]].couple].f += now;
            i = fa[i];
        }
    }
    return cost;
}

void bad() {
	cout << -1 << endl;
	exit(0);
}

const int N = 510;
int n, x, y, a[N], q1, q2, r1[N], r2[N], src, tgt, wts = 0;
vector<int> g1[N], g2[N];

void dfs1(int s, int p = 0, int pre = src) {
	if(pre) {
		if(r1[s]) {
			// Add (pre, n+s) with capacity r1[s], cost {-1, 0}
			newnode(pre, n+s, r1[s], {-1, 0});
			// Add (n+s, s) with capacity 1, cost {0, -a[s]}
			newnode(n+s, s, 1, {0, -a[s]});
		} else {
			if(pre != src) {
				// If pre != src, add (pre, s) with capacity 1, cost {0, -a[s]}
				newnode(pre, s, 1, {0, -a[s]});
			}
		}
	}
	for(int v : g1[s]) {
		if(v != p) {
			dfs1(v, s, r1[s] ? n+s : pre);
		}
	}
}

void dfs2(int s, int p = 0, int pre = tgt) {
	if(pre) {
		if(r2[s]) {
			// Add (n+n+s, pre) with capacity r2[s], cost {-1, 0}
			newnode(n+n+s, pre, r2[s], {-1, 0});
			// Add (s, n+n+s) with capacity 1, cost {0, 0}
			newnode(s, n+n+s, 1, {0, 0});
		} else {
			if(pre != tgt) {
				// If pre != tgt, add (s, pre) with capacity 1, cost {0, 0}
				newnode(s, pre, 1, {0, 0});
			}
		}
	}
	for(int v : g2[s]) {
		if(v != p) {
			dfs2(v, s, r2[s] ? n+n+s : pre);
		}
	}
}

signed main() {
    ios::sync_with_stdio(false);
    //cout << fixed << setprecision(7);
    cin >> n >> x >> y;
    FOR(i, 1, n) cin >> a[i];
    REP(_, n-1) {
    	int u, v;
    	cin >> u >> v;
    	g1[u].push_back(v);
    	g1[v].push_back(u);
    }
    REP(_, n-1) {
    	int u, v;
    	cin >> u >> v;
    	g2[u].push_back(v);
    	g2[v].push_back(u);
    }
    cin >> q1;
    REP(_, q1) {
    	int k, num;
    	cin >> k >> num;
    	if(r1[k]) bad();
    	r1[k] = num;
    	wts += num;
    }
    cin >> q2;
    REP(_, q2) {
    	int k, num;
    	cin >> k >> num;
    	if(r2[k]) bad();
    	r2[k] = num;
    	wts += num;
    }
    src = n+n+n+1;
    tgt = n+n+n+1+1;
    dfs1(x);
    dfs2(y);
    // FOR(i, 1, 3*n+2) {
    // 	for(auto& p : g[i]) {
    // 		cout << "{" << p.v << ", " << p.f << ", {" << p.w.first << ", " << p.w.second << "}} ";
    // 	}
    // 	cout << endl;
    // }
    pii ans = min_cost(src, tgt, tgt);
    if(ans.first != -wts) {
    	bad();
    } else {
    	cout << -ans.second << endl;
    }

    return 0;
}

// Easy things to check:
// - LONG LONGS (WA)
// - const int N is correct (WA, RTE)
// - .size()-k underflow (WA, RTE, TLE)
// - small n edge cases (WA)

// Rare mistakes made in the past:
// - division by 0 (WA)
// - integer division (WA)
// - setprecision (WA)
// - INF not big enough using ll (WA)
// - setting min to 0 instead of -INF (WA)
// - outputting debug (WA)
// - allocating too much memory (locRTE, MLE)
// - stack size (locRTE)