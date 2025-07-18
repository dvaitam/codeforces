#include <bits/stdc++.h>

using namespace std;

typedef pair<int, int> PA;
typedef long long LL;

#define MAXN 100001
#define MAXL 21
#define FST first
#define SCD second

int n, m, x;
struct Segment {
	int cL, cR;
	PA vL[MAXL], vR[MAXL];
	LL sum;
} seg[MAXN << 1 | 1];

void seg_up(Segment &rt, Segment &lft, Segment &rht) {
	rt.cL = lft.cL;
	memcpy(rt.vL, lft.vL, rt.cL * sizeof(PA));
	for(int i = 0; i < rht.cL; ++i) {
		PA pre = rt.vL[rt.cL - 1];
		PA cur = rht.vL[i];
		cur.FST |= pre.FST;
		if(pre.FST < cur.FST)
			rt.vL[rt.cL++] = cur;
	}
	rt.cR = rht.cR;
	memcpy(rt.vR, rht.vR, rt.cR * sizeof(PA));
	for(int i = 0; i < lft.cR; ++i) {
		PA pre = rt.vR[rt.cR - 1];
		PA cur = lft.vR[i];
		cur.FST |= pre.FST;
		if(pre.FST < cur.FST)
			rt.vR[rt.cR++] = cur;
	}
	rt.sum = lft.sum + rht.sum;
	for(int i = lft.cR - 1, j = 0; i >= 0; --i) {
		for( ; j < rht.cL && (lft.vR[i].FST | rht.vL[j].FST) < x; ++j);
		if(j < rht.cL)
			rt.sum += (lft.vR[i].SCD - (i + 1 < lft.cR ? lft.vR[i + 1].SCD : lft.vL[0].SCD - 1)) * (rht.vR[0].SCD - rht.vL[j].SCD + 1LL);
	}
}

inline int seg_idx(int L, int R) {
	return (L + R) | (L < R);
}

void seg_build(int L, int R) {
	Segment &rt = seg[seg_idx(L, R)];
	if(L == R) {
		int val; 
		scanf("%d", &val);
		rt.cL = rt.cR = 1;
		rt.vL[0] = rt.vR[0] = (PA){val, L};
		rt.sum = val >= x;
	} else {
		int M = (L + R) >> 1;
		seg_build(L, M);
		seg_build(M + 1, R);
		seg_up(rt, seg[seg_idx(L, M)], seg[seg_idx(M + 1, R)]);
	}
}

void seg_upd(int L, int R, int u, int v) {
	Segment &rt = seg[seg_idx(L, R)];
	if(L == R) {
		rt.cL = rt.cR = 1;
		rt.vL[0] = rt.vR[0] = (PA){v, L};
		rt.sum = v >= x;
	} else {
		int M = (L + R) >> 1;
		if(u <= M)
			seg_upd(L, M, u, v);
		else
			seg_upd(M + 1, R, u, v);
		seg_up(rt, seg[seg_idx(L, M)], seg[seg_idx(M + 1, R)]);
	}
}

Segment seg_que(int L, int R, int u, int v) {
	Segment ret = {};
	if(u <= L && R <= v) {
		ret = seg[seg_idx(L, R)];
	} else {
		int M = (L + R) >> 1;
		if(v <= M) {
			ret = seg_que(L, M, u, v);
		} else if(u > M) {
			ret = seg_que(M + 1, R, u, v);
		} else {
			Segment lft = seg_que(L, M, u, v);
			Segment rht = seg_que(M + 1, R, u, v);
			seg_up(ret, lft, rht);
		}
	}
	return ret;
}

int main() {
	#ifndef ONLINE_JUDGE
	freopen("in.txt", "r", stdin);
	#endif
	scanf("%d%d%d", &n, &m, &x);
	seg_build(1, n);
	while(m--) {
		int typ, x, y;
		scanf("%d%d%d", &typ, &x, &y);
		if(typ == 1) {
			seg_upd(1, n, x, y);
		} else {
			Segment res = seg_que(1, n, x, y);
			printf("%lld\n", res.sum);
		}
	}
	return 0;
}