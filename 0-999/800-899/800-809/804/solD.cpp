#ifdef _MSC_VER
#define _CRT_SECURE_NO_DEPRECATE
#pragma comment(linker, "/STACK:66777216")
#include <cstdio>
#else
#pragma GCC optimize("O3")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx")
#include <stdio.h>
#endif
#include <algorithm>
#include <string>
#include <set>
#include <map>
#include <vector>
#include <queue>
#include <iostream>
#include <iterator>
#include <cmath>
#include <cstdlib>
#include <sstream>
#include <fstream>
#include <ctime>
#include <cstring>
#include <functional>
#include <unordered_map>
#include <chrono>
using namespace std;
#define pb push_back
#define ppb pop_back
#define pi 3.1415926535897932384626433832795028841971
#define mp make_pair
#define x first
#define y second
#define pii pair<int,int>
#define pdd pair<double,double>
#define INF 1000000000
#define FOR(i,a,b) for (int _n(b), i(a); i <= _n; i++)
#define FORD(i,a,b) for(int i=(a),_b=(b);i>=_b;i--)
#define all(c) (c).begin(), (c).end()
#define SORT(c) sort(all(c))
#define rep(i,n) FOR(i,1,(n))
#define rept(i,n) FOR(i,0,(n)-1)
#define L(s) (int)((s).size())
#define C(a) memset((a),0,sizeof(a))
#define VI vector <int>
#define ll long long

int a, b, c, d, n, m, k;
VI sm[100002];
int num[100002], tot[100002];
bool cool[100002];
int down[100002], ds[100002];

vector<int> sr[100002];
vector<ll> ssr[100002];

pii dfs(int v, int c) {
	num[v] = c;
	pii t(1, L(sm[v]));
	rept(i, L(sm[v])) {
		int w = sm[v][i];
		if (num[w] != -1) continue;
		auto tt = dfs(w, c);
		t.x += tt.x;
		t.y += tt.y;
	}
	return t;
}


void dfs2(int v, int pr) {
	down[v] = 0;
	rept(i, L(sm[v])) {
		int w = sm[v][i];
		if (w == pr) continue;
		dfs2(w, v);
		down[v] = max(down[v], down[w] + 1);
	}
}
void dfs3(int v, int pr, int sof, int h, VI& res) {
	ds[v] = max(down[v], sof + h);
	ds[v] = max(ds[v], h);
	res.pb(ds[v]);

	int mx = -INF, mx2 = -INF;
	rept(i, L(sm[v])) {
		int w = sm[v][i];
		if (w == pr) continue;
		if (down[w] > mx) {
			mx2 = mx;
			mx = down[w];
		}
		else if (down[w] > mx2) {
			mx2 = down[w];
		}
	}

	rept(i, L(sm[v])) {
		int w = sm[v][i];
		if (w == pr) continue;
		int nx = sof;
		if (down[w] != mx) {
			nx = max(nx, mx - h + 1);
		}
		else {
			nx = max(nx, mx2 - h + 1);
		}
		dfs3(w, v, nx, h + 1, res);
	}
}
inline void prepare(int s) {
	dfs2(s, -1);
	VI cur;
	dfs3(s, -1, -INF, 0, cur);
	int t = *max_element(all(cur));
	int c = num[s];
	tot[c] = L(cur);
	sr[c].resize(t + 1);
	ssr[c].resize(t + 1);
	rept(i, L(cur)) {
		++sr[c][cur[i]];
	}
	FORD(i, L(sr[c]) - 2, 0) {
		sr[c][i] += sr[c][i + 1];
	}
	FORD(i, L(sr[c]) - 1, 0) {
		ssr[c][i] = sr[c][i];
		if (i < L(sr[c]) - 1) {
			ssr[c][i] += ssr[c][i + 1];
		}
	}
}

unordered_map<ll, double> mem;
double solve(int a, int b) {
	if (a > b) swap(a, b);
	ll h = (ll)a * 1234567 + (a ^ b);
	auto it = mem.find(h);
	if (it != mem.end()) return it->second;

	if (L(sr[a]) > L(sr[b])) swap(a, b);
	int da = L(sr[a]) - 1;
	int db = L(sr[b]) - 1;
	int r = max(da, db);

	long double sum = 0.0;
	rept(i, L(sr[a])) {
		int need = max(r - i, 0);
		int cc = sr[a][i];
		if (i < L(sr[a]) - 1) cc -= sr[a][i + 1];
		if (!cc) continue;
		if (need >= L(sr[b])) {
			sum += (long double)tot[b] * cc * r;
		}
		else {
			int cnt = sr[b][need];
			int lf = tot[b] - cnt;
			sum += (long double)lf * cc * r;
			int d0 = need + i + 1;

			sum += (long double)d0 * cc * sr[b][need];
			if (need < L(sr[b]) - 1) {
				sum += (long double)cc * ssr[b][need + 1];
			}
		}
	}

	sum /= (ll)tot[a] * tot[b];
	return mem[h] = sum;
}
int main() {
	//freopen("input.txt", "r", stdin);
	//freopen("output.txt", "w", stdout);
	scanf("%d%d%d", &n, &m, &k);
	rept(i, m) {
		scanf("%d%d", &a, &b); --a; --b;
		sm[a].pb(b);
		sm[b].pb(a);
	}

	memset(num, -1, sizeof(num));
	c = 0;
	rept(i, n) {
		if (num[i] != -1) continue;
		pii t = dfs(i, c);
		if (t.y == t.x * 2 - 2) {
			cool[c] = 1;
			prepare(i);
		}
		++c;
	}

	rept(i, k) {
		scanf("%d%d", &a, &b); --a; --b;
		int va = num[a], vb = num[b];
		if (va == vb || !cool[va] || !cool[vb]) {
			printf("-1\n");
			continue;
		}
		double ans = solve(va, vb);
		printf("%.9lf\n", ans);
	}
}