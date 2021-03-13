package storage

import "fmt"

//1:实现两个线程，使之交替打印1-100,
//如：两个线程分别为：Printer1和Printer2,
//最后输出结果为： Printer1 — 1 Printer2 一 2 Printer1 一 3 Printer2 一 4

func ParallelPrint() {
	print1Ch := make(chan int)
	print2Ch := make(chan int)
	// goroutine1
	go func(beginNum int) {
		printNum := beginNum
		fmt.Print(printNum)
		print2Ch <- printNum + 1
		for {
			printNum = <-print1Ch
			fmt.Print(printNum)
			print2Ch <- printNum + 1
			if printNum == 99 {
				close(print1Ch)
			}
		}
	}(1)
	// goroutine2
	go func(beginNum int) {
		for {
			printNum := <-print2Ch
			fmt.Print(printNum)
			print1Ch <- printNum + 1
			if printNum == 100 {
				close(print2Ch)
			}
		}
	}(2)
}

//
//2:实现函数,给定一个字符串数组,求该数组的连续非空子集，分別打印出来各子集 ，举例数组为[abc]，输出[a],[b],[c],[ab],[bc],[abc]
func PrintSubset(s string) {
	length := len(s)
	subsets := make(map[string]interface{})

	for beginIndex := range s {
		for endIndex := beginIndex; endIndex <= length; endIndex++ {
			subsets[s[beginIndex:endIndex]] = nil
		}
	}

	for subset := range subsets {
		fmt.Print(subset + " ")
	}
	fmt.Println()
}

//
//3:给定任意一个正整数（int 表示范围内）n，求n!（n的阶乘）结果的末尾有几个连续的零，如：3!=6，有0个连续的0；12!=479001600，有2个连续的0；

func GetZeroNums(num int) (zeroNum int) {
	if num == 0 {
		return 1
	}
	var fiveNum, twoNum int
	for cur := 2; cur <= num; cur++ {
		fiveNum += getNumExpr(5, cur)
		twoNum += getNumExpr(2, cur)
		if fiveNum > 0 && twoNum > 0 {
			var newZero int
			if fiveNum >= twoNum {
				newZero = twoNum
				fiveNum -= twoNum
				twoNum = 0
			} else {
				newZero = fiveNum
				twoNum -= fiveNum
				fiveNum = 0
			}
			zeroNum += newZero
		}
	}
	return
}

func getNumExpr(dividen, originNum int) (expr int) {
	for originNum%dividen == 0 {
		expr++
		originNum /= dividen
	}
	return
}
