# include <stdio.h>
# include <math.h>
int main(void){
double a,b,c;
double r1,r2,d,q;
scanf("%lf",&a);
scanf("%lf",&b);
scanf("%lf",&c);
if(a==0 && b==0 && (c!=0)){
printf("0");
return 0;
}
if(a==0 && b==0 && c==0){
printf("-1");
return 0;
}

if (a==0){
    r1=(-1*c)/b;
printf("1\n");

printf("%lf",r1);
return 0;
}
q=b*b-(4*a*c);
if(q<0){
printf("0");
return 0;
}
d=sqrt(q);
if(d==0){
printf("1\n");
r1=(-1*b)/(2*a);
printf("%lf",r1);
return 0;
}
if(d!=0){
printf("2\n");
r1=((-1*b)-d)/(2*a);
r2=((-1*b)+d)/(2*a);
if(r1>r2){
printf("%lf\n",r2);
printf("%lf\n",r1);
}
if(r1<r2){
printf("%lf\n",r1);
printf("%lf\n",r2);
}
}
}