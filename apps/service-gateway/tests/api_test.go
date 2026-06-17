package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	v1 "openerp/apps/service-gateway/api/v1"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
)

func setupRouter() *gin.Engine {
	// Initialize in-memory SQLite for testing
	database, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to testing database: " + err.Error())
	}
	db.DB = database

	// Auto-migrate models
	db.DB.AutoMigrate(
		&models.SysTask{},
		&models.SysOrganization{},
		&models.SysAccountDetermination{},
		&models.SysGLAccount{},
		&models.SysJournalEntry{},
		&models.SysJournalEntryLine{},
		&models.SysExpenseClaim{},
		&models.SysMaterial{},
		&models.SysGoodsMovement{},
		&models.SysPurchaseOrder{},
		&models.SysPurchaseOrderLine{},
		&models.SysPurchaseReceipt{},
		&models.SysPurchaseReceiptLine{},
		&models.SysBusinessPartner{},
		&models.SysSalesOrder{},
		&models.SysSalesOrderLine{},
		&models.SysProductionOrder{},
	)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Organizations API
	r.POST("/openerp/v1/organizations", v1.CreateOrganization)
	r.GET("/openerp/v1/organizations", v1.GetOrganizations)
	r.PUT("/openerp/v1/organizations/:id", v1.UpdateOrganization)
	r.DELETE("/openerp/v1/organizations/:id", v1.DeleteOrganization)

	// GLAccounts API
	r.POST("/openerp/v1/gl-accounts", v1.CreateGLAccount)
	r.GET("/openerp/v1/gl-accounts", v1.GetGLAccounts)
	r.PUT("/openerp/v1/gl-accounts/:id", v1.UpdateGLAccount)
	r.DELETE("/openerp/v1/gl-accounts/:id", v1.DeleteGLAccount)

	// ExpenseClaims API
	r.POST("/openerp/v1/expense-claims", v1.CreateExpenseClaim)
	r.GET("/openerp/v1/expense-claims", v1.GetExpenseClaims)
	r.PUT("/openerp/v1/expense-claims/:id", v1.UpdateExpenseClaim)
	r.DELETE("/openerp/v1/expense-claims/:id", v1.DeleteExpenseClaim)

	// Purchasing
	r.POST("/openerp/v1/purchase-orders", v1.CreatePurchaseOrder)
	r.GET("/openerp/v1/purchase-orders", v1.GetPurchaseOrders)
	r.POST("/openerp/v1/purchase-receipts", v1.CreatePurchaseReceipt)
	r.GET("/openerp/v1/purchase-receipts", v1.GetPurchaseReceipts)

	// Business Partners
	r.POST("/openerp/v1/business-partners", v1.CreateBusinessPartner)
	r.GET("/openerp/v1/business-partners", v1.GetBusinessPartners)
	r.PUT("/openerp/v1/business-partners/:id", v1.UpdateBusinessPartner)
	r.DELETE("/openerp/v1/business-partners/:id", v1.DeleteBusinessPartner)

	// Sales Orders
	r.POST("/openerp/v1/sales-orders", v1.CreateSalesOrder)
	r.GET("/openerp/v1/sales-orders", v1.GetSalesOrders)
	r.PUT("/openerp/v1/sales-orders/:id", v1.UpdateSalesOrder)
	r.DELETE("/openerp/v1/sales-orders/:id", v1.DeleteSalesOrder)

	// Goods Movement
	r.POST("/openerp/v1/goods-movement", v1.CreateGoodsMovement)
	r.GET("/openerp/v1/goods-movement", v1.GetGoodsMovements)

	// Production Orders
	r.POST("/openerp/v1/production-orders", v1.CreateProductionOrder)
	r.GET("/openerp/v1/production-orders", v1.GetProductionOrders)
	r.PUT("/openerp/v1/production-orders/:id", v1.UpdateProductionOrder)
	r.DELETE("/openerp/v1/production-orders/:id", v1.DeleteProductionOrder)

	return r
}

func TestOrganizationCRUD(t *testing.T) {
	r := setupRouter()

	// 1. Create
	org := models.SysOrganization{
		OrgCode: "ORG-001",
		OrgName: "HQ",
		OrgType: "Legal",
	}
	body, _ := json.Marshal(org)
	req, _ := http.NewRequest("POST", "/openerp/v1/organizations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(200), resp["code"])
	orgID := resp["data"].(float64)

	// 2. Read
	req, _ = http.NewRequest("GET", "/openerp/v1/organizations", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	dataList := resp["data"].([]interface{})
	assert.Len(t, dataList, 1)

	// 3. Update
	orgUpdate := models.SysOrganization{
		OrgName: "Global HQ",
	}
	bodyUpdate, _ := json.Marshal(orgUpdate)
	req, _ = http.NewRequest("PUT", "/openerp/v1/organizations/1", bytes.NewBuffer(bodyUpdate))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Delete
	req, _ = http.NewRequest("DELETE", "/openerp/v1/organizations/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_ = orgID // avoid unused variable error
}

func TestGLAccountCRUD(t *testing.T) {
	r := setupRouter()

	account := models.SysGLAccount{
		AccountCode: "1001",
		AccountName: "Cash",
		AccountType: "Asset",
	}
	body, _ := json.Marshal(account)
	req, _ := http.NewRequest("POST", "/openerp/v1/gl-accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/openerp/v1/gl-accounts", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("DELETE", "/openerp/v1/gl-accounts/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestExpenseClaimCRUD(t *testing.T) {
	r := setupRouter()

	claim := models.SysExpenseClaim{
		ClaimNo:    "EXP-1001",
		EmployeeID: "EMP001",
		Amount:     1500.50,
		Reason:     "Travel",
		Status:     "Draft",
	}
	body, _ := json.Marshal(claim)
	req, _ := http.NewRequest("POST", "/openerp/v1/expense-claims", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/openerp/v1/expense-claims", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("DELETE", "/openerp/v1/expense-claims/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBusinessPartnerCRUD(t *testing.T) {
	r := setupRouter()

	partner := models.SysBusinessPartner{
		BpCode: "BP001",
		BpName: "Test Partner",
		BpType: "Customer",
	}
	body, _ := json.Marshal(partner)
	req, _ := http.NewRequest("POST", "/openerp/v1/business-partners", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/openerp/v1/business-partners", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("DELETE", "/openerp/v1/business-partners/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSalesOrderCRUD(t *testing.T) {
	r := setupRouter()

	order := models.SysSalesOrder{
		DocNo:       "SO001",
		SoldToBpID:  "BP001",
		ShipToBpID:  "BP001",
		BillToBpID:  "BP001",
		PayerBpID:   "BP001",
		TotalAmount: 100.0,
	}
	body, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/openerp/v1/sales-orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/openerp/v1/sales-orders", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("DELETE", "/openerp/v1/sales-orders/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGoodsMovementReceiptUpdatesStockAndPostsJournalEntry(t *testing.T) {
	r := setupRouter()

	material := models.SysMaterial{
		MaterialCode:      "MAT-RECEIPT",
		MaterialName:      "Receipt Material",
		Category:          "RAW",
		Unit:              "EA",
		Stock:             10,
		MovingAverageCost: 5,
	}
	assert.NoError(t, db.DB.Create(&material).Error)
	assert.NoError(t, db.DB.Create(&models.SysAccountDetermination{
		TransactionType: v1.PostingEventGoodsReceipt,
		ItemCategory:    "RAW",
		DebitAccount:    "1405",
		CreditAccount:   "2202",
	}).Error)

	body, _ := json.Marshal(map[string]interface{}{
		"mvmtNo":   "GM-RECEIPT-001",
		"mvmtType": "Receipt",
		"itemCode": "MAT-RECEIPT",
		"qty":      5,
		"unitCost": 8,
	})
	req, _ := http.NewRequest("POST", "/openerp/v1/goods-movement", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/openerp/v1/goods-movement", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var updated models.SysMaterial
	assert.NoError(t, db.DB.Where("material_code = ?", "MAT-RECEIPT").First(&updated).Error)
	assert.Equal(t, 15.0, updated.Stock)
	assert.InDelta(t, 6.0, updated.MovingAverageCost, 0.0001)

	var entry models.SysJournalEntry
	assert.NoError(t, db.DB.Preload("Lines").Where("total_amount = ?", 40).First(&entry).Error)
	assert.Equal(t, 40.0, entry.TotalAmount)
	assert.Len(t, entry.Lines, 2)
	assertJournalLine(t, entry.Lines, "1405", 40, 0)
	assertJournalLine(t, entry.Lines, "2202", 0, 40)
}

func TestGoodsMovementIssueUpdatesStockAndPostsJournalEntry(t *testing.T) {
	r := setupRouter()

	material := models.SysMaterial{
		MaterialCode:      "MAT-ISSUE",
		MaterialName:      "Issue Material",
		Category:          "RAW",
		Unit:              "EA",
		Stock:             10,
		MovingAverageCost: 6,
	}
	assert.NoError(t, db.DB.Create(&material).Error)
	assert.NoError(t, db.DB.Create(&models.SysAccountDetermination{
		TransactionType: v1.PostingEventGoodsIssue,
		ItemCategory:    "RAW",
		DebitAccount:    "6401",
		CreditAccount:   "1405",
	}).Error)

	body, _ := json.Marshal(map[string]interface{}{
		"mvmtNo":   "GM-ISSUE-001",
		"mvmtType": "Issue",
		"itemCode": "MAT-ISSUE",
		"qty":      3,
	})
	req, _ := http.NewRequest("POST", "/openerp/v1/goods-movement", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var updated models.SysMaterial
	assert.NoError(t, db.DB.Where("material_code = ?", "MAT-ISSUE").First(&updated).Error)
	assert.Equal(t, 7.0, updated.Stock)
	assert.Equal(t, 6.0, updated.MovingAverageCost)

	var entry models.SysJournalEntry
	assert.NoError(t, db.DB.Preload("Lines").Where("total_amount = ?", 18).First(&entry).Error)
	assert.Equal(t, 18.0, entry.TotalAmount)
	assert.Len(t, entry.Lines, 2)
	assertJournalLine(t, entry.Lines, "6401", 18, 0)
	assertJournalLine(t, entry.Lines, "1405", 0, 18)
}

func TestPurchaseOrderReceiptPostsInventoryAndJournalWithTraceability(t *testing.T) {
	r := setupRouter()

	material := models.SysMaterial{
		MaterialCode:      "MAT-GRPO",
		MaterialName:      "GRPO Material",
		Category:          "RAW",
		Unit:              "EA",
		Stock:             10,
		MovingAverageCost: 5,
	}
	assert.NoError(t, db.DB.Create(&material).Error)
	assert.NoError(t, db.DB.Create(&models.SysAccountDetermination{
		TransactionType: v1.PostingEventGoodsReceipt,
		ItemCategory:    "RAW",
		DebitAccount:    "1405",
		CreditAccount:   "2202",
	}).Error)

	orderBody, _ := json.Marshal(map[string]interface{}{
		"docNo":      "PO-GRPO-001",
		"title":      "GRPO integration test",
		"soldToBpId": "VEND-001",
		"shipToBpId": "PLANT-001",
		"billToBpId": "VEND-INV",
		"payerBpId":  "VEND-PAY",
		"lines": []map[string]interface{}{
			{
				"itemCode":  "MAT-GRPO",
				"qty":       4,
				"unitPrice": 8,
			},
		},
	})
	req, _ := http.NewRequest("POST", "/openerp/v1/purchase-orders", bytes.NewBuffer(orderBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var orderResponse struct {
		Data models.SysPurchaseOrder `json:"data"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &orderResponse))
	assert.NotZero(t, orderResponse.Data.ID)

	receiptBody, _ := json.Marshal(map[string]interface{}{
		"docNo":           "GRPO-001",
		"purchaseOrderId": orderResponse.Data.ID,
		"costCenter":      "WH01",
		"lines": []map[string]interface{}{
			{
				"baseLine": 0,
				"qty":      4,
			},
		},
	})
	req, _ = http.NewRequest("POST", "/openerp/v1/purchase-receipts", bytes.NewBuffer(receiptBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var updated models.SysMaterial
	assert.NoError(t, db.DB.Where("material_code = ?", "MAT-GRPO").First(&updated).Error)
	assert.Equal(t, 14.0, updated.Stock)
	assert.InDelta(t, 5.8571, updated.MovingAverageCost, 0.0001)

	var receipt models.SysPurchaseReceipt
	assert.NoError(t, db.DB.Preload("Lines").Where("doc_no = ?", "GRPO-001").First(&receipt).Error)
	assert.Equal(t, orderResponse.Data.ID, receipt.PurchaseOrderID)
	assert.Equal(t, "VEND-PAY", receipt.PayerBpID)
	assert.Equal(t, 32.0, receipt.TotalAmount)
	assert.Len(t, receipt.Lines, 1)
	assert.Equal(t, models.DocumentTypePurchaseOrder, receipt.Lines[0].BaseType)
	assert.Equal(t, orderResponse.Data.ID, receipt.Lines[0].BaseEntry)
	assert.Equal(t, uint(0), receipt.Lines[0].BaseLine)

	var entry models.SysJournalEntry
	assert.NoError(t, db.DB.Preload("Lines").Where("total_amount = ?", 32).First(&entry).Error)
	assertJournalLine(t, entry.Lines, "1405", 32, 0)
	assertJournalLine(t, entry.Lines, "2202", 0, 32)
	assert.Equal(t, "VEND-PAY", entry.Lines[0].PayerBpID)
}

func TestProductionOrderCRUD(t *testing.T) {
	r := setupRouter()

	order := models.SysProductionOrder{
		OrderNo:    "PO001",
		ItemCode:   "ITEM001",
		PlannedQty: 10,
	}
	body, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/openerp/v1/production-orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/openerp/v1/production-orders", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("DELETE", "/openerp/v1/production-orders/1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func assertJournalLine(t *testing.T, lines []models.SysJournalEntryLine, accountCode string, debit float64, credit float64) {
	t.Helper()
	for _, line := range lines {
		if line.AccountCode == accountCode {
			assert.Equal(t, debit, line.DebitAmount)
			assert.Equal(t, credit, line.CreditAmount)
			return
		}
	}
	t.Fatalf("journal line for account %s not found", accountCode)
}
