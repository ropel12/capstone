package service_test

import (
	"context"
	"errors"
	"mime/multipart"
	"os"
	"testing"

	entity "github.com/education-hub/BE/app/entities/user"
	mocks "github.com/education-hub/BE/app/features/user/mocks/repository"
	user "github.com/education-hub/BE/app/features/user/service"
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

var _ = Describe("user", func() {
	var Mock *mocks.UserRepo
	var UserService user.UserService
	var Depend dependcy.Depend
	var ctx context.Context
	BeforeEach(func() {
		Depend.Db = config.GetConnectionTes()
		log := logrus.New()
		Depend.Log = log
		Mock = mocks.NewUserRepo(GinkgoT())
		UserService = user.NewUserService(Mock, Depend)

	})
	Context("User Login", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				err, _, _ := UserService.Login(ctx, entity.LoginReq{})
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Username Tidak terdaftar", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, "1321321ewqewq").Return(nil, errors.New("Username not registered")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Username not registered'", func() {
				_, _, err := UserService.Login(ctx, entity.LoginReq{Username: "1321321ewqewq", Password: "123"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Username not registered"))
			})
		})
		When("Password Salah", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, "satrio123").Return(&entity.User{Username: "satrio2@gmail.com", Password: "321"}, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'wrong password' ", func() {
				_, _, err := UserService.Login(ctx, entity.LoginReq{Username: "satrio123", Password: "123"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Wrong password"))
			})
		})
		When("Email Belum Diverifikasi", func() {
			BeforeEach(func() {
				data := &entity.User{Email: "satrio2@gmail.com", Password: "$2a$10$vu7o2Wl9LKyzTFkRDp7tc.VyoBB48nj97qyQjlgGCeQXJ067KZGQu", IsVerified: false}
				Mock.On("FindByUsername", mock.Anything, "satrio").Return(data, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Email Not Verified'", func() {
				_, _, err := UserService.Login(ctx, entity.LoginReq{Username: "satrio", Password: "123"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email Not Verified"))
			})
		})
		When("Berhasil Login", func() {
			BeforeEach(func() {
				data := &entity.User{Email: "satrio2@gmail.com", Password: "$2a$10$vu7o2Wl9LKyzTFkRDp7tc.VyoBB48nj97qyQjlgGCeQXJ067KZGQu", IsVerified: true}
				data.ID = 1
				data.Role = "student"
				Mock.On("FindByUsername", mock.Anything, "satrio").Return(data, nil).Once()
			})
			It("Akan Mengembalikan error", func() {
				uid, role, err := UserService.Login(ctx, entity.LoginReq{Username: "satrio", Password: "123"})
				Expect(err).Should(BeNil())
				Expect(uid).To(Equal(1))
				Expect(role).To(Equal("student"))
			})
		})

	})
	Context("User Register", func() {
		When("Request body kosong", func() {
			It("Akan Mengembalikan error", func() {
				err := UserService.Register(ctx, entity.RegisterReq{})
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Username sudah terdaftar", func() {
			BeforeEach(func() {
				data := &entity.User{Email: "satrio2@gmail.com"}
				Mock.On("FindByUsername", mock.Anything, "satrio").Return(data, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Username already registered'", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", FirstName: "satrio", LastName: "w", Password: "123", Address: "bogor ct", Username: "satrio", Role: "student"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Username already registered"))
			})
		})
		When("Email sudah terdaftar", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, "satrio").Return(nil, errors.New("error")).Once()
				data := &entity.User{Email: "satrio2@gmail.com", IsVerified: true}
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(data, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'email already registered'", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", FirstName: "satrio", LastName: "w", Password: "123", Address: "bogor ct", Username: "satrio", Role: "student"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email already registered"))
			})
		})

		When("Password Terlalu panjang (melebihi 72 char)", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, "satrio").Return(nil, errors.New("error")).Once()
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(nil, errors.New("email not registered")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'email already registered'", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", FirstName: "satrio", LastName: "w", Password: "12332222222222222222322222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222223232", Address: "bogor ct", Username: "satrio", Role: "student"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Register failed"))

			})
		})
		When("Query database Salah", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, "satrio").Return(nil, errors.New("error")).Once()
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(nil, errors.New("email not registered")).Once()
				Mock.On("Create", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Internal Server Error'", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", FirstName: "satrio", LastName: "w", Password: "123", Address: "bogor ct", Username: "satrio", Role: "student"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))

			})
		})
		When("Berhasil membuat user", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, "satrio").Return(nil, errors.New("error")).Once()
				Mock.On("FindByEmail", mock.Anything, "satrio2@gmail.com").Return(nil, errors.New("email not registered")).Once()
				Mock.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan error dengan nilai null", func() {
				err := UserService.Register(ctx, entity.RegisterReq{Email: "satrio2@gmail.com", FirstName: "satrio", LastName: "w", Password: "123", Address: "bogor ct", Username: "satrio", Role: "student"})
				Expect(err).Should(BeNil())

			})
		})
		When("Terdapat kesalahan query pada saat memverifikasi email user", func() {
			BeforeEach(func() {
				Mock.On("VerifyEmail", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Internal Server Error'", func() {
				err := UserService.VerifyEmail(ctx, "yewquei31231231======")
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil pada saat memverifikasi email user", func() {
			BeforeEach(func() {
				Mock.On("VerifyEmail", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan error dengan nilai nil", func() {
				err := UserService.VerifyEmail(ctx, "yewquei31231231======")
				Expect(err).Should(BeNil())
			})
		})

	})

	Context("Lupa Password", func() {
		When("Email tidak terdaftar", func() {
			BeforeEach(func() {
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("Email not registered")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Email not registered'", func() {
				err := UserService.ForgetPass(ctx, "satrio2@gmail.com")
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email not registered"))
			})
		})
		When("Email belum diverifikasi", func() {
			BeforeEach(func() {
				data := &entity.User{Email: "satrio2@gmail.com", IsVerified: false}
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(data, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Internal Server Error'", func() {
				err := UserService.ForgetPass(ctx, "satrio2@gmail.com")
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email not verified"))
			})
		})
		When("Terdapat kesalahan query pada saaat memasukan data", func() {
			BeforeEach(func() {
				data := &entity.User{Email: "satrio2@gmail.com", IsVerified: true}
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(data, nil).Once()
				Mock.On("InsertForgotPassToken", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Internal Server Error'", func() {
				err := UserService.ForgetPass(ctx, "satrio2@gmail.com")
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Terdapat kesalahan query pada saaat memasukan data", func() {
			BeforeEach(func() {
				data := &entity.User{Email: "satrio2@gmail.com", IsVerified: true}
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(data, nil).Once()
				Mock.On("InsertForgotPassToken", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Internal Server Error'", func() {
				err := UserService.ForgetPass(ctx, "satrio2@gmail.com")
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil melakukan permintaan lupa password", func() {
			BeforeEach(func() {
				data := &entity.User{Email: "satrio2@gmail.com", IsVerified: true}
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(data, nil).Once()
				Mock.On("InsertForgotPassToken", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan error dengan nilai nil", func() {
				err := UserService.ForgetPass(ctx, "satrio2@gmail.com")
				Expect(err).Should(BeNil())
			})
		})
	})

	When("Terdapat kesalahan query pada saat memasukan data password baru", func() {
		BeforeEach(func() {
			Mock.On("ResetPass", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
		})
		It("Akan Mengembalikan error dengan pesan 'Internal Server Error'", func() {
			err := UserService.ResetPass(ctx, "satrio123", "ewqe12312321=312======")
			Expect(err).ShouldNot(BeNil())
			Expect(err.Error()).To(Equal("Internal Server Error"))
		})
	})
	When("Berhasil Merubah password", func() {
		BeforeEach(func() {
			Mock.On("ResetPass", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		})
		It("Akan Mengembalikan error dengan nilai nil", func() {
			err := UserService.ResetPass(ctx, "satrio123", "ewqe12312321=312======")
			Expect(err).Should(BeNil())
		})
	})

	Context("User Update", func() {
		When("Password baru terlalu panjang (melebihi 72 char)", func() {
			It("Akan Mengembalikan error", func() {
				var file multipart.File
				_, err := UserService.Update(ctx, entity.UpdateReq{Password: "wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww"}, file)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Username baru sudah terdaftar pada database", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, mock.Anything).Return(&entity.User{Username: "satrio"}, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Username already registered'", func() {
				var file multipart.File
				_, err := UserService.Update(ctx, entity.UpdateReq{Username: "satrio"}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Username already registered"))
			})
		})

		When("Email baru sudah terdaftar pada database", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(&entity.User{Email: "satrio@gmail.com", IsVerified: true}, nil).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Email already registered'", func() {
				var file multipart.File
				_, err := UserService.Update(ctx, entity.UpdateReq{Username: "satrio", Password: "wwwwwwww", Email: "satrio@gmail.com"}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Email already registered"))
			})
		})
		When("Terdapat request gambar yang tidak sesuai dengan format gambar", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()
			})
			It("Akan Mengembalikan error dengan pesan 'Failed to upload image' ", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				_, err := UserService.Update(ctx, entity.UpdateReq{Username: "satrio", Password: "satrio2323223", Email: "satrio3@gmail.com", Image: "gambar.php"}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Failed to upload image"))
			})
		})

		When("Kesalahan Database", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error dengan pesan 'Internal Server Error'", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				_, err := UserService.Update(ctx, entity.UpdateReq{Username: "satrio", Password: "satrio2323223", Email: "satrio44@gmail.com", Image: "gambar.jpg"}, file)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil mengupdate profile", func() {
			BeforeEach(func() {
				Mock.On("FindByUsername", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()
				Mock.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(&entity.User{Email: "satrio44@gmail.com"}, nil).Once()
			})
			It("Akan Mengembalikan data user terbaru", func() {
				var file multipart.File
				file = os.NewFile(uintptr(2), "2")
				udata, err := UserService.Update(ctx, entity.UpdateReq{Username: "satrio", Password: "satrio2323223", Email: "satrio44@gmail.com", Image: "gambar.jpg"}, file)
				Expect(err).Should(BeNil())
				Expect(udata.Email).To(Equal("satrio44@gmail.com"))
			})
		})

	})
	Context("User Delete", func() {
		When("Terjadi kesalahan pada database", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				err := UserService.Delete(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil menghapus akun", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan nil error", func() {
				err := UserService.Delete(ctx, 1)
				Expect(err).Should(BeNil())
			})
		})

	})
	Context("User Profile", func() {
		When("Id user tidak ditemukan", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(nil, errors.New("id not found")).Once()
			})
			It("Akan Mengembalikan Error dengan pesan 'id not found'", func() {
				user, err := UserService.GetProfile(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("id not found"))
				Expect(user).To(BeNil())
			})
		})
		When("Server error", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan data user", func() {
				user, err := UserService.GetProfile(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
				Expect(user).Should(BeNil())
			})
		})
		When("Id user ditemukan", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(&entity.User{Email: "satrio@gmail.com"}, nil).Once()
			})
			It("Akan Mengembalikan data user", func() {
				user, err := UserService.GetProfile(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(user.Email).To(Equal("satrio@gmail.com"))
			})
		})

	})
})
