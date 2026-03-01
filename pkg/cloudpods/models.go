package cloudpods

// 区域信息
// Region 代表云平台的区域
// json标签用于API序列化
// 可根据实际API返回字段调整

type Region struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

// VPC信息
// VPC 代表虚拟私有云

type VPC struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	RegionID    string `json:"region_id"`
	CIDR        string `json:"cidr"`
	Description string `json:"description,omitempty"`
}

// Server信息
// Server 代表虚拟机

type Server struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	RegionID string   `json:"region_id"`
	VPCID    string   `json:"vpc_id"`
	IPs      []string `json:"ips,omitempty"`
	CPU      int      `json:"cpu,omitempty"`
	MemoryMB int      `json:"memory_mb,omitempty"`
}

// ListServerOptions 用于查询虚拟机列表的参数

type ListServerOptions struct {
	RegionID string `json:"region_id"`
	Status   string `json:"status"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}
