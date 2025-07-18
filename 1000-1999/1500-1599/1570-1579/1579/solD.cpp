/*input
1
3
1 2 2
*/
/*inhuman strength
#pragma GCC optimize("O2,unroll-loops")
#pragma GCC target("avx,avx2")
*/
#include<bits/stdc++.h>
#define foR(i,a,b) for(int i=a;i<=b;i++)
#define roF(i,a,b) for(int i=a;i>=b;i--)
#define _ff exit(0);
#define Sort(i,a,b) sort(i+a,i+a+b)
#define SORT(i,a,b) sort(i+a,i+a+b,greater<int>())
#define code_chef for(int i=1,x=0;i<=100000;i++){x++;cerr<<x;}
#define TLE cerr<<"\n" <<  1.0*clock()/CLOCKS_PER_SEC << "\n";
#define F first
#define S second
#define sperm(a) next_permutation(a.begin(),a.end());
#define all(a) a.begin(),a.end()
#define nxhieu_speed cin.tie(0)->sync_with_stdio(false);
#define i2 pair<int,int>
#define getbit(x,i) ((x>>(i))&1)
#define batbit(x,i) (x|(1ll<<(i)))
#define tatbit(x,i) (x&~(1<<(i)))
#define RTE cout<<-1;_ff
using namespace std;
mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());
int rnd(int l, int r) {return l + rng() % (r - l + 1);}
int dx[] = {0, 0, 1, -1};
int dy[] = {1, -1, 0, 0};
struct Triple
{
	int Fi;
	int Se;
	int Th;
};
bool operator >(Triple A, Triple B)
{
	if (A.Fi == B.Fi && A.Se == B.Se) return (A.Th > B.Th);
	if (A.Fi == B.Fi && A.Se != B.Se) return (A.Se > B.Se);
	return (A.Fi > B.Fi);
}
bool operator <(Triple A, Triple B)
{
	if (A.Fi == B.Fi && A.Se == B.Se) return (A.Th < B.Th);
	if (A.Fi == B.Fi && A.Se != B.Se) return (A.Se < B.Se);
	return (A.Fi < B.Fi);
}
bool operator ==(Triple A, Triple B)
{
	return (A.Fi == B.Fi && A.Se == B.Se && A.Th == B.Th);
}
const int N = 2e5 + 1e2;
int a[N], b[N], pref[N], n;
void solve()
{
	int cnt = 0;
	cin >> n;
	for (int i = 1; i <= n; i++)
	{
		cin >> a[i];
		while (a[i]--) b[++cnt] = i;
	}
	int START = 1;
	int X = (cnt) / 2 + 1;
	while (b[X] == b[1]) X++;
	vector<i2> GAY;
	int ST_X = 0;
	while (X <= cnt&&START!=ST_X)
	{
		while (b[START] == b[X] && X < cnt) X++;
		if (b[START] != b[X])
		{
			GAY.push_back(make_pair(b[START], b[X]));
			if (!ST_X) ST_X = X;
		}
		START++; ++X;
	}
	cout << GAY.size() << "\n";
	for (auto h : GAY) cout << h.F << " " << h.S << "\n";
}
int main()
{
	nxhieu_speed
	int t;
	cin >> t;
	while (t-- > 0) solve();
	return 0;
}