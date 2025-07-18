#include <cstdio>
#include <cstring>
#include <algorithm>
#include <queue>
#include <vector>
#define LL long long
#define MEM(x, y) memset(x, y, sizeof(x));
using namespace std;

const int INF = 0x3f3f3f3f;
const int P = 1e9+7;
const int maxn = 1e3+10;

struct Edge {
	int to, next;
} edges[maxn<<1];
int first[maxn], mm = 0;

void AddEdge(int u, int v) {
	Edge &e = edges[++mm];
	e.to = v, e.next = first[u];
	first[u] = mm;
}

struct Node {
	int u, d;
	bool operator < (const Node &x) const {
		return d < x.d;
	}
};

int n, m, ans = 0;
//ÿ�������������ӵ�����ĳ��� 
priority_queue<Node> vex;
//int vis[maxn];
vector<pair<int, int> > E;

int fa[maxn], vis[maxn], dp[maxn], f[maxn];
int Find(int x) {
	return x == fa[x] ? x : fa[x] = Find(fa[x]);
}
void Union(int x, int y) {
	fa[Find(x)] = Find(y);
}

int dfs(int u, int f) {
//	vis[u] = 1;
	int res = 1, mx = 0, t = 0;
	for (int i = first[u]; i; i = edges[i].next) {
		Edge &e = edges[i];
		if (e.to != f) {
			t = dfs(e.to, u);
			mx = max(mx, t);
		}
	}
	return res + mx;
}
int temp;
void dfs_d(int u, int f, int d) {
    if (d > ans) {
        ans = d;
        temp = u;
    }
    for (int i = first[u]; i; i = edges[i].next) {
        Edge &e = edges[i];
        if (e.to != f)
            dfs_d(e.to, u, d + 1);
    }
}

int main() {
	scanf("%d%d", &n, &m);
	for (int i = 1; i <= n; i++) {
		fa[i] = i;
		dp[i] = INF;
	}
	for (int i = 0; i < m; i++) {
		int u, v;
		scanf("%d%d", &u, &v);
		AddEdge(u, v);
		AddEdge(v, u);
		Union(u, v);
	}
	for (int i = 1; i <= n; i++) {
		int t = dfs(i, 0);
		if (t < dp[Find(i)]) {
			dp[Find(i)] = t;
			f[Find(i)] = i;
		}
//		dp[Find(i)] = min(t, dp[Find(i)]);
//		printf("%d fa %d, t = %d, dp[%d] = %d\n", i, Find(i), t, Find(i), dp[Find(i)]);
	}
//		if (!vis[i]) {
//			vex.push(Node{i, dfs(i, 0)});
//		}
	for (int i = 1; i <= n; i++) {
		if (!vis[Find(i)]) {
			vis[Find(i)] = 1;
			vex.push(Node{f[Find(i)], dp[Find(i)]});
		}
	}
	Node x = vex.top(); vex.pop();
//	int ans = x.d - 1;
//	if (!vex.empty())
//		ans += vex.top().d;
//	printf("%d\n", ans);
//	while (!vex.empty()) {
//		Node y = vex.top(); vex.pop();
//		printf("%d %d\n", x.u, y.u);
//	}
	while (!vex.empty()) {
//		Node x = vex.top(); vex.pop();
		Node y = vex.top(); vex.pop();
		E.push_back(make_pair(x.u, y.u));
		AddEdge(x.u, y.u);
		AddEdge(y.u, x.u);
//		x.d = max(x.d, y.d + 1);
	}
	dfs_d(1, 0, 0);
	dfs_d(temp, 0, 0);
	printf("%d\n", ans);
	for (int i = 0; i < E.size(); i++)
		printf("%d %d\n", E[i].first, E[i].second);
	
	return 0;
}