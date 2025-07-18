#include<stdio.h>

#include<bits/stdc++.h>

#define fir first

#define sec second

#define all(x) begin(x),end(x)

using namespace std;

typedef long long ll;

typedef unsigned uint;

typedef unsigned long long ull;

typedef double db;

typedef long double ldb;

typedef __int128 int128;

typedef __uint128_t uint128;

typedef pair<int,int> pii;

template<typename type>

void chmin(type &x,const type &y)

{

	if(x>y)

		x=y;

}

template<typename type>

void chmax(type &x,const type &y)

{

	if(x<y)

		x=y;

}

constexpr int Max=(1<<17)+10,Size=3e6;

int a[Max],n,l,r,ch[Size][2],tot=1;

void insert(const int &x)

{

	int p=1;

	for(int i=16;i>=0;--i)

	{

		const int k=x>>i&1;

		if(!ch[p][k])

			ch[p][k]=++tot;

		p=ch[p][k];

	}

}

int query1(const int &x)

{

	int p=1,ans=0;

	for(int i=16;i>=0;--i)

	{

		const int k=x>>i&1;

		if(ch[p][k])

			p=ch[p][k];

		else

			p=ch[p][!k],ans|=1<<i;

	}

	return ans;

}

int query2(const int &x)

{

	int p=1,ans=0;

	for(int i=16;i>=0;--i)

	{

		const int k=x>>i&1;

		if(ch[p][!k])

			p=ch[p][!k],ans|=1<<i;

		else

			p=ch[p][k];

	}

	return ans;

}

void solve()

{

	cin>>l>>r,n=r-l+1;

	for(int i=1;i<=n;++i)

		cin>>a[i],insert(a[i]);

	for(int i=1;i<=n;++i)

	{

		const int x=a[i]^l;

		if(query1(x)==l&&query2(x)==r)

		{

			cout<<x<<"\n";

			return;

		}

	}

}

void clear()

{

	for(int i=1;i<=tot;++i)

		ch[i][0]=ch[i][1]=0;

	tot=1;

}

signed main()

{

	ios::sync_with_stdio(false);

	cin.tie(nullptr),cout.tie(nullptr);

	int t;

	cin>>t;

	while(t--)

		solve(),clear();

	return 0;

}