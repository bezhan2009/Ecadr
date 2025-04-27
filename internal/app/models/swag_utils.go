package models

import (
	"time"
)

// TokenResponse represents the response with access token and user ID
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       uint   `json:"user_id"`
}

// RefreshTokenResponse represents the response with access token and user ID
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// ErrorResponse represents an error message response
type ErrorResponse struct {
	Error string `json:"error"`
}

// DefaultResponse represents a default message response
type DefaultResponse struct {
	Message string `json:"message"`
}

type UserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	RoleID    uint   `json:"role_id"`

	TitleCompany       string   `json:"title_company"`
	DescriptionCompany string   `json:"description_company"`
	Criteria           []string `json:"criteria"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required"` // Название категории, обязательное поле
	ParentID     uint   `json:"parent_id,omitempty"`              // Идентификатор родительской категории, необязательное поле
	Description  string `json:"description,omitempty"`            // Описание категории, необязательное поле
}

type AddressRequest struct {
	AddressName string `json:"address_name"`
}

type AccountRequest struct {
	AccountName string `json:"account_number"`
}

type FillAccountRequest struct {
	AccountName string `json:"account_number"`
	Balance     uint   `json:"balance"`
}

type AccountsResponse struct {
	AccountName string `json:"account_name"`
}

type FeaturedProductsRequest struct {
	ProductID uint `json:"product_id"`
}

type ReviewRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Rating  uint   `json:"rating"`
}

type OrderRequest struct {
	StatusID  uint `json:"status_id"`
	AddressID uint `json:"address_id"`
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type OrderStatusRequest struct {
	StatusName  string `json:"status_name"`
	Description string `json:"description"`
}

type PaymentRequest struct {
	AccountID uint `json:"account_id"`
	OrderID   uint `json:"order_id"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ProductRequest struct {
	StoreID       uint     `json:"store_id"`
	CategoryID    uint     `json:"category_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Price         uint     `json:"price"`
	Amount        uint     `json:"amount"`
	ProductImages []string `json:"product_images"`
}

type ProductResponse struct {
	StoreID       uint     `json:"store_id"`
	CategoryID    uint     `json:"category_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Price         uint     `json:"price"`
	Amount        uint     `json:"amount"`
	ProductImages []string `json:"product_images"`
}

type StoreRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StoreReviewRequest struct {
	Rating  uint   `json:"rating"`
	Comment string `json:"comment"`
}

type CommentRequest struct {
	ParentID    uint   `json:"parent_id"`
	CommentText string `json:"text" binding:"required"`
}

type CourseResponse struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	Tags        []string  `json:"tags"`
	WorkerID    int       `json:"worker_id"`
	Worker      User      `json:"-"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type CourseReq struct {
	Message  string `json:"message"`
	CourseID uint   `json:"course_id"`
}

type VacancyResponse struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Subject     string   `json:"subject"`
	Tags        []string `json:"tags"`
	WorkerID    int      `json:"worker_id"`
	Worker      User     `json:"-"`
	Contact     string   `json:"contact"`
	Salary      Salary   `json:"salary"`
	Location    string   `json:"location"`
	Experience  string   `json:"experience"`
}

type VacancyReq struct {
	Message   string `json:"message"`
	VacancyID uint   `json:"vacancy_id"`
}

type RecommendReq struct {
	Message     string `json:"message"`
	RecommendID uint   `json:"recommend_id"`
}

type CriteriaReq struct {
	Title     string `json:"title"`
	VacancyID uint   `json:"vacancy_id"`
}

type CriteriaResp struct {
	Message string `json:"message"`
}
