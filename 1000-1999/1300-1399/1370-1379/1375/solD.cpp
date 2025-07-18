#include<bits/stdc++.h>

int a[1005],c[1005],x[5005];
int main(){
int T,n,i,j,s;
for(scanf("%d",&T);T--;){
memset(c,0,sizeof(c));
for(scanf("%d",&n),i=0;i<n;i++)scanf("%d",a+i),c[a[i]]++;
for(s=0;;){
for(i=1;i<n&&a[i]>=a[i-1];i++);
if(i>=n)break;
for(i=0;c[i];i++);
x[s++]=j=i-(i>=n);
n-=i>=n;
c[a[j]]--,c[i]++,a[j]=i;
}
for(printf("%d\n",s),i=0;i<s;i++)printf("%d ",x[i]+1);puts("");
}
}