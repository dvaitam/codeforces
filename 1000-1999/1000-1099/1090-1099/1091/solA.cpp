#include <iostream>

using namespace std;

int main()
{
    int b,r,y;
    cin>>y>>b>>r;
    cout<<min(3*r-3,min(3*b,3*y+3));
    return 0;
}