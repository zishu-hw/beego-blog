/*
 * ==============================================================
 *
 *       FileName: test.cpp
 *    Description:
 *
 *         Author: zhaiyu, zishuzy@gmail.com
 *        Created: 2018-12-21 09:52:21
 *  Last Modified: 2019-01-08 20:42:19
 *
 *  Copyright (C) 2018 zhaiyu. All rights reserved.
 *
 * ==============================================================
 */

#include <iostream>
#include <stdio.h>
#include <vector>
#include <limits>

using namespace std;

int main(void)
{
    cout << "double: \t"
         << "所占字节数：" << sizeof(double);
    cout << "\t最大值：" << (numeric_limits<double>::max)();
    cout << "\t最小值：" << (numeric_limits<double>::min)() << endl;
    // double t = 0;
    // while (true) {
    //     t += 100000000000;
    //     std::cout << t << std::endl;
    // }
    return 0;
}
