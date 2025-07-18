/*

    Author : SharpnessV & SharpnessX & 7KByte

    Right Output ! & Accepted !

*/

#include<bits/stdc++.h>

//#include<atcoder/all>

//#define int long long



#define rep(i, a, b) for(int i = (a);i <= (b);i++)

#define pre(i, a, b) for(int i = (a);i >= (b);i--)

#define rp(i, a) for(int i = 1; i <= (a); i++)

#define pr(i, a) for(int i = (a); i >= 1; i--)

#define go(i, x) for(auto i : x)



#define mp make_pair

#define pb emplace_back

#define pf push_front

#define fi first

#define se second

#define si(x) (int)(x).size()

#define pc putchar

#define gc getchar

#define el putchar('\n')



using namespace std;

const double eps = 1e-15, pi = 3.1415926535897932385;

typedef long long LL;

typedef pair<int,int> Pr;

const int dx[4] = {1, 0, -1, 0}, dy[4] = {0, 1, 0, -1}, inf = 0x7fffffff, inf_ = 0x80000000;

const LL Inf = 0x7fffffffffffffffLL, Inf_ = 0x8000000000000000LL;

//char buf[1<<22],*p1=buf,*p2=buf;

//#define getchar() (p1==p2&&(p2=(p1=buf)+fread(buf,1,1<<21,stdin),p1==p2)?EOF:*p1++)

template <typename T> inline void read(T &x) {

    x = 0;bool flag = false; char ch = getchar();

    while (ch < '0' || ch > '9')flag = ch == '-' ? true : false, ch = getchar();

    while (ch >= '0' && ch <= '9')x = (x << 3) + (x << 1) + (ch & 15), ch = getchar();

    if(flag) x = -x;

}

template <typename T,typename... Args> inline void read(T &t,Args&... args){read(t);read(args...);}

int gcd(int x,int y) {	return y ? gcd(y, x % y) : x;}

int lcm(int x,int y) {	return x / gcd(x, y) * y;}

#define P 1000000007

//#define P 998244353

#define bas 229

template<typename T> void cmx(T &x, T y){if(y > x) x = y;}

template<typename T> void cmn(T &x, T y){if(y < x) x = y;}

template<typename T> void ad(T &x, T y) {x += y; if(x >= P) x -= P;}

template<typename T> void su(T &x, T y) {x -= y; if(x < 0) x += P;}

int Pow(int x, int y){

	int now = 1 ;

	for(; y; y >>= 1, x = 1LL * x * x % P)if(y & 1) now = 1LL * now * x % P;

	return now;

}



/***************************************************************************************************************************/

/*                                                                                                                         */

/***************************************************************************************************************************/

#define N 600005

int n, m, idx, u[N], v[N], w[N], vs[N], c[N], l[N], r[N], f[N][3], in[N], ed[N];

queue<int>q; vector<Pr>e[N];

void dfs(int x){

	c[x] = 1;

	go(y, e[x])if(!c[y.fi]){

		ed[y.se] = u[y.se] == x ? 1 : 2;

		//cout << "ss " << y.se << " " << ed[y.se] << endl;

		dfs(y.fi);

	}

}

void dfs(int x,int pv){

	if(c[x])return;

	c[x] = 1;

	go(y, e[x])if(y.se != pv){

		ed[y.se] = u[y.se] == x ? 1 : 2;

		//cout << "tt " << y.se << " " << ed[y.se] << endl;

		dfs(y.fi, y.se); break;

	}

}

void calc(int x){

	if(!l[x])return;

	if(ed[x] == 2)swap(u[x], v[x]), swap(l[x], r[x]);

	//cout << "ff " << x << " " << u[x] << " " << v[x] << " " << w[x] << " " << l[x] << " " << r[x] << endl;

	//cout << "pp " << u[l[x]] << " " << v[r[x]] << endl;

	if(u[l[x]] == u[x])ed[l[x]] = 1; else ed[l[x]] = 2;

	if(v[r[x]] == v[x])ed[r[x]] = 1; else ed[r[x]] = 2;

	calc(l[x]), calc(r[x]);

}

int main(){

	read(n, m), idx = m;

	rp(i, m)read(u[i], v[i], w[i]), q.push(i);

	while(!q.empty()){

		int x = q.front(), y = 0; q.pop();

		if(f[u[x]][w[x]])y = f[u[x]][w[x]];

		else if(f[v[x]][w[x]])y = f[v[x]][w[x]];

		if(y){

			f[u[y]][w[y]] = f[v[y]][w[y]] = 0;

			int z = ++idx; l[z] = x, r[z] = y;

			if(u[x] == u[y])u[z] = v[x], v[z] = v[y];

			else if(u[x] == v[y])u[z] = v[x], v[z] = u[y];

			else if(v[x] == u[y])u[z] = u[x], v[z] = v[y];

			else u[z] = u[x], v[z] = u[y];

			w[z] = w[x], vs[x] = vs[y] = 1, q.push(z);

			//cout << "ss " << x << " " << y << " " << z << " " << u[z] << " " << v[z] << " " << w[z] << endl;

		}

		else f[u[x]][w[x]] = f[v[x]][w[x]] = x;

	}

	rp(i, idx)if(!vs[i])e[u[i]].pb(mp(v[i], i)), e[v[i]].pb(mp(u[i], i)), in[u[i]]++, in[v[i]]++;

	rp(i, n)if(!c[i] && in[i] <= 1)dfs(i);

	rp(i, n)if(!c[i])dfs(i, 0);

	rp(i, idx)if(!vs[i])calc(i);

	memset(c, 0, sizeof(c));

	rp(i, m)if(ed[i] == 1)c[u[i]] += w[i], c[v[i]] -= w[i]; else c[v[i]] += w[i], c[u[i]] -= w[i];

	int ans = 0 ;

	rp(i, n)ans += c[i] == 1 || c[i] == -1;

	printf("%d\n", ans);

	rp(i, m)printf("%d", ed[i]);

	return 0;

}