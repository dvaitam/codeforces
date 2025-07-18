#include <iostream>
using namespace std;
int main()
{
    int cases;
    cin>>cases;
    while(cases--)
    {
        int x=0;
        int n;
        cin>>n;
        int arr[n];
        for(int i=0;i<n;i++)
        {
            cin>>arr[i];
            x^=arr[i];
        }
        
        if(n%2)
        {
            cout<<x<<endl;
        }
        else
        {
            if(x==0)
            {
                cout<<arr[n-1]<<endl;
            }
            else
            {
                cout<<"-1"<<endl;
            }
        }
    }
}