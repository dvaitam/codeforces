#include <iostream>
using namespace std;
int main()
{
    long long da,db;
    cin >> da >> db;
    if (db-da==1)
        cout << db*10-1 << " " << db*10;
    else if (da==db)
        cout << da*10 << " " << db*10+1;
    else if (da==9 && db==1)
        cout << "99" << " " << "100";
    else cout << "-1";
    return 0;
}