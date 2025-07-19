#include<bits/stdc++.h>
int T,n,m,ans;
int main(){
	scanf("%d",&T);
	while(T--){
		scanf("%d%d",&n,&m);
		ans=0;
		for(int i=1;i<=n&&i<=m;++i)
			for(int j=1;i+j<m&&i*j<n;++j)
				ans+=std::min((n-i*j)/(i+j),m-i-j);
		printf("%d\n",ans);
	}
	return 0;
}
