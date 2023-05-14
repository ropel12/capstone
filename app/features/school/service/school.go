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
			return
		}
		data.Image = filename
		image.Close()
		errchan <- nil
		return

	}()
	go func() {
		defer wg.Done()
		filename := fmt.Sprintf("%s_%s_%s", "School_", req.Npsn, req.Pdf)
		if err1 := s.dep.Gcp.UploadFile(pdf, filename); err1 != nil {
			s.dep.Log.Errorf("Error Service : %v", err1)
			errchan <- err1
			return
		}
		data.Pdf = filename
		image.Close()
		errchan <- nil
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
