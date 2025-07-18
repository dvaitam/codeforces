#include <bits/stdc++.h>
using namespace std;
#define MP make_pair
#define PB push_back
#define fo(i,a,b) for(int i=(a);i<(b);i++)
using ll = long long;

int N, a[100100];
ll totgood, good;

struct tree {
	ll fen[100100];
	void u (int i, ll v) { while (i <= N) fen[i] += v, i += i & (-i); }
	ll q (int i) {
		ll r = 0;
		while (i >= 1) r += fen[i], i -= i & (-i);
		return r;
	}
} g, k;

int main () {
	scanf("%d", &N);
	fo(i, 0, N) {
		scanf("%d", &a[i]);
		good += g.q(a[i]-1), g.u(a[i], 1);
		totgood += k.q(a[i]-1) * (N-i), k.u(a[i], i+1);
	}

	long double n = N;
	long double dv = n * (n+1.);
	long double v = (n * (n-1) / 2. - ((long double) good)) * dv;

	long double f = (n-1) * n * (n+1) * (n+2) / 24.;
	long double t = ((long double) totgood*2) - f;

	printf("%.16lf\n", double((v+t) / dv));

	/*
	long double n = N;
	long double dv = n * (n+1.) / 2.;
	long double v = (n * (n-1) / 2. - ((long double) good)) * dv;

	long double f = (n-1) * n * (n+1) * (n+2) / 48.;
	long double t = ((long double) totgood) - f;

	printf("%.12lf\n", double((v+t) / dv));
	*/

	return 0;
}