#include<stdio.h>

double n,m,k;
double ans;
int main(){
    while(scanf("%lf %lf %lf",&n,&m,&k)!=EOF){
        ans = 1;
        if(n+k < m) printf("0\n");
        else{
            for(double i=0;i<=k;i++)
                ans *= (m-i)/(n+1+i);
            printf("%.6lf\n",1.0-ans);
        }
    }
    return 0;
}