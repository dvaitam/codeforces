#include <bits/stdc++.h>



using namespace std;



#ifdef SG

	#include <debug.h>

#else

	#define show(...)

	#define debug(...)

	#define deepen(...)

	#define timer(...)

#endif



#define ARG4(_1,_2,_3,_4,...) _4



#define forn3(i,l,r) for (int i = int(l); i < int(r); ++i)

#define forn2(i,n) forn3 (i, 0, n)

#define forn(...) ARG4(__VA_ARGS__, forn3, forn2) (__VA_ARGS__)



#define ford3(i,l,r) for (int i = int(r) - 1; i >= int(l); --i)

#define ford2(i,n) ford3 (i, 0, n)

#define ford(...) ARG4(__VA_ARGS__, ford3, ford2) (__VA_ARGS__)



#define ve vector

#define pa pair

#define tu tuple

#define mp make_pair

#define mt make_tuple

#define pb push_back

#define fs first

#define sc second

#define all(a) (a).begin(), (a).end()

#define sz(a) ((int)(a).size())



typedef long double ld;

typedef long long ll;

typedef unsigned long long ull;

typedef unsigned int ui;

typedef unsigned char uc;

typedef pa<int, int> pii;

typedef pa<int, ll> pil;

typedef pa<ll, int> pli;

typedef pa<ll, ll> pll;

typedef ve<int> vi;



const ld pi = 3.1415926535897932384626433832795l;



template<typename T> inline auto sqr (T x) -> decltype(x * x) {return x * x;}

template<typename T1, typename T2> inline bool umx (T1& a, T2 b) {if (a < b) {a = b; return 1;} return 0;}

template<typename T1, typename T2> inline bool umn (T1& a, T2 b) {if (b < a) {a = b; return 1;} return 0;}



const int N = 18;

double p[N][N];

bool use[N];

double d[2][N];

int n;



double calc(vi tek) {

	reverse(all(tek));

	memset(d, 0, sizeof(d));

	d[0][tek[0]] = 1.0;

	forn(i, sz(tek) - 1) {

		int i0 = (i & 1), i1 = (i0 ^ 1);

		int nw = tek[i + 1];

		double sum = 0;

		forn(j, n) {

			d[i1][j] = d[i0][j] * p[j][nw]; 	

			sum += d[i1][j];

		}

		d[i1][nw] = 1 - sum;

	}

	return d[(sz(tek) - 1) & 1][0];

}



int main () {

	cout.setf(ios::showpoint | ios::fixed);

	cout.precision(20);

#ifdef SG

	freopen("e.in", "r", stdin);

//	freopen((problemname + ".out").c_str(), "w", stdout);

#endif

	cin >> n;

	forn(i, n)

		forn(j, n)

			cin >> p[i][j];	

	vi tek;

	tek.pb(0);

	memset(use, 0, sizeof(use));

	forn(qqq, n - 1) {

		double maxx = -1;

		int maxi = -1;

		forn(j, 1, n) {

			if (use[j])

				continue;		

			tek.pb(j);

			if (umx(maxx, calc(tek))) {

				maxi = j;

			}

			tek.pop_back();

		}

		tek.pb(maxi);

		use[maxi] = 1;

	}

	cout << calc(tek) << endl;

	return 0;

}