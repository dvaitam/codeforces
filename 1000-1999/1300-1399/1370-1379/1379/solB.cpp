#include<stdio.h>

int main(){
int t;scanf("%d",&t);
while(t--){
int l,r;
long long m;
scanf("%d%d%lld",&l,&r,&m);
for(int a=l;a<=r;++a){
int d=(m-l+r)%a+l-r;
if(d>r-l)continue;
int c=d>0?l:(l-d);
int b=c+d;
printf("%d %d %d\n",a,b,c);
break;
}
}
}