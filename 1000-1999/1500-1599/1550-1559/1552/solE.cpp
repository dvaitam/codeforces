#include<bits/stdc++.h>

using namespace std;

int n,k,x,L[105],R[105],p[105],a[105][2],vis[105];

int main(){

	cin>>n>>k;

	for (int i=1; i<=n; i+=k-1)

		for (int j=i; j<=min(n,i+k-2); j++) L[j]=i,R[j]=min(n,i+k-2);

	for (int i=1; i<=n*k; i++){

		scanf("%d",&x);

		if (vis[x]) continue;

		if (p[x]){

			a[x][0]=p[x],a[x][1]=i,vis[x]=1;

			for (int j=L[x]; j<=R[x]; j++) p[j]=0;

		} else p[x]=i;

	}

	for (int i=1; i<=n; i++) printf("%d %d\n",a[i][0],a[i][1]);

	return 0;

}