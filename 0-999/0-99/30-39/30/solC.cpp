#include <cstdio>
#include <cmath>
#include <iostream>
#include <algorithm>

using namespace std;

const int MAXN = 1010;

int n;

struct node{
	int x,y,t;
	double p;
	friend bool operator < (const node &a,const node &b){
		return a.t<b.t;
	}
}ar[MAXN];

double dn[MAXN];

int main(){

	scanf("%d",&n);

	for(int i=1;i<=n;i++)
		scanf("%d %d %d %lf",&ar[i].x,&ar[i].y,&ar[i].t,&ar[i].p);

	double res=0;

	sort(ar+1,ar+n+1);

	for(int i=1;i<=n;i++)
	{
		dn[i]=ar[i].p;
		for(int j=1;j<i;j++)
			if(pow(ar[j].t-ar[i].t,2)-pow(ar[i].x-ar[j].x,2)-pow(ar[i].y-ar[j].y,2)>-1e-9)
				dn[i]=max(dn[i],dn[j]+ar[i].p);
		res=max(res,dn[i]);
	}

	printf("%.9lf\n",res);

	return 0;

}