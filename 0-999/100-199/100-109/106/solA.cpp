#include <iostream>

using namespace std;

int main()
{
    char m;
    string a, b, t="6789TJQKA";
    cin >> m >> a >> b;
    if (((a[1]==m)&&(b[1]!=m)) || ((a[1]==b[1])&&(t.find(a[0])>t.find(b[0]))))
        cout << "YES";
    else
        cout << "NO";
    return 0;
}