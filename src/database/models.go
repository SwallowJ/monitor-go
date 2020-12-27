package database

import "time"

//CPUInfo 信息
type CPUInfo struct {
	Cors       int       `json:"cors"`      //内核
	Logical    int       `json:"logical"`   //逻辑处理器
	Percent    float64   `json:"percent"`   //利用率
	CPU        string    `json:"cpu"`       //cpu
	User       float64   `json:"user"`      //用户态空间运行时间
	System     float64   `json:"system"`    //内核空间运行时间
	Idle       float64   `json:"idle"`      //空闲时间
	Nice       float64   `json:"nice"`      //用户空间进程的CPU的调度优先级
	Iowait     float64   `json:"iowait"`    //读写等待状态时间
	Irq        float64   `json:"irq"`       //硬中断的CPU时间
	Softirq    float64   `json:"softirq"`   //软中断的CPU时间
	Steal      float64   `json:"steal"`     //其他虚拟机占用的CPU时间
	Guest      float64   `json:"guest"`     //拟化运行其他操作系统的时间
	GuestNice  float64   `json:"guestNice"` //低优先级运行虚拟机的时间
	CreateTime time.Time `json:"createTime"`
}

//DiskInfo 磁盘信息
type DiskInfo struct {
	Path        string    `json:"path"`
	Fstype      string    `json:"fstype"`
	Total       uint64    `json:"total"`
	Free        uint64    `json:"free"`
	Used        uint64    `json:"used"`
	Speed       float64   `json:"speed"`
	UsedPercent float64   `json:"usedPercent"`
	CreateTime  time.Time `json:"createTime"`
}

//MemoryInfo 内存信息
type MemoryInfo struct {
	Total        uint64    `json:"total"`       //系统共可用内存
	Available    uint64    `json:"available"`   //可以内存
	Used         uint64    `json:"used"`        //已用内存
	UsedPercent  float64   `json:"usedPercent"` //使用百分比
	Free         uint64    `json:"free"`        //空闲内存
	Active       uint64    `json:"active"`
	Inactive     uint64    `json:"inactive"`
	Wired        uint64    `json:"wired"`
	Buffers      uint64    `json:"buffers"`
	Cached       uint64    `json:"cached"`
	ActiveFile   uint64    `json:"activefile"`
	InactiveFile uint64    `json:"inactivefile"`
	ActiveAnon   uint64    `json:"activeanon"`
	InactiveAnon uint64    `json:"inactiveanon"`
	Unevictable  uint64    `json:"unevictable"`
	CreateTime   time.Time `json:"createTime"`
}

//NetInfo 网络信息
type NetInfo struct {
	Name        string    `json:"name"`        // name
	BytesSent   uint64    `json:"bytesSent"`   // 上传数据
	BytesRecv   uint64    `json:"bytesRecv"`   // 下载数据
	PacketsSent uint64    `json:"packetsSent"` // 上传数据包
	PacketsRecv uint64    `json:"packetsRecv"` // 下载数据包
	Errin       uint64    `json:"errin"`       // 接收到的错误
	Errout      uint64    `json:"errout"`      // 发送时错误
	Dropin      uint64    `json:"dropin"`      // 被丢弃的入站数据包总数
	Dropout     uint64    `json:"dropout"`     // 被丢弃的传出数据包总数（在OSX和BSD上始终为0）
	Fifoin      uint64    `json:"fifoin"`      // 接收时FIFO缓冲区错误总数
	Fifoout     uint64    `json:"fifoout"`     // 发送时FIFO缓冲区错误总数
	SentSpeed   float64   `json:"sentSpeed"`   //上传速率
	RecvSpeed   float64   `json:"recvSpeed"`   //下载速率
	CreateTime  time.Time `json:"createTime"`
}
