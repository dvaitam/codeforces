#include<cstdio>
#include<cmath>

int main(int argc,char **argv){

    int hh,mm;
    int H,D,C,N;

    double cost1;
    double cost2;

    scanf("%d%d",&hh,&mm);
    scanf("%d%d%d%d",&H,&D,&C,&N);
    
    int burg=ceil((double)H/N);

    if(hh>=20){

      cost1=((double)(4*C)/5)*burg;
      printf("%0.5f",cost1);
      return(0);   
    } 

    cost1=((double)burg)*C;

    H+=((20-hh)*60-mm)*D;

    burg=ceil((double)H/N);
    cost2=((double)(4*C)/5)*burg;

    if(cost1<cost2)
       printf("%0.5f",cost1);
    else
       printf("%0.5f",cost2);       
          
    return(0);       
}