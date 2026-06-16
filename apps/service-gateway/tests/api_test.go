package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	v1 "openerp/apps/service-gateway/api/v1"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
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
		&models.SysOrganization{},
		&models.SysGLAccount{},
		&models.SysExpenseClaim{},
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
		ClaimNo: "EXP-1001",
		EmployeeID: "EMP001",
		Amount: 1500.50,
		Reason: "Travel",
		Status: "Draft",
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
		DocNo: "SO001",
		BpID: "BP001",
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

func TestProductionOrderCRUD(t *testing.T) {
	r := setupRouter()

	order := models.SysProductionOrder{
		OrderNo: "PO001",
		ItemCode: "ITEM001",
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

