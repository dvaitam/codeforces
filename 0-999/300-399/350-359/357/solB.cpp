#include<iostream>
#include<cstdio>
using namespace std;

int main(){
	int i,j,n,num,t[4],m,k;
	int a[100001]={0};bool b[4];
	scanf("%d%d",&n,&m);
	for(i=1;i<=m;i++){
		scanf("%d%d%d",&t[1],&t[2],&t[3]);
		for(j=1;j<=3;j++)b[j]=false;
		b[a[t[1]]]=true;
		b[a[t[2]]]=true;
		b[a[t[3]]]=true;
		for(j=1;j<=3;j++)
			for(k=1;k<=3;k++)
				if(a[t[j]]==0&&b[k]==false){
					a[t[j]]=k;b[k]=true;}
		}
	printf("%d",a[1]);
	for(i=2;i<=n;i++)
		printf(" %d",a[i]);
	cout<<endl;
}