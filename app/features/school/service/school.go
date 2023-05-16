package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"sync"

	entity "github.com/education-hub/BE/app/entities"
	"github.com/education-hub/BE/app/features/school/repository"
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
	}
	SchoolService interface {
		Create(ctx context.Context, req entity.ReqCreateSchool, image multipart.File, pdf multipart.File) (int, error)
		Update(ctx context.Context, req entity.ReqUpdateSchool, image multipart.File, pdf multipart.File) (*entity.ResUpdateSchool, error)
		Search(searchval string) any
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
	}
)

func NewSchoolService(repo repository.SchoolRepo, dep dependency.Depend) SchoolService {
	return &school{repo: repo, dep: dep, validator: validator.New()}
}

func (s *school) Create(ctx context.Context, req entity.ReqCreateSchool, image multipart.File, pdf multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE CREATE SCHOOL REQ, Error: %v", err)
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}

	if err := s.repo.FindByNPSN(s.dep.Db.WithContext(ctx), req.Npsn); err == nil {
		return 0, errorr.NewBad("School Already Registered")
	}
	if err := helper.CheckNPSN(req.Npsn, s.dep.Log); err != nil {
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
			return 0, err
		}
	}
	data.UserID = uint(req.UserId)
	id, err2 := s.repo.Create(s.dep.Db.WithContext(ctx), data)
	if err2 != nil {
		return 0, err2
	}
	return id, nil
}
func (s *school) Update(ctx context.Context, req entity.ReqUpdateSchool, image multipart.File, pdf multipart.File) (*entity.ResUpdateSchool, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE REQUPDATE")
		return nil, errorr.NewBad("Missing Or Invalid Request Body")
	}
	if req.Npsn != "" {
		if err := s.repo.FindByNPSN(s.dep.Db.WithContext(ctx), req.Npsn); err == nil {
			return nil, errorr.NewBad("School Already Registered")
		}
		if err := helper.CheckNPSN(req.Npsn, s.dep.Log); err != nil {
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
			return nil, err
		}
		data.Image = filename
		image.Close()
	}
	if pdf != nil {
		filename := fmt.Sprintf("%s_%s_%s", "School_", req.Npsn, req.Pdf)
		if err := s.dep.Gcp.UploadFile(pdf, filename); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			return nil, err
		}
		data.Pdf = filename
		pdf.Close()
	}
	resdata, err := s.repo.Update(s.dep.Db.WithContext(ctx), data)
	if err != nil {
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
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Achv_", req.SchoolID, req.Image)
	if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
		s.dep.Log.Errorf("Error Service : %v", err)
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
		return 0, err
	}
	return res, err
}

func (s *school) DeleteAchievement(ctx context.Context, id int) error {
	if err := s.repo.DeleteAchievement(s.dep.Db.WithContext(ctx), id); err != nil {
		return err
	}
	return nil
}

func (s *school) UpdateAchievement(ctx context.Context, req entity.ReqUpdateAchievemnt, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Achievement REQ, Error: %v", err)
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Achv_", req.Id, req.Image)
	data := entity.Achievement{
		Description: req.Description,
		Image:       filename,
		Title:       req.Title,
	}
	data.ID = uint(req.Id)
	res, err := s.repo.UpdateAchievement(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		return 0, err
	}
	if image != nil {
		if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
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
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Extra_", req.SchoolID, req.Image)
	if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
		s.dep.Log.Errorf("Error Service : %v", err)
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
		return 0, err
	}
	return res, err
}

func (s *school) DeleteExtracurricular(ctx context.Context, id int) error {
	if err := s.repo.DeleteExtracurricular(s.dep.Db.WithContext(ctx), id); err != nil {
		return err
	}
	return nil
}

func (s *school) UpdateExtracurricular(ctx context.Context, req entity.ReqUpdateExtracurricular, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Extracurricular REQ, Error: %v", err)
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Extra_", req.Id, req.Image)
	data := entity.Extracurricular{
		Description: req.Description,
		Image:       filename,
		Title:       req.Title,
	}
	data.ID = uint(req.Id)
	res, err := s.repo.UpdateExtracurricular(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		return 0, err
	}
	if image != nil {
		if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			image.Close()
			return 0, err
		}
		image.Close()
	}
	return int(res.SchoolID), nil
}
func (s *school) GetByUid(ctx context.Context, uid int) (*entity.ResDetailSchool, error) {
	data, err := s.repo.GetByUid(s.dep.Db.WithContext(ctx), uid)
	if err != nil {
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
		}
		res.Achievements = append(res.Achievements, achivement)
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
func (s *school) GetByid(ctx context.Context, id int) (*entity.ResDetailSchool, error) {

	data, err := s.repo.GetByUid(s.dep.Db.WithContext(ctx), id)
	if err != nil {
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
		}
		res.Achievements = append(res.Achievements, achivement)
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
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	data := entity.Faq{
		SchoolID: uint(req.SchoolId),
		Question: req.Question,
		Answer:   req.Answer,
	}
	res, err := s.repo.AddFaq(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		return 0, err
	}
	return res, err
}

func (s *school) DeleteFaq(ctx context.Context, id int) error {
	if err := s.repo.DeleteFaq(s.dep.Db.WithContext(ctx), id); err != nil {
		return err
	}
	return nil
}

func (s *school) UpdateFaq(ctx context.Context, req entity.ReqUpdateFaq) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Faq REQ, Error: %v", err)
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	data := entity.Faq{
		Question: req.Question,
		Answer:   req.Answer,
	}
	data.ID = uint(req.Id)
	res, err := s.repo.UpdateFaq(s.dep.Db.WithContext(ctx), data)
	if err != nil {
		return 0, err
	}
	return int(res.SchoolID), nil
}

func (s *school) AddPayment(ctx context.Context, req entity.ReqAddPayment, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Payment REQ, Error: %v", err)
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Payment_", req.SchoolID, req.Image)
	if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
		s.dep.Log.Errorf("Error Service : %v", err)
		image.Close()
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
		return 0, err
	}
	return res, err
}

func (s *school) DeletePayment(ctx context.Context, id int) error {
	if err := s.repo.DeletePayment(s.dep.Db.WithContext(ctx), id); err != nil {
		return err
	}
	return nil
}

func (s *school) UpdatePayment(ctx context.Context, req entity.ReqUpdatePayment, image multipart.File) (int, error) {
	if err := s.validator.Struct(req); err != nil {
		s.dep.Log.Errorf("[ERROR] WHEN VALIDATE Add Payment REQ, Error: %v", err)
		return 0, errorr.NewBad("Missing or Invalid Request Body")
	}
	filename := fmt.Sprintf("%s_%d_%s", "Payment_", req.ID, req.Image)
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
		return 0, err
	}
	if image != nil {
		if err := s.dep.Gcp.UploadFile(image, filename); err != nil {
			s.dep.Log.Errorf("Error Service : %v", err)
			image.Close()
			return 0, err
		}
		image.Close()
	}
	return int(res.SchoolID), nil
}
