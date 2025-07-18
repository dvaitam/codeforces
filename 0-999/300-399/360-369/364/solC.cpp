#include <cstdio>
#include <iostream>
#include <cstdlib>
#include <cmath>
#include <cassert>
#include <cstring>
#include <algorithm>
#include <string>
#include <vector>
#include <list>
#include <set>
#include <map>
#include <sstream>
using namespace std;
#pragma comment(linker, "/STACK:255000000")

typedef long long ll;

#define rep(i, a, b) for(i = (a); i < (b); ++i)
#define repb(i, a, b) for(i = (a) - 1; i >= (b); --i)
#define repd(i, a, b, d) for(i = (a); i < (b); i += (d))
#define repbd(i, a, b, d) for(i = (a) - 1; i >= (b); i -= (d))
#define reps(i, s) for(i = 0; (s)[i]; ++i)
#define repl(i, l) for(i = l.begin(); i != l.end(); ++i)

#define in(f, a) scanf("%"#f, &(a))

bool firstout = 1;

#define out(f, a) printf("%"#f, (a))
#define outf(f, a) printf((firstout) ? "%"#f : " %"#f, (a)), firstout = 0
#define nl printf("\n"), firstout = 1

#define all(x) (x).begin(),(x).end()
#define sqr(x) ((x) * (x))
#define mp make_pair

template<class T>
T &minn(T &a, T b)
{
	if(b < a) a = b;
	return a;
}

template<class T>
T &maxx(T &a, T b)
{
	if(a < b) a = b;
	return a;
}

#define inf 1012345678
#define eps 1e-9


#ifdef XDEBUG
#define mod 23
#else
#define mod 1000000007
#endif

int &madd(int &a, int b)
{
	a += b;
	if(a >= mod) a -= mod;
	return a;
}

int &msub(int &a, int b)
{
	a -= b;
	if(a < 0) a += mod;
	return a;
}

int &mmult(int &a, int b)
{
	return a = (ll)a * b % mod;
}

int mdiv(ll a, ll b, ll m)
{
	a = (a % m + m) % m;
	b = (b % m + m) % m;
	if(a % b == 0) return a / b;
	return (a + m * mdiv(-a, m, b)) / b;
}

#define N 5012
#define M 1012

int n, m, q;
int A[N], AA[N];
int P[15] = {2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47};

void test(int n)
{
	int i, j, k;
	int t = 2 * n * n;
	m = 1;
	A[0] = 1;
	for(k = 0; m < n; ++k)
	{
		if(k == 15) 
		{
			cerr << "error k" << endl;
			break;
		}
		int mm = 0;
		ll d = P[k];
		int l = m;
		int p = 0;
		for(; mm < n && d <= t; d = (d == P[k]) ? 1 : (d == 1) ? sqr(P[k]) : d * P[k])
		{
			rep(i, 0, l) 
			{
				if(d * A[i] <= t) AA[mm++] = A[i] * d;
				else break;
				if(d > 1) ++p;
				if(mm == n) break;
			}
			if(d == P[k]) l = i;
		}
		for(; p >= (mm + 2) / 2 && l < m && mm < n;) AA[mm++] = A[l++];
		if(m == mm)
		{
			cerr << "error " << n << endl;
			break;
		}
		m = mm;
		rep(i, 0, m) A[i] = AA[i];
		sort(A, A + m);		
	}
	if(m != n) cerr << "error not found " << n << endl;
	/*sort(A, A + n);
	rep(i, 1, n) if(A[i] == A[i - 1]) break;
	if(i < n) cerr << "error " << n << endl;
	if(A[0] <= 0) cerr << "error " << n << endl;
	if(A[n - 1] > 2 * n * n) cerr << "error" << endl;
	rep(k, 0, 15)
	{
		int r = 0;
		rep(i, 0, n) if(A[i] % P[k] == 0) ++r;
		if(r && r < (n + 1) / 2) cerr << "error " << n << endl;
	}*/
}

int main()
{
#ifdef XDEBUG
	freopen("in.txt", "rt", stdin);
	//freopen("out.txt", "wt", stdout);
#else
#endif

	int i, j, k;
	char c;
	int a, d;

	//test(725);
	//rep(k, 10, 5001) test(k);

	int ts;	
#if 0
	int tss;
	in(d, tss);
	rep(ts, 1, tss + 1)
#else
	for(ts = 1; in(d, n) > 0; ++ts)
#endif
	{
		test(n);
		rep(i, 0, n) outf(d, A[i]); nl;
	}

	return 0;
}