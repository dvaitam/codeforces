#include<cstdio>
int main(){
    int n,ans,i,j,k,u,v,x,y;
    int a[404];
    scanf("%d",&n);
    for (i=0;i<n;i++) scanf("%d",a+i);
    ans=0x7fffffff;
    for (i=n-1;i>=0;i--){
        u=a[i]-1,x=1;
        if (a[i]>=ans) break;
        for (j=i+1;j<n;j++){
            x+=(v=u/a[j]);
            u-=v*a[j];
            v=a[i]+a[j]-u-1;
            if (v>=ans) continue;
            for (y=0,k=0;v && y<=x;k++){
                y+=v/a[k];
                v%=a[k];
            }
            if (y>x) ans=a[i]+a[j]-u-1;
        }
    }
    printf("%d\n",ans==0x7fffffff?-1:ans);
    return 0;
}