#include<bits/stdc++.h>

#define ll long long

#define gmax(x,y) x=max(x,y)

#define gmin(x,y) x=min(x,y)

#define F first

#define S second

#define P pair

#define FOR(i,a,b) for(int i=a;i<=b;i++)

#define rep(i,a,b) for(int i=a;i<b;i++)

#define V vector

#define RE return

#define ALL(a) a.begin(),a.end()

#define MP make_pair

#define PB emplace_back

#define PF push_front

#define FILL(a,b) memset(a,b,sizeof(a))

#define lwb lower_bound

#define upb upper_bound

#define sz(x) ((int)x.size())

#define pc putchar

using namespace std;

int lc[200005],rc[200005];

int n;

P<int,int> p[200005];

int ans[200005];

void dfs(int x){

	if(!lc[x]){

		p[x]=MP(1,0);RE ;

	}

	if(!rc[x]){

		dfs(lc[x]);

		if(p[lc[x]].S==0){

			ans[lc[x]]^=1;

			p[x]=MP(1-p[lc[x]].F,1);

		}else{

			if(p[lc[x]].F==2){

				ans[lc[lc[x]]]^=1;ans[lc[x]]^=1;

				p[x]=MP(1,1);

			}else{

				ans[lc[x]]^=1;

				p[x]=MP(1-p[lc[x]].F,1);

			}

		}

		RE ;

	}

	dfs(lc[x]);dfs(rc[x]);

	if(p[lc[x]].F==2){

		ans[lc[lc[x]]]^=1;

		p[lc[x]]=MP(0,0);

	}

	if(p[rc[x]].F==2){

		ans[lc[rc[x]]]^=1;

		p[rc[x]]=MP(0,0);

	}

	if(p[lc[x]].S==0&&p[rc[x]].S==0){

		ans[lc[x]]^=1;ans[rc[x]]^=1;

		p[x]=MP(1-p[lc[x]].F-p[rc[x]].F,1);

		RE;

	}

	rep(f1,0,2)rep(f2,0,2){

		if(p[lc[x]].S==0&&f1==0)continue;

		if(p[rc[x]].S==0&&f2==0)continue;

		int val=1;

		if(f1)val-=p[lc[x]].F;else val+=p[lc[x]].F;

		if(f2)val-=p[rc[x]].F;else val+=p[rc[x]].F;

		if((val==0||val==1)&&(f1||f2)){

			ans[lc[x]]^=f1;ans[rc[x]]^=f2;

			p[x]=MP(val,1);RE;

		}

	}

}

void solve(){

	cin>>n;

	FOR(i,1,n)lc[i]=rc[i]=0,ans[i]=0;

	FOR(i,2,n){

		int x;

		cin>>x;

		if(lc[x])rc[x]=i;else lc[x]=i;

		if(i==4&&lc[2]&&rc[2]){

			cout<<2<<'\n';

		}else cout<<(i&1)<<'\n';

	}

	dfs(1);

	FOR(i,1,n){

		if(lc[i])ans[lc[i]]^=ans[i];

		if(rc[i])ans[rc[i]]^=ans[i];

		if(ans[i])cout<<'w';else cout<<'b';

	}

	cout<<'\n';

}

signed main(){

	ios::sync_with_stdio(0);

	cin.tie(0);

	int T;

	cin>>T;

	while(T--)solve();

	RE 0;

}