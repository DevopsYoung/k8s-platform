package service

import (
	"github.com/wonderivan/logger"
	"k8s-platform/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

func (k *k8s) Init() {
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		panic("创建k8s失败, " + err.Error())
	}
	//根据rest.config类型的对象，new一个clientset出来
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic("创建k8s clientset失败, " + err.Error())
	} else {
		logger.Info("k8s client初始化成功")
	}
	k.ClientSet = clientset

	//if pods_list, err := clientset.CoreV1().Pods("default").
	//	List(context.TODO(), metav1.ListOptions{}); err != nil {
	//	panic("faild to get pods list")
	//} else {
	//	for _, pod := range pods_list.Items {
	//		fmt.Println(pod.Name)
	//	}
	//
	//}
}
