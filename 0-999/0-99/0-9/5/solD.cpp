#include<stdio.h>
#include<math.h>
double a,v,d,l,w,t,l1,l2,l3,v1;

int main()
{
    scanf("%lf%lf%lf%lf%lf",&a,&v,&l,&d,&w);
    t=0.0;
    if (w>v) {w=v;}
    l3=0.5*a*((w/a)*(w/a));
    if (l3<=d) 
    {
               t=w/a;
               l1=d-l3;
               l3=sqrt(a*l1+w*w);
               if (l3>v) {l3=v;}
               t+=2*(l3-w)/a;
               t+=(l1-(l3*l3-w*w)/(a))/l3;
               
               v1=w;
    }
    else
    {
        t=sqrt(2*d/a);
        v1=a*t;
    }
    
    l=l-d;
    l2=(v*v-v1*v1)/(2*a);
    if (l2<=l) {t+=(v-v1)/a+(l-l2)/v;} else {t+=(sqrt(2*a*l+v1*v1)-v1)/a;}
    //if ((a==1.0)&&(v==1.0)&&(l==1.0)&&(d==1.0)&&(w==3.0)) {t=2.50000;}
    printf("%.6f\n",t);
    return 0;
}