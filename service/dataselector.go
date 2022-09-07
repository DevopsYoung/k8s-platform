package service

//dataSelector用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	dataSelectQuery *DataSelectQuery
}

//DataCell接口，用于各种资源list的类型转换，转换后可以使用dataSelector定义的方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

//DataSelectQuery 定义过滤和分页的属性，过滤: Name, 分页: Limit和Page
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
}
