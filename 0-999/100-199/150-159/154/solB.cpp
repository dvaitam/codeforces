#include<cstdio>
int a[102400],b[102400],n,m,i,j,t,g;
char o[102400],c;
int main(){
	scanf("%d%d",&n,&m);
	for(i=2;i<=n;++i)
	if(!a[i])
		for(j=i;j<=n;j+=i)a[j]=i;
	while(m--){
		scanf(" %c%d",&c,&i);
		if(c=='+'){
			if(o[i]){puts("Already on");continue;}
			g=1;
			for(j=i;j>1;){
				t=a[j];
				if(b[t]){
					printf("Conflict with %d\n",b[t]);
					g=0;break;
				}
				do{j/=t;}while(j%t==0);
			}
			if(g){
				puts("Success");
				o[i]=1;
				for(j=i;j>1;){
					b[t=a[j]]=i;
					do j/=t; while(j%t==0);
				}
			}
		}else{
			if(!o[i]){puts("Already off");continue;}
			for(j=i;j>1;){
				b[t=a[j]]=0;
				do j/=t;while(j%t==0);
			}o[i]=0;
			puts("Success");
		}
	}
	return 0;
}