package features

import (
	schoolrepo "github.com/education-hub/BE/app/features/school/repository"
	schoolserv "github.com/education-hub/BE/app/features/school/service"
	userrepo "github.com/education-hub/BE/app/features/user/repository"
	userserv "github.com/education-hub/BE/app/features/user/service"
	"go.uber.org/dig"
)

func RegisterRepo(C *dig.Container) error {
	if err := C.Provide(userrepo.NewUserRepo); err != nil {
		return err
	}
	if err := C.Provide(schoolrepo.NewSchoolRepo); err != nil {
		return err
	}
	return nil
}

func RegisterService(C *dig.Container) error {
	if err := C.Provide(userserv.NewUserService); err != nil {
		return err
	}
	if err := C.Provide(schoolserv.NewSchoolService); err != nil {
		return err
	}

	return nil
}
