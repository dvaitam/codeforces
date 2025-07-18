#include<bits/stdc++.h>

using namespace std;

#define F first

#define S second

#define in(n) scanf("%d",&n)

int main()

{

    int a,b,x,y,c,d,r1,r2;

    double r;

    in(a),in(b),in(r1);

    in(x),in(y),in(r2);

    c=(a-x)*(a-x),d=(b-y)*(b-y);

    cout << setprecision(12) << fixed;

    r=sqrt(c+d);

    if(r>r1+r2)r=r-r1-r2,r=r/2;

    else if(r+min(r1,r2)<max(r1,r2))r=max(r1,r2)-min(r1,r2)-r,r=r/2;

    else

        r=0;

    cout << r << endl;

}