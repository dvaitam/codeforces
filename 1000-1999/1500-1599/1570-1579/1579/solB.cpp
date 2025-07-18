#include<bits/stdc++.h>
using namespace std;

int t,n,a[51],tmp,book1[51],book2[51],sum,flag;
int minn,mini;

void find_minn(int step){
	for(int i=0;i<=step;i++)
		if(a[i]>minn){
			minn=a[i];
			mini=i;
		}
	return;
}

void solve(){
	for(int i=n-1;i>=0;i--){
		minn=-1000000001;
		mini=-1;
		find_minn(i);
		if(a[i]==minn){
			continue;
		}else{
			//printf("%d %d 1\n",i,mini);
			book1[sum]=i;
			book2[sum]=mini;
			sum++;
			/*/tmp=a[i];
			a[i]=minn;
			a[mini]=tmp;*/
			tmp=a[i];
			for(int j=mini;j<i;j++){
				a[j]=a[j+1];
			}
			a[i]=tmp;
		}
	}
	printf("%d\n",sum);
	for(int i=0;i<sum;i++){
		printf("%d %d 1\n",book2[i]+1,book1[i]+1);
	}
}

int main(){
	scanf("%d",&t);
	while(t--){
		scanf("%d",&n);
		for(int i=0;i<n;i++){
			scanf("%d",&a[i]);
		}
		sum=0;
		flag=0;
		for(int i=0;i<n-1;i++){
			if(a[i]>a[i+1]){
				flag=1;
				break;
			}
		}
		if(flag==0){
			printf("0\n");
			continue;
		}
		solve();
	}
	return 0;
}