#include <bits/stdc++.h>
using namespace std;
int n,q[1005],cev,k[1005];
int main(){
	scanf("%d",&n);
	for(int i=1;i<=n;i++)	
		scanf("%d",&q[i]);
		q[n+1]=1;
	for(int i=1;i<=n+1;i++){
		if(q[i]==1){
			cev++;
			}				
	}	
	cev-=1;
	printf("%d\n",cev);
							//cout<<cev;
							//return 0;
	for(int i=2;i<=n+1;i++){						
		if(q[i]==1)
		printf("%d ", q[i-1] );	
	}
	return 0;	
}