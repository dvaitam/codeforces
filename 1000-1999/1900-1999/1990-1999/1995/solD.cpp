#include<bits/stdc++.h>
using namespace std;
int n,c,k,ans,i,T,sum[18][1<<19],j,t;
bool tag[1<<18];
string s;
int main(){
	cin>>T; 
	for (;T--;){
		cin>>n>>c>>k>>s;
		for (i=0;i<1<<c;i++) tag[i]=0;
		ans=1e9;
		for (j=0;j<c;j++)
			for (i=1;i<=n;i++) sum[j][i]=sum[j][i-1]+(s[i-1]=='A'+j);
		for (i=1;i<=n-k+1;i++){
			t=0;
			for (j=0;j<c;j++)
				if (sum[j][i+k-1]-sum[j][i-1]==0) t|=1<<j;
			tag[t]=1;
		}
		for (i=(1<<c)-1;i;i--) if (tag[i])
			for (j=0;j<c;j++)
				if (i>>j&1) tag[i^(1<<j)]=1;
		for (i=0;i<1<<c;i++)
			if (!tag[i] && i>>(s[s.size()-1]-'A')&1) ans=min(ans,__builtin_popcount(i));
		cout<<ans<<endl;
	}
}