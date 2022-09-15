package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct {
}

//定义列表的返回内容，Items是pod元素的列表，Total是元素的数量
type PodResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}

//获取pod列表，支持过滤、分页、排序

func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodResp, err error) {
	//context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时，这里的常用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label、field等
	//kubectl get services --all-namespaces --field-selector metadata.namespace != default
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取pod列表失败," + err.Error())
		//返回给上一层，最终返回给前端，前端打印出来这个error
		return nil, errors.New("获取pod列表失败, " + err.Error())
	}
	//实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	//再排序和分页
	data := filtered.Sort().Paginate()
	fmt.Printf("分页后的数据长度：%d，数据：%s\n", len(data.GenericDataList), data.GenericDataList)
	//将Data转换成podCell类型
	pods := p.fromCells(data.GenericDataList)

	//数据处理后的数据和原始数据的比较
	//处理后的数据
	fmt.Println("处理后的数据")
	for _, pod := range pods {
		fmt.Println(pod.Name, pod.CreationTimestamp.Time)
	}
	//原始数据
	fmt.Println("原始数据")
	for _, pod := range podList.Items {
		fmt.Println(pod.Name, pod.CreationTimestamp.Time)
	}
	return &PodResp{
		Total: total,
		Items: pods,
	}, nil
}

//toCells方法用于将pod类型数组，转换成DataCell类型数组
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	//定义一个切片
	cells := make([]DataCell, len(std))
	for i := range std {
		//目的：给列表添加元素，每个元素值包含两个属性：1、time 2、name
		//dataselecter.go中podCell已经实现了DataCell接口，那这里podCell等价与接口（DataCell）
		// v, ok := podCell.(DataCell)
		cells[i] = podCell(std[i])

	}
	fmt.Printf("cells => %#v", cells)
	return cells

}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {

	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//cells[i].(podCell)是将DataCell类型转化成podCell类型
		pods[i] = corev1.Pod(cells[i].(podCell))
	}

	return pods
}

//interface{} ()
