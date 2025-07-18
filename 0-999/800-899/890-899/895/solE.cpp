#include <cstdio>
#include <algorithm>

struct node {
	double x, t, s;
	int sz;
	node *l, *r;
	node() : x(0), t(1), s(0), sz(1), l(nullptr), r(nullptr) {}
	void upd(double tt, double ts) {
		x = tt * x + ts * sz;
		t *= tt; s = s * tt + ts;
	}
	void pull() {
		x = l->x + r->x;
		sz = l->sz + r->sz;
	}
	void push() {
		l->upd(t, s); r->upd(t, s);
		t = 1, s = 0;
	}
} *root;

int D[100007];

void build(node*& nd, int l, int r)
{
	nd = new node;
	if (l + 1 == r) { nd->x = D[l]; return; }
	int m = (l + r) / 2;
	build(nd->l, l, m); build(nd->r, m, r);
	nd->pull();
}

double query(node* nd, int l, int r, int ql, int qr)
{
	if (ql <= l && r <= qr) return nd->x;
	nd->push();
	int m = (l + r) / 2;
	double ans = 0;
	if (!(qr <= l || m <= ql)) ans = query(nd->l, l, m, ql, qr);
	if (!(qr <= m || r <= ql)) ans += query(nd->r, m, r, ql, qr);
	return ans;
}

void modify(node* nd, int l, int r, int ql, int qr, double t, double s)
{
	if (ql <= l && r <= qr) {
		nd->upd(t, s);
		return;
	}
	nd->push();
	int m = (l + r) / 2;
	if (!(qr <= l || m <= ql)) modify(nd->l, l, m, ql, qr, t, s);
	if (!(qr <= m || r <= ql)) modify(nd->r, m, r, ql, qr, t, s);
	nd->pull();
}

int main()
{
	int N, Q, a, b, c, d;
	scanf("%d%d", &N, &Q);
	for (int i = 0; i < N; i++) scanf("%d", D + i);
	build(root, 0, N);
	while (Q--) {
		scanf("%d%d%d", &d, &a, &b);
		if (d == 1) {
			scanf("%d%d", &c, &d);
			a--, c--;
			double s1 = query(root, 0, N, c, d) / (d - c) / (b - a);
			double s2 = query(root, 0, N, a, b) / (d - c) / (b - a);
			double t1 = 1 - 1. / (b - a);
			double t2 = 1 - 1. / (d - c);
			//printf("%lf %lf %lf %lf\n", s1,s2,t1,t2);
			modify(root, 0, N, a, b, t1, s1);
			modify(root, 0, N, c, d, t2, s2);
		}
		else printf("%.7lf\n", query(root, 0, N, a - 1, b));
	}
}