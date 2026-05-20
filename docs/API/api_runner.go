package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const baseURL = "http://127.0.0.1:8080"

// TestResult holds a single API call outcome for the recap report.
type TestResult struct {
	Phase    string `json:"phase"`
	Name     string `json:"name"`
	Method   string `json:"method"`
	URL      string `json:"url"`
	Status   int    `json:"status"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
	Response string `json:"response,omitempty"`
	Duration string `json:"duration"`
}

// RecapData is the root object written to the JSON/Markdown report files.
type RecapData struct {
	Date      string            `json:"date"`
	Total     int               `json:"total"`
	Passed    int               `json:"passed"`
	Failed    int               `json:"failed"`
	Results   []TestResult      `json:"results"`
	Variables map[string]string `json:"variables"`
}

var (
	results   []TestResult
	variables = make(map[string]string)
	client    = &http.Client{Timeout: 10 * time.Second}
)

// ─────────────────────────────────────────────────────────────────────────────
// main — Automated test sequence
// Roles covered: manager, cashier, kitchen
// Payment model: normalized — splits[] carries per-method detail
// ─────────────────────────────────────────────────────────────────────────────
func main() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("   POS KopiTiam — Automated API Test Runner")
	fmt.Println("   Roles: Manager | Cashier | Kitchen")
	fmt.Println("═══════════════════════════════════════════════════")

	// ── Phase 1: Health & Auth ──────────────────────────────────────────────
	runTest("1.1", "Health Check", "GET", "/health", "", nil)

	// Manager login (email + password)
	loginManager := runTest("1.2", "Login Manager", "POST", "/api/v1/auth/login",
		`{"email":"manager@kopitiam.id","password":"password123"}`, nil)
	extractToken(loginManager, "manager_token")
	extractUserData(loginManager, "manager_id", "manager_branch_id")

	// Cashier login (PIN — fast terminal login)
	loginCashier := runTest("1.3", "Login Cashier (PIN)", "POST", "/api/v1/auth/login-pin",
		`{"pin":"222222"}`, nil)
	extractToken(loginCashier, "cashier_token")
	extractUserData(loginCashier, "cashier_id", "cashier_branch_id")

	// Kitchen login (PIN)
	loginKitchen := runTest("1.4", "Login Kitchen (PIN)", "POST", "/api/v1/auth/login-pin",
		`{"pin":"333333"}`, nil)
	extractToken(loginKitchen, "kitchen_token")

	// ── Phase 2: Menus ─────────────────────────────────────────────────────
	menus := runTest("2.1", "Get Active Menus (M2M categories)", "GET", "/api/v1/menus", "",
		header("manager_token"))
	extractMenuIDs(menus)

	// ── Phase 3: Create & Confirm Order (Dine In) ──────────────────────────
	orderPayload := fmt.Sprintf(`{
		"table_id": "",
		"customer_id": "",
		"order_type": "dine_in",
		"notes": "Less sugar, extra ice",
		"discount_amount": 0,
		"items": [
			{
				"menu_id": "%s",
				"menu_name": "%s",
				"quantity": 2,
				"price": %s,
				"notes": "Hot"
			},
			{
				"menu_id": "%s",
				"menu_name": "%s",
				"quantity": 1,
				"price": %s,
				"notes": ""
			}
		]
	}`,
		variables["menu0_id"], variables["menu0_name"], variables["menu0_price"],
		variables["menu4_id"], variables["menu4_name"], variables["menu4_price"],
	)

	createOrder := runTest("3.1", "Create Order (Dine In)", "POST", "/api/v1/orders",
		orderPayload, header("cashier_token"))
	extractOrderData(createOrder)

	// Guard: if order creation failed, skip all dependent tests
	if variables["order_id"] == "" {
		skipTest("3.2", "Get All Orders — SKIPPED (order creation failed)")
		skipTest("3.3", "Get Order by ID — SKIPPED")
		skipTest("3.4", "Get Order by Number — SKIPPED")
		skipTest("3.5", "Confirm Order — SKIPPED")
		skipTest("3.6", "Update Order Notes — SKIPPED")
		skipTest("4.1", "Process Payment (Cash) — SKIPPED")
		skipTest("6.1", "Process Return — SKIPPED")
		skipTest("12.1", "Generate Receipt — SKIPPED")
		skipTest("12.2", "Generate Kitchen Ticket — SKIPPED")
		fmt.Println("\n ⚠️  WARNING: Order creation failed. Skipping all order-dependent tests.")
		fmt.Println("   Tip: Run 'go run cmd/migrate/main.go && go run seeds/seed.go' then retry.")
	} else {
		runTest("3.2", "Get All Orders", "GET", "/api/v1/orders", "", header("manager_token"))
		runTest("3.3", "Get Order by ID", "GET", "/api/v1/orders/"+variables["order_id"], "", header("manager_token"))
		runTest("3.4", "Get Order by Number", "GET", "/api/v1/orders/by-number/"+variables["order_number"], "", header("cashier_token"))

		// Confirm: branches on menu.IsRecipeBased — uses pessimistic lock if recipe-based
		confirm := runTest("3.5", "Confirm Order (stock deduction)", "PATCH",
			"/api/v1/orders/"+variables["order_id"]+"/confirm", "", header("cashier_token"))
		extractField(confirm, "status", "confirmed_status")

		runTest("3.6", "Update Order Notes", "PATCH", "/api/v1/orders/"+variables["order_id"]+"/notes",
			`{"notes":"Customer requests extra napkins"}`, header("cashier_token"))
	}

	// ── Phase 4: Payment (normalized — splits array) ──────────────────────
	if variables["order_id"] != "" {
		orderTotal := parseFloat(variables["order_total"])
		cashTendered := math.Ceil(orderTotal/10000) * 10000 // round up to nearest 10k
		change := cashTendered - orderTotal

		// 4.1 — Single method: full cash
		cashPayload := fmt.Sprintf(`{
		"change_amount": %.0f,
		"splits": [
			{
				"payment_method": "cash",
				"amount": %.0f
			}
		]
	}`, change, cashTendered)
		runTest("4.1", "Process Payment — Cash (single split)", "POST",
			"/api/v1/payments?order_id="+variables["order_id"],
			cashPayload, header("cashier_token"))
	} else {
		skipTest("4.1", "Process Payment (Cash) — SKIPPED (no order_id)")
	}

	// ── Phase 5: Takeaway Order → Cancel / Void ───────────────────────────
	time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)

	takeawayPayload := fmt.Sprintf(`{
		"table_id": "",
		"customer_id": "",
		"order_type": "takeaway",
		"notes": "Packaging separate",
		"discount_amount": 5000,
		"items": [
			{
				"menu_id": "%s",
				"menu_name": "%s",
				"quantity": 1,
				"price": %s,
				"notes": ""
			}
		]
	}`, variables["menu1_id"], variables["menu1_name"], variables["menu1_price"])

	createTakeaway := runTest("5.1", "Create Order (Takeaway)", "POST", "/api/v1/orders",
		takeawayPayload, header("cashier_token"))
	extractTakeawayData(createTakeaway)

	// If takeaway creation failed, skip dependent tests
	if variables["takeaway_id"] == "" {
		skipTest("5.2", "Process Payment — Split (Cash + QRIS) — SKIPPED (takeaway creation failed)")
		skipTest("5.3", "Cancel Order (cashier) — SKIPPED (no takeaway_id)")
		if variables["order_id"] != "" {
			voidPayload := `{"reason": "Customer walked away"}`
			runTest("5.4", "Void Order (manager only)", "PATCH",
				"/api/v1/orders/"+variables["order_id"]+"/void",
				voidPayload, header("manager_token"))
		} else {
			skipTest("5.4", "Void Order — SKIPPED (no order_id)")
		}
	} else {
		// 5.2 — Split payment on takeaway order (Cash + QRIS)
		takeawayTotal := parseFloat(variables["takeaway_total"])
		splitPayload := fmt.Sprintf(`{
			"change_amount": 0,
			"splits": [
				{
					"payment_method": "cash",
					"amount": %.0f
				},
				{
					"payment_method": "qris",
					"amount": %.0f,
					"reference_number": "QRIS-%s-001"
				}
			]
		}`, takeawayTotal/2, takeawayTotal-takeawayTotal/2,
			time.Now().Format("20060102"))
		runTest("5.2", "Process Payment — Split (Cash + QRIS)", "POST",
			"/api/v1/payments?order_id="+variables["takeaway_id"],
			splitPayload, header("cashier_token"))

		// 5.3 — Cancel the takeaway (cashier)
		time.Sleep(100 * time.Millisecond)
		cancelOrder := runTest("5.3a", "Create Order for Cancel", "POST", "/api/v1/orders",
			takeawayPayload, header("cashier_token"))

		cancelOrderID := ""
		if cancelOrder != nil && cancelOrder["data"] != nil {
			if dataMap, ok := cancelOrder["data"].(map[string]interface{}); ok {
				if id, ok := dataMap["id"].(string); ok {
					cancelOrderID = id
				}
			}
		}

		if cancelOrderID != "" {
			runTest("5.3", "Cancel Order (cashier)", "PATCH",
				"/api/v1/orders/"+cancelOrderID+"/cancel", "", header("cashier_token"))
		} else {
			skipTest("5.3", "Cancel Order (cashier) — SKIPPED (no cancel_order_id)")
		}

		// 5.4 — Void the takeaway order (manager only)
		time.Sleep(100 * time.Millisecond)
		voidOrder := runTest("5.4a", "Create Order for Void", "POST", "/api/v1/orders",
			takeawayPayload, header("cashier_token"))

		voidOrderID := ""
		if voidOrder != nil && voidOrder["data"] != nil {
			if dataMap, ok := voidOrder["data"].(map[string]interface{}); ok {
				if id, ok := dataMap["id"].(string); ok {
					voidOrderID = id
				}
			}
		}

		if voidOrderID != "" {
			voidPayload := `{"reason": "Customer walked away"}`
			runTest("5.4", "Void Order (manager only)", "PATCH",
				"/api/v1/orders/"+voidOrderID+"/void",
				voidPayload, header("manager_token"))
		} else {
			skipTest("5.4", "Void Order (manager only) — SKIPPED (no void_order_id)")
		}
	}

	// ── Phase 6: Returns ──────────────────────────────────────────────────
	returnPayload := `{"return_amount": 18000, "reason": "Wrong order — customer wanted tea not coffee"}`
	if variables["order_id"] != "" {
		runTest("6.1", "Process Return (Cashier)", "POST",
			"/api/v1/returns?order_id="+variables["order_id"]+"&processed_by="+variables["cashier_id"],
			returnPayload, header("cashier_token"))
	} else {
		skipTest("6.1", "Process Return — SKIPPED (no order_id)")
	}

	// ── Phase 7: Shifts ───────────────────────────────────────────────────
	openShift := runTest("7.1", "Open Shift", "POST", "/api/v1/shifts/open",
		`{"opening_cash": 500000}`, header("cashier_token"))
	extractShiftData(openShift)

	runTest("7.2", "Get Current Shift", "GET", "/api/v1/shifts/current", "", header("cashier_token"))

	closeShift := runTest("7.3", "Close Shift", "PATCH", "/api/v1/shifts/close",
		`{"actual_closing_cash": 525000}`, header("cashier_token"))
	extractCloseShiftData(closeShift)

	// ── Phase 8: Raw Materials (Manager Only) ────────────────────────────
	rmPayload := `{
		"name": "Arabica Coffee Beans",
		"unit": "gram",
		"current_stock": 5000,
		"min_stock_level": 500,
		"cost_per_unit": 150
	}`
	createRM := runTest("8.1", "Create Raw Material", "POST", "/api/v1/inventory/raw-materials",
		rmPayload, header("manager_token"))
	extractField(createRM, "id", "rm_id")
	extractField(createRM, "name", "rm_name")

	runTest("8.2", "Get Raw Material", "GET",
		"/api/v1/inventory/raw-materials/"+variables["rm_id"], "", header("manager_token"))

	updateRMPayload := `{
		"name": "Arabica Coffee Beans (Premium)",
		"unit": "gram",
		"current_stock": 4500,
		"min_stock_level": 500,
		"cost_per_unit": 175
	}`
	runTest("8.3", "Update Raw Material", "PUT",
		"/api/v1/inventory/raw-materials/"+variables["rm_id"], updateRMPayload, header("manager_token"))

	// Create a separate raw material for deletion test
	deleteRMPayload := `{
		"name": "Test Delete Material",
		"unit": "gram",
		"current_stock": 100,
		"min_stock_level": 10,
		"cost_per_unit": 50
	}`
	createDeleteRM := runTest("8.4a", "Create Raw Material for Deletion", "POST", "/api/v1/inventory/raw-materials",
		deleteRMPayload, header("manager_token"))
	extractField(createDeleteRM, "id", "rm_delete_id")

	runTest("8.4", "Delete Raw Material", "DELETE",
		"/api/v1/inventory/raw-materials/"+variables["rm_delete_id"], "", header("manager_token"))

	// ── Phase 9: Recipes / BOM (Manager Only) ────────────────────────────
	recipePayload := fmt.Sprintf(`{
		"menu_id": "%s",
		"ingredients": [
			{
				"raw_material_id": "%s",
				"quantity": 20
			}
		]
	}`, variables["menu0_id"], variables["rm_id"])

	// 9.0 — Create Recipe (always create)
	createRecipe := runTest("9.0", "Create Recipe (BOM)", "POST", "/api/v1/inventory/recipes",
		recipePayload, header("manager_token"))
	extractField(createRecipe, "id", "recipe_id")

	// Guard: if recipe creation failed, skip dependent tests
	if variables["recipe_id"] == "" {
		skipTest("9.1", "Get Recipe by Menu ID — SKIPPED (recipe creation failed)")
		skipTest("9.2", "Update Recipe — SKIPPED (no recipe_id)")
		skipTest("9.3", "Delete Recipe — SKIPPED (no recipe_id)")
		fmt.Println("\n ⚠️  WARNING: Recipe creation failed. Skipping recipe-dependent tests.")
	} else {
		// 9.1 — Get Recipe by Menu ID
		runTest("9.1", "Get Recipe by Menu ID", "GET",
			"/api/v1/inventory/recipes/"+variables["menu0_id"], "", header("manager_token"))

		// 9.2 — Update Recipe
		updateRecipePayload := fmt.Sprintf(`{
			"menu_id": "%s",
			"ingredients": [
				{
					"raw_material_id": "%s",
					"quantity": 25
				}
			]
		}`, variables["menu0_id"], variables["rm_id"])
		runTest("9.2", "Update Recipe", "PUT",
			"/api/v1/inventory/recipes/"+variables["recipe_id"], updateRecipePayload, header("manager_token"))

		// 9.3 — Delete Recipe
		runTest("9.3", "Delete Recipe", "DELETE",
			"/api/v1/inventory/recipes/"+variables["recipe_id"], "", header("manager_token"))
	}

	// ── Phase 10: Employees (Manager Only) ───────────────────────────────
	runTest("10.1", "Get All Employees", "GET", "/api/v1/employees", "", header("manager_token"))

	empPayload := fmt.Sprintf(`{
		"branch_id": "%s",
		"name": "Test Cashier API",
		"email": "testcashier_api_%d@kopitiam.id",
		"password": "password123",
		"pin_code": "555555",
		"role": "cashier",
		"is_active": true
	}`, variables["manager_branch_id"], time.Now().Unix())

	createEmp := runTest("10.2", "Create Employee (cashier)", "POST", "/api/v1/employees",
		empPayload, header("manager_token"))
	extractField(createEmp, "id", "emp_id")

	updateEmpPayload := `{
		"name": "Test Cashier (Updated)",
		"role": "cashier",
		"is_active": true
	}`
	runTest("10.3", "Update Employee", "PUT",
		"/api/v1/employees/"+variables["emp_id"], updateEmpPayload, header("manager_token"))
	runTest("10.4", "Delete Employee", "DELETE",
		"/api/v1/employees/"+variables["emp_id"], "", header("manager_token"))

	// ── Phase 11: Analytics (Manager Only) ───────────────────────────────
	runTest("11.1", "Sales Summary (this month)", "GET",
		"/api/v1/analytics/sales-summary?start_date=2026-05-01&end_date=2026-05-31", "",
		header("manager_token"))
	runTest("11.2", "Best Sellers (top 5)", "GET",
		"/api/v1/analytics/best-sellers?limit=5", "", header("manager_token"))
	runTest("11.3", "Return Impact", "GET",
		"/api/v1/analytics/return-impact", "", header("manager_token"))
	runTest("11.4", "COGS", "GET",
		"/api/v1/analytics/cogs?start_date=2026-05-01&end_date=2026-05-31", "", header("manager_token"))

	// ── Phase 12: Printers ────────────────────────────────────────────────
	if variables["order_id"] != "" {
		runTest("12.1", "Generate Receipt (splits-aware)", "GET",
			"/api/v1/printer/receipt/"+variables["order_id"], "", header("cashier_token"))
		runTest("12.2", "Generate Kitchen Ticket", "GET",
			"/api/v1/printer/kitchen-ticket/"+variables["order_id"], "", header("kitchen_token"))
	} else {
		skipTest("12.1", "Generate Receipt — SKIPPED (no order_id)")
		skipTest("12.2", "Generate Kitchen Ticket — SKIPPED (no order_id)")
	}

	// ── Phase 13: Catalog self-service (Basic Auth) ─────────────────────
	catalogHeaders := map[string]string{
		"Authorization": "Basic Y2F0YWxvZzprb3BpdGlhbTEyMw==", // "catalog:kopitiam123"
	}
	catalogMenus := runTest("13.1", "Get Catalog Menus (Basic Auth)", "GET", "/catalog/api/menus", "", catalogHeaders)
	extractCatalogMenuID(catalogMenus)

	if variables["catalog_menu_id"] != "" {
		checkoutPayload := fmt.Sprintf(`{
			"items": [
				{
					"menu_id": "%s",
					"quantity": 1,
					"notes": "no sugar"
				}
			],
			"payment_method": "qris",
			"table_number": "5"
		}`, variables["catalog_menu_id"])
		runTest("13.2", "Catalog Self-Checkout (QRIS)", "POST", "/catalog/api/checkout", checkoutPayload, catalogHeaders)
	} else {
		skipTest("13.2", "Catalog Self-Checkout — SKIPPED (no catalog_menu_id)")
	}

	// ── Final Report ─────────────────────────────────────────────────────
	generateRecap()
}

// ─────────────────────────────────────────────────────────────────────────────
// HTTP helpers
// ─────────────────────────────────────────────────────────────────────────────

func runTest(phase, name, method, path, body string, headers map[string]string) map[string]interface{} {
	url := baseURL + path
	start := time.Now()

	var req *http.Request
	var err error

	if body != "" {
		req, err = http.NewRequest(method, url, bytes.NewBufferString(body))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		recordResult(phase, name, method, path, 0, false, err.Error(), "", "0s")
		return nil
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		recordResult(phase, name, method, path, 0, false, err.Error(), "", duration.String())
		return nil
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var data map[string]interface{}
	json.Unmarshal(respBody, &data)

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		errMsg := string(respBody)
		if msg, ok := data["message"].(string); ok {
			errMsg = msg
		}
		recordResult(phase, name, method, path, resp.StatusCode, false, errMsg, string(respBody), duration.String())
	} else {
		recordResult(phase, name, method, path, resp.StatusCode, true, "", string(respBody), duration.String())
	}

	return data
}

func recordResult(phase, name, method, url string, status int, success bool, errMsg, response, duration string) {
	icon := "✅"
	if !success {
		icon = "❌"
	}
	fmt.Printf("%s [%s] %s %s (HTTP %d) %s\n", icon, phase, method, url, status, duration)
	if errMsg != "" {
		fmt.Printf("   Error: %s\n", truncate(errMsg, 120))
	}
	results = append(results, TestResult{
		Phase:    phase,
		Name:     name,
		Method:   method,
		URL:      url,
		Status:   status,
		Success:  success,
		Error:    errMsg,
		Response: truncate(response, 200),
		Duration: duration,
	})
}

func header(tokenKey string) map[string]string {
	token := variables[tokenKey]
	if token == "" {
		return nil
	}
	return map[string]string{"Authorization": "Bearer " + token}
}

// skipTest records a skipped test — it neither passes nor fails, shown as ⏭ in the report.
func skipTest(phase, name string) {
	fmt.Printf("⏭  [%s] SKIP: %s\n", phase, name)
	results = append(results, TestResult{
		Phase:   phase,
		Name:    name,
		Method:  "SKIP",
		URL:     "-",
		Status:  0,
		Success: false,
		Error:   "skipped: dependency failed",
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Variable extractors
// ─────────────────────────────────────────────────────────────────────────────

func extractToken(data map[string]interface{}, key string) {
	if data == nil {
		return
	}
	if d, ok := data["data"].(map[string]interface{}); ok {
		if token, ok := d["access_token"].(string); ok {
			variables[key] = token
		}
	}
}

func extractUserData(data map[string]interface{}, idKey, branchKey string) {
	if data == nil {
		return
	}
	if d, ok := data["data"].(map[string]interface{}); ok {
		if user, ok := d["user"].(map[string]interface{}); ok {
			if id, ok := user["id"].(string); ok {
				variables[idKey] = id
			}
			if bid, ok := user["branch_id"].(string); ok {
				variables[branchKey] = bid
			}
		}
	}
}

// extractMenuIDs fetches the first 5 menus (indices 0–4) into variables.
func extractMenuIDs(data map[string]interface{}) {
	if data == nil {
		return
	}
	if d, ok := data["data"].([]interface{}); ok && len(d) > 0 {
		for i := 0; i < len(d) && i < 5; i++ {
			if menu, ok := d[i].(map[string]interface{}); ok {
				if id, ok := menu["id"].(string); ok {
					variables[fmt.Sprintf("menu%d_id", i)] = id
				}
				if name, ok := menu["name"].(string); ok {
					variables[fmt.Sprintf("menu%d_name", i)] = name
				}
				if price, ok := menu["price"].(float64); ok {
					variables[fmt.Sprintf("menu%d_price", i)] = fmt.Sprintf("%.0f", price)
				}
			}
		}
	}
}

func extractOrderData(data map[string]interface{}) {
	if data == nil {
		return
	}
	if d, ok := data["data"].(map[string]interface{}); ok {
		extractFieldFromMap(d, "id", "order_id")
		extractFieldFromMap(d, "order_number", "order_number")
		extractFieldFromMap(d, "total", "order_total")
		extractFieldFromMap(d, "status", "order_status")
	}
}

func extractTakeawayData(data map[string]interface{}) {
	if data == nil {
		return
	}
	if d, ok := data["data"].(map[string]interface{}); ok {
		extractFieldFromMap(d, "id", "takeaway_id")
		extractFieldFromMap(d, "total", "takeaway_total")
	}
}

func extractShiftData(data map[string]interface{}) {
	if data == nil {
		return
	}
	if d, ok := data["data"].(map[string]interface{}); ok {
		extractFieldFromMap(d, "id", "shift_id")
		extractFieldFromMap(d, "status", "shift_status")
	}
}

func extractCloseShiftData(data map[string]interface{}) {
	if data == nil {
		return
	}
	if d, ok := data["data"].(map[string]interface{}); ok {
		extractFieldFromMap(d, "variance", "shift_variance")
		extractFieldFromMap(d, "expected_cash", "shift_expected")
		extractFieldFromMap(d, "total_sales", "shift_sales")
	}
}

func extractField(data map[string]interface{}, field, key string) {
	if data == nil {
		return
	}
	if d, ok := data["data"].(map[string]interface{}); ok {
		extractFieldFromMap(d, field, key)
	}
}

func extractFieldFromMap(data map[string]interface{}, field, key string) {
	if val, ok := data[field]; ok {
		switch v := val.(type) {
		case string:
			variables[key] = v
		case float64:
			variables[key] = fmt.Sprintf("%.2f", v)
		case bool:
			variables[key] = fmt.Sprintf("%t", v)
		}
	}
}

// parseFloat safely parses a string stored in variables to float64.
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// ─────────────────────────────────────────────────────────────────────────────
// Recap / Report generation
// ─────────────────────────────────────────────────────────────────────────────

func generateRecap() {
	passed, failed := 0, 0
	for _, r := range results {
		if r.Success {
			passed++
		} else {
			failed++
		}
	}

	recap := RecapData{
		Date:      time.Now().Format("2006-01-02 15:04:05"),
		Total:     len(results),
		Passed:    passed,
		Failed:    failed,
		Results:   results,
		Variables: variables,
	}

	// Save JSON recap
	jsonData, _ := json.MarshalIndent(recap, "", "  ")
	os.WriteFile("api_test_recap.json", jsonData, 0644)

	// Save Markdown recap
	md := generateMarkdown(recap)
	os.WriteFile("API_TEST_RECAP.md", []byte(md), 0644)

	// Print summary
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("                  TESTING COMPLETE")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("Total:     %d\n", recap.Total)
	fmt.Printf("✅ Passed: %d\n", recap.Passed)
	fmt.Printf("❌ Failed: %d\n", recap.Failed)
	passRate := 0.0
	if recap.Total > 0 {
		passRate = float64(recap.Passed) / float64(recap.Total) * 100
	}
	fmt.Printf("Pass Rate: %.1f%%\n", passRate)
	fmt.Println("\nOutput files:")
	fmt.Println("  → api_test_recap.json")
	fmt.Println("  → API_TEST_RECAP.md")
}

func generateMarkdown(r RecapData) string {
	var sb strings.Builder
	sb.WriteString("# 📊 POS KopiTiam — API Test Recap\n\n")
	sb.WriteString(fmt.Sprintf("**Date:** %s\n\n", r.Date))
	sb.WriteString("## Summary\n\n")
	passRate := 0.0
	if r.Total > 0 {
		passRate = float64(r.Passed) / float64(r.Total) * 100
	}
	sb.WriteString(fmt.Sprintf(
		"| Metric | Value |\n|--------|-------|\n| Total | %d |\n| ✅ Passed | %d |\n| ❌ Failed | %d |\n| Pass Rate | %.1f%% |\n",
		r.Total, r.Passed, r.Failed, passRate,
	))

	sb.WriteString("\n## Results\n\n")
	sb.WriteString("| Phase | Name | Method | Status | Duration | Error |\n")
	sb.WriteString("|-------|------|--------|--------|----------|-------|\n")
	for _, res := range r.Results {
		icon := "✅"
		if !res.Success {
			icon = "❌"
		}
		errCell := "-"
		if res.Error != "" {
			errCell = truncate(res.Error, 60)
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s %d | %s | %s |\n",
			res.Phase, res.Name, res.Method, icon, res.Status, res.Duration, errCell))
	}

	sb.WriteString("\n## Variables Captured\n\n")
	sb.WriteString("| Key | Value |\n|-----|-------|\n")
	for k, v := range r.Variables {
		display := v
		if strings.Contains(k, "token") {
			maxLen := 40
			if len(display) > maxLen {
				display = display[:maxLen] + "..."
			}
		}
		sb.WriteString(fmt.Sprintf("| %s | %s |\n", k, display))
	}

	return sb.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func extractCatalogMenuID(data map[string]interface{}) {
	if data == nil {
		return
	}
	if d, ok := data["data"].([]interface{}); ok {
		for _, item := range d {
			if menu, ok := item.(map[string]interface{}); ok {
				if recipeBased, ok := menu["is_recipe_based"].(bool); ok && !recipeBased {
					if id, ok := menu["id"].(string); ok {
						variables["catalog_menu_id"] = id
						return
					}
				}
			}
		}
	}
}
