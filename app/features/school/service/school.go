package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"strings"
	"sync"
	"time"

	entity "github.com/education-hub/BE/app/entities"
	"github.com/education-hub/BE/app/features/school/repository"
	user "github.com/education-hub/BE/app/features/user/repository"
	"github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/errorr"
	"github.com/education-hub/BE/helper"
	"github.com/education-hub/BE/pkg"
	"github.com/go-playground/validator"
)

type (
	school struct {
		repo      repository.SchoolRepo
		validator *validator.Validate
		dep       dependency.Depend
		userrepo  user.UserRepo
	}
	SchoolService interface {
		Create(ctx context.Context, req entity.ReqCreateSchool, image multipart.File, pdf multipart.File) (int, error)
		Update(ctx context.Context, req entity.ReqUpdateSchool, image multipart.File, pdf multipart.File) (*entity.ResUpdateSchool, error)
		Delete(ctx context.Context, id int, uid int) error
		Search(searchval string) any
		GetAll(ctx context.Context, page, limit int, search string) (*entity.Response, error)
		AddAchievement(ctx context.Context, req entity.ReqAddAchievemnt, image multipart.File) (int, error)
		DeleteAchievement(ctx context.Context, id int) error
		UpdateAchievement(ctx context.Context, req entity.ReqUpdateAchievemnt, image multipart.File) (int, error)
		GetByUid(ctx context.Context, uid int) (*entity.ResDetailSchool, error)
		GetByid(ctx context.Context, id int) (*entity.ResDetailSchool, error)
		AddExtracurricular(ctx context.Context, req entity.ReqAddExtracurricular, image multipart.File) (int, error)
		DeleteExtracurricular(ctx context.Context, id int) error
		UpdateExtracurricular(ctx context.Context, req entity.ReqUpdateExtracurricular, image multipart.File) (int, error)
		AddFaq(ctx context.Context, req entity.ReqAddFaq) (int, error)
		DeleteFaq(ctx context.Context, id int) error
		UpdateFaq(ctx context.Context, req entity.ReqUpdateFaq) (int, error)
		AddPayment(ctx context.Context, req entity.ReqAddPayment, image multipart.File) (int, error)
		DeletePayment(ctx context.Context, id int) error
		UpdatePayment(ctx context.Context, req entity.ReqUpdatePayment, image multipart.File) (int, error)
		CreateSubmission(ctx context.Context, req entity.ReqCreateSubmission, studentph, signstudent, signparent multipart.File) (int, error)
		UpdateProgressByid(ctx context.Context, id int, status string) (int, error)
		GetAllProgressByUid(ctx context.Context, uid int) ([]entity.ResAllProgress, error)
		GetProgressById(ctx context.Context, id int) (*entity.ResDetailProgress, error)
		GetAllProgressAndSubmissionByuid(ctx context.Context, uid int) ([]entity.ResAllProgressSubmission, error)
		GetSubmissionByid(ctx context.Context, id int) (*entity.ResDetailSubmission, error)
		AddReview(ctx context.Context, req entity.Reviews) (int, error)
		CreateQuiz(ctx context.Context, req []entity.ReqAddQuiz) error
		GetTestResult(ctx context.Context, uid int) ([]pkg.TestResult, error)
	}
)

func NewSchoolService(repo repository.SchoolRepo, dep dependency.Depend, user user.UserRepo) SchoolService {
	return &school{repo: repo, dep: dep, validator: validator.New(), userrepo: user}
}

func (s *school) Create(ctx context.Context, req entity.ReqCreateSchool, image multipart.File, pdf multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE CREATE SCHOOL REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}

	if err := s.repo.FindByNPSN(s.dep.Db.WithContext(ctx), req.Npsn); err == nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("School Already Registered")
	}
	if err := helper.CheckNPSN(req.Npsn, s.dep.Log); err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	data := entity.School{
		Npsn:          req.Npsn,
		Name:          req.Name,
		Description:   req.Description,
		Image:         req.Image,
		Video:         req.Video,
		Pdf:           req.Pdf,
		Web:           req.Web,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		Village:       req.Village,
		Detail:        req.Detail,
		ZipCode:       req.ZipCode,
		Students:      req.Students,
		Teachers:      req.Teachers,
		Staff:         req.Staff,
		Accreditation: req.Accreditation,
	}
	if image != nil && pdf != nil {
		wg := &sync.WaitGroup{}
		errchan := make(chan error, 2)
		wg.Add(2)
		go func() {
			defer wg.Done()
			filename := fmt.Sprintf("%s_%s_%s", "School_", req.Npsn, req.Image)
			if err1 := s.dep.Gcp.UploadFile(image, filename); err1 != nil {
				s.dep.Log.Errorf("Error Service : %v", err1)
				errchan <- err1
				image.Close()
				return
			}
			data.Image = filename
			errchan <- nil
			image.Close()
			return

		}()
		go func() {
			defer wg.Done()
			filename := fmt.Sprintf("%s_%s_%s", "School_", req.Npsn, req.Pdf)
			if err1 := s.dep.Gcp.UploadFile(pdf, filename); err1 != nil {
				s.dep.Log.Errorf("Error Service : %v", err1)
				errchan <- err1
				pdf.Close()
				return
			}
			data.Pdf = filename
			errchan <- nil
			pdf.Close()
			return
		}()
		wg.Wait()
		close(errchan)
		for err := range errchan {
			if err != nil {
				s.dep.PromErr["error"] = err.Error()
				return 0, err
			}
		}
	}
	data.UserID = uint(req.UserId)
	id, err2 := s.repo.Create(s.dep.Db.WithContext(ctx), data)
	if err2 != nil {
		s.dep.PromErr["error"] = err2.Error()
		return 0, err2
	}
	return id, nil
}
func (s *school) Delete(ctx context.Context, id int, uid int) error {
	if err := s.repo.Delete(s.dep.Db.WithContext(ctx), id, uid); err != nil {
		s.dep.PromErr["error"] = err.Error()
		return err
	}
	return nil
}
func (s *school) Update(ctx context.Context, req entity.ReqUpdateSchool, image multipart.File, pdf multipart.File) (*entity.ResUpdateSchool, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE REQUPDATE")
		s.dep.PromErr["error"] = err.Error()
		return nil, errorr.NewBad("Missing Or Invalid Request Body")
	}
	if req.Npsn != "" {
		if err := s.repo.FindByNPSN(s.dep.Db.WithContext(ctx), req.Npsn); err == nil {
			s.dep.PromErr["error"] = "School Already Registered"
			return nil, errorr.NewBad("School Already Registered")
		}
		if err := helper.CheckNPSN(req.Npsn, s.dep.Log); err != nil {
			s.dep.PromErr["error"] = err.Error()
			return nil, err
		}
	}
	data := entity.School{
		Npsn:            req.Npsn,
		Name:            req.Name,
		Description:     req.Description,
		Image:           req.Image,
		Video:           req.Video,
		Pdf:             req.Pdf,
		Web:             req.Web,
		Province:        req.Province,
		City:            req.City,
		District:        req.District,
		Village:         req.Village,
		Detail:          req.Detail,
		ZipCode:         req.ZipCode,
		Students:        req.Students,
		Teachers:        req.Teachers,
		Staff:           req.Staff,
		Accreditation:   req.Accreditation,
		Gmeet:           req.Gmeet,
		QuizLinkPub:     req.QuizLinkPub,
		QuizLinkPreview: req.QuizLinkPreview,
	}
	data.ID = uint(req.Id)
	if image != nil {
		filename := fmt.Sprintf("%s_%s_%s", "School_", req.Npsn, req.Image)
		if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			s.dep.PromErr["error"] = err.Error()
			return nil, err
		}
		data.Image = filename
		image.Close()
	}
	if pdf != nil {
		filename := fmt.Sprintf("%s_%s_%s", "School_", req.Npsn, req.Pdf)
		if err := s.dep.Gcp.UploadFile(pdf, filename); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			s.dep.PromErr["error"] = err.Error()
			return nil, err
		}
		data.Pdf = filename
		pdf.Close()
	}
	resdata, err := s.repo.Update(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	res := entity.ResUpdateSchool{
		Id:            int(resdata.ID),
		Npsn:          resdata.Npsn,
		Name:          resdata.Name,
		Description:   resdata.Description,
		Image:         resdata.Image,
		Video:         resdata.Video,
		Pdf:           resdata.Pdf,
		Web:           resdata.Web,
		Students:      resdata.Students,
		Teachers:      resdata.Teachers,
		Staff:         resdata.Staff,
		Accreditation: resdata.Accreditation,
		Location: entity.Location{
			Province: resdata.Province,
			City:     resdata.City,
			District: resdata.District,
			Village:  resdata.Village,
			Detail:   resdata.Detail,
			ZipCode:  resdata.ZipCode,
		},
	}
	return &res, nil
}
func (s *school) Search(searchval string) any {
	return pkg.NewClientGmaps(s.dep.Config.GmapsKey, s.dep.Log).Search(searchval)
}

func (s *school) AddAchievement(ctx context.Context, req entity.ReqAddAchievemnt, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Achievement REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Achv_", req.SchoolID, req.Image)
	if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
		s.dep.Log.Errorf("Error Service : %v", err)
		s.dep.PromErr["error"] = err.Error()
		image.Close()
		return 0, err
	}
	image.Close()
	data := entity.Achievement{
		SchoolID:    req.SchoolID,
		Description: req.Description,
		Image:       filename,
		Title:       req.Title,
	}
	res, err := s.repo.AddAchievement(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	return res, err
}

func (s *school) DeleteAchievement(ctx context.Context, id int) error {
	if err := s.repo.DeleteAchievement(s.dep.Db.WithContext(ctx), id); err != nil {
		s.dep.PromErr["error"] = err.Error()
		return err
	}
	return nil
}

func (s *school) UpdateAchievement(ctx context.Context, req entity.ReqUpdateAchievemnt, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Achievement REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Achv_", req.Id, req.Image)
	if image != nil {
		req.Image = filename
	}
	data := entity.Achievement{
		Description: req.Description,
		Image:       req.Image,
		Title:       req.Title,
	}
	data.ID = uint(req.Id)
	res, err := s.repo.UpdateAchievement(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	if image != nil {
		if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			s.dep.PromErr["error"] = err.Error()
			image.Close()
			return 0, err
		}
		image.Close()
	}
	return int(res.SchoolID), nil
}

func (s *school) AddExtracurricular(ctx context.Context, req entity.ReqAddExtracurricular, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Extracurricular REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Extra_", req.SchoolID, req.Image)
	if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
		s.dep.Log.Errorf("Error Service : %v", err)
		s.dep.PromErr["error"] = err.Error()
		image.Close()
		return 0, err
	}
	image.Close()
	data := entity.Extracurricular{
		SchoolID:    req.SchoolID,
		Description: req.Description,
		Image:       filename,
		Title:       req.Title,
	}
	res, err := s.repo.AddExtracurricular(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	return res, err
}

func (s *school) DeleteExtracurricular(ctx context.Context, id int) error {
	if err := s.repo.DeleteExtracurricular(s.dep.Db.WithContext(ctx), id); err != nil {
		s.dep.PromErr["error"] = err.Error()
		return err
	}
	return nil
}

func (s *school) UpdateExtracurricular(ctx context.Context, req entity.ReqUpdateExtracurricular, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Extracurricular REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Extra_", req.Id, req.Image)
	if image != nil {
		req.Image = filename
	}
	data := entity.Extracurricular{
		Description: req.Description,
		Image:       req.Image,
		Title:       req.Title,
	}
	data.ID = uint(req.Id)
	res, err := s.repo.UpdateExtracurricular(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	if image != nil {
		if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			image.Close()
			s.dep.PromErr["error"] = err.Error()
			return 0, err
		}
		image.Close()
	}
	return int(res.SchoolID), nil
}
func (s *school) GetByUid(ctx context.Context, uid int) (*entity.ResDetailSchool, error) {
	data, err := s.repo.GetByUid(s.dep.Db.WithContext(ctx), uid)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	res := entity.ResDetailSchool{
		Id:              int(data.ID),
		Npsn:            data.Npsn,
		Name:            data.Name,
		Description:     data.Description,
		Image:           data.Image,
		Video:           data.Video,
		Pdf:             data.Pdf,
		Web:             data.Web,
		Province:        data.Province,
		City:            data.City,
		District:        data.District,
		Village:         data.Village,
		Detail:          data.Detail,
		ZipCode:         data.ZipCode,
		Students:        data.Students,
		Teachers:        data.Teachers,
		Staff:           data.Staff,
		Accreditation:   data.Accreditation,
		Gmeet:           data.Gmeet,
		QuizLinkPub:     data.QuizLinkPub,
		QuizLinkPreview: fmt.Sprintf("https://go-event.online/quiz/%s?preview=1", data.QuizLinkPub),
	}
	for _, val := range data.Achievements {
		achivement := entity.ResAddItems{
			Id:          int(val.ID),
			Name:        val.Title,
			Img:         val.Image,
			Description: val.Description,
		}
		res.Achievements = append(res.Achievements, achivement)
	}
	for _, val := range data.Faqs {
		faq := entity.ResFaq{
			Id:       int(val.ID),
			Question: val.Question,
			Answer:   val.Answer,
		}
		res.Faqs = append(res.Faqs, faq)
	}

	for _, val := range data.Extracurriculars {
		extracurricular := entity.ResAddItems{
			Name:        val.Title,
			Img:         val.Image,
			Description: val.Description,
			Id:          int(val.ID),
		}
		res.Extracurriculars = append(res.Extracurriculars, extracurricular)
	}
	for _, val := range data.Payments {
		if val.Type == "one" {
			onetime := entity.ResPaymentType{
				Id:          int(val.ID),
				Img:         val.Image,
				Description: val.Description,
				Price:       val.Price,
			}
			res.ResPayment.OneTime = append(res.ResPayment.OneTime, onetime)
		} else {
			interval := entity.ResPaymentType{
				Id:          int(val.ID),
				Img:         val.Image,
				Description: val.Description,
				Price:       val.Price,
			}
			res.ResPayment.Interval = append(res.ResPayment.OneTime, interval)
		}
	}
	return &res, nil
}
func (s *school) GetAll(ctx context.Context, page, limit int, search string) (*entity.Response, error) {
	offset := (page - 1) * limit
	data, total, err := s.repo.GetAll(s.dep.Db.WithContext(ctx), limit, offset, search)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	schools := []entity.ResAllSchool{}
	for _, val := range data {
		school := entity.ResAllSchool{
			ID:            int(val.ID),
			Name:          val.Name,
			Location:      fmt.Sprintf("%s, %s", val.City, val.Province),
			AdminName:     val.User.Username,
			Accreditation: val.Accreditation,
			Image:         val.Image,
		}
		schools = append(schools, school)
	}
	res := entity.Response{
		Limit:     limit,
		Page:      page,
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		TotalData: total,
		Data:      schools,
	}
	return &res, nil
}
func (s *school) GetByid(ctx context.Context, id int) (*entity.ResDetailSchool, error) {

	data, err := s.repo.GetById(s.dep.Db.WithContext(ctx), id)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	res := entity.ResDetailSchool{
		Id:              int(data.ID),
		Npsn:            data.Npsn,
		Name:            data.Name,
		Description:     data.Description,
		Image:           data.Image,
		Video:           data.Video,
		Pdf:             data.Pdf,
		Web:             data.Web,
		Province:        data.Province,
		City:            data.City,
		District:        data.District,
		Village:         data.Village,
		Detail:          data.Detail,
		ZipCode:         data.ZipCode,
		Students:        data.Students,
		Teachers:        data.Teachers,
		Staff:           data.Staff,
		Accreditation:   data.Accreditation,
		Gmeet:           data.Gmeet,
		QuizLinkPub:     data.QuizLinkPub,
		QuizLinkPreview: data.QuizLinkPreview,
	}

	for _, val := range data.Achievements {
		achivement := entity.ResAddItems{
			Name:        val.Title,
			Img:         val.Image,
			Description: val.Description,
			Id:          int(val.ID),
		}
		res.Achievements = append(res.Achievements, achivement)
	}
	for _, val := range data.Reviews {
		review := entity.ResReview{
			UserImage: val.User.Image,
			Review:    val.Review,
		}
		res.Reviews = append(res.Reviews, review)
	}
	for _, val := range data.Faqs {
		faq := entity.ResFaq{
			Id:       int(val.ID),
			Question: val.Question,
			Answer:   val.Answer,
		}
		res.Faqs = append(res.Faqs, faq)
	}
	for _, val := range data.Extracurriculars {
		extracurricular := entity.ResAddItems{
			Name:        val.Title,
			Img:         val.Image,
			Description: val.Description,
			Id:          int(val.ID),
		}
		res.Extracurriculars = append(res.Extracurriculars, extracurricular)
	}
	for _, val := range data.Payments {
		if val.Type == "one" {
			onetime := entity.ResPaymentType{
				Id:          int(val.ID),
				Img:         val.Image,
				Description: val.Description,
				Price:       val.Price,
			}
			res.ResPayment.OneTime = append(res.ResPayment.OneTime, onetime)
		} else {
			interval := entity.ResPaymentType{
				Id:          int(val.ID),
				Img:         val.Image,
				Description: val.Description,
				Price:       val.Price,
			}
			res.ResPayment.Interval = append(res.ResPayment.OneTime, interval)
		}
	}

	return &res, nil
}

func (s *school) AddFaq(ctx context.Context, req entity.ReqAddFaq) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Faq REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	data := entity.Faq{
		SchoolID: uint(req.SchoolId),
		Question: req.Question,
		Answer:   req.Answer,
	}
	res, err := s.repo.AddFaq(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	return res, err
}

func (s *school) DeleteFaq(ctx context.Context, id int) error {
	if err := s.repo.DeleteFaq(s.dep.Db.WithContext(ctx), id); err != nil {
		s.dep.PromErr["error"] = err.Error()
		return err
	}
	return nil
}

func (s *school) UpdateFaq(ctx context.Context, req entity.ReqUpdateFaq) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Faq REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	data := entity.Faq{
		Question: req.Question,
		Answer:   req.Answer,
	}
	data.ID = uint(req.Id)
	res, err := s.repo.UpdateFaq(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	return int(res.SchoolID), nil
}

func (s *school) AddPayment(ctx context.Context, req entity.ReqAddPayment, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Payment REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()

		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Payment_", req.SchoolID, req.Image)
	if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
		s.dep.Log.Errorf("Error Service : %v", err)
		image.Close()
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	image.Close()
	typee := "one"
	if *req.Interval != 0 {
		typee = "interval"
	}
	data := entity.Payment{
		SchoolID:    uint(req.SchoolID),
		Description: req.Description,
		Price:       req.Price,
		Interval:    *req.Interval,
		Image:       filename,
		Type:        typee,
	}
	res, err := s.repo.AddPayment(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	return res, err
}

func (s *school) DeletePayment(ctx context.Context, id int) error {
	if err := s.repo.DeletePayment(s.dep.Db.WithContext(ctx), id); err != nil {
		s.dep.PromErr["error"] = err.Error()
		return err
	}
	return nil
}

func (s *school) UpdatePayment(ctx context.Context, req entity.ReqUpdatePayment, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Payment REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Payment_", req.ID, req.Image)
	if image != nil {
		req.Image = filename
	}
	typee := ""
	if req.Interval != nil {
		if *req.Interval != 0 {
			typee = "interval"
		} else {
			typee = "one"
		}
	}
	Interval := -1
	if req.Interval != nil {
		Interval = *req.Interval
	}
	data := entity.Payment{
		Description: req.Description,
		Image:       req.Image,
		Price:       req.Price,
		Interval:    Interval,
		Type:        typee,
	}
	data.ID = uint(req.ID)
	res, err := s.repo.UpdatePayment(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	if image != nil {
		if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
			s.dep.PromErr["error"] = err.Error()
			s.dep.Log.Errorf("Error Service : %v", err)
			image.Close()
			return 0, err
		}
		image.Close()
	}
	return int(res.SchoolID), nil
}

func (s *school) CreateSubmission(ctx context.Context, req entity.ReqCreateSubmission, studentph, studentsign, parentsign multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR]WHEN VALIDATE CREATESUBMISSION REQ, err : %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing Or Invalid Req Body")
	}
	studentphoname := fmt.Sprintf("%s_%d_%s", "Student_", req.UserID, req.StudentPhoto)
	studentsignname := fmt.Sprintf("%s_%d_%s", "StudentSign_", req.UserID, req.StudentSignature)
	parentsignname := fmt.Sprintf("%s_%d_%s", "ParentSign_", req.UserID, req.ParentSignature)
	wg := &sync.WaitGroup{}
	wg.Add(3)
	errchan := make(chan error)
	go func() {
		defer wg.Done()
		if err := s.dep.Gcp.UploadFile(studentph, studentphoname); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			studentph.Close()
			errchan <- err
			return
		}
		studentph.Close()
		errchan <- nil
	}()
	go func() {
		defer wg.Done()
		if err := s.dep.Gcp.UploadFile(studentsign, studentsignname); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			studentsign.Close()
			errchan <- err
			return
		}
		studentsign.Close()
		errchan <- nil
	}()
	go func() {
		defer wg.Done()
		if err := s.dep.Gcp.UploadFile(parentsign, parentsignname); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			parentsign.Close()
			errchan <- err
			return
		}
		parentsign.Close()
		errchan <- nil
	}()
	wg.Wait()
	close(errchan)
	for err := range errchan {
		if err != nil {
			s.dep.PromErr["error"] = err.Error()
			return 0, err
		}
	}
	parentaddres := entity.ReqAdressSubmission{
		Province: req.ParentProvince,
		District: req.ParentDistrict,
		Village:  req.ParentVillage,
		ZipCode:  req.ParentZipCode,
		City:     req.ParentCity,
		Detail:   req.ParentDetail,
	}
	studentaddres := entity.ReqAdressSubmission{
		Province: req.StudentProvince,
		District: req.StudentDistrict,
		Village:  req.StudentVillage,
		ZipCode:  req.StudentZipCode,
		City:     req.StudentCity,
		Detail:   req.StudentDetail,
	}
	studentadd, _ := json.Marshal(studentaddres)
	parentadd, _ := json.Marshal(parentaddres)
	data := entity.Submission{
		SchoolID:         uint(req.SchoolID),
		UserID:           req.UserID,
		StudentPhoto:     studentphoname,
		StudentName:      req.StudentName,
		ParentName:       req.ParentName,
		ParentJob:        req.ParentJob,
		Religion:         req.Religion,
		ParentReligion:   req.Religion,
		PlaceDate:        req.PlaceDate,
		Gender:           req.Gender,
		GraduationFrom:   req.GraduationFrom,
		NISN:             req.NISN,
		ParentPhone:      req.ParentPhone,
		ParentSignature:  parentsignname,
		StudentSignature: studentsignname,
		Date:             time.Now().Format("2006-01-02"),
		ParentAddress:    string(parentadd),
		StudentAddress:   string(studentadd),
	}
	res, err := s.repo.CreateSubmission(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	return res, nil
}

func (s *school) UpdateProgressByid(ctx context.Context, id int, status string) (int, error) {
	if status != "Check File Registration" && status != "File Approved" && status != "Send Detail Costs Registration" &&
		status != "Failed File Approved" && status != "Failed Test Result" &&
		status != "Send Test Link" && status != "Check Test Result" &&
		status != "Test Result" && status != "Send Detail Costs Her-Registration" &&
		status != "Finish" {
		s.dep.PromErr["error"] = "Status Not Available"
		return 0, errorr.NewBad("Invalid Request Body")
	}
	res, err := s.repo.UpdateProgress(s.dep.Db.WithContext(ctx), id, status)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	if status == "File Approved" {
		user, err := s.userrepo.GetById(s.dep.Db.WithContext(ctx), int(res.SchoolID))
		if err != nil {
			s.dep.Log.Errorf("[ERROR]WHEN GETTING USER DATA: %v", err)
		}
		school, err := s.repo.GetById(s.dep.Db.WithContext(ctx), int(res.SchoolID))
		if err != nil {
			s.dep.Log.Errorf("[ERROR]WHEN GETTING SCHOOL DATA: %v", err)
		}
		if err := s.dep.Pusher.Publish(map[string]string{"username": user.Username, "type": "admission", "school_name": school.Name, "status": "File Approved"}, 2); err != nil {
			s.dep.Log.Errorf("Failed to publish to PusherJs: %v", err)
		}
	}
	if status == "Send Test Link" {
		schooldata, _ := s.repo.GetById(s.dep.Db.WithContext(ctx), int(res.SchoolID))
		userdata, _ := s.userrepo.GetById(s.dep.Db.WithContext(ctx), int(res.UserID))
		encodeddata, _ := json.Marshal(map[string]any{"email": userdata.Email, "name": userdata.FirstName + " " + userdata.SureName, "school": schooldata.Name, "test": schooldata.QuizLinkPub})
		if err := s.dep.Pusher.Publish(map[string]string{"username": userdata.Username, "type": "admission", "school_name": schooldata.Name, "status": "Send Test Link"}, 2); err != nil {
			s.dep.Log.Errorf("Failed to publish to PusherJs: %v", err)
		}
		go func() {
			if err := s.dep.Nsq.Publish("8", encodeddata); err != nil {
				s.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
			}
		}()

	}
	return int(res.ID), nil
}

func (s *school) GetAllProgressByUid(ctx context.Context, uid int) ([]entity.ResAllProgress, error) {
	data, err := s.repo.GetAllProgressByuid(s.dep.Db.WithContext(ctx), uid)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	res := []entity.ResAllProgress{}
	for _, val := range data {
		progrss := entity.ResAllProgress{
			SchoolName:  val.School.Name,
			SchoolImage: val.School.Image,
			SchoolWeb:   val.School.Web,
			ProgressId:  int(val.ID),
		}
		res = append(res, progrss)
	}
	return res, nil
}
func (s *school) GetProgressById(ctx context.Context, id int) (*entity.ResDetailProgress, error) {
	data, err := s.repo.GetProgressByid(s.dep.Db.WithContext(ctx), id)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	return &entity.ResDetailProgress{Id: int(data.ID), Status: data.Status}, nil
}

func (s *school) GetAllProgressAndSubmissionByuid(ctx context.Context, uid int) ([]entity.ResAllProgressSubmission, error) {
	data, err := s.repo.GetAllProgressAndSubmissionByuid(s.dep.Db.WithContext(ctx), uid)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	res := []entity.ResAllProgressSubmission{}
	for id, val := range data.Progresses {
		progresssubmission := entity.ResAllProgressSubmission{
			UserId:       int(val.User.ID),
			UserImage:    val.User.Image,
			UserName:     val.User.FirstName + " " + val.User.SureName,
			SubmissionId: int(data.Submissions[id].ID),
			ProgressId:   int(val.ID),
		}
		res = append(res, progresssubmission)
	}
	return res, nil
}
func (s *school) GetSubmissionByid(ctx context.Context, id int) (*entity.ResDetailSubmission, error) {
	data, err := s.repo.GetSubmissionByid(s.dep.Db.WithContext(ctx), id)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	studentaddress := entity.ReqAdressSubmission{}
	Parentaddress := entity.ReqAdressSubmission{}
	json.Unmarshal([]byte(data.StudentAddress), &studentaddress)
	json.Unmarshal([]byte(data.ParentAddress), &Parentaddress)
	res := entity.ResDetailSubmission{
		ParentSignature:  data.ParentSignature,
		StudentSignature: data.StudentSignature,
		DatePlace:        fmt.Sprintf("%s, %s", studentaddress.City, data.Date),
		StudentData: entity.StudentData{
			Photo:          data.StudentPhoto,
			Name:           data.StudentName,
			PlaceDate:      data.PlaceDate,
			NISN:           data.NISN,
			GraduationFrom: data.GraduationFrom,
			Religion:       data.Religion,
			Gender:         data.Gender,
			Adress:         studentaddress,
		},
		ParentData: entity.ParentData{
			Name:     data.ParentName,
			Job:      data.ParentJob,
			Religion: data.ParentReligion,
			Phone:    data.ParentPhone,
			Adress:   Parentaddress,
		},
	}
	return &res, nil
}

func (s *school) AddReview(ctx context.Context, req entity.Reviews) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Review REQ, Error: %v", err)
		s.dep.PromErr["error"] = err.Error()
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	if !s.dep.Validation.Validate(req.Review) {
		s.dep.PromErr["error"] = "comment contains bad words"
		return 0, errorr.NewBad("Your comment contains bad words")
	}
	res, err := s.repo.AddReview(s.dep.Db.WithContext(ctx), req)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return 0, err
	}
	return res, nil
}

func (s *school) CreateQuiz(ctx context.Context, req []entity.ReqAddQuiz) error {

	if len(req) == 0 {
		s.dep.PromErr["error"] = "the length of the data is 0"
		return errorr.NewBad("Missing or Invalid Request Body")
	}
	data, err := s.repo.GetById(s.dep.Db.WithContext(ctx), req[0].SchoolID)
	if err != nil {
		s.dep.PromErr["error"] = "school id doesn't exist"
		return err
	}
	quizlink, prev, result, err := s.dep.Quiz.CreateQuiz(data.Name, s.dep.Log)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return err
	}
	newdata := entity.School{
		QuizLinkPub:     quizlink,
		QuizLinkPreview: prev,
		QuizLinkResult:  result,
	}
	reqdata := entity.ReqDataQuiz{
		PubLink:    strings.ReplaceAll(quizlink, "https://www.flexiquiz.com/SC/N/", ""),
		Prevlink:   prev,
		ResultLink: result,
		Data:       req,
	}
	newdata.ID = uint(req[0].SchoolID)
	_, err = s.repo.Update(s.dep.Db.WithContext(ctx), newdata)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return err
	}
	encodeddata, _ := json.Marshal(reqdata)
	go func() {
		if err := s.dep.Nsq.Publish("9", encodeddata); err != nil {
			s.dep.PromErr["error"] = err.Error()
			s.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
		}
	}()

	return err
}

func (s *school) GetTestResult(ctx context.Context, uid int) ([]pkg.TestResult, error) {

	schooldata, err := s.repo.GetByUid(s.dep.Db.WithContext(ctx), uid)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	res, err := s.dep.Quiz.GetResult(schooldata.QuizLinkResult, s.dep.Log)
	if err != nil {
		s.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	if len(res) == 0 {
		s.dep.PromErr["error"] = "Data Not Found"
		return nil, errorr.NewBad("Data Not Found")
	}
	return res, nil
}
