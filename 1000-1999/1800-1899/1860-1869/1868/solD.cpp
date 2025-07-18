#include<bits/stdc++.h>
#define ll unsigned long long
#define db double
#define pii pair<int,int>
//#define int ll
#define mp make_pair
#define fi first
#define se second
using namespace std;
ll read(){
	char st=getchar();
	ll sum=0,f=1;
	while(st<'0'||st>'9'){
		if(st=='-')f=-1;
		st=getchar();
	}
	while('0'<=st&&st<='9'){
		sum=(sum<<3)+(sum<<1)+st-'0';
		st=getchar();
	}
	return sum*f;
}
const int maxn=1000010,maxm=20300,mod=1e9+7;

int n,m,c[maxn];
pii d[maxn];
vector<pii> ans;
void add(int x,int y){
	ans.push_back(mp(d[x].se,d[y].se));
	--d[x].fi,--d[y].fi;
}
int solve(){
	int n=read();
	int su=0;
	for(int i=1;i<=n;++i)
		su+=d[i].fi=c[i]=read(),d[i].se=i;
	if(su!=n+n)return 0;
	sort(d+1,d+1+n);
	if(d[1].fi==2){
//		puts("Yes");
		for(int i=1;i<=n;++i)
			ans.push_back(mp(i,i%n+1));
//			printf("%d %d\n",i,i%n+1);
		return 1;
	}
	if(d[2].fi!=1||d[n-1].fi<=2)return 0;
	int top=1,s=1;
	while(d[top].fi==1)++top;
	if(d[top].fi>2){
		add(top,n);
		for(int i=top;i<n;++i)add(i,i+1);
		for(int i=top;i<=n;++i)
			while(d[i].fi)add(s++,i);
		return 1;
	}
//	int A=d[top].se,B=d[top+1].se;
//	top+=2;
	while(1){
//		cout<<"top="<<top<<" d[top-2].se="<<d[top-2].se<<" "
		if(top+1==n){
			add(top-2,top),add(top-1,top+1);
			add(top,top+1);add(top,top+1);
			while(d[top].fi)add(s++,top);
			while(d[top+1].fi)add(s++,top+1);
			break;
		}
		if(top+2==n){
			if(c[d[top-1].se]==1)return 0;
//			cout<<"d"
			if(d[top+2].fi<=3)return 0;
			add(top,top+2);
			while(d[top].fi)add(s++,top);
			swap(d[top],d[top-2]);
			++top;
			continue;
		}
		if(c[d[top-1].se]!=1&&top+4==n&&d[top+2].fi>=3){
			add(top,top+2);
			while(d[top].fi)add(s++,top);
			add(top-2,top+1),add(top-1,top+2);
			while(d[top+1].fi>1)add(s++,top+1);
			while(d[top+2].fi>1)add(s++,top+2);
			top+=3;
			continue;
		}
		add(top-2,top),add(top-1,top+1);
		while(d[top].fi>1)add(s++,top);
		while(d[top+1].fi>1)add(s++,top+1);
		top+=2;
	}
	return 1;
}
signed main(){
//	freopen("1.in","r",stdin);
//	freopen("1.out","w",stdout);
	int TT=read();
	while(TT--){
		if(solve()){
			puts("Yes");
			for(auto x:ans)
				printf("%d %d\n",x.fi,x.se);
		}
		else puts("No");
		ans.clear();
		
	}
	return 0;
}