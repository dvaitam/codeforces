#include <bits/stdc++.h>

using namespace std;



#define getchar() (p1==p2&&(p2=(p1=buf)+fread(buf,1,1<<21,stdin),p1==p2)?EOF:*p1++)

char buf[1 << 21], *p1 = buf, *p2 = buf;



inline int qread() {

	register char c = getchar();

	register int x = 0, f = 1;

	while (c < '0' || c > '9') {

		if (c == '-') f = -1;

		c = getchar();

	}

	while (c >= '0' && c <= '9') {

		x = (x << 3) + (x << 1) + c - 48;

		c = getchar();

	}

	return x * f;

}



const int N = 300005;

struct Edge {

	int to, nxt;

	Edge() {

		nxt = -1;

	}

};

Edge e[N << 1];

int n, hd[N], pnt, m, f[N], cnt, vis[N];



inline void AddEdge(int u, int v) {

	e[++pnt].to = v;

	e[pnt].nxt = hd[u];

	hd[u] = pnt;

}



inline int GetRoot(int v) {

	if (f[v] == v) return v;

	return f[v] = GetRoot(f[v]);

}



inline bool Merge(int x, int y) {

	int u = GetRoot(x), v = GetRoot(y);

	if (u != v) {

		f[v] = u;

		return 1;

	}

	return 0;

}



inline void Read() {

	n = qread(); m = qread();

	for (int i = 1;i <= n;i++) {

		hd[i] = -1;

		f[i] = i;

	}

	pnt = 0;

	cnt = n;

	for (int i = 1;i <= m;i++) {

		int u = qread(), v = qread();

		AddEdge(u, v);

		AddEdge(v, u);

		cnt -= Merge(u, v);

	}

}



inline void Solve() {

	if (cnt > 1) {

		cout << "NO\n";

		return;

	}

	for (int i = 1;i <= n;i++) vis[i] = -1;

	queue <int> que;

	vis[1] = 1;

	for (int i = hd[1];~i;i = e[i].nxt) {

		que.push(e[i].to);

		vis[e[i].to] = 0;

	}

	while (!que.empty()) {

		int u = que.front();

		que.pop();

		/*

		bool flag = 0;

		for (int i = hd[u];~i;i = e[i].nxt) {

			if (vis[e[i].to] == 1) {

				flag = 1;

				break;

			}

		}

		if (!flag) {*/

			for (int i = hd[u];~i;i = e[i].nxt) {

				if (vis[e[i].to] == -1) {

					vis[e[i].to] = 1;

					for (int j = hd[e[i].to];~j;j = e[j].nxt) {

						if (vis[e[j].to] == -1) {

							vis[e[j].to] = 0;

							que.push(e[j].to);

						}

					}

					// break;

				}

			}

		//}

	}

	int sum = 0;

	for (int i = 1;i <= n;i++) sum += vis[i];

	cout << "YES\n" << sum << "\n";

	for (int i = 1;i <= n;i++) {

		if (vis[i] == 1) cout << i << " ";

	}

	cout << "\n";

}



int main() {

	int t = qread();

	while (t--) {

		Read();

		Solve();

	}

	return 0;

}