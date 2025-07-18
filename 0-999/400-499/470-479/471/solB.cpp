#include <iostream>
#include <string>
#include <cmath>
using namespace std;

int main ()
{
    int n, *tasks, *pos;
    cin >> n;
    tasks = new int[n];
    pos = new int[n];
    for(int i = 0; i < n; i++)
    {
        cin >> tasks[i];
        pos[i] = i + 1;
    }

    for(int i = 0; i < n; i++)
        for(int j = 0; j < n - i - 1; j++)
            if(tasks[j] > tasks[j + 1])
            {
                swap(tasks[j], tasks[j + 1]);
                swap(pos[j], pos[j + 1]);
            }

    int l1 = -1, l2 = -1;
    for(int i = 0; i < n - 1; i++)
        if(tasks[i] == tasks[i + 1])
        {
            if(l1 == -1)
                l1 = i;
            else
                l2 = i;
        }
    if(l2 == -1)
    {
        cout << "NO";
        return 0;
    }
    
    cout << "YES\n";
    for(int i = 0; i < n; i++)
        cout << pos[i] << " ";
    cout << endl;
    for(int i = 0; i < n; i++)
    {
        if(i == l1)
            cout << pos[l1 + 1];
        else if (i == l1 + 1)
            cout << pos[l1];
        else
            cout << pos[i];
        cout << " ";
            
    }
    cout << endl;
    for(int i = 0; i < n; i++)
    {
        if(i == l2)
            cout << pos[l2 + 1];
        else if (i == l2 + 1)
            cout << pos[l2];
        else
            cout << pos[i];
        cout << " ";
        
    }
    cout << endl;
    
    delete [] tasks;
    delete [] pos;
    return 0;
}