package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/internal/service"
	"github.com/naseej/server/pkg/database"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://naseej:naseej@localhost:5432/naseej?sslmode=disable"
	}

	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	log.Info().Msg("seeding development data...")

	// Create test merchant
	merchant := &models.Merchant{
		BaseModel:      models.BaseModel{ID: uuid.New()},
		BusinessName:   "Naseej Demo Store",
		Slug:           "naseej-demo",
		Domain:         "demo.naseej.local",
		PrimaryColor:   "#1B2A4A",
		SecondaryColor: "#D4A574",
		TaxRate:        0.15,
		Currency:       "SAR",
		Timezone:       "Asia/Riyadh",
		IsActive:       true,
	}

	if err := db.Where("slug = ?", merchant.Slug).FirstOrCreate(merchant).Error; err != nil {
		log.Fatal().Err(err).Msg("failed to create merchant")
	}
	log.Info().Str("merchant_id", merchant.ID.String()).Msg("merchant created")

	// Create permissions
	permissions := createPermissions(db)
	log.Info().Int("count", len(permissions)).Msg("permissions created")

	// Create roles
	roles := createRoles(db, merchant.ID, permissions)
	log.Info().Int("count", len(roles)).Msg("roles created")

	// Create admin user
	adminPassword, _ := service.HashPassword("admin123")
	adminUser := &models.User{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchant.ID,
		},
		Email:        "admin@naseej.demo",
		PasswordHash: adminPassword,
		FirstName:    "Admin",
		LastName:     "User",
		IsActive:     true,
	}
	if err := db.Where("merchant_id = ? AND email = ?", merchant.ID, adminUser.Email).FirstOrCreate(adminUser).Error; err != nil {
		log.Fatal().Err(err).Msg("failed to create admin user")
	}

	// Assign super_admin role
	if superAdmin, ok := roles["super_admin"]; ok {
		db.Model(adminUser).Association("Roles").Append(&superAdmin)
	}
	log.Info().Str("email", adminUser.Email).Msg("admin user created")

	// Create cashier user
	cashierPassword, _ := service.HashPassword("cashier123")
	cashierPin, _ := service.HashPin("1234")
	cashierUser := &models.User{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchant.ID,
		},
		Email:        "cashier@naseej.demo",
		PasswordHash: cashierPassword,
		FirstName:    "Sara",
		LastName:     "Cashier",
		PinCode:      cashierPin,
		IsActive:     true,
	}
	if err := db.Where("merchant_id = ? AND email = ?", merchant.ID, cashierUser.Email).FirstOrCreate(cashierUser).Error; err != nil {
		log.Fatal().Err(err).Msg("failed to create cashier user")
	}
	if cashierRole, ok := roles["cashier"]; ok {
		db.Model(cashierUser).Association("Roles").Append(&cashierRole)
	}
	log.Info().Str("email", cashierUser.Email).Msg("cashier user created")

	// Create location
	location := &models.Location{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchant.ID,
		},
		Name:     "Main Store",
		Address:  "123 Fashion Ave",
		City:     "Riyadh",
		Country:  "SA",
		IsActive: true,
	}
	if err := db.Where("merchant_id = ? AND name = ?", merchant.ID, location.Name).FirstOrCreate(location).Error; err != nil {
		log.Fatal().Err(err).Msg("failed to create location")
	}
	log.Info().Str("name", location.Name).Msg("location created")

	// Create size groups
	sizeGroups := createSizeGroups(db, merchant.ID)
	log.Info().Int("count", len(sizeGroups)).Msg("size groups created")

	// Create colors
	colors := createColors(db, merchant.ID)
	log.Info().Int("count", len(colors)).Msg("colors created")

	// Create categories
	createCategories(db, merchant.ID, sizeGroups)
	log.Info().Msg("categories created")

	log.Info().Msg("seed complete!")
	fmt.Printf("\n--- Development Credentials ---\n")
	fmt.Printf("Merchant ID: %s\n", merchant.ID)
	fmt.Printf("Admin:   admin@naseej.demo / admin123\n")
	fmt.Printf("Cashier: cashier@naseej.demo / cashier123 (PIN: 1234)\n")
	fmt.Printf("-------------------------------\n")
	_ = colors
}

func createPermissions(db *gorm.DB) []models.Permission {
	permDefs := []struct {
		Resource string
		Action   string
	}{
		{"products", "view"}, {"products", "create"}, {"products", "edit"}, {"products", "delete"}, {"products", "publish"},
		{"inventory", "view"}, {"inventory", "edit"}, {"inventory", "transfer"}, {"inventory", "price"},
		{"orders", "view"}, {"orders", "create"}, {"orders", "edit"}, {"orders", "cancel"}, {"orders", "refund"}, {"orders", "discount"}, {"orders", "export"},
		{"pos", "access"}, {"pos", "open_drawer"}, {"pos", "void_item"}, {"pos", "override_price"}, {"pos", "apply_coupon"}, {"pos", "park_sale"}, {"pos", "cash_in"}, {"pos", "cash_out"}, {"pos", "end_shift"}, {"pos", "end_shift_other"},
		{"customers", "view"}, {"customers", "create"}, {"customers", "edit"}, {"customers", "delete"}, {"customers", "export"}, {"customers", "loyalty"},
		{"coupons", "view"}, {"coupons", "create"}, {"coupons", "edit"}, {"coupons", "delete"},
		{"reports", "sales"}, {"reports", "inventory"}, {"reports", "customers"}, {"reports", "financial"},
		{"staff", "view"}, {"staff", "create"}, {"staff", "edit"}, {"staff", "delete"}, {"staff", "assign_role"},
		{"settings", "view"}, {"settings", "edit"}, {"settings", "locations"},
	}

	var permissions []models.Permission
	for _, p := range permDefs {
		perm := models.Permission{
			BaseModel: models.BaseModel{ID: uuid.New()},
			Resource:  p.Resource,
			Action:    p.Action,
		}
		db.Where("resource = ? AND action = ?", p.Resource, p.Action).FirstOrCreate(&perm)
		permissions = append(permissions, perm)
	}
	return permissions
}

func createRoles(db *gorm.DB, merchantID uuid.UUID, permissions []models.Permission) map[string]models.Role {
	roles := make(map[string]models.Role)

	roleDefs := []struct {
		Name        string
		Description string
	}{
		{"super_admin", "Full access to everything"},
		{"admin", "Manages full operation"},
		{"store_manager", "Runs day-to-day operations"},
		{"senior_cashier", "Trusted front-line staff"},
		{"cashier", "Standard POS staff"},
		{"ecommerce_manager", "Manages online storefront"},
		{"inventory_manager", "Warehouse/stockroom staff"},
		{"viewer", "Read-only access"},
	}

	for _, rd := range roleDefs {
		role := models.Role{
			TenantModel: models.TenantModel{
				ID:         uuid.New(),
				MerchantID: merchantID,
			},
			Name:        rd.Name,
			Description: rd.Description,
			IsSystem:    true,
		}
		db.Where("merchant_id = ? AND name = ?", merchantID, rd.Name).FirstOrCreate(&role)

		// Assign all permissions to super_admin
		if rd.Name == "super_admin" {
			db.Model(&role).Association("Permissions").Replace(permissions)
		}

		roles[rd.Name] = role
	}

	return roles
}

func createSizeGroups(db *gorm.DB, merchantID uuid.UUID) map[string]models.SizeGroup {
	groups := make(map[string]models.SizeGroup)

	groupDefs := map[string][]string{
		"Women's Shoes":  {"36", "37", "38", "39", "40", "41"},
		"Men's Shoes":    {"40", "41", "42", "43", "44", "45"},
		"Shirts":         {"S", "M", "L", "XL", "XXL"},
		"Kids Clothes":   {"2Y", "4Y", "6Y", "8Y", "10Y", "12Y"},
	}

	for name, sizes := range groupDefs {
		group := models.SizeGroup{
			TenantModel: models.TenantModel{
				ID:         uuid.New(),
				MerchantID: merchantID,
			},
			Name: name,
		}
		db.Where("merchant_id = ? AND name = ?", merchantID, name).FirstOrCreate(&group)

		for i, sizeName := range sizes {
			size := models.Size{
				TenantModel: models.TenantModel{
					ID:         uuid.New(),
					MerchantID: merchantID,
				},
				SizeGroupID: group.ID,
				Name:        sizeName,
				SortOrder:   i,
			}
			db.Where("size_group_id = ? AND name = ?", group.ID, sizeName).FirstOrCreate(&size)
		}

		groups[name] = group
	}

	return groups
}

func createColors(db *gorm.DB, merchantID uuid.UUID) []models.Color {
	colorDefs := []struct {
		Name    string
		HexCode string
	}{
		{"Black", "#000000"},
		{"White", "#FFFFFF"},
		{"Navy Blue", "#1B2A4A"},
		{"Red", "#DC2626"},
		{"Beige", "#D4A574"},
		{"Gray", "#6B7280"},
		{"Green", "#059669"},
		{"Pink", "#EC4899"},
	}

	var colors []models.Color
	for _, c := range colorDefs {
		color := models.Color{
			TenantModel: models.TenantModel{
				ID:         uuid.New(),
				MerchantID: merchantID,
			},
			Name:    c.Name,
			HexCode: c.HexCode,
		}
		db.Where("merchant_id = ? AND name = ?", merchantID, c.Name).FirstOrCreate(&color)
		colors = append(colors, color)
	}

	return colors
}

func createCategories(db *gorm.DB, merchantID uuid.UUID, sizeGroups map[string]models.SizeGroup) {
	type catDef struct {
		Name      string
		Slug      string
		SizeGroup string
	}

	categories := []catDef{
		{"Women's Footwear", "womens-footwear", "Women's Shoes"},
		{"Men's Footwear", "mens-footwear", "Men's Shoes"},
		{"Women's Tops", "womens-tops", "Shirts"},
		{"Men's Shirts", "mens-shirts", "Shirts"},
		{"Kids", "kids", "Kids Clothes"},
	}

	for _, c := range categories {
		cat := models.Category{
			TenantModel: models.TenantModel{
				ID:         uuid.New(),
				MerchantID: merchantID,
			},
			Name: c.Name,
			Slug: c.Slug,
		}

		if sg, ok := sizeGroups[c.SizeGroup]; ok {
			cat.SizeGroupID = &sg.ID
		}

		db.Where("merchant_id = ? AND slug = ?", merchantID, c.Slug).FirstOrCreate(&cat)
	}
}
