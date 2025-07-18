#include <iostream>
#include <string>

using namespace std;

int main( int argc, char* argv[] )
{
    int n;
    cin >> n;
    double ySum = 0;
    for( int i = 0; i < n; i++ ) {
        double x, y;
        cin >> x >> y;
        ySum += y;
    }

    cout << 5 + ySum / n;

    return 0;
}