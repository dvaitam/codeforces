#include <stdio.h>

int main() {
	int n,p,a,b,d,to[1001][2],ans=0,tank[501];
	bool out[1001]={},in[1001]={};
	scanf("%d%d",&n,&p);
	while(p-->0){
		scanf("%d%d%d",&a,&b,&d);
		to[a][0]=b;
		to[a][1]=d;
		out[a]=1;
		in[b]=1;
	}
	for(int i=1;i<=n;++i)
		if(out[i]&&!in[i])
			tank[ans++]=i;
	printf("%d",ans);
	for(int i=0;i<ans;++i){
		int min=to[tank[i]][1],at=to[tank[i]][0];
		while(out[at]){
			if(min>to[at][1])
				min=to[at][1];
			at=to[at][0];
		}
		printf("\n%d %d %d",tank[i],at,min);
	}
	return 0;
}