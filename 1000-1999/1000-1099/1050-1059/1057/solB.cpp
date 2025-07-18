#include<iostream>
#include<cstdio>
#include<cstring>
#include<cmath>
#include<algorithm>
using namespace std;
int n,ans;
int a[5010],sum[5010];
int main()
{
	scanf("%d",&n);
	for(int i=1;i<=n;i++){
		scanf("%d",&a[i]);
		sum[i]=sum[i-1]+a[i];
	}
	for(int i=1;i<=n;i++){
		for(int j=0;j<i;j++){
			if((sum[i]-sum[j])>(i-j)*100) {
				ans=max(ans,i-j);
			}
		}
	}
	printf("%d",ans);
	return 0;
}