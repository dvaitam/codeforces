#include <iostream>
using namespace std;

int main()
{
    int a, b, x;
    cin >> a >> b >> x;

    string s0;
    string s1;

    for(int i = 0; i <= x; ++i)
    {
        if(i % 2 == 0)
        {
            s0.push_back('0');
            s1.push_back('1');
        }
        else
        {
            s0.push_back('1');
            s1.push_back('0');
        }
    }

    int atemp, btemp;
    if(x % 2)
    {
        atemp = a - (x + 1)/2;
        btemp = b - (x + 1)/2;
        if(btemp >= 0 and atemp >= 0)
        {
            s0 = string(atemp, '0') + s0 + string(btemp, '1');

            cout << s0 << "\n";
            return 0;
        }
    }
    else
    {
        // s0
        atemp = a - 1 - (x/2);
        btemp = b - (x)/2;
        if(btemp >= 0 and atemp >= 0)
        {
            s0 = string(atemp, '0') + s0.substr(0, s0.length() - 1) + string(btemp, '1') + s0.substr(s0.length() - 1, 1);

            cout << s0 << "\n";
            return 0;
        }

        // s1
        atemp = a - (x)/2;
        btemp = b - 1 - (x/2);
        if(btemp >= 0 and atemp >= 0)
        {
            s1 = string(btemp, '1') + s1.substr(0, s1.length() - 1) + string(atemp, '0') + s1.substr(s1.length() - 1, 1);

            cout << s1 << "\n";
            return 0;
        }
    }
}