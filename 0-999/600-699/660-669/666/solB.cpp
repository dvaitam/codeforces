#ifdef _MSC_VER
#define _CRT_SECURE_NO_DEPRECATE
#pragma comment(linker, "/STACK:66777216")
#else
#pragma GCC optimize("O3")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx")
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
#include <cstdio>
#include <cstdlib>
#include <sstream>
#include <fstream>
#include <ctime>
#include <cstring>
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
VI sm[3002];
int dout[3002][3002];
int q[3002];
int res[4];
pii mdout[3002][3], mdin[3002][3];
inline void bfs(int s, VI sm[], int ds[]) {
	rept(i, n) ds[i] = INF;
	ds[s] = 0;
	int a = 0, b = 0;
	q[b++] = s;
	while (a < b) {
		int v = q[a++];
		rept(i, L(sm[v])) {
			if (ds[sm[v][i]] > ds[v] + 1) {
				ds[sm[v][i]] = ds[v] + 1;
				q[b++] = sm[v][i];
			}
		}
	}
}
inline void upd(pii md[], pii t) {
	if (t > md[0]) {
		md[2] = md[1];
		md[1] = md[0];
		md[0] = t;
	}
	else if (t > md[1]) {
		md[2] = md[1];
		md[1] = t;
	}
	else if (t > md[2]) {
		md[2] = t;
	}
}

int main() {
	//freopen("input.txt", "r", stdin);
	//freopen("output.txt", "w", stdout);
	scanf("%d%d", &n, &m);
	rept(i, m) {
		scanf("%d%d", &a, &b); --a; --b;
		if (a == b) continue;
		sm[a].pb(b);
	}
	rept(i, n) {
		SORT(sm[i]);
		sm[i].resize(unique(all(sm[i])) - sm[i].begin());
	}

	
	memset(mdout, -1, sizeof(mdout));
	memset(mdin, -1, sizeof(mdin));
	rept(i, n) {
		bfs(i, sm, dout[i]);

		rept(j, n) {
			if (dout[i][j] < INF && j != i) {
				pii t = mp(dout[i][j], j);
				upd(mdout[i], t);
				
				t = mp(dout[i][j], i);
				upd(mdin[j], t);
			}
		}
	}

	int ans = -1;
	rept(i, n) {
		rept(j, n) {
			if (i == j || dout[i][j] >= INF) continue;
			rept(p0, 3) {
				if (mdin[i][p0].x == -1 || mdin[i][p0].y == j) continue;
				rept(p1, 3) {
					if (mdout[j][p1].x == -1 || mdout[j][p1].y == mdin[i][p0].y || mdout[j][p1].y == i) continue;
					int cur = dout[i][j] + mdin[i][p0].x + mdout[j][p1].x;
					if (cur > ans) {
						ans = cur;
						res[0] = mdin[i][p0].y;
						res[1] = i;
						res[2] = j;
						res[3] = mdout[j][p1].y;
					}
				}
			}
		}
	}

	printf("%d %d %d %d\n", res[0] + 1, res[1] + 1, res[2] + 1, res[3] + 1);
}