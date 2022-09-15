package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/service"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//初始化k8s client
	service.K8s.Init()
	//if pods_list, err := service.K8s.ClientSet.CoreV1().Pods("kube-system").
	//	List(context.TODO(), metav1.ListOptions{}); err != nil {
	//	panic("faild to get pods list")
	//} else {
	//	for _, pod := range pods_list.Items {
	//		fmt.Println(pod.Name)
	//	}
	//
	//}
	r := gin.Default()
	controller.Router.InitApiRouter(r)
	r.Run(config.ListenAddr)

}

func test() {
	conf, err := clientcmd.BuildConfigFromFlags("", "kube-config")
	if err != nil {
		panic("error")
	}
	//根据rest.config类型的对象，new一个clientset出来
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic("error")
	}

	if pods_list, err := clientset.CoreV1().Pods("kube-system").
		List(context.TODO(), metav1.ListOptions{}); err != nil {
		panic("faild to get pods list")
	} else {
		for _, pod := range pods_list.Items {
			fmt.Println(pod.Name)
		}

	}
}
