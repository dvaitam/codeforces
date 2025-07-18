#include<stdio.h>
#include<stdlib.h>
struct vertices
{
    double x;
    double y;
}V[10];

int main()
{
    double x1,x2,x3,y1,y2,y3,one,two;
    int i,last=-1;
    scanf("%lf %lf %lf %lf %lf %lf",&x1,&y1,&x2,&y2,&x3,&y3);
    one = (x1 + x2)/2.0;
    two = (y1 + y2)/2.0;
    one = one + (one - x3);
    two = two + (two - y3);
    for(i=0;i<=last;i++)
        if(one==V[i].x && two==V[i].y)
        break;
    if(i==last+1 && one == (int)one && two == (int)two)
    {
        V[++last].x = one;
        V[last].y = two;
    }

    one = (x2 + x3)/2.0;
    two = (y2 + y3)/2.0;
    one = one + (one - x1);
    two = two + (two - y1);
    for(i=0;i<=last;i++)
        if(one==V[i].x && two==V[i].y)
        break;
    if(i==last+1 && one == (int)one && two == (int)two)
    {
        V[++last].x = one;
        V[last].y = two;
    }

    one = (x1 + x3)/2.0;
    two = (y1 + y3)/2.0;
    one = one + (one - x2);
    two = two + (two - y2);
    for(i=0;i<=last;i++)
        if(one==V[i].x && two==V[i].y)
        break;
    if(i==last+1 && one == (int)one && two == (int)two)
    {
        V[++last].x = one;
        V[last].y = two;
    }

    printf("%d\n",last+1);
    for(i=0;i<=last;i++)
        printf("%.0lf %.0lf\n",V[i].x,V[i].y);
    return 0;
}