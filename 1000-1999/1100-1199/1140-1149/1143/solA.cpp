//
#include <bits/stdc++.h>
#define ll long long
#define pb push_back
#define mk make_pair
#define pii pair<int,int>
#define fst first
#define scd second
using namespace std;
inline ll read() {
	ll f=1,x=0;char ch=getchar();
	while(!isdigit(ch)){if(ch=='-')f=-1;ch=getchar();}
	while(isdigit(ch)){x=x*10+ch-'0';ch=getchar();}
	return x*f;
}
const int MAXN=2e5+5;
int n,a[MAXN],x[MAXN],y[MAXN];
int main() {
//	freopen("","r",stdin);
//	freopen("","w",stdout);
	ios::sync_with_stdio(0);cin.tie(0);/*syn加速*/
	n=read();
	for(int i=1;i<=n;++i) a[i]=read();
	for(int i=n;i>=1;--i) if(a[i]) x[i]=x[i+1]+1,y[i]=y[i+1]; else x[i]=x[i+1],y[i]=y[i+1]+1;
	for(int i=1;i<=n;++i) if(x[i]==0||y[i]==0) {
		cout<<i-1<<endl;return 0;
	}
	return 0;
}