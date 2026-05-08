package bigcapital

type JournalEntry struct {
	Date        string               `json:"date"` // YYYY-MM-DD
	Description string               `json:"description"`
	Reference   string               `json:"reference"`
	Entries     []JournalEntryDetail `json:"entries"`
}

type JournalEntryDetail struct {
	AccountId   string  `json:"accountId"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Description string  `json:"description"`
}

type Item struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CostPrice   float64 `json:"costPrice"`
	Active      bool    `json:"active"`
}

type Customer struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}
