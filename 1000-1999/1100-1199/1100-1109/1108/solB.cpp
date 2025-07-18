#include<iostream>
#include<vector>
#include<unordered_set>
#include<algorithm>
using namespace std;

int main()
{
    long long int n,gc=0,c,b=0,t,t1,ans=0;
    cin>>n;

    vector<int>num;
    unordered_set<int>s;
    unordered_set<int>:: iterator it;

    for(int i=0;i<n;i++)
    {
        cin>>t;
        if(t>b)b=t;
        num.push_back(t);
    }
    int si = num.size();
    sort(num.begin(),num.end());

    for(int i=2;i<b;i++)
    {
        if(b%i==0)
        {
            for(int j=0;j<num.size();j++)
            {
                if(num[j]==i)
                {
                    num[j]=0;
                    break;
                }
            }
        }
    }
    for(int i=num.size()-2;i>=0;i--)
    {
        if(num[i]!=0)
        {
            cout<<num[i]<<" "<<b<<endl; return 0;
        }
    }






}