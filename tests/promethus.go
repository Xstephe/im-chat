package main

import (
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/prometheus"
	"time"
)

// 使用go-zero内置promethus
func main() {
	prometheus.StartAgent(prometheus.Config{
		Host: "0.0.0.0",
		Port: 1234,
		Path: "/metrics",
	})

	gague := metric.NewGaugeVec(&metric.GaugeVecOpts{
		Name: "tests_go_zerogauge",
		Help: "this is a go_zero gugue",
	})

	var i int
	for {
		i++
		if i%2 == 0 {
			gague.Inc()
		}
		time.Sleep(time.Second)
	}
}

//func main() {
//	//1. go语言运行的相关数据指标
//	//2. 自己定义的采集信息
//	gague := prometheus.NewGauge(prometheus.GaugeOpts{
//		Name: "tests_gauge",
//		Help: "this is a gugue",
//	})
//	prometheus.MustRegister(gague)
//
//	var i int
//	go func() {
//		for {
//			i++
//			if i%2 == 0 {
//				gague.Inc()
//			}
//			time.Sleep(time.Second)
//		}
//	}()
//
//	http.Handle("/metrics", promhttp.Handler())
//	http.ListenAndServe(":1234", nil)
//}
