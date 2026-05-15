package models

type Role string

// TASK 5: Only 3 core roles are supported.
const (
	RoleManager Role = "manager"
	RoleCashier Role = "cashier"
	RoleKitchen Role = "kitchen"
)

type Station string

const (
	StationKitchen Station = "kitchen"
	StationBar     Station = "bar"
	StationPastry  Station = "pastry"
	StationNone    Station = "none"
)

type TableStatus string

const (
	TableAvailable TableStatus = "available"
	TableOccupied  TableStatus = "occupied"
	TableReserved  TableStatus = "reserved"
	TableCleaning  TableStatus = "cleaning"
)

type OrderType string

const (
	OrderDineIn   OrderType = "dine_in"
	OrderTakeaway OrderType = "takeaway"
	OrderDelivery OrderType = "delivery"
	OrderOnline   OrderType = "online"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderConfirmed OrderStatus = "confirmed"
	OrderCooking   OrderStatus = "cooking"
	OrderReady     OrderStatus = "ready"
	OrderServed    OrderStatus = "served"
	OrderPaid      OrderStatus = "paid"
	OrderCancelled OrderStatus = "cancelled"
	OrderRefunded  OrderStatus = "refunded"
)

type KDSStatus string

const (
	KDSPending KDSStatus = "pending"
	KDSCooking KDSStatus = "cooking"
	KDSReady   KDSStatus = "ready"
	KDSServed  KDSStatus = "served"
)

type PaymentMethod string

const (
	PayCash       PaymentMethod = "cash"
	PayQRIS       PaymentMethod = "qris"
	PayDebitCard  PaymentMethod = "debit_card"
	PayCreditCard PaymentMethod = "credit_card"
	PayTransfer   PaymentMethod = "transfer"
	PayEWallet    PaymentMethod = "e_wallet"
	PaySplit      PaymentMethod = "split"
)

type DiscountType string

const (
	DiscountPercentage  DiscountType = "percentage"
	DiscountFixedAmount DiscountType = "fixed_amount"
	DiscountBuyXGetY    DiscountType = "buy_x_get_y"
)
