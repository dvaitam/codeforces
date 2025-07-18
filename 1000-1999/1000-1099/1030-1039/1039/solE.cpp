#include <iostream>
#include <cstdio>
#include <cmath>
#include <vector>
#include <cstring>
#include <algorithm>
using namespace std;
#define debug(x) cout << #x << "=" << x << endl 
const int N = 1e5 + 5;

namespace LCT {
	struct Node {
		int ch[2], sz, fa;
		Node (int x = 0) {
			ch[0] = ch[1] = 0, sz = x, fa = 0;
		}
	} a[N];
	#define ls(x) a[x].ch[0]
	#define rs(x) a[x].ch[1]
	#define fa(x) a[x].fa
	#define sz(x) a[x].sz
	bool isroot(int x) {
		return ls(fa(x)) != x && rs(fa(x)) != x;
	}
	bool get(int x) {
		return rs(fa(x)) == x; 
	}
	void pushup(int x) {
		sz(x) = sz(ls(x)) + sz(rs(x)) + 1;
	} 
	void rotate(int x) {
		int y = fa(x), z = fa(y), k = get(x);
		if (!isroot(y))
			a[z].ch[get(y)] = x;
		fa(a[x].ch[1 - k]) = y;
		a[y].ch[k] = a[x].ch[1 - k];
		a[x].ch[1 - k] = y;
		fa(y) = x, fa(x) = z;
		pushup(y), pushup(x); 
	} 
	void splay(int x) {
		while (!isroot(x)) {
			if (!isroot(fa(x)) && get(x) == get(fa(x)))
				rotate(fa(x));
			rotate(x);
		} 
	}
	void access(int x) {
		int p = 0;
		while (x) {
			splay(x);
			rs(x) = p;
			pushup(x);
			p = x, x = fa(x); 
		}
	}
	int findroot(int x) {
		access(x);
		splay(x);
		while (ls(x))
			x = ls(x);
		splay(x);
		return x;
	}
	void link(int x, int y) {
		splay(x);
		fa(x) = y; 
	}
	void cut(int x) {//x和父节点砍掉 
		access(x);
		splay(x);
		ls(x) = fa(ls(x)) = 0; 
	}
	int ask(int x) {//查询 x 到根的距离 
		access(x);
		splay(x);
		return a[x].sz;
	}
}

int n, q;
int a[N] = {0};
typedef pair<int, int> pii;
pii c[N];
int b;

int logN[N] = {0};
pii st[N][22];
pii mrg(pii x, pii y) {
	return make_pair(max(x.first, y.first), min(x.second, y.second));
}
void initST() {
	logN[1] = 0;
	for (int i = 2; i <= n; i++)
		logN[i] = logN[i / 2] + 1;
	for (int i = 1; i <= n; i++)
		st[i][0] = make_pair(a[i], a[i]);
	for (int j = 1; (1 << j) <= n; j++)
		for (int i = 1; i + (1 << j) - 1 <= n; i++)	
			st[i][j] = mrg(st[i][j - 1], st[i + (1 << (j - 1))][j - 1]);
}
int qry(int l, int r) {
	int k = logN[r - l + 1];
	pii res = mrg(st[l][k], st[r - (1 << k) + 1][k]);
	return res.first - res.second;
}
int fnd(int pos, int w) {
	int l = pos, r = n + 1;
	while (l + 1 < r) {
		int mid = (l + r) / 2;
		if (qry(pos, mid) > w)
			r = mid;
		else
			l = mid;
	}
	return r;
}

vector<pii> add[N];

int res[N] = {0};

int main() {
	int w;
	scanf("%d%d%d", &n, &w, &q);
	for (int i = 1; i <= n; i++)
		scanf("%d", &a[i]), LCT::a[i] = LCT::Node(1);
	initST();
	for (int i = 1; i <= q; i++)
		scanf("%d", &c[i].first), c[i].first = w - c[i].first, c[i].second = i;
	b = min(100.0, sqrt(n) + 0.5); 
	sort(c + 1, c + q + 1);
	for (int i = 1; i <= n; i++) {
		int dif = 0;
		int nxt = fnd(i, 0);
	//	debug(i);
	//	cout << fnd(1, 1) << endl;
		while (nxt <= min(n, i + b)) {
			int pos = lower_bound(c + 1, c + q + 1, make_pair(dif, 0)) - c;
	//		debug(pos), debug(nxt);
			add[pos].push_back(make_pair(i, nxt));
			dif = qry(i, nxt);
			nxt = fnd(i, dif);
		}
	//	debug(nxt);
	//	debug(dif);
		int pos = lower_bound(c + 1, c + q + 1, make_pair(dif, 0)) - c;
		add[pos].push_back(make_pair(i, 0));
	}
	for (int i = 1; i <= q; i++) {
	//	debug(i);
	//	debug(c[i].first);
	//	debug(c[i].second);
		for (auto j: add[i]) {
	//		debug(j.first);
	//		debug(j.second);
			LCT::cut(j.first);
			if (j.second)
				LCT::link(j.first, j.second);
		}
		int cur = 1, ans = 0;
		while (cur <= n) {
			ans += LCT::ask(cur);
			cur = LCT::findroot(cur);
			cur = fnd(cur, c[i].first);
		}
		res[c[i].second] = ans - 1;
	}
	for (int i = 1; i <= q; i++)
		printf("%d\n", res[i]);
	return 0;
}