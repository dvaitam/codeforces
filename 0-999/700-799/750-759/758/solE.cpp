#include <bits/stdc++.h>

using namespace std;

#ifdef WIN32
	#define I64 "%I64d"
#else
	#define I64 "%lld"
#endif

typedef long long ll;

#define f first
#define s second
#define mp make_pair
#define pb push_back
#define all(s) s.begin(), s.end()
#define sz(s) (int(s.size()))
#define fname "a"
#define MAXN 500005
#define INF 100000000000000000LL

struct edge {
	int v1, v2;
	int w, p;
	ll d;
};

int n;
edge ed[MAXN];
vector <int> g[MAXN];
bool ok;
ll s[MAXN];

void dfs(int v = 0) {
	for (const int& i : g[v]) {
		edge& e = ed[i];
		e.d = min(e.w - 1, e.p);
		e.w -= e.d;
		e.p -= e.d;
		dfs(e.v2);
		if (s[e.v2] > e.p) {
			ll t = min(s[e.v2] - e.p, e.d);
			e.p += t;
			e.w += t;
			e.d -= t;
			if (s[e.v2] > e.p) ok = 0;
		}
		s[v] += s[e.v2] + e.w;
	}
}

ll go(int v = 0, ll add = INF) {
	ll res = 0;
	for (const int& i : g[v]) {
		edge& e = ed[i];
		ll t = min(add, e.d);
		e.w += t;
		e.p += t;
		e.d -= t;
		add -= t;
		res += t;
		ll tt = go(e.v2, min(add, e.p - s[e.v2]));
		add -= tt;
		res += tt;
	}
	return res;
}

int main()
{
	#ifdef LOCAL
	freopen(fname".in", "r", stdin);
	freopen(fname".out", "w", stdout);
	#endif

	scanf("%d", &n);
	for (int i = 0; i < n - 1; ++i) {
		scanf("%d%d%d%d", &ed[i].v1, &ed[i].v2, &ed[i].w, &ed[i].p);
		--ed[i].v1, --ed[i].v2;
		ed[i].d = 0;
		g[ed[i].v1].pb(i);
	}

	ok = 1;
	dfs();
	if (!ok) {
		puts("-1");
		return 0;
	}

	go();

	printf("%d\n", n);
	for (int i = 0; i < n - 1; ++i) {
		printf("%d %d %d %d\n", ed[i].v1 + 1, ed[i].v2 + 1, ed[i].w, ed[i].p);
	}

	return 0;
}