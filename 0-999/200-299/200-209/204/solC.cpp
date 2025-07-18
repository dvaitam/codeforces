#include <bits/stdc++.h>
using namespace std;
//typedefs
typedef long long ll;
typedef unsigned long long ull;
typedef pair <int, int> pii;
typedef pair <int, pii> piii;
typedef vector <int> vi;
typedef vector <ll> vl;
typedef pair <ll, ll> pll;
const double PI = acos(-1);
//defines
#define MP make_pair
#define PB push_back
#define F first
#define S second
#define mem(a, b) memset(a, b, sizeof(a))
#define gcd(a,b) __gcd(a,b)
#define lcm(a,b) (a*(b/gcd(a,b)))
#define sqr(a) ((a)*(a))
#define inf 100000000
#define mod 1000000007
#define mod1 1000000007
#define mod2 1000000009
#define b1 43
#define b2 41
#define EPS 1e-9
//define harmonic(n) 0.57721566490153286l+log(n)
#define nl puts("")
#define odd(n) ((n)&1)
#define even(n) (!((n)&1))
#define vsort(v) sort(v.begin(), v.end())
#define lc (node<<1)
#define rc ((node<<1)|1)
//loop
#define rep(i, n) for(int i = 0; i < n; ++i)
#define REP(i, n) for(int i = 1; i <= n; ++i)
//input
#define si(a) scanf("%d", &a)
#define sii(a, b) scanf("%d%d", &a, &b)
#define siii(a, b, c) scanf("%d%d%d", &a, &b, &c)
#define sl(a) scanf("%lld", &a)
#define sll(a, b) scanf("%lld%lld", &a, &b)
#define slll(a, b, c) scanf("%lld%lld%lld", &a, &b, &c)
#define sd(a) scanf("%lf", &a)
#define sc(a) scanf("%c", &a)
#define sst(a) scanf("%s", a)
inline bool EQ(double a, double b) { return fabs(a-b) < 1e-9; }
//debug
#ifdef tahsin
template < typename F, typename S >
ostream& operator << ( ostream& os, const pair< F, S > & p ) {
	return os << "(" << p.first << ", " << p.second << ")";
}

template < typename T >
ostream &operator << ( ostream & os, const vector< T > &v ) {
	os << "{";
	for(auto it = v.begin(); it != v.end(); ++it) {
		if( it != v.begin() ) os << ", ";
		os << *it;
	}
    return os << "}";
}

template < typename T >
ostream &operator << ( ostream & os, const set< T > &v ) {
	os << "[";
	for(auto it = v.begin(); it != v.end(); ++it) {
		if( it != v.begin()) os << ", ";
		os << *it;
	}
    return os << "]";
}

template < typename F, typename S >
ostream &operator << ( ostream & os, const map< F, S > &v ) {
	os << "[";
	for(auto it = v.begin(); it != v.end(); ++it) {
		if( it != v.begin() ) os << ", ";
		os << it -> first << " = " << it -> second ;
	}
    return os << "]";
}

#define dbg(args...) do {cerr << #args << " : "; faltu(args); } while(0)

clock_t tStart = clock();
#define timeStamp dbg("Execution Time: ", (double)(clock() - tStart)/CLOCKS_PER_SEC)

void faltu () { cerr << endl; }

template <typename T>
void faltu( T a[], int n ) {
	for(int i = 0; i < n; ++i) cerr << a[i] << ' ';
	cerr << endl;
}

template <typename T, typename ... hello>
void faltu( T arg, const hello &... rest) { cerr << arg << ' '; faltu(rest...); }

#else
#define dbg(args...)
#endif
ll add(ll a, ll b) {
	ll ret = a+b;
	if(ret >= mod) ret -= mod;
	return ret;
}

ll subtract(ll a, ll b) {
	ll ret = a-b;
	if(ret < 0) ret += mod;
	return ret;
}

ll mult(ll a, ll b) {
	ll ret = a*b;
	if(ret >= mod) ret %= mod;
	return ret;
}

ll bigmod(ll a, ll b) {
	ll ret = 1;
	while(b) {
		if(b&1) ret = mult(ret, a);
		b >>= 1; a = mult(a, a);
	}
	return ret;
}

ll inverse(ll n) { return bigmod(n, mod-2); }

bool base[1000010];
vi primes;
void sieve(int mx) {
	mx += 10;
	int x = sqrt(mx);
	for(int i = 3; i <= x; i += 2) if(base[i] == 0) for(int j = i*i, k = i<<1; j < mx; j += k) base[j] = 1;
	primes.PB(2);
	for(int i = 3; i < mx; i += 2) if(base[i] == 0) primes.PB(i);
}


//Direction Array 
//int fx[]={1, -1, 0, 0}; int fy[]={0, 0, 1, -1};
//int fx[]={0, 0, 1, -1, -1, 1, -1, 1}; int fy[]={-1, 1, 0, 0, 1, 1, -1, -1};

//bit manipulation
bool checkBit(int n, int i) { return (n&(1<<i)); }
int setBit(int n, int i) { return (n|(1<<i)); }
int resetBit(int n, int i) { return (n&(~(1<<i))); }
//end of template

#define MX 200010
char a[MX], b[MX];
vi v1[26], v2[26];
vector <double> sum1[26], sum2[26];

int main () {
#ifdef tahsin
//	freopen("in", "r", stdin);
//	freopen("out", "w", stdout);
#endif
	int n;

	si(n);
	sst(a);
	sst(b);

	rep(i, n) {
		v1[a[i]-'A'].PB(i);
		v2[b[i]-'A'].PB(i);
	}

	double res = 0;
	rep(i, 26) {
		int sz = v2[i].size();
		sum1[i].resize(sz+1);
		sum2[i].resize(sz+1);

		rep(j, sz) {
			if(j == 0) sum1[i][0] = v2[i][0] + 1;
			else sum1[i][j] = sum1[i][j-1] + v2[i][j] + 1;
		}

		for(int j = sz-1; j >= 0; --j) {
			if(j == sz-1) sum2[i][j] = n - v2[i][j];
			else sum2[i][j] = sum2[i][j+1] + n - v2[i][j];
		}

		sz = v1[i].size();

		rep(j, sz) {
			int idx = lower_bound(v2[i].begin(), v2[i].end(), v1[i][j]) - v2[i].begin();
			if(idx) res += 1.0 * (n - v1[i][j]) * sum1[i][idx-1];
			res += 1.0 * (v1[i][j]+1) * sum2[i][idx];
		}
	}

	double div = 0;
	REP(i, n) div += 1.0 * i * i;

	printf("%.10lf\n", 1.0 * res / div);

//	timeStamp;
	return 0;
}