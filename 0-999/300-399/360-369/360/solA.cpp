#include<iostream>
#include<cstring>
#include<algorithm>
#include<cstdio>
#include<set>
#include<vector>
using namespace std;
int plu[5005],mx[5005];
vector < int > col[5005];
set < int > st;
int n,m;
void solve(){
	int i,j,s=0,a,b,c,d;
	scanf("%d %d",&n,&m);
	for(i=1;i<=n;i++)
		mx[i]=1000000000;
	for(j=1;j<=m;j++){
		scanf("%d %d %d %d",&a,&b,&c,&d);
		if(a==1){
			for(i=b;i<=c;i++)
				plu[i]+=d;
		}
		else{
			s++;
			for(i=b;i<=c;i++){
				if(mx[i]>d-plu[i]){
					mx[i]=d-plu[i];
					col[i].clear();
				}
				if(mx[i]==d-plu[i])
					col[i].push_back(s);
			}
		}
	}
	for(i=1;i<=n;i++)
		for(j=0;j<col[i].size();j++)
			st.insert(col[i][j]);
	if(st.size()==s){
		printf("YES\n");
		for(i=1;i<=n;i++)
			printf("%d ",mx[i]);
		printf("\n");
	}
	else{
		printf("NO\n");
	}
}
int main(){
	solve();
	return 0;
}