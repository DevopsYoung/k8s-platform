package service

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
)

/*
最终目的是将数据展示，可以使用排序、过滤、分页的方法进行展示
	如何调用这些方法？
		定义一个结构体A，作为这些方法的接收者
	需要展示的数据从哪里来？
		通过k8s的接口获取。
			那定义的结构A的数据字段(GenericDataList)为什么不使用与k8s接口返回数据类型相同？
				因为获取的数据种类比较多，那么就要定义多个结构体，代码冗余。
				所以在这里引入了DataCell接口，接口在一个结构体实现它的情况下与该结构体等价。
				所以k8s返回的不同的数据类型都可以与DataCell接口等价。
				那数据类型转换的写法是怎样的呢？
					cells[i] = podCell(std[i]) 该如何理解呢？
*/

//dataSelector用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	dataSelectQuery *DataSelectQuery
}

//DataCell接口，用于各种资源List的类型转换，转换后可以使用dataSelector定义的排序、过滤、分页方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

/*DataSelectQuery 定义过滤和分页的属性，
过滤: Name, 分页: Limit和Page
		Limit: 单页的数据条数
		Page： 是第几页
*/
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

/*
排序
	实现自定义结构的排序，需要重写Len、Swap、Less方法
*/
//Len方法用于获取数组长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

//Swap方法用于数组中的元素在比较大小后交换位置，可定义升序或者降序
func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

//Less方法用于定义数组中元素排序‘大小’的比较方式
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

//重写上述三个方法适用sort.Sort进行排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

/*
过滤
*/

//Filter方法用于过滤元素，比较数据的Name属性，若包含，再返回
func (d *dataSelector) Filter() *dataSelector {
	//判断入参是否为空，若为空，则返回所有数据
	if d.dataSelectQuery.FilterQuery.Name == "" {
		return d
	}

	//若Name的传参不为空，则返回元素名中包含Name的所有元素
	filterdList := []DataCell{}
	for _, value := range d.GenericDataList {
		//定义是否匹配的标签变量，默认是匹配
		//注意：其实这里matches没必要使用
		matches := true
		objName := value.GetName()
		if !strings.Contains(objName, d.dataSelectQuery.FilterQuery.Name) {
			matches = false
			continue
		}

		if matches {
			filterdList = append(filterdList, value)
		}
	}

	d.GenericDataList = filterdList
	return d
}

/*
分页
*/

//Paginate方法用于数组分页，根据Limit和Page参数，返回数据
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectQuery.PaginateQuery.Limit
	page := d.dataSelectQuery.PaginateQuery.Page
	//验证参数是否合法，若不合法，则返回所有数据
	if limit <= 0 || page <= 0 {
		return d
	}
	//举例：25个元素的数组，limit是10，page是3，startIndex是20，endIndex是29（实际上endIndex是24）
	startIndex := limit * (page - 1)
	endIndex := limit * page

	//处理最后一页，这时候就把endIndex由29改为24了
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}

	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	fmt.Println("分页后的数据数量：", len(d.GenericDataList))
	return d
}

//测试：定义一个podcell类型，实现GetCreation和GetName方法后，可进行数据转换
//corev1.Pod -> podCell -> DataCell
//appsv1.Deployment -> deployCel -> DataCell
type podCell corev1.Pod

//通过podCell实现了DadaCell接口
func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}
