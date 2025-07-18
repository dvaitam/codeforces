/*

宣告——
汝身听吾号令，托付吾之命运于汝之剑，
应圣杯之召，若汝遵从此意志此理，回应吧。
在此起誓，吾愿成就世间一切之善行，吾愿诛尽世间一切之恶行。
汝为身缠三大言灵之七天，从抑止之轮显现吧，天秤之守护者！

*/
#include<iostream>
#include<cstdio>
#include<cstring>
#include<algorithm>
#include<cassert>
#include<cmath>
#include<vector>
#include<queue>
#define ll long long
using namespace std;
inline int read(){
	int re=0,flag=1;char ch=getchar();
	while(!isdigit(ch)){
		if(ch=='-') flag=-1;
		ch=getchar();
	}
	while(isdigit(ch)) re=(re<<1)+(re<<3)+ch-'0',ch=getchar();
	return re*flag;
}
int n,m,cnt[1000010],dp[1000010][3][3];bool vis[1000010][3][3];
int main(){
	n=read();m=read();int ans=0,i,j,k,tmp,cur;
	for(i=1;i<=n;i++) cnt[read()]++;
//	for(i=1;i<=m;i++) while(cnt[i]>5) ans++,cnt[i]-=3;
	dp[0][0][0]=0;vis[0][0][0]=1;
	for(i=1;i<=m;i++){
		for(j=0;j<=2;j++){
			for(k=0;k<=2;k++){
				if(!vis[i-1][j][k]) continue;
//				cout<<"do "<<i-1<<' '<<j<<' '<<k<<' '<<dp[i-1][j][k]<<'\n';
				for(cur=0;cur<=2;cur++){
					if(cnt[i]<cur) break;
					if(cnt[i-1]-k<cur) break;
					if(cnt[i-2]-k-j<cur) break;
					dp[i][k][cur]=max(dp[i][k][cur],dp[i-1][j][k]+cur+(cnt[i-2]-k-j-cur)/3);
					vis[i][k][cur]=1;
				}
			}
		}
	}
	tmp=0;
	for(j=0;j<=2;j++){
		for(k=0;k<=2;k++){
			if(!vis[m][j][k]) continue;
			tmp=max(tmp,dp[m][j][k]+(cnt[m-1]-j-k)/3+(cnt[m]-k)/3);
		}
	}
	cout<<ans+tmp<<'\n';
}