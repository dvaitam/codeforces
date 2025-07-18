#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
#define pii pair<int,int>
#define pll pair<ll,ll>
#define pdd pair<double,double>
#define X first
#define Y second
#define REP(i,a) for(int i=0;i<a;++i)
#define REPP(i,a,b) for(int i=a;i<b;++i)
#define FILL(a,x) memset(a,x,sizeof(a))
#define	foreach( gg,itit )	for( typeof(gg.begin()) itit=gg.begin();itit!=gg.end();itit++ )
#define	mp make_pair
#define	pb push_back
#define all(s) s.begin(),s.end()
#define present(c,x) ((c).find(x) != (c).end())
const double EPS = 1e-8;
const int mod = 1e9+7;
const int N = 1e6+10;
const ll INF = 1e18;

//#define DEBUG
ll power(ll x,ll y){
  ll t=1;
  while(y>0){
    if(y%2) y-=1,t=t*x%mod;
    else y/=2,x=x*x%mod;
  }
  return t;
}
#ifdef DEBUG
#define dprintf(fmt,...) fprintf(stderr,fmt,__VA_ARGS__)
#else
#define dprintf(fmt,...)
#endif

#define ld long double
int arr[N];
int ans[N],c[N];
int main(){
	int n,k;
	scanf("%d%d",&n,&k);
	REP(i,n){
		scanf("%d",&arr[i]);
	}
	REP(i,n){
		 scanf("%d",&c[i]);
	}
	for(int i=n-1;i>=0;i--){
		if(100*arr[i]>c[i]){
			arr[i]=c[i];
		}else{
			arr[i]=100*arr[i];
			if(arr[i+1]>arr[i]) arr[i]=min(c[i],arr[i+1]);
		} 
	}	sort(arr,arr+n);
	ld ans=arr[0],tot=1;

	REP(i,n-k){
		tot*=(n-k-i); tot/=(n-i);
		if(tot<1e-18) break;
		ans+=tot*(arr[i+1]-arr[i]); 
	}
	printf("%lf\n",(double)ans);
  return 0;
}