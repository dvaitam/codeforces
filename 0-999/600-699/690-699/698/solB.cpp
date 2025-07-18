#include <stdio.h>
#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
#define pii pair<int,int>
#define pll pair<ll,ll>
#define pdd pair<double,double>
#define FILL(a,x) memset(a,x,sizeof(a))
#define foreach( gg,ii ) for( typeof(gg.begin()) ii=gg.begin();ii!=gg.end();ii++)
#define mp make_pair
#define pb push_back
#define X first
#define Y second
#define sz(a) int((a).size())
#define N 1000010
#define MAX 30
#define mod 1000000007
#define REP(i,a) for(int i=0;i<a;++i)
#define REPP(i,a,b) for(int i=a;i<b;++i)
#define all(a) a.begin(),a.end()
const ll INF = 1e18+1;

inline ll input(void)
{
	char t;
	ll x=0;
	int neg=0;
	t=getchar();
	while((t<48 || t>57) && t!='-')
		t=getchar();
	if(t=='-')
		{neg=1; t=getchar();}
	while(t>=48 && t<=57)
	{
		x=(x<<3)+(x<<1)+t-48;
		t=getchar();
	}
	if (neg) x=-x;
	return x;
}

inline void output(ll x)
{
	char a[20];
	int i=0,j;
	a[0]='0';
	if (x<0) {putchar('-'); x=-x;}
	if (x==0) putchar('0');
	while(x)
	{
		a[i++]=x%10+48;
		x/=10;
	}
	for(j=i-1;j>=0;j--)
	{
		putchar(a[j]);
	}
	putchar(' ');
}

ll a[N];

int p[N],ranker[N];

inline void create(int x)
{
	p[x]=x;
	ranker[x]=1;
}

inline int findset(int x)
{
	if (x != p[x]) p[x] = findset(p[x]);
	return p[x];
}

inline void merge(int x,int y)
{
	int px=findset(x),py=findset(y);
	if(px==py) return;
	if(ranker[px]>ranker[py])
	{
		p[py]=px;
		ranker[px]=ranker[px]+ranker[py];
	}
	else
	{
		p[px]=py;
		ranker[py]=ranker[py]+ranker[px];
	}
}

vector <int> selfloops;

int main()
{
	ll n,r=-1,x,y,ans=0;
	n=input();
	REP(i,n)
		a[i]=input()-1;
	REP(i,n) create(i);
	REP(i,n)
	{
		if (a[i]==i&&r==-1) r=i;
		else
		{
			x=findset(i);
			y=findset(a[i]);
			if (x==y) selfloops.pb(i);
			else merge(x,y);
		}
	}
	ans=sz(selfloops);
	if (r==-1)
	{
		r=selfloops.back();
		selfloops.pop_back();
	}
	a[r]=r;
	REP(i,sz(selfloops))
		a[selfloops[i]]=r;
	output(ans);
	putchar('\n');
	REP(i,n) output(a[i]+1);
	putchar('\n');
return 0;
}