#include <bits/stdc++.h> 
using namespace std;
#define ll long long
const int mod=1e9+7;
#define mymaxn 400400
int fpow(int a,int b)
{
	int s=1;if(a==1)return 1;
	while(b)
	{
		if(b&1)
		{
			s=1ll*s*a%mod;
		}
		a=1ll*a*a%mod;
		b>>=1;
	}
	return s;
}
ll qaqaq[mymaxn];
bool vis[mymaxn];int prime[mymaxn],cnt;
int n,Q,a[mymaxn];
int t[mymaxn<<2];
ll qwqwqwqwqwq[mymaxn<<2];
int tagt[mymaxn<<2];
ll pwpwpw[mymaxn<<2];
void Build(int now,int l,int r)
{
	t[now]=1;
	qwqwqwqwqwq[now]=0;
	tagt[now]=1;
	if(l==r)
	{
		t[now]=a[l];
		qwqwqwqwqwq[now]=qaqaq[a[l]];
		return;
	}
	int mid=(l+r)>>1;
	Build((now<<1),l,mid);
	Build((now<<1|1),mid+1,r);
	t[now]=1ll*t[(now<<1)]*t[(now<<1|1)]%mod;
	qwqwqwqwqwq[now]=qwqwqwqwqwq[(now<<1)]|qwqwqwqwqwq[(now<<1|1)];
}
void modify(int now,int l,int r,int L,int R,int w,ll W)
{
	if(L==l&&r==R)
	{
		tagt[now]=1ll*tagt[now]*w%mod;pwpwpw[now]|=W;
		return;
	}
	int mid=(l+r)>>1;
	if(L>mid)
	{
		modify((now<<1|1),mid+1,r,L,R,w,W);
	}
	else if(R<=mid)
	{
		modify((now<<1),l,mid,L,R,w,W);
	}
	else 
	{
		modify((now<<1),l,mid,L,mid,w,W);
		modify((now<<1|1),mid+1,r,mid+1,R,w,W);
	}
	qwqwqwqwqwq[now]|=W;
	t[now]=1ll*t[now]*fpow(w,R-L+1)%mod;
}
int queryqwq(int now,int l,int r,int L,int R)
{
	if(L==l&&r==R)
	{
		return 1ll*t[now]*fpow(tagt[now],r-l+1)%mod;
	}
	int mid=(l+r)>>1;
	if(R<=mid)
	{
		return 1ll*queryqwq((now<<1),l,mid,L,R)*fpow(tagt[now],R-L+1)%mod;
	}
	if(L>mid)
	{
		return 1ll*queryqwq((now<<1|1),mid+1,r,L,R)*fpow(tagt[now],R-L+1)%mod;
	}
	return 1ll*queryqwq((now<<1),l,mid,L,mid)*queryqwq((now<<1|1),mid+1,r,mid+1,R)%mod*fpow(tagt[now],R-L+1)%mod;
}
ll queryqaq(int now,int l,int r,int L,int R)
{
	if(L<=l&&r<=R)
	{
		return qwqwqwqwqwq[now]|pwpwpw[now];
	}
	ll ret=pwpwpw[now];
	int mid=(l+r)>>1;
	if(L<=mid)
	{
		ret|=queryqaq((now<<1),l,mid,L,R);
	}
	if(R>mid)
	{
		ret|=queryqaq((now<<1|1),mid+1,r,L,R);
	}
	return ret;
}
int inv[500];
int main()
{
	cin>>n>>Q;
	for(int i=2;i<=300;++i)
	{
		if(!vis[i])
		{
			prime[++cnt]=i;
			inv[cnt]=fpow(i,mod-2);
		}
		for(int j=i+i;j<=300;j+=i)
		{
			vis[j]=true;
		}
	}
	for(int i=2;i<=300;++i)
	{
		for(int j=1;j<=cnt;++j)
		{
			if(i%prime[j]==0)
			{
				qaqaq[i]|=1ll<<(j-1);
			}
		}
	}
		
	for(int i=1;i<=n;++i)
	{
		scanf("%d",&a[i]);
	}
	Build(1,1,n);
	char ch[50];
	while(Q--)
	{
		scanf("%s",ch);
		int l,r;
		scanf("%d%d",&l,&r);
		if(ch[0]=='M')
		{
			int x;
			scanf("%d",&x);
			modify(1,1,n,l,r,x,qaqaq[x]);
		}
		else
		{
			ll S=queryqaq(1,1,n,l,r);
			int v=queryqwq(1,1,n,l,r);
			for(int i=0;i<cnt;++i)
			{
				if(S&(1ll<<i))
				{
					v=1ll*v*inv[i+1]%mod*(prime[i+1]-1)%mod;
				}	
			}		
			printf("%d\n",v);
		}
	}
	return 0;
}