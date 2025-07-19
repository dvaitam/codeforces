#include<cstdio>
int n,q,a[200001],x,y,req[200001],t[200001];
void ad(int x){for(;x<200001;x+=x&-x)++t[x];}
int main(){
	scanf("%d%d",&n,&q);
	for(int i=1;i<=n;++i)scanf("%d",a+i);
	for(int i=1;i<=n;++i){
		x=y=0;
		for(int j=17;~j;--j)if(1ll*a[i]*(x|1<<j)<=y+t[x|1<<j])x|=1<<j,y+=t[x];
		ad(++x),req[i]=x;
	}
	while(q--)scanf("%d%d",&x,&y),printf("%s\n",y<req[x]?"NO":"YES");
	return 0;
}