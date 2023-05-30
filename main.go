package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	N        = 5   //传输带上最多可以放_盘寿司
	m        = 3   //有_位做寿司的师傅
	n        = 15  //当时有_位顾客在店
	totalNum = 100 //寿司原材料全部制作，可完成_份
)

type Cooker struct {
	name       string
	totalNumP  int           //需完成的总量
	presentNum int           //目前完成的份数
	perTime    time.Duration //做一个寿司需要的时间
}

type Customer struct {
	name       string
	buyNum     int           //购买的份数
	perEatTime time.Duration //吃一个寿司需要的时间
}

var Ch = make(chan string, N)
var Total = make(chan int, totalNum)
var (
	logFilename = flag.String("log", "diary.log", "helpmessage")
)

func (cook *Cooker) Produce(in chan<- string) {
	for i := 0; i < cook.totalNumP; i++ {
		select {
		case v := <-Total:
			fmt.Printf("第%v份材料\n", v)
			log.Printf("第%v份材料\n", v)
			cook.presentNum++
			s1 := cook.name
			s2 := strconv.FormatInt(int64(cook.presentNum), 10)
			s3 := s1 + "做的第" + s2 + "个寿司"
			in <- s3
			fmt.Println(cook.name+"制作了寿司", s3)
			log.Println(cook.name+"制作了寿司", s3)
			time.Sleep(cook.perTime)
		default:
			fmt.Println("原材料用完了。")
		}
	}

}

func (eat *Customer) Buy(out <-chan string) {
	for i := 0; i < eat.buyNum; i++ {
		j := <-out
		time.Sleep(eat.perEatTime)
		fmt.Println(eat.name+"吃了寿司", j)
		log.Println(eat.name+"吃了寿司", j)
		time.Sleep(eat.perEatTime)
	}
}

func diary() {
	log.Printf("\n")
	log.Printf("\n")
	log.Printf("工作完成，闭店\n")
	//log.Printf("已做寿司：%v，传送带现有寿司数：%v\n", count, len(Ch))
	log.Printf("\n")
	log.Printf("\n")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	logFile, Err := os.OpenFile(*logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if Err != nil {
		fmt.Println("Can not open", *logFile)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("程序正常启动")

	defer diary()

	for i := 1; i <= 100; i++ {
		Total <- i
	}

	cookers := [m]Cooker{
		{"厨师A", 50, 0, 2e9},
		{"厨师B", 30, 0, 3e9},
		{"厨师C", 35, 0, 4e9},
	}

	customers := [n]Customer{
		{"顾客1", 5, 3e9},
		{"顾客2", 3, 4e9},
		{"顾客3", 6, 3e9},
		{"顾客4", 2, 2e9},
		{"顾客5", 1, 6e9},
		{"顾客6", 7, 2e9},
		{"顾客7", 4, 2e9},
		{"顾客8", 1, 4e9},
		{"顾客9", 5, 3e9},
		{"顾客10", 2, 3e9},
		{"顾客11", 4, 6e9},
		{"顾客12", 3, 4e9},
		{"顾客13", 5, 5e9},
		{"顾客14", 1, 2e9},
		{"顾客15", 2, 3e9},
	}
	for i := 0; i < m; i++ {
		go cookers[i].Produce(Ch)
	}
	for i := 0; i < n; i++ {
		go customers[i].Buy(Ch)
	}
	time.Sleep(1e11)
}
