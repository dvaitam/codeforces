#include <cstdio>
#include <iostream>
using namespace std;
int main()
{
    double a,b,c,d;
    while(cin>>a>>b>>c>>d)
    {
        double p1 = a/b;
        double p2 = c/d;

        double ans = p1*(1/(1-((1-p2)*(1-p1))));
        printf("%.10lf\n",ans);
    }
    return 0;
}