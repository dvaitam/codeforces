#include <cstdio>
double L,D,V,G,R,T,ans=0;
int main()
{
    scanf("%lf%lf%lf%lf%lf",&L,&D,&V,&G,&R);
    T=D/V;
    while(T>G+R)T-=G+R;
    if(T>=G)ans+=G+R-T;
    ans+=L/V;
    if(D>=L)ans=L/V;
    printf("%lf\n",ans);
}