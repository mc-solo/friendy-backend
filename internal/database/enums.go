package database

// only male and female genders
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

// invalid gender value check
func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale:
		return true
	}
	return false
}

// String returns the string representation
func (g Gender) String() string {
	return string(g)
}

// edu bg
type EducatoinalLevel string

const (
	// formal
	EduHighSchool EducatoinalLevel = "high_school"
	EduBachelor   EducatoinalLevel = "bachelor"
	EduMaster     EducatoinalLevel = "master"
	EduPhd        EducatoinalLevel = "phd"

	// freaks
	UniDropout        EducatoinalLevel = "university_drop_out"
	HighSchoolDropout EducatoinalLevel = "high_school_dropout"
	HomeSchooled      EducatoinalLevel = "home_schooled"
	EduOther          EducatoinalLevel = "other"
)

func (e EducatoinalLevel) IsValid() bool {
	switch e {
	case EduHighSchool, EduBachelor, EduMaster, EduPhd, UniDropout, HighSchoolDropout, EduOther:
		return true
	}
	return false
}

type BodyType string

const (
	BodyTypeSlim     BodyType = "slim"
	BodyTypeAthletic BodyType = "athletic"
	BodyTypeChubby   BodyType = "chubby"
	BodyTypeCurvy    BodyType = "curvy"
	BodyTypeFit      BodyType = "fit"
	BodyTypeOther    BodyType = "other"
)

func (b BodyType) IsValid() bool {
	switch b {
	case BodyTypeSlim, BodyTypeAthletic, BodyTypeChubby, BodyTypeCurvy, BodyTypeFit, BodyTypeOther:
		return true
	}
	return false
}
