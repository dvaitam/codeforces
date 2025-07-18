#include <cstdio>
#include <algorithm>
#include <set>
using namespace std;
typedef long long ll;

inline void R(int &x) {
	char ch = getchar(); x = 0;
	for (; ch<'0'; ch=getchar());
	for (; ch>='0'; ch=getchar()) x = x*10+ch-'0';
}
const int N = 100005;
int n, m = 0;
struct pt {
	int l, r, x;
	inline void read() {
		R(l); R(x); R(r);
	}
} a[N], *px[N], *pr[N];
inline bool cmpx(pt *x, pt *y) {return x->x < y->x;}
inline bool cmpr(pt *x, pt *y) {return x->r < y->r;}
namespace seg {
	struct node {
		int mx, mf, sa;
		inline void SA(int x) {
			mx += x, sa += x;
		}
	} a[1048577];
	int tsz, tlv;
	inline void D(int w) {
		if (a[w].sa != 0) {
			a[w<<1].SA(a[w].sa);
			a[(w<<1)|1].SA(a[w].sa);
			a[w].sa = 0;
		}
	}
	void D(int l, int r) {
		for (int i=tlv; i; --i) {
			D(l>>i);
			if ((r>>i)!=(l>>i)) D(r>>i);
		}
	}
	inline void nupd(int w) {
		if (a[w<<1].mx > a[(w<<1)|1].mx) {
			a[w].mx = a[w<<1].mx,
			a[w].mf = a[w<<1].mf;
		}
		else {
			a[w].mx = a[(w<<1)|1].mx,
			a[w].mf = a[(w<<1)|1].mf;
		}
	}
	void B(int n) {
		for (++n,tlv=0,tsz=1; tsz<=n; tsz<<=1,++tlv);
		for (int i=0; i<tsz; ++i) a[i|tsz].mf = i;
		for (int i=tsz-1; i; --i) nupd(i);
	}
	void A(int l, int r, int x) {
		l|=tsz, r|=tsz, --l, ++r;
		D(l, r);
		for (; (l^r)>1; l>>=1, r>>=1) {
			if (!(l&1)) a[l^1].SA(x);
			if (  r&1 ) a[r^1].SA(x);
			nupd(l>>1), nupd(r>>1);
		}
		for (l>>=1; l; l>>=1) nupd(l);
	}
}
int main() {
	R(n);
	for (int i=1; i<=n; ++i) {
		a[i].read(), px[i] = pr[i] = &a[i];
		m = max(m, a[i].r);
	}
	sort(px+1, px+n+1, cmpx);
	sort(pr+1, pr+n+1, cmpr);
	int cpx = 1, cpr = 1, ans = 0, al = 0, ar = 0;
	seg::B(m);
	for (int i=1; i<=m; ++i) {
		for (; cpx<=n && px[cpx]->x<=i; ++cpx) seg::A(px[cpx]->l, i, 1);
		if (seg::a[1].mx > ans) {
			ans = seg::a[1].mx;
			al = seg::a[1].mf;
			ar = i;
		}
		for (; cpr<=n && pr[cpr]->r<=i; ++cpr) seg::A(pr[cpr]->l, pr[cpr]->x, -1);
	}
	printf("%d\n", ans);
	for (int i=1; i<=n; ++i)
		if (a[i].x >= al && a[i].x <= ar && a[i].l <= al && a[i].r >= ar)
			printf("%d ", i);
	return 0;
}