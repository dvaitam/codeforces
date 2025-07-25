#include <cstdio>
#include <cstring>
#include <iostream>
#include <algorithm>
#include <vector>
using namespace std;

#define LY(p) freopen (p".in", "r", stdin); freopen (p".out", "w", stdout)
#define LL long long
#define dbl double
#define ld long double
#define LLD "%I64d"
#define pb push_back
#define N 200010

int n, m, x, y, h[N], ent = 1;
int tim, cnt, vis[N], dep[N], fa[N], U[N], V[N], W[N], R[N];

struct edge {
	int v, n;
	edge (int y = 0, int t = 0): v(y), n(t) {}
} e[N << 1];

void link (int x, int y) {e[++ ent] = edge (y, h[x]), h[x] = ent;}

vector<int> get (int x, int y) {
	vector<int> a;
	for (; x != y; x = fa[x]) a.pb (x);
	a.pb (y);
	return a;
}

void print (vector<int> p) {
	printf ("%d ", p.size());
	for (int i = 0, t = p.size(); i < t; i++)
		printf ("%d ", p[i]);
	printf ("\n");
}

void print (int st, int en, int x, int y) {
	vector<int> a, b;
	a = get (x, st);
	reverse (a.begin(), a.end());
	b = get (en, y);
	reverse (b.begin(), b.end());
	a.insert (a.end(), b.begin(), b.end());
	print (a);
}

void print (int x, int y) {
	printf ("YES\n");
	int st = W[cnt];
	while (dep[st] != V[cnt]) st = fa[st];
	int en = dep[y] > U[cnt]? y : R[cnt];
	vector<int> a = get (st, en);
	print (a);
	print (st, en, x, y);
	print (st, en, W[cnt], R[cnt]);
	exit(0);
}

void dfs (int o, int ft, int d) {
	vis[o] = ++ tim, dep[o] = d;
	for (int x = h[o], y; y = e[x].v, x; x = e[x].n)
		if (x ^ ft ^ 1)
			if (! vis[y]) {
				fa[y] = o, dfs (y, x, d + 1);
				while (cnt && U[cnt] == d) cnt --;
				if (cnt && V[cnt] > d) V[cnt] = d;
			}
			else if (vis[y] < vis[o]) {
				if (cnt && dep[y] < V[cnt])
					print (o, y);
				while (cnt && U[cnt] > dep[y]) cnt --;
				cnt ++;
				U[cnt] = dep[y], V[cnt] = dep[o], R[cnt] = y, W[cnt] = o;
			}
}

int main()
{
#ifndef ONLINE_JUDGE
	LY("E");
#endif
	scanf ("%d %d", &n, &m);
	for (int i = 1; i <= m; i++)
		scanf ("%d %d", &x, &y), link (x, y), link (y, x);

	for (int i = 1; i <= n; i++)
		if (! vis[i])
			dfs (i, -1, 1);
	printf ("NO");
	return 0;
}