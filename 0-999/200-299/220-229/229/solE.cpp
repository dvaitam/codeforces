//229E
#include <iostream>
#include <cstdio>
#include <vector>
#include <algorithm>
using namespace std;
const int MAXN=1000;
double f[MAXN+10][MAXN+10];
int n,m,k[MAXN+10],a[MAXN+10];
vector<pair<int,int> > v;
int main(){
	scanf("%d%d",&n,&m);
	for(int i=1;i<=m;i++){
		scanf("%d",&k[i]);
		for(int j=0;j<k[i];j++){
			int x;
			scanf("%d",&x);
			v.push_back(make_pair(x,i));
		}
	}
	sort(v.begin(),v.end(),greater<pair<int,int> >());
	int keng=v[n-1].first;
	double p=1;
	int done=0;
	for(int i=0;v[i].first!=keng;i++){
		int cur=v[i].second;
		p*=(a[cur]+1)/double(max(k[cur],1));
		a[cur]++;
		k[cur]--;
		done++;
	}
	vector<int> ids;
	for(int i=0;v[i].first>=keng;i++)if(v[i].first==keng){
		ids.push_back(v[i].second);
	}
	int s=ids.size();
	static double tmp[MAXN+10];
	for(int i=0;i<s;i++)
		tmp[i]=(a[ids[i]]+1)/double(max(k[ids[i]],1));
	f[0][0]=p;
	for(int j=0;j<s;j++)
		for(int i=0;i<=j;i++){
			f[i+1][j+1]+=f[i][j]*tmp[j]*(i+1)/(j+1);
			f[i][j+1]+=f[i][j]*(j-i+1)/(j+1);
		}
	printf("%.10f\n",f[n-done][s]);
	return 0;
}