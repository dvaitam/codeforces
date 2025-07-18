#include <bits/stdc++.h>
using namespace std;
const int N = 1e5+5;
int n;
#define pdd pair<double,double>
#define x first
#define y second
pdd t[N*2];
pdd merge(pdd x,pdd y){
	double a1=x.x,a2=y.x,b1= x.y,b2=y.y,r1,r2;
	if(a1==-1)return y;
	if(a2==-1)return x;
	r1= a1*a2/(1-b1*(1-a2));
	r2= b2+ ( (1-b2)*b1*a2/(1-(1-a2)*b1));
	return make_pair(r1,r2);
}
void build(){
	int i;pdd x;
	for(i=n-1;i>=1;--i){
		t[i]= merge(t[i<<1],t[i<<1|1]);
	}
}
void update(int ind,double p){
	int i= ind+n-1;pdd x;
	for(t[i].x=p,t[i].y=p,i>>=1;i>1;i>>=1)
		t[i]= merge(t[i<<1],t[i<<1|1]);
}
pdd query(int l,int r){
	pdd x,y;x.x=-1,y.x=-1;
	for(l+=n-1,r+=n;l<r;l>>=1,r>>=1){
		if(l&1)x=merge(x,t[l++]);
		if(r&1)y=merge(t[--r],y);
	}
	return merge(x,y);
}
int main(){
	int a,b,i,q,type,l,r;double p;
	cin>>n>>q;
	for(i=0;i<n;++i)
		scanf("%d%d",&a,&b),p= (double)a/b,t[i+n].x= p,t[i+n].y=p;
	build();
	pdd ret;
	while(q--){
		scanf("%d",&type);
		if(type==1){
			scanf("%d%d%d",&i,&a,&b);p= (double)a/b;
			update(i,p);
		}
		else{
			scanf("%d%d",&l,&r);
			ret= query(l,r);
			printf("%0.9lf\n",ret.x);
		}
	}
}