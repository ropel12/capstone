package service_test

import (
	"context"
	"errors"
	"mime/multipart"
	"os"
	"testing"

	entity "github.com/education-hub/BE/app/entities"
	mocks "github.com/education-hub/BE/app/features/school/mocks/repository"
	school "github.com/education-hub/BE/app/features/school/service"
	"github.com/education-hub/BE/config"
	dependcy "github.com/education-hub/BE/config/dependency"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("school", func() {
	var Mock *mocks.SchoolRepo
	var SchoolService school.SchoolService
	var Depend dependcy.Depend
	var ctx context.Context
	BeforeEach(func() {
		Depend.Db = config.GetConnectionTes()
		log := logrus.New()
		Depend.Log = log
		ctx = context.Background()
		Mock = mocks.NewSchoolRepo(GinkgoT())
		SchoolService = school.NewSchoolService(Mock, Depend)
	})

	Context("Menambah Sekolah Baru", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.Create(ctx, entity.ReqCreateSchool{}, image, pdf)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("NPSN sudah terdaftar", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "20100251",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.jpg",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.pdf",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      21,
					Teachers:      21,
					Staff:         21,
					Accreditation: "A"}
				_, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("School Already Registered"))
			})
		})
		When("NPSN tidak valid", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "201002512323",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.jpg",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.pdf",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      21,
					Teachers:      21,
					Staff:         21,
					Accreditation: "A"}
				_, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("NPSN not registered"))
			})
		})
		When("Tipe file bukan merupakan gambar atau pdf", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "20100251",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.js",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.java",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      21,
					Teachers:      21,
					Staff:         21,
					Accreditation: "A"}
				_, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("File type not allowed"))
			})
		})
		When("Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
				Mock.On("Create", mock.Anything, mock.Anything).Return(0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "20100251",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.jpg",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.pdf",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      21,
					Teachers:      21,
					Staff:         21,
					Accreditation: "A"}
				_, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menambahkan Data Sekolah", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
				Mock.On("Create", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan id  dan error bernailai nil", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "20100251",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.jpg",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.pdf",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      21,
					Teachers:      21,
					Staff:         21,
					Accreditation: "A"}
				id, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).Should(BeNil())
				Expect(id).To(Equal(1))
			})
		})
	})

})
