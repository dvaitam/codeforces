#include<bits/stdc++.h>
using namespace std;
bool cmp(int a,int b){
	return a>b;
}
int main(){
	int n,h;
	while(cin>>n>>h){
		int data[1005];
		int max=0;
		for(int i=1;i<=n;i++){
			scanf("%d",&data[i]);
		}
		int ans[1005];
		memset(ans,0,sizeof(ans));
		for(int i=1;i<=n;i++){
			for(int j=1;j<=i;j++){
				ans[j]=data[j];
			}
			sort(ans+1,ans+1+i,cmp);
			int sum=0;
			for(int j=1;j<=i;j+=2){
				sum+=ans[j];
			}
			if(sum<=h){
				max++;
			}
			else break;
		}
		cout<<max<<endl;
	}
	return 0;
}