package model

type DashboardStatsResponse struct {
	TotalProducts       int `json:"total_products"`
	TotalCustomers      int `json:"total_customers"`
	TotalInventoryItems int `json:"total_inventory_items"`
	LowStockProducts    int `json:"low_stock_products"`
}
