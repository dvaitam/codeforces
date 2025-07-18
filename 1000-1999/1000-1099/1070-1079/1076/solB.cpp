#include <bits/stdc++.h>
using namespace std;
bool ngto(long long n)
{
    if(n<2) return 0;
    for(long long i=2;i<=sqrt(n);i++)
        if(n%i==0) return 0;
    return 1;
}
int main()
{
    long long n;
    cin>>n;
    if(ngto(n))
    {
        cout<<1;
        return 0;
    }
    else if(n%2==0)
        cout<<n/2;
    else
    {
    for(int i=3;i<=n;i++)
    {
        if(n%i==0)
        {
            cout<<(n-i)/2+1;
            return 0;
        }
    }
    }
}