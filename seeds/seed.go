package main

import (
	"os"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/config"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	_ = config.LoadConfig()

	dbCfg := config.DefaultDBConfig()
	db, err := config.InitDB(dbCfg)
	if err != nil {
		color.Red("InitDB Failed: %v", err)
		os.Exit(1)
	}

	color.Cyan("Seeding data...")

	seedBranch(db)
	seedUsers(db)
	seedCategories(db)
	seedMenus(db)
	seedRawMaterials(db) // FIX: Tambah raw materials
	seedRecipes(db)      // FIX: Tambah recipes untuk IsRecipeBased menus
	seedTaxConfigs(db)
	seedHistoricalAnalytics(db) // Seed 30-day comprehensive historical sales data
	seedAuditLogs(db)

	color.Green("Seeding completed successfully!")
}

func seedBranch(db *gorm.DB) {
	branch := models.Branch{
		Name:    "KopiTiam Pusat",
		Address: "Jl. Sudirman No. 1, Jakarta",
		Phone:   "021-12345678",
	}
	if err := db.FirstOrCreate(&branch, models.Branch{Name: "KopiTiam Pusat"}).Error; err != nil {
		color.Red("Failed to seed branch: %v", err)
		os.Exit(1)
	}
	color.Blue("Branch seeded: %v", branch.ID)
}

func seedUsers(db *gorm.DB) {
	var branch models.Branch
	db.First(&branch)

	users := []models.User{
		{
			BranchID:     branch.ID,
			Name:         "Manager Satu",
			Email:        "manager@kopitiam.id",
			PasswordHash: hashPassword("password123"),
			PIN:          hashPIN("111111"),
			Role:         models.RoleManager,
			IsActive:     true,
		},
		{
			BranchID:     branch.ID,
			Name:         "Kasir Satu",
			Email:        "cashier@kopitiam.id",
			PasswordHash: hashPassword("password123"),
			PIN:          hashPIN("222222"),
			Role:         models.RoleCashier,
			IsActive:     true,
		},
		{
			BranchID:     branch.ID,
			Name:         "Dapur Satu",
			Email:        "kitchen@kopitiam.id",
			PasswordHash: hashPassword("password123"),
			PIN:          hashPIN("333333"),
			Role:         models.RoleKitchen,
			IsActive:     true,
		},
	}

	for _, usr := range users {
		if err := db.FirstOrCreate(&usr, models.User{Email: usr.Email}).Error; err != nil {
			color.Red("Failed to seed user %s: %v", usr.Email, err)
			os.Exit(1)
		}
		color.Blue("User seeded: %s", usr.Email)
	}
}

func seedCategories(db *gorm.DB) {
	categories := []models.Category{
		{Name: "Minuman", SortOrder: 1},
		{Name: "Coffee", SortOrder: 2},
		{Name: "Tea", SortOrder: 3},
		{Name: "Japanese Tea", SortOrder: 4},
		{Name: "Makanan Ringan", SortOrder: 5},
		{Name: "Makanan Berat", SortOrder: 6},
		{Name: "Dessert", SortOrder: 7},
	}

	for _, cat := range categories {
		if err := db.FirstOrCreate(&cat, models.Category{Name: cat.Name}).Error; err != nil {
			color.Red("Failed to seed category: %v", err)
			os.Exit(1)
		}
		color.Blue("Category seeded: %s", cat.Name)
	}
}

func seedMenus(db *gorm.DB) {
	var categories []models.Category
	db.Find(&categories)

	menus := []models.Menu{
		{Name: "Kopi Susu", Price: 18000, DailyStock: 100, IsActive: true, IsRecipeBased: false, Station: string(models.StationBar)},
		{Name: "Espresso", Price: 15000, DailyStock: 100, IsActive: true, IsRecipeBased: false, Station: string(models.StationBar)},
		{Name: "Teh Tarik", Price: 16000, DailyStock: 80, IsActive: true, IsRecipeBased: false, Station: string(models.StationBar)},
		{Name: "Nasi Lemak", Price: 35000, DailyStock: 50, IsActive: true, IsRecipeBased: true, Station: string(models.StationKitchen)},
		{Name: "Roti Bakar", Price: 12000, DailyStock: 60, IsActive: true, IsRecipeBased: true, Station: string(models.StationPastry)},
		{Name: "Kaya Toast", Price: 14000, DailyStock: 60, IsActive: true, IsRecipeBased: true, Station: string(models.StationPastry)},
		{Name: "Mee Goreng", Price: 28000, DailyStock: 40, IsActive: true, IsRecipeBased: true, Station: string(models.StationKitchen)},
		{Name: "Ice Lemon Tea", Price: 14000, DailyStock: 80, IsActive: true, IsRecipeBased: false, Station: string(models.StationBar)},
		{Name: "Kopi O", Price: 12000, DailyStock: 100, IsActive: true, IsRecipeBased: false, Station: string(models.StationBar)},
		{Name: "Chendol", Price: 20000, DailyStock: 30, IsActive: true, IsRecipeBased: true, Station: string(models.StationKitchen)},
		{Name: "Matcha Latte", Price: 22000, DailyStock: 50, IsActive: true, IsRecipeBased: false, Station: string(models.StationBar)},
		{Name: "Oolong Tea", Price: 15000, DailyStock: 60, IsActive: true, IsRecipeBased: false, Station: string(models.StationBar)},
	}

	for i := range menus {
		if err := db.FirstOrCreate(&menus[i], models.Menu{Name: menus[i].Name}).Error; err != nil {
			color.Red("Failed to seed menu %s: %v", menus[i].Name, err)
			os.Exit(1)
		}
		color.Blue("Menu seeded: %s", menus[i].Name)
	}

	// Assign categories (M2M)
	for i := range menus {
		var cats []models.Category
		switch menus[i].Name {
		case "Kopi Susu", "Espresso", "Kopi O":
			db.Where("name IN ?", []string{"Minuman", "Coffee"}).Find(&cats)
		case "Teh Tarik", "Ice Lemon Tea":
			db.Where("name IN ?", []string{"Minuman", "Tea"}).Find(&cats)
		case "Matcha Latte":
			db.Where("name IN ?", []string{"Minuman", "Tea", "Japanese Tea"}).Find(&cats)
		case "Oolong Tea":
			db.Where("name IN ?", []string{"Minuman", "Tea"}).Find(&cats)
		case "Nasi Lemak", "Mee Goreng":
			db.Where("name IN ?", []string{"Makanan Berat"}).Find(&cats)
		case "Roti Bakar", "Kaya Toast":
			db.Where("name IN ?", []string{"Makanan Ringan"}).Find(&cats)
		case "Chendol":
			db.Where("name IN ?", []string{"Minuman", "Dessert"}).Find(&cats)
		}
		if len(cats) > 0 {
			db.Model(&menus[i]).Association("Categories").Replace(cats)
		}
	}
}

// FIX: Tambah seed raw materials
func seedRawMaterials(db *gorm.DB) {
	rawMaterials := []models.RawMaterial{
		{Name: "Beras Nasi Lemak", Unit: "gram", CostPerUnit: 0.05, CurrentStock: 10000, MinStockLevel: 1000},
		{Name: "Ikan Teri", Unit: "gram", CostPerUnit: 0.8, CurrentStock: 5000, MinStockLevel: 500},
		{Name: "Kacang Tanah", Unit: "gram", CostPerUnit: 0.3, CurrentStock: 3000, MinStockLevel: 300},
		{Name: "Telur", Unit: "pcs", CostPerUnit: 2000, CurrentStock: 500, MinStockLevel: 50},
		{Name: "Roti Tawar", Unit: "slice", CostPerUnit: 1500, CurrentStock: 200, MinStockLevel: 20},
		{Name: "Kaya Jam", Unit: "gram", CostPerUnit: 0.05, CurrentStock: 5000, MinStockLevel: 500},
		{Name: "Mentega", Unit: "gram", CostPerUnit: 0.1, CurrentStock: 3000, MinStockLevel: 300},
		{Name: "Mie Kuning", Unit: "gram", CostPerUnit: 0.08, CurrentStock: 8000, MinStockLevel: 800},
		{Name: "Minyak Goreng", Unit: "ml", CostPerUnit: 0.02, CurrentStock: 10000, MinStockLevel: 1000},
		{Name: "Sayur Tauge", Unit: "gram", CostPerUnit: 0.05, CurrentStock: 2000, MinStockLevel: 200},
		{Name: "Santan", Unit: "ml", CostPerUnit: 0.03, CurrentStock: 5000, MinStockLevel: 500},
		{Name: "Gula Melaka", Unit: "gram", CostPerUnit: 0.04, CurrentStock: 3000, MinStockLevel: 300},
		{Name: "Cendol", Unit: "gram", CostPerUnit: 0.06, CurrentStock: 4000, MinStockLevel: 400},
	}

	for _, rm := range rawMaterials {
		if err := db.FirstOrCreate(&rm, models.RawMaterial{Name: rm.Name}).Error; err != nil {
			color.Red("Failed to seed raw material %s: %v", rm.Name, err)
			os.Exit(1)
		}
		color.Blue("Raw material seeded: %s", rm.Name)
	}
}

// FIX: Tambah seed recipes untuk menu IsRecipeBased
func seedRecipes(db *gorm.DB) {
	// Map menu name → recipe ingredients
	recipes := []struct {
		MenuName    string
		Ingredients []struct {
			RawMaterialName string
			Quantity        float64
		}
	}{
		{
			MenuName: "Nasi Lemak",
			Ingredients: []struct {
				RawMaterialName string
				Quantity        float64
			}{
				{"Beras Nasi Lemak", 200},
				{"Ikan Teri", 30},
				{"Kacang Tanah", 20},
				{"Telur", 1},
			},
		},
		{
			MenuName: "Roti Bakar",
			Ingredients: []struct {
				RawMaterialName string
				Quantity        float64
			}{
				{"Roti Tawar", 2},
				{"Mentega", 10},
			},
		},
		{
			MenuName: "Kaya Toast",
			Ingredients: []struct {
				RawMaterialName string
				Quantity        float64
			}{
				{"Roti Tawar", 2},
				{"Kaya Jam", 20},
				{"Mentega", 5},
			},
		},
		{
			MenuName: "Mee Goreng",
			Ingredients: []struct {
				RawMaterialName string
				Quantity        float64
			}{
				{"Mie Kuning", 150},
				{"Minyak Goreng", 30},
				{"Sayur Tauge", 50},
				{"Telur", 1},
			},
		},
		{
			MenuName: "Chendol",
			Ingredients: []struct {
				RawMaterialName string
				Quantity        float64
			}{
				{"Santan", 100},
				{"Gula Melaka", 30},
				{"Cendol", 80},
			},
		},
	}

	for _, r := range recipes {
		// Find menu
		var menu models.Menu
		if err := db.Where("name = ?", r.MenuName).First(&menu).Error; err != nil {
			color.Red("Menu %s not found for recipe: %v", r.MenuName, err)
			continue
		}

		// Check if recipe already exists
		var existing models.Recipe
		if err := db.Where("menu_id = ?", menu.ID).First(&existing).Error; err == nil {
			color.Yellow("Recipe for %s already exists, skipping", r.MenuName)
			continue
		}

		// Build ingredients
		var ingredients []models.RecipeIngredient
		for _, ing := range r.Ingredients {
			var rm models.RawMaterial
			if err := db.Where("name = ?", ing.RawMaterialName).First(&rm).Error; err != nil {
				color.Red("Raw material %s not found: %v", ing.RawMaterialName, err)
				continue
			}
			ingredients = append(ingredients, models.RecipeIngredient{
				RawMaterialID: rm.ID,
				Quantity:      ing.Quantity,
			})
		}
		// Create recipe (omitting Ingredients to avoid GORM double-saving association)
		recipe := models.Recipe{
			MenuID:      menu.ID,
			Ingredients: ingredients,
		}

		if err := db.Omit("Ingredients").Create(&recipe).Error; err != nil {
			color.Red("Failed to seed recipe for %s: %v", r.MenuName, err)
			continue
		}

		// Then create ingredients manually with recipe ID set
		if len(recipe.Ingredients) > 0 {
			for i := range recipe.Ingredients {
				recipe.Ingredients[i].RecipeID = recipe.ID
			}
			if err := db.Create(&recipe.Ingredients).Error; err != nil {
				color.Red("Failed to seed ingredients for %s: %v", r.MenuName, err)
				continue
			}
		}

		color.Blue("Recipe seeded for: %s (%d ingredients)", r.MenuName, len(ingredients))
	}
}

func seedTaxConfigs(db *gorm.DB) {
	taxes := []models.TaxConfig{
		{Name: "PPN", TaxType: models.TaxTypeTax, Percentage: 0.11, IsActive: true, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "Service Charge", TaxType: models.TaxTypeService, Percentage: 0.05, IsActive: true, EffectiveFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tax := range taxes {
		if err := db.FirstOrCreate(&tax, models.TaxConfig{Name: tax.Name}).Error; err != nil {
			color.Red("Failed to seed tax config: %v", err)
			os.Exit(1)
		}
		color.Blue("Tax config seeded: %s", tax.Name)
	}
}

func hashPassword(password string) string {
	hashed, _ := utils.HashPassword(password)
	return hashed
}

func hashPIN(pin string) string {
	hashed, _ := utils.HashPIN(pin)
	return hashed
}

func seedHistoricalAnalytics(db *gorm.DB) {
	color.Cyan("Seeding historical analytics data for the last 30 days...")

	var branch models.Branch
	if err := db.First(&branch).Error; err != nil {
		color.Red("No branch found for seeding historical analytics")
		return
	}

	var cashier models.User
	if err := db.Where("role = ?", models.RoleCashier).First(&cashier).Error; err != nil {
		color.Red("No cashier found for seeding historical analytics")
		return
	}

	var menus []models.Menu
	if err := db.Find(&menus).Error; err != nil || len(menus) == 0 {
		color.Red("No menus found for seeding historical analytics")
		return
	}

	now := time.Now()

	for day := 30; day >= 0; day-- {
		// Target date
		targetDate := now.AddDate(0, 0, -day)
		
		// 1. Create a Shift
		openingTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 8, 0, 0, 0, targetDate.Location())
		closingTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 22, 0, 0, 0, targetDate.Location())
		
		shift := models.Shift{
			BranchID:          branch.ID,
			OpenedBy:          cashier.ID,
			OpeningTime:       openingTime,
			OpeningCash:       500000,
			ClosedBy:          &cashier.ID,
			ClosingTime:       &closingTime,
			ActualClosingCash: 500000, // will compute below
			Status:            "closed",
		}
		// Set created_at/updated_at to target date
		shift.CreatedAt = openingTime
		shift.UpdatedAt = closingTime
		
		if err := db.Create(&shift).Error; err != nil {
			color.Red("Failed to create shift for day -%d: %v", day, err)
			continue
		}

		// Random number of orders per day (between 5 and 15)
		numOrders := 5 + (day % 11)
		totalSalesCash := 0.0

		for oIdx := 0; oIdx < numOrders; oIdx++ {
			orderTime := openingTime.Add(time.Duration(oIdx) * time.Hour)
			
			// Decide order status: 90% paid, 10% cancelled
			status := models.OrderPaid
			if oIdx == 0 && day%3 == 0 {
				status = models.OrderCancelled
			}

			// Generate order number
			orderNumber := targetDate.Format("20060102") + "-00" + targetDate.Format("05") + string(rune('0'+oIdx))

			order := models.Order{
				OrderNumber:    orderNumber,
				BranchID:       branch.ID,
				EmployeeID:     cashier.ID,
				ShiftID:        &shift.ID,
				OrderType:      models.OrderDineIn,
				Status:         status,
				Notes:          "Seeded historical order",
			}
			order.CreatedAt = orderTime
			order.UpdatedAt = orderTime

			// Randomly select 1 to 3 items
			numItems := 1 + (oIdx % 3)
			var details []models.OrderDetail
			subtotal := 0.0

			for iIdx := 0; iIdx < numItems; iIdx++ {
				menu := menus[(oIdx+iIdx)%len(menus)]
				qty := 1 + (iIdx % 2)
				itemPrice := menu.Price
				costPrice := itemPrice * 0.4
				itemSubtotal := itemPrice * float64(qty)

				subtotal += itemSubtotal

				detail := models.OrderDetail{
					MenuID:    menu.ID,
					MenuName:  menu.Name,
					Quantity:  qty,
					Price:     itemPrice,
					CostPrice: costPrice,
					Subtotal:  itemSubtotal,
					Station:   menu.Station,
					IsVoided:  false,
				}
				detail.CreatedAt = orderTime
				detail.UpdatedAt = orderTime
				details = append(details, detail)
			}

			// Tax calculation
			taxAmount := subtotal * 0.11
			serviceCharge := subtotal * 0.05
			total := subtotal + taxAmount + serviceCharge

			order.Subtotal = subtotal
			order.TaxAmount = taxAmount
			order.ServiceCharge = serviceCharge
			order.Total = total
			order.OrderDetails = details

			if err := db.Create(&order).Error; err != nil {
				color.Red("Failed to create order for day -%d: %v", day, err)
				continue
			}

			// Create payment if paid
			if status == models.OrderPaid {
				payment := models.Payment{
					OrderID:      order.ID,
					ChangeAmount: 0,
					Status:       "success",
				}
				payment.CreatedAt = orderTime.Add(5 * time.Minute)
				payment.UpdatedAt = orderTime.Add(5 * time.Minute)

				method := models.PayCash
				if oIdx%3 == 1 {
					method = models.PayQRIS
				} else if oIdx%3 == 2 {
					method = models.PayDebitCard
				}

				split := models.PaymentSplit{
					PaymentMethod: method,
					Amount:        total,
				}
				split.CreatedAt = payment.CreatedAt
				split.UpdatedAt = payment.UpdatedAt
				payment.Splits = []models.PaymentSplit{split}

				if err := db.Create(&payment).Error; err != nil {
					color.Red("Failed to create payment for day -%d: %v", day, err)
					continue
				}

				if method == models.PayCash {
					totalSalesCash += total
				}
			}
		}

		// Update shift actual cash closing
		shift.ActualClosingCash = shift.OpeningCash + totalSalesCash
		db.Save(&shift)
	}

	color.Green("Historical analytics data seeded successfully!")
}

func seedAuditLogs(db *gorm.DB) {
	var manager models.User
	db.Where("role = ?", models.RoleManager).First(&manager)

	var cashier models.User
	db.Where("role = ?", models.RoleCashier).First(&cashier)

	logs := []models.AuditLog{
		{
			UserID:    manager.ID,
			Action:    "CREATE",
			Entity:    "menus",
			OldValue:  "null",
			NewValue:  `{"name":"Matcha Latte","price":22000,"daily_stock":100}`,
			Timestamp: time.Now().Add(-2 * time.Hour),
		},
		{
			UserID:    manager.ID,
			Action:    "UPDATE",
			Entity:    "menus",
			OldValue:  `{"name":"Matcha Latte","price":22000,"daily_stock":100}`,
			NewValue:  `{"name":"Matcha Latte","price":24000,"daily_stock":100}`,
			Timestamp: time.Now().Add(-1 * time.Hour - 30 * time.Minute),
		},
		{
			UserID:    manager.ID,
			Action:    "CREATE",
			Entity:    "employees",
			OldValue:  "null",
			NewValue:  `{"name":"Kasir Baru","role":"cashier","email":"kasir.baru@kopitiam.id"}`,
			Timestamp: time.Now().Add(-1 * time.Hour),
		},
		{
			UserID:    cashier.ID,
			Action:    "CREATE",
			Entity:    "orders",
			OldValue:  "null",
			NewValue:  `{"order_id":"146ee955-419b-4bfa-a5aa-76c8a6cf193b","total":58000,"table_number":"Meja 5"}`,
			Timestamp: time.Now().Add(-15 * time.Minute),
		},
	}

	for _, log := range logs {
		log.ID = uuid.New()
		log.CreatedAt = log.Timestamp
		log.UpdatedAt = log.Timestamp
		db.Create(&log)
	}

	color.Blue("Audit logs seeded: %d logs", len(logs))
}
