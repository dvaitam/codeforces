#include<bits/stdc++.h>

using namespace std;

#define ll long long
#define pb push_back
#define rep(i,a,b) for(int (i) = (a); (i) < (b); (i)++)
#define all(v) (v).begin(),(v).end()
#define S(x) scanf("%d",&(x))
#define S2(x,y) scanf("%d%d",&(x),&(y))
#define SL(x) scanf("%lld",&(x))
#define SL2(x) scanf("%lld%lld",&(x),&(y))
#define P(x) printf("%d\n",(x))
#define FF first
#define SS second

int main(){
	//freopen("in.txt","r",stdin);
	int tc;
	S(tc);
	while(tc--){
		int L,v,l,r;
		S2(L,v);
		S2(l,r);
		int cnt = 0;
		//cout<<L<<v<<l<<r<<endl;
		/*for(int i=v; i<=L; i+=v){
			if(i>=l and i<=r)
				i = r- (r%v);
			else
				cnt++;
		}*/
		//int cnt1 = max((int)round((L/v)) +1,0);
		//int cnt2 = max((int)round((r/v)) - (int)round((l/v)) +1,0);
		//cout<<cnt1<<" "<<cnt2<<endl;
		int cnt1 = L/v;
		int cnt2 = floor(r/v)-floor((l-1)/v);
		printf("%d\n",cnt1-cnt2 );
	}
 
return 0;
}