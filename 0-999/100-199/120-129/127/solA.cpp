#include<iostream>
#include<math.h>
using namespace std;
int main()
{
    int n,k,x,y,i,a,b;double d=0;cin >> n >> k >> a >>b;
    for(i=0;i<n-1;i++){cin >> x >> y;d=d+sqrt((a-x)*(a-x)+(b-y)*(b-y));a=x;b=y;}
    cout.precision(8);
    cout << d*k/50;
}