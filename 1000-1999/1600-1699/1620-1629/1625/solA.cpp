#include<bits/stdc++.h>
using namespace std;
int main()
{
	int T,n,l,a[110];
	cin>>T;
	while(T--){
		int ans=0;
		cin>>n>>l;
		for(int i=1;i<=n;i++)cin>>a[i];
		for(int i=0;i<30;i++){
			int cnt=0;
			for(int j=1;j<=n;j++)if(a[j]>>i&1)cnt++;
			if(cnt*2>=n)ans|=1<<i;
		}
		cout<<ans<<endl;
	}
}