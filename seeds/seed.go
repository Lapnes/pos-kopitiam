package main

import (
	"log"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/google/uuid"
)

func main() {
	cfg := config.LoadConfig()
	db := config.InitDB(cfg)

	log.Println("Seeding Database...")

	// -------------------------------------------------------------------------
	// 1. Branch
	// -------------------------------------------------------------------------
	mainBranch := models.Branch{
		Name:    "Main Branch",
		Address: "Jl. Sudirman No. 1",
		Phone:   "08123456789",
		IsMain:  true,
	}
	db.FirstOrCreate(&mainBranch, models.Branch{Name: "Main Branch"})

	// -------------------------------------------------------------------------
	// 2. Employees — PIN codes aligned with api_runner.go expectations
	// -------------------------------------------------------------------------
	// FIX: PIN codes must match api_runner.go expectations:
	// manager: 111111, cashier: 222222, kitchen: 333333
	hashedPassword, _ := utils.HashPassword("password123")
	users := []models.Employee{
		{BranchID: mainBranch.ID, Name: "Manager", Email: "manager@kopitiam.com", Password: hashedPassword, PINCode: "111111", Role: models.RoleManager},
		{BranchID: mainBranch.ID, Name: "Cashier 1", Email: "cashier1@kopitiam.com", Password: hashedPassword, PINCode: "222222", Role: models.RoleCashier},
		{BranchID: mainBranch.ID, Name: "Kitchen 1", Email: "kitchen1@kopitiam.com", Password: hashedPassword, PINCode: "333333", Role: models.RoleKitchen},
	}
	for _, u := range users {
		// FIX: Use Count to avoid "record not found" logs
		var count int64
		db.Model(&models.Employee{}).Where("email = ?", u.Email).Count(&count)
		
		if count == 0 {
			// Create new
			if err := db.Create(&u).Error; err != nil {
				log.Printf("  [error] failed to create user %s: %v", u.Email, err)
			} else {
				log.Printf("  [ok]   created user: %s (PIN: %s)", u.Name, u.PINCode)
			}
		} else {
			// Update existing to ensure PIN and password are correct
			var existing models.Employee
			db.Where("email = ?", u.Email).First(&existing)
			
			updates := map[string]interface{}{
				"password": u.Password,
				"pin_code": u.PINCode,
				"role":     u.Role,
				"name":     u.Name,
			}
			if err := db.Model(&existing).Updates(updates).Error; err != nil {
				log.Printf("  [error] failed to update user %s: %v", u.Email, err)
			} else {
				log.Printf("  [ok]   updated user: %s (PIN: %s)", u.Name, u.PINCode)
			}
		}
	}

	// -------------------------------------------------------------------------
	// 3. Areas & Tables
	// -------------------------------------------------------------------------
	area := models.Area{BranchID: mainBranch.ID, Name: "Indoor"}
	db.FirstOrCreate(&area, models.Area{Name: "Indoor", BranchID: mainBranch.ID})

	for _, t := range []models.Table{
		{AreaID: area.ID, Name: "Table 1", Capacity: 4},
		{AreaID: area.ID, Name: "Table 2", Capacity: 2},
	} {
		db.FirstOrCreate(&t, models.Table{Name: t.Name, AreaID: area.ID})
	}

	// -------------------------------------------------------------------------
	// 4. Categories — pre-seed ALL categories once, build a lookup map.
	// -------------------------------------------------------------------------
	categoryNames := []string{
		"Coffee", "Tea", "Food", "Cold Beverages", "Snacks", "Desserts", "Juices",
	}
	categoriesMap := make(map[string]models.Category, len(categoryNames))
	for _, name := range categoryNames {
		cat := models.Category{Name: name}
		db.FirstOrCreate(&cat, models.Category{Name: name})
		categoriesMap[name] = cat
	}

	// -------------------------------------------------------------------------
	// 5. Menus — define all menu seeds up front, then loop once.
	// -------------------------------------------------------------------------
	type menuSeed struct {
		Name          string
		Price         float64
		CostPrice     float64
		DailyStock    int
		Station       models.Station
		CatNames      []string
		IsRecipeBased bool // NEW: flag for recipe-based items
	}

	menuSeeds := []menuSeed{
		{"Milk Coffee", 18000, 10000, 48, models.StationBar, []string{"Coffee"}, true},
		{"Black Coffee", 15000, 8000, 50, models.StationBar, []string{"Coffee"}, false},
		{"Pulled Tea", 15000, 8000, 39, models.StationBar, []string{"Tea"}, false},
		{"Plain Tea", 8000, 4000, 40, models.StationBar, []string{"Tea"}, false},
		{"Fried Rice", 35000, 20000, 20, models.StationKitchen, []string{"Food"}, false},
		{"Fried Noodles", 30000, 15000, 20, models.StationKitchen, []string{"Food"}, false},
		{"Iced Milk Coffee", 20000, 12000, 30, models.StationBar, []string{"Cold Beverages"}, true},
		{"French Fries", 18000, 10000, 25, models.StationKitchen, []string{"Snacks"}, false},
		{"Chocolate Pudding", 15000, 8000, 15, models.StationKitchen, []string{"Desserts"}, false},
		{"Avocado Juice", 22000, 12000, 20, models.StationBar, []string{"Juices"}, false},
	}

	menuMap := make(map[string]models.Menu) // Track created menus for recipe linking
	for _, ms := range menuSeeds {
		var count int64
		db.Model(&models.Menu{}).Where("name = ?", ms.Name).Count(&count)
		if count > 0 {
			log.Printf("  [skip] menu already exists: %s", ms.Name)
			// Still load into menuMap for recipe linking
			var existing models.Menu
			db.Where("name = ?", ms.Name).First(&existing)
			menuMap[ms.Name] = existing
			continue
		}

		cats := make([]models.Category, 0, len(ms.CatNames))
		for _, cn := range ms.CatNames {
			cats = append(cats, categoriesMap[cn])
		}

		menu := models.Menu{
			Name:          ms.Name,
			Price:         ms.Price,
			CostPrice:     ms.CostPrice,
			DailyStock:    ms.DailyStock,
			IsActive:      true,
			IsRecipeBased: ms.IsRecipeBased,
			Station:       ms.Station,
			Categories:    cats,
		}

		if err := db.Create(&menu).Error; err != nil {
			log.Printf("  [error] failed to create menu %s: %v", ms.Name, err)
		} else {
			log.Printf("  [ok]   created menu: %s", ms.Name)
			menuMap[ms.Name] = menu
		}
	}

	// -------------------------------------------------------------------------
	// 6. Raw Materials & Recipes (BOM) — NEW SECTION
	// -------------------------------------------------------------------------
	// Pre-seed a raw material and recipe for Milk Coffee so API tests work
	// without needing to create raw materials first.

	coffeeBeans := models.RawMaterial{
		Name:         "Arabica Coffee Beans",
		Unit:         "gram",
		CurrentStock: 5000,
		MinimumStock: 500,
		CostPerUnit:  150,
	}

	var existingRM models.RawMaterial
	var count int64
	db.Model(&models.RawMaterial{}).Where("name = ?", coffeeBeans.Name).Count(&count)
	
	if count == 0 {
		// Create new
		if err := db.Create(&coffeeBeans).Error; err != nil {
			log.Printf("  [error] failed to create raw material: %v", err)
		} else {
			log.Printf("  [ok]   created raw material: %s (ID: %s)", coffeeBeans.Name, coffeeBeans.ID)
			existingRM = coffeeBeans
		}
	} else {
		db.Where("name = ?", coffeeBeans.Name).First(&existingRM)
		log.Printf("  [skip] raw material already exists: %s", coffeeBeans.Name)
	}

	// Create recipe for Milk Coffee if it doesn't exist
	milkCoffeeMenu := menuMap["Milk Coffee"]
	if milkCoffeeMenu.ID != uuid.Nil && milkCoffeeMenu.IsRecipeBased {
		var recipeCount int64
		db.Model(&models.Recipe{}).Where("menu_id = ?", milkCoffeeMenu.ID).Count(&recipeCount)

		if recipeCount == 0 {
			// Create recipe with ingredients
			recipe := models.Recipe{
				MenuID:       milkCoffeeMenu.ID,
				Instructions: "1. Grind 20g Arabica beans\n2. Brew double espresso\n3. Steam 100ml milk\n4. Combine",
				Ingredients: []models.RecipeItem{
					{
						RawMaterialID: existingRM.ID,
						Quantity:      20,
					},
				},
			}

			if err := db.Create(&recipe).Error; err != nil {
				log.Printf("  [error] failed to create recipe: %v", err)
			} else {
				log.Printf("  [ok]   created recipe for: %s (ID: %s)", milkCoffeeMenu.Name, recipe.ID)
			}
		} else {
			log.Printf("  [skip] recipe already exists for: %s", milkCoffeeMenu.Name)
		}
	}

	// Create recipe for Iced Milk Coffee if it doesn't exist
	icedMilkCoffeeMenu := menuMap["Iced Milk Coffee"]
	if icedMilkCoffeeMenu.ID != uuid.Nil && icedMilkCoffeeMenu.IsRecipeBased {
		var recipeCount int64
		db.Model(&models.Recipe{}).Where("menu_id = ?", icedMilkCoffeeMenu.ID).Count(&recipeCount)

		if recipeCount == 0 {
			recipe := models.Recipe{
				MenuID:       icedMilkCoffeeMenu.ID,
				Instructions: "1. Grind 20g Arabica beans\n2. Brew double espresso\n3. Add ice\n4. Pour milk",
				Ingredients: []models.RecipeItem{
					{
						RawMaterialID: existingRM.ID,
						Quantity:      20,
					},
				},
			}

			if err := db.Create(&recipe).Error; err != nil {
				log.Printf("  [error] failed to create recipe for iced milk coffee: %v", err)
			} else {
				log.Printf("  [ok]   created recipe for: %s (ID: %s)", icedMilkCoffeeMenu.Name, recipe.ID)
			}
		} else {
			log.Printf("  [skip] recipe already exists for: %s", icedMilkCoffeeMenu.Name)
		}
	}

	log.Println("Seeding Completed!")
}
