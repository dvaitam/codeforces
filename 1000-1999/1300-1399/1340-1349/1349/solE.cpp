#include<bits/stdc++.h>

#define R(i,a,b) for(int i=(a),i##E=(b);i<=i##E;i++)

#define L(i,a,b) for(int i=(b),i##E=(a);i>=i##E;i--)

using namespace std;

typedef long long ll;

int n,m;

ll t[222222];

pair<int,int>rg[222222];

int dp[222222][2],to[222222][2];

int res[222222];



inline ll calc1(ll x)

{

	return 1ll*x*(x-1)/2;

}

inline ll clac2(ll x,ll y)

{

	return 1ll*x*y+calc1(y);

}

inline int check(int L,int R,ll s)

{

	if(s<0) return 0;

	int l=0,r=R-L+1,ans=0;

	while(r-l>5)

	{

		int mid=(l+r)>>1;

		if(clac2(L,mid)<=s) l=mid,ans=mid;

		else r=mid;

	}

	L(i,l,r) if(clac2(L,i)<=s) {ans=i;break;}

	//cout<<"ans:"<<L<<" "<<R<<" "<<ans<<endl;

	ll sl=clac2(L,ans),sr=clac2(R-ans+1,ans);

	return sl<=s&&s<=sr;

}

void dfs(int l,int r,ll s)

{

	if(l>r) return;

	if(check(l+1,r,s)) dfs(l+1,r,s);

	else 

	{

		res[l]=1;

		dfs(l+1,r,s-l);

	}

}

signed main()

{

	ios::sync_with_stdio(false);

	cin.tie(NULL);

	cin>>n;

	R(i,1,n) cin>>t[i];

	reverse(t+1,t+n+1);

	L(i,1,n) if(t[i])

	{

		if(m&&t[i]<=t[rg[m].second]) rg[m].first=i;

		else rg[++m]={i,i};

	}

	if(!m)

	{

		R(i,1,n) cout<<0;

		cout<<endl;

		exit(0);

	}

	reverse(rg+1,rg+m+1);

	int cnt=0,cntt=0;

	R(i,rg[1].first,rg[1].second) if(t[i])

	{

		if(t[i]==t[rg[1].second]) ++cnt;

		else ++cntt;

	}

	dp[m+1][1]=n+1;

	L(i,2,m) R(j1,0,1) if(dp[i+1][j1]) R(j2,0,1) if(j2||t[rg[i].first]==t[rg[i].second])

	{

		ll lst=t[rg[i+1].first]+j1-1;

		for(int k=rg[i].first;k>rg[i-1].second&&k>dp[i][j2];k--) if(j2||k!=rg[i].first)

		{

			if(j2==1&&k<rg[i].first) break;

			if(check(rg[i].second+1,dp[i+1][j1]-1,t[rg[i].first]+j2-1-k-lst)) 

				{dp[i][j2]=k,to[i][j2]=j1;}

		}

	}

	R(i,0,1) if(dp[2][i])

	{

		ll lst=t[rg[2].first]+i-1;

		R(j,0,rg[1].second) 

		{

			if(cnt!=1&&t[j]==t[rg[1].second]) continue;

			if(cntt!=0&&t[j]!=t[rg[1].second]-1) continue;

			ll now=t[j]?t[j]:(t[rg[1].second]-1);

			//cout<<"lst&now"<<lst<<" "<<now<<" "<<i<<" "<<j<<endl;

			if(check(rg[1].second+1,dp[2][i]-1,now-lst-j))

			{

				res[j]=1;

				dfs(rg[1].second+1,dp[2][i]-1,now-lst-j);

				int cur=i;

				R(k,2,m)

				{

					int nxt=to[k][cur];

					ll lstt=t[rg[k+1].first]+nxt-1,noww=t[rg[k].first]+cur-1;

					res[dp[k][cur]]=1;

					dfs(rg[k].second+1,dp[k+1][nxt]-1,noww-lstt-dp[k][cur]);

					cur=nxt;

				}

				L(i,1,n) cout<<res[i];

				cout<<endl;

				exit(0);

			}

		}

	}

	cout<<dp[2][0]<<" "<<dp[2][1]<<endl;

}