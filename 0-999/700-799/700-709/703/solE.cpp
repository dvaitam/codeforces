#include <iostream>

#include <cstdio>

#include <ctime>

#include <vector>

#include <algorithm>

#include <cmath>

#include <cstring>

#include <algorithm>



#define rep(i,a,b) for(int i=(a),_b=(b);i<_b;i++)

#define repn(i,a,b) for(int i=(a),_b=(b);i<=_b;i++)

#define pern(i,b,a) for(int i=(b),_a=(a);i>=_a;i--)

#define fi first

#define se second

#define pb push_back

#define mp make_pair

#define sz(x) int((x).size())



using namespace std;



typedef long long ll;

typedef unsigned long long ull;

typedef pair<int, int> pii;

typedef vector<int> vi;

typedef vector<pii> vii;



template <typename T> inline void read(T &x){

	x = 0; char c; int positive = 1;

	while (!isdigit(c = getchar())) if (c == '-') positive = 0;

	do x = x * 10 + c - '0'; while (isdigit(c = getchar()));

	if (!positive) x = -x;

}



template <typename T> inline void write(T x){

	if (x > 9) write(x / 10);

	putchar(x % 10 + 48);

}



const int maxn = 1010;



int f[1<<20];

ll g[1<<20];

ll a[maxn];

ll k;

int n;

vi cnt;

vector<ll> pr, List[1<<20];

int represent[1<<20][17];

int current[17];



void input(){

	read(n); read(k);

	rep(i,0,n) read(a[i]);

} 



void init(){

}



void solve(){

	if (k == 1){

		f[0] = 1;

		int x = n;

		a[n] = 1e12 + 1;

		rep(i,0,n) if (a[x] > a[i]) x = i;

		List[0].pb(x);

		return;

	}

	

	repn(i,2,1<<20)

		if (k % i == 0){

			int x = 0;

			while (k % i == 0) k /= i, x++;

			cnt.pb(x);

			pr.pb(i);

		}		

	if (k > 1) cnt.pb(k), pr.pb(1);

	

	fill(f, f + (1 << 20), n + 1);

	

	int maxstate = 1;

	rep(i,0,sz(pr)) maxstate *= (cnt[i] + 1);

	rep(i,0,sz(pr)) represent[maxstate - 1][i] = cnt[i];

	

	f[maxstate-1] = 0;

	rep(i,0,n){

		ll x = a[i];

		rep(j,0,sz(pr)){

			current[j] = 0;

			while (x % pr[j] == 0)

				x /= pr[j], current[j]++;

		}

		

		rep(curr,0,maxstate) if (f[curr] != n + 1){

			int neXt = 0;

			rep(j,0,sz(pr))

				neXt = neXt * (cnt[j] + 1) + max(0, represent[curr][j] - current[j]);

			if (f[neXt] > f[curr] + 1 || (f[neXt] == f[curr] + 1 && g[neXt] > g[curr] + a[i])){

				f[neXt] = f[curr] + 1;

				g[neXt] = g[curr] + a[i];

				rep(j,0,sz(pr))

					represent[neXt][j] = max(0, represent[curr][j] - current[j]);

				List[neXt] = List[curr];

				List[neXt].pb(i);

			}

		}

	}

}



void output(){

	if (f[0] == n + 1) printf("-1");

	else{

		write(f[0]);

		putchar('\n');

		rep(i,0,sz(List[0]))

			write(List[0][i]+1), putchar(' ');

	}

}



int main(){

	clock_t timestart, timeend;

	timestart = clock();

//	freopen("input.inp","r",stdin);

//	freopen("output.out","w",stdout);

	input();

	init();

	solve();

	output();

	timeend = clock();

	//system("pause");

	//printf("\nRun time : %lf",(timeend-timestart)*1.0/CLK_TCK);

	return 0;

}