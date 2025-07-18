#include<bits/stdc++.h>
#define pb push_back
#define mp make_pair
#define fi first
#define se second
#define SZ(x) ((int)x.size())
#define FOR(i,a,b) for (int i=a;i<=b;++i)
#define FORD(i,a,b) for (int i=a;i>=b;--i)
using namespace std;
typedef long long LL;
typedef pair<int,int> pa;
typedef vector<int> vec;
void getint(int &v){
    char ch,fu=0;
    for(ch='*'; (ch<'0'||ch>'9')&&ch!='-'; ch=getchar());
    if(ch=='-') fu=1, ch=getchar();
    for(v=0; ch>='0'&&ch<='9'; ch=getchar()) v=v*10+ch-'0';
    if(fu) v=-v;
}
pa mx1,mx2;
int n,x,ans;
int main(){
	cin>>n;
	mx1=mp(0,0);
	mx2=mp(0,0);
	FOR(i,1,n){
		getint(x);
		if (mx1.se && mx1.se!=x){
			ans=max(ans,i-mx1.fi);
		}
		if (mx2.se && mx2.se!=x){
			ans=max(ans,i-mx2.fi);
		}
		if (mx1.se==0){
			mx1=mp(i,x);
			continue;
		}
		if (mx2.se==0 && x!=mx1.se){
			mx2=mp(i,x);
		}
	}
	cout<<ans<<endl;
	return 0;
}