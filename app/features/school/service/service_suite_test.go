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
		Depend.Config = &config.Config{GmapsKey: os.Getenv("GMAPS")}
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

	Context("Memperbaharui Data Sekolah", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.Update(ctx, entity.ReqUpdateSchool{}, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Missing Or Invalid Request Body"))
			})
		})

		When("Npsn sudah terdaftar pada database", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "motivasion letter.pdf",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("School Already Registered"))
			})
		})
		When("Npsn tidak terdaftar pada data kementrian pendidikan", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "2010025112",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "motivasion letter.pdf",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("NPSN not registered"))
			})
		})
		When("Format gambar tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.php",
					Pdf:           "motivasion letter.pdf",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("File type not allowed"))
			})
		})
		When("Format pdf tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "brochure.php",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("File type not allowed"))
			})
		})
		When("Terjadi kesalahn qury database", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "brochure.pdf",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil memperbaharui data", func() {
			BeforeEach(func() {
				res := entity.School{
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "brochure.php",
					Accreditation: "A"}
				res.ID = uint(1)
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(&res, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "brochure.pdf",
					Accreditation: "A"}
				res, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).Should(BeNil())
				Expect(res.Npsn).To(Equal("20100251"))
			})
		})
	})
	Context("Mencari Detail Alamat Sekolah", func() {
		When("Data Sekolah Tidak Ditemukan", func() {
			It("Akan Mengembalikan nil", func() {
				res := SchoolService.Search("exsdwqeqewqxqwxwqxqwxqxwwqxwwxwqxqxwqxwqxwqxwqxwwwwwwwwwwwwwwwwwwwwwwwwwwxwxwq")
				Expect(res).Should(BeEmpty())
			})
		})
		When("Data Sekolah Tidak Ditemukan", func() {
			It("Akan Mengembalikan nil", func() {
				res := SchoolService.Search("smpn 6 jakarta")
				Expect(res).ShouldNot(BeEmpty())
			})
		})
	})
	Context("Add Achievement", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.AddAchievement(ctx, entity.ReqAddAchievemnt{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddAchievemnt{SchoolID: 1, Description: "test", Image: "gambar.php", Title: "tes"}
				_, err := SchoolService.AddAchievement(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Format gambar tidak sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddAchievemnt{SchoolID: 1, Description: "test", Image: "gambar.php", Title: "tes"}
				_, err := SchoolService.AddAchievement(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("AddAchievement", mock.Anything, mock.Anything).Return(0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddAchievemnt{SchoolID: 1, Description: "test", Image: "gambar.jpg", Title: "tes"}
				_, err := SchoolService.AddAchievement(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menambahakan Prestasi", func() {
			BeforeEach(func() {
				Mock.On("AddAchievement", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan Id Sekolah", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddAchievemnt{SchoolID: 1, Description: "test", Image: "gambar.jpg", Title: "tes"}
				res, err := SchoolService.AddAchievement(ctx, req, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})

	})

	Context("Update Achievement", func() {
		When("id tidak ada", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateAchievement(ctx, entity.ReqUpdateAchievemnt{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("UpdateAchievement", mock.Anything, mock.Anything).Return(&entity.Achievement{}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateAchievement(ctx, entity.ReqUpdateAchievemnt{Id: 1, Image: "backdoor.aspx"}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("UpdateAchievement", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateAchievement(ctx, entity.ReqUpdateAchievemnt{Id: 1, Image: "img.jpg"}, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil memperbahrui data achievement", func() {
			BeforeEach(func() {
				Mock.On("UpdateAchievement", mock.Anything, mock.Anything).Return(&entity.Achievement{SchoolID: 1}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				res, err := SchoolService.UpdateAchievement(ctx, entity.ReqUpdateAchievemnt{Id: 1, Image: "img.jpg"}, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})
	})

	Context("Delete Achievement", func() {
		When("Id tidak ditemukan", func() {
			BeforeEach(func() {
				Mock.On("DeleteAchievement", mock.Anything, mock.Anything).Return(errors.New("Id not found")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteAchievement(ctx, 9999)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Id not found"))
			})
		})
		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("DeleteAchievement", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteAchievement(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menghapus data achievement", func() {
			BeforeEach(func() {
				Mock.On("DeleteAchievement", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteAchievement(ctx, 1)
				Expect(err).Should(BeNil())
			})
		})

	})
	Context("Add Extracurricular", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.AddExtracurricular(ctx, entity.ReqAddExtracurricular{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddExtracurricular{SchoolID: 1, Description: "test", Image: "gambar.php", Title: "tes"}
				_, err := SchoolService.AddExtracurricular(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Format gambar tidak sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddExtracurricular{SchoolID: 1, Description: "test", Image: "gambar.php", Title: "tes"}
				_, err := SchoolService.AddExtracurricular(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("AddExtracurricular", mock.Anything, mock.Anything).Return(0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddExtracurricular{SchoolID: 1, Description: "test", Image: "gambar.jpg", Title: "tes"}
				_, err := SchoolService.AddExtracurricular(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menambahakan Prestasi", func() {
			BeforeEach(func() {
				Mock.On("AddExtracurricular", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan Id Sekolah", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddExtracurricular{SchoolID: 1, Description: "test", Image: "gambar.jpg", Title: "tes"}
				res, err := SchoolService.AddExtracurricular(ctx, req, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})

	})

	Context("Update Extracurricular", func() {
		When("id tidak ada", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateExtracurricular(ctx, entity.ReqUpdateExtracurricular{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("UpdateExtracurricular", mock.Anything, mock.Anything).Return(&entity.Extracurricular{}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateExtracurricular(ctx, entity.ReqUpdateExtracurricular{Id: 1, Image: "backdoor.aspx"}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("UpdateExtracurricular", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateExtracurricular(ctx, entity.ReqUpdateExtracurricular{Id: 1, Image: "img.jpg"}, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil memperbahrui data Extracurricular", func() {
			BeforeEach(func() {
				Mock.On("UpdateExtracurricular", mock.Anything, mock.Anything).Return(&entity.Extracurricular{SchoolID: 1}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				res, err := SchoolService.UpdateExtracurricular(ctx, entity.ReqUpdateExtracurricular{Id: 1, Image: "img.jpg"}, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})
	})

	Context("Delete Extracurricular", func() {
		When("Id tidak ditemukan", func() {
			BeforeEach(func() {
				Mock.On("DeleteExtracurricular", mock.Anything, mock.Anything).Return(errors.New("Id not found")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteExtracurricular(ctx, 9999)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Id not found"))
			})
		})
		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("DeleteExtracurricular", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteExtracurricular(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menghapus data Extracurricular", func() {
			BeforeEach(func() {
				Mock.On("DeleteExtracurricular", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteExtracurricular(ctx, 1)
				Expect(err).Should(BeNil())
			})
		})

	})

})
