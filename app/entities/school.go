package entities

import (
	"gorm.io/gorm"
)

type (
	School struct {
		gorm.Model
		UserID           uint
		Npsn             string `gorm:"type:varchar(12);not null"`
		Name             string `gorm:"type:varchar(150);not null"`
		Description      string `gorm:"type:varchar(255);not null"`
		Image            string `gorm:"type:varchar(150);not null"`
		Video            string `gorm:"type:varchar(150);not null"`
		Pdf              string `gorm:"type:varchar(150);not null"`
		Web              string `gorm:"type:varchar(150);not null"`
		Province         string `gorm:"type:varchar(150);not null"`
		City             string `gorm:"type:varchar(150);not null"`
		District         string `gorm:"type:varchar(150);not null"`
		Village          string `gorm:"type:varchar(150);not null"`
		Detail           string `gorm:"type:varchar(150);not null"`
		ZipCode          string `gorm:"type:varchar(150);not null"`
		Students         string `gorm:"not null"`
		Teachers         string `gorm:"not null"`
		Staff            string `gorm:"not null"`
		Phone            string `gorm:"type:varchar(15);not null"`
		Accreditation    string `gorm:"type:varchar(3);not null"`
		Gmeet            string `gorm:"type:varchar(70);default: "`
		GmeetDate        string `gorm:"type:varchar(25)"`
		QuizLinkPub      string `gorm:"type:varchar(150);default:"`
		QuizLinkPreview  string `gorm:"type:varchar(150);default:"`
		QuizLinkResult   string `gorm:"type:varchar(150);default:"`
		Achievements     []Achievement
		Extracurriculars []Extracurricular
		Faqs             []Faq
		Payments         []Payment
		User             *User
		Submissions      []Submission
		Progresses       []Progress
		Reviews          []Reviews
		Carts            []Carts
	}

	Submission struct {
		ID               uint `gorm:"primaryKey;autoIncrement;not null"`
		SchoolID         uint
		UserID           uint
		StudentPhoto     string `gorm:"type:varchar(255);not null"`
		StudentName      string `gorm:"type:varchar(255);not null"`
		PlaceDate        string `gorm:"type:varchar(255);not null"`
		Gender           string `gorm:"type:varchar(255);not null"`
		Religion         string `gorm:"type:varchar(255);not null"`
		GraduationFrom   string `gorm:"type:varchar(255);not null"`
		NISN             string `gorm:"type:varchar(255);not null"`
		StudentAddress   string `gorm:"type:varchar(255);not null"`
		ParentName       string `gorm:"type:varchar(255);not null"`
		ParentJob        string `gorm:"type:varchar(255);not null"`
		ParentReligion   string `gorm:"type:varchar(255);not null"`
		ParentGender     string `gorm:"type:varchar(255);not null"`
		ParentAddress    string `gorm:"type:varchar(255);not null"`
		ParentPhone      string `gorm:"type:varchar(255);not null"`
		ParentSignature  string `gorm:"type:varchar(255);not null"`
		StudentSignature string `gorm:"type:varchar(255);not null"`
		Date             string `gorm:"type:varchar(255);not null"`
		School           School
		User             User
	}
	ReqAdressSubmission struct {
		Province string `json:"province" `
		District string `json:"district" `
		Village  string `json:"village" `
		ZipCode  string `json:"zip_code" `
		City     string `json:"city"`
		Detail   string `json:"detail"`
	}
	ReqAddQuiz struct {
		SchoolID int    `json:"school_id" validate:"required"`
		Question string `json:"question" validate:"required"`
		Option1  string `json:"option1" validate:"required"`
		Option2  string `json:"option2" validate:"required"`
		Option3  string `json:"option3" validate:"required"`
		Option4  string `json:"option4" validate:"required"`
		Answer   int    `json:"answer" validate:"required"`
	}
	ReqDataQuiz struct {
		PubLink    string       `json:"pub_link"`
		Prevlink   string       `json:"prev_link"`
		ResultLink string       `json:"result_link"`
		Data       []ReqAddQuiz `json:"data"`
	}
	ReqCreateSubmission struct {
		UserID           uint
		SchoolID         int    `form:"school_id" validate:"required"`
		StudentPhoto     string `form:"student_photo" validate:"required"`
		StudentName      string `form:"student_name" validate:"required"`
		PlaceDate        string `form:"place_date" validate:"required"`
		Gender           string `form:"gender" validate:"required"`
		Religion         string `form:"religion" validate:"required"`
		GraduationFrom   string `form:"graduation_from" validate:"required"`
		NISN             string `form:"nisn" validate:"required"`
		StudentProvince  string `form:"student_province" validate:"required"`
		StudentDistrict  string `form:"student_district" validate:"required"`
		StudentVillage   string `form:"student_village" validate:"required"`
		StudentZipCode   string `form:"student_zip_code" validate:"required"`
		StudentCity      string `form:"student_city" validate:"required"`
		StudentDetail    string `form:"student_detail" validate:"required"`
		ParentProvince   string `form:"parent_province" validate:"required"`
		ParentDistrict   string `form:"parent_district" validate:"required"`
		ParentVillage    string `form:"parent_village" validate:"required"`
		ParentZipCode    string `form:"parent_zip_code" validate:"required"`
		ParentCity       string `form:"parent_city" validate:"required"`
		ParentDetail     string `form:"parent_detail" validate:"required"`
		ParentName       string `form:"parent_name" validate:"required"`
		ParentGender     string `form:"parent_gender" validate:"required"`
		ParentJob        string `form:"parent_job" validate:"required"`
		ParentReligion   string `form:"parent_religion" validate:"required"`
		ParentPhone      string `form:"parent_phone" validate:"required"`
		ParentSignature  string `form:"parent_signature" validate:"required"`
		StudentSignature string `form:"student_signature" validate:"required"`
		Date             string `form:"date"`
	}
	Progress struct {
		ID        uint `gorm:"primaryKey;autoIncrement;not null"`
		UserID    uint
		SchoolID  uint
		Status    string
		DeletedAt gorm.DeletedAt `gorm:"index"`
		School    School
		User      User
	}
	ResAllProgress struct {
		SchoolName  string `json:"school_name"`
		SchoolImage string `json:"school_image"`
		SchoolWeb   string `json:"school_web"`
		ProgressId  int    `json:"progress_id"`
	}
	ResDetailProgress struct {
		Id     int    `json:"progress_id"`
		Status string `json:"progress_status"`
	}
	ResAllProgressSubmission struct {
		UserId         int    `json:"user_id"`
		UserImage      string `json:"user_image"`
		UserName       string `json:"user_name"`
		SubmissionId   int    `json:"submission_id"`
		ProgressId     int    `json:"progress_id"`
		ProgressStatus string `json:"progress_status"`
	}

	StudentData struct {
		Photo          string              `json:"photo"`
		Name           string              `json:"name"`
		PlaceDate      string              `json:"place_date"`
		Gender         string              `json:"gender"`
		Religion       string              `json:"religion"`
		GraduationFrom string              `json:"graduation_from"`
		NISN           string              `json:"nisn"`
		Adress         ReqAdressSubmission `json:"address"`
	}
	ParentData struct {
		Name     string              `json:"name"`
		Job      string              `json:"job"`
		Religion string              `json:"religion"`
		Phone    string              `json:"phone"`
		Gender   string              `json:"gender"`
		Adress   ReqAdressSubmission `json:"address"`
	}
	ResDetailSubmission struct {
		StudentData      StudentData `json:"student_data"`
		ParentData       ParentData  `json:"parent_data"`
		ParentSignature  string      `json:"parent_signature"`
		StudentSignature string      `json:"student_signature"`
		DatePlace        string      `json:"date_place"`
		SchoolName       string      `json:"school_name"`
	}
	Achievement struct {
		gorm.Model
		SchoolID    uint
		Description string `gorm:"type:varchar(255);not null"`
		Image       string `gorm:"type:varchar(50);not null"`
		Title       string `gorm:"type:varchar(50);not null"`
	}
	Extracurricular struct {
		gorm.Model
		SchoolID    uint
		Description string `gorm:"type:varchar(255);not null"`
		Image       string `gorm:"type:varchar(50);not null"`
		Title       string `gorm:"type:varchar(50);not null"`
	}
	Faq struct {
		gorm.Model
		SchoolID uint
		Question string `gorm:"type:varchar(255);not null"`
		Answer   string `gorm:"type:varchar(255);not null"`
	}
	ResFaq struct {
		Id       int    `json:"id"`
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}
	Payment struct {
		gorm.Model
		SchoolID    uint
		Description string `gorm:"type:varchar(255);not null"`
		Image       string `gorm:"type:varchar(70);not null"`
		Type        string `gorm:"type:varchar(15);not null"`
		Price       int
		Interval    int
	}
	ReqAddPayment struct {
		SchoolID    uint   `form:"school_id" validate:"required"`
		Description string `form:"description" validate:"required"`
		Price       int    `form:"price" validate:"required"`
		Image       string `form:"image" validate:"required"`
		Interval    *int   `form:"interval" validate:"required"`
	}
	ReqUpdatePayment struct {
		ID          int    `form:"id" validate:"required"`
		Description string `form:"description"`
		Price       int    `form:"price"`
		Image       string `form:"image"`
		Interval    *int   `form:"interval"`
	}
	ReqAddFaq struct {
		SchoolId int    `json:"school_id" validate:"required"`
		Question string `json:"question" validate:"required"`
		Answer   string `json:"answer" validate:"required"`
	}
	ReqUpdateFaq struct {
		Id       int    `json:"id" validate:"required"`
		Question string `json:"question"`
		Answer   string `json:"answer" `
	}
	ReqAddExtracurricular struct {
		SchoolID    uint   `form:"school_id" validate:"required"`
		Description string `form:"description" validate:"required"`
		Image       string `form:"image" validate:"required"`
		Title       string `form:"title" validate:"required"`
	}
	ReqUpdateExtracurricular struct {
		Id          int    `form:"id" validate:"required"`
		Description string `form:"description" `
		Image       string `form:"image" `
		Name        string `form:"name" `
	}
	ReqAddAchievemnt struct {
		SchoolID    uint   `form:"school_id" validate:"required"`
		Description string `form:"description" validate:"required"`
		Image       string `form:"image" validate:"required"`
		Title       string `form:"title" validate:"required"`
	}
	ReqUpdateAchievemnt struct {
		Id          int    `form:"id" validate:"required"`
		Description string `form:"description" `
		Image       string `form:"image" `
		Name        string `form:"name" `
	}
	ReqCreateGmeet struct {
		StartDate string `json:"start_time" validate:"required"`
		EndDate   string `json:"end_time" validate:"required"`
		SchoolId  int    `json:"school_id" validate:"required"`
	}
	ResAddItems struct {
		Id          int    `json:"id,omitempty"`
		Name        string `json:"name"`
		Img         string `json:"image"`
		Description string `json:"description"`
	}
	ResPaymentType struct {
		Id          int    `json:"id"`
		Img         string `json:"image"`
		Description string `json:"description"`
		Price       int    `json:"price"`
		Interval    string `json:"interval,omitempty"`
	}
	ResPayment struct {
		OneTime  []ResPaymentType `json:"onetime"`
		Interval []ResPaymentType `json:"interval"`
	}
	Response struct {
		Limit     int `json:"limit,omitempty"`
		Page      int `json:"page,omitempty"`
		TotalPage int `json:"total_page,omitempty"`
		TotalData int `json:"total_data,omitempty"`
		Data      any `json:"data"`
	}
	ResAllSchool struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		AdminName     string `json:"admin_name"`
		Image         string `json:"image"`
		Accreditation string `json:"accreditation"`
		Location      string `json:"location"`
	}
	ResReview struct {
		UserImage string `json:"image"`
		Review    string `json:"review"`
	}
	ResDetailSchool struct {
		Id               int           `json:"id"`
		Npsn             string        `json:"npsn"`
		Name             string        `json:"name"`
		Description      string        `json:"description"`
		Image            string        `json:"image"`
		Video            string        `json:"video"`
		Pdf              string        `json:"pdf"`
		Web              string        `json:"web"`
		Province         string        `json:"province"`
		City             string        `json:"city"`
		District         string        `json:"district"`
		Village          string        `json:"village"`
		Detail           string        `json:"detail"`
		ZipCode          string        `json:"zipCode"`
		Students         string        `json:"students"`
		Teachers         string        `json:"teachers"`
		Staff            string        `json:"staff"`
		Accreditation    string        `json:"accreditation"`
		Gmeet            string        `json:"gmeet"`
		GmeetDate        string        `json:"gmeet_date"`
		QuizLinkPub      string        `json:"quizLinkPub,omitempty"`
		QuizLinkPreview  string        `json:"quizLinkPreview,omitempty"`
		WaLink           string        `json:"wa_link,omitempty"`
		Phone            string        `json:"phone"`
		Achievements     []ResAddItems `json:"achievements"`
		Extracurriculars []ResAddItems `json:"extracurriculars"`
		ResPayment       ResPayment    `json:"payments"`
		Reviews          []ResReview   `json:"reviews"`
		Faqs             []ResFaq      `json:"faqs"`
	}
	ReqCreateSchool struct {
		UserId        int
		Npsn          string `form:"npsn" validate:"required"`
		Name          string `form:"name" validate:"required"`
		Description   string `form:"description" validate:"required"`
		Image         string `form:"image"`
		Video         string `form:"video" validate:"required"`
		Pdf           string `form:"pdf"`
		Web           string `form:"web" validate:"required"`
		Province      string `form:"province" validate:"required"`
		City          string `form:"city" validate:"required"`
		District      string `form:"district" validate:"required"`
		Village       string `form:"village" validate:"required"`
		Detail        string `form:"detail" validate:"required"`
		ZipCode       string `form:"zipcode" validate:"required"`
		Students      string `form:"students" validate:"required"`
		Teachers      string `form:"teachers" validate:"required"`
		Staff         string `form:"staff" validate:"required"`
		Accreditation string `form:"accreditation" validate:"required"`
		Phone         string `form:"phone" validate:"required"`
	}
	ReqUpdateSchool struct {
		Id              int    `form:"id" validate:"required"`
		Npsn            string `form:"npsn" `
		Name            string `form:"name" `
		Description     string `form:"description" `
		Image           string `form:"image" `
		Video           string `form:"video" `
		Pdf             string `form:"pdf" `
		Web             string `form:"web" `
		Province        string `form:"province" `
		City            string `form:"city" `
		District        string `form:"district" `
		Village         string `form:"village" `
		Detail          string `form:"detail" `
		ZipCode         string `form:"zipcode" `
		Students        string `form:"students" `
		Teachers        string `form:"teachers" `
		Staff           string `form:"staff" `
		Phone           string `form:"phone"`
		Accreditation   string `form:"accreditation"`
		Gmeet           string `json:"gmeet"`
		GmeetDate       string `json:"gmeet_date"`
		QuizLinkPub     string `json:"quizLinkPub"`
		QuizLinkPreview string `json:"quizLinkPreview"`
	}
	Location struct {
		Province string `json:"province"`
		City     string `json:"city"`
		District string `json:"district"`
		Village  string `json:"village"`
		Detail   string `json:"detail"`
		ZipCode  string `json:"zipcode"`
	}
	ResUpdateSchool struct {
		Id            int      `json:"id"`
		Npsn          string   `json:"npsn"`
		Name          string   `json:"school_name"`
		Description   string   `json:"description"`
		Image         string   `json:"image"`
		Video         string   `json:"video"`
		Pdf           string   `json:"pdf"`
		Web           string   `json:"school_web"`
		Students      string   `json:"students"`
		Teachers      string   `json:"teachers"`
		Staff         string   `json:"staff"`
		Accreditation string   `json:"accreditation"`
		Location      Location `json:"location"`
	}
	Reviews struct {
		SchoolID uint   `json:"school_id" validate:"required"`
		UserID   uint   `json:"user_id"`
		Review   string `json:"review"  validate:"required"`
		User     User   `json:"-"`
	}
	StaticData struct {
		Title   string
		Content string
	}
)
