#include <bits/stdc++.h>
using namespace std;

int n, w;
int ans;
struct cha{
    int b;
    int index;
    int c;
};

cha a[200];

bool cmp1(cha a, cha b)
{
    return a.b > b.b;
}

bool cmp2(cha a, cha b)
{
    return a.index < b.index;
}

int main()
{
    cin >> n >> w;
    for(int i = 0; i < n; i++)
    {
        int aa;
        cin >> aa;
        a[i].b = aa;
        a[i].index = i;
        ans += (aa+1)/2;
    }
    if(w < ans)
        cout << "-1" << endl;
    else
    {
        sort(a,a+n,cmp1);
        int sum = 0;
        for(int i = 0; i < n; i++)
        {
            a[i].c = (a[i].b+1)/2;
            sum += a[i].c;
        }
        //cout << sum << endl;
        for(int i = 0; i < n; i++)
        {
            //cout << a[i].b << endl;
            if(w-(sum - a[i].c) > a[i].b)
            {
                sum = sum - a[i].c;
                a[i].c = a[i].b;
                w = w-a[i].c;
                //cout << sum << " " << w << endl;
            }
            else
            {
                //cout << "yes" << endl;
                a[i].c = w-(sum - a[i].c);
                break;
            }
        }
        sort(a,a+n,cmp2);
        for(int i = 0; i <n; i++)
            cout << a[i].c << " ";
    }
    return 0;
}