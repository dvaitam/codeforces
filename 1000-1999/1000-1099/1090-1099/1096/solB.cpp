#include<bits/stdc++.h>
#define rep(i, l, r) for(int i = (l); i <= (r); i++)
#define per(i, r, l) for(int i = (r); i >= (l); i--)
#define foreach(i, x) for(auto i = x.begin(); i != x.end(); i++)
#define forg(i) for(int i = head[u]; i; i = g[i].nxt)
#define INF (ll)0x3f3f3f3f
#define EPS 1e-8
#define X first
#define Y second
#define mkp(x, y) make_pair(x, y)
#define pb(x) push_back(x)
#define lbound lower_bound
#define ubound upper_bound
#define rnd(x) rand()%(x)
#define squ(x) (x)*(x)
#define disp0(A){foreach(i,A)cout<<A[i]<<" ";cout<endl;}
#define disp(A, l, r) {rep(i,l,r)cout<<A[i]<<" ";cout<<endl;}
#define disp2(A, l, r, b, e){									\
		rep(i,l,r){rep(j,b,e)cout<<A[i][j]<<"\t";cout<<endl;}	\
		cout<<endl;												\
	}
using namespace std;
typedef long long ll;
typedef unsigned long long ull;
typedef double db;
typedef long double ld;
typedef pair<int, int> P;
typedef pair<ll, int> Pli;
typedef pair<int, ll> Pil;
typedef pair<ll, ll> Pll;
template<class T> inline T lowbit(T x) {return x&(-x);}
template<class T> T gcd(T a, T b) {return b?gcd(b,a%b):a;}
inline ll read()
{
	ll x=0,f=1;char ch=getchar();
	while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}
	while(ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getchar();}
	return x*f;
}

const ll Z = 998244353; const int N = 2e5+10;
char s[N];
int main()
{
	//freopen("std.in", "r", stdin);
	//freopen("std.out", "w", stdout);
	int n = read();
	scanf("%s", s);
	ll l = 1, r = 1;
	rep(i, 0, n-1) if(s[i]==s[0]) l++; else break;
	per(i, n-1, 0) if(s[i]==s[n-1]) r++; else break;

	if(s[0] == s[n-1]) printf("%lld", l*r%Z);
	else printf("%lld", (l+r-1)%Z);
		
}