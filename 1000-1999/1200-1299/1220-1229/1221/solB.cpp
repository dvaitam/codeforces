#include<iostream>
using namespace std;

int main()
{
    int n, i, j;
    char a[100][100];
    cin>>n;
    for (i=0; i<n; i++)
    {    
        for (j=0; j<n; j++)
            if ((i+j)%2==0)
                a[i][j]='W';
            else
                a[i][j]='B';
    }
    for (i=0; i<n; i++)
    {
         for (j=0; j<n; j++)
            cout<<a[i][j];
    
        cout<<endl;
    }
    return 0;
}