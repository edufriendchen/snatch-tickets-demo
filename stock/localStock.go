package stock

// LocalStock 本地仓库
type LocalStock struct {
	LocalServerNumber string //本地库存的唯一编号 （本地服务器唯一编号）
	LocalTicketStock  int64  //本地票库存量
	LocalSalesVolume  int64  //本地售出量
}

// LocalDeductionStock 本地扣库存,返回bool值
func (spike *LocalStock) LocalDeductionStock() bool {
	spike.LocalSalesVolume = spike.LocalSalesVolume + 1
	return spike.LocalSalesVolume <= spike.LocalTicketStock
}
