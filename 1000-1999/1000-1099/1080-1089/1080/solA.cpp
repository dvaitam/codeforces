#include <iostream>

using namespace std;

int main()
{
    int n,k;

    cin>>n>>k;

    int a=0,v=0,r=0;

    a=n*8;

    if(a%k==0)
    {
        a=a/k;
    }
    else
    {
        a=(a/k)+1;
    }

    v=n*5;

    if(v%k==0)
    {
        v=v/k;
    }
    else
    {
        v=(v/k)+1;
    }

    r=n*2;

    if(r%k==0)
    {
        r=r/k;
    }
    else
    {
        r=(r/k)+1;
    }

    cout<<r+a+v;

    return 0;
}