#include <iostream>
using namespace std;
int main()
{
    int n, t[111], z=0, kol=1, max=0;
    cin >> n;
    for(int i=0;i<n;i++)
    {
        cin >> t[i];
    }
    for(int i=0;i<n;i++)
    {
        kol=1;
        for(int y=0;y<n;y++)
        {
            if(i!=y)
            {
                if(t[i]==t[y])
                {
                    kol++;
                }
            }
        }
        if(kol>max){max=kol;}
    }
    cout<<max<<endl;
    return 0;
}