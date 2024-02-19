package db

type Reference interface {
	GetID() int64
	GetName() string
}

var ReferencesMap = map[string]Reference{
	"College":        new(College),
	"Funnel":         new(Funnel),
	"LessonLocation": new(LessonLocation),
	"LessonSubject":  new(LessonSubject),
	"PaymentMethod":  new(PaymentMethod),
}

func (ref *College) GetID() int64 {
	return ref.CollegeID
}

func (ref *College) GetName() string {
	return ref.Name
}

func (ref *Funnel) GetID() int64 {
	return ref.FunnelID
}

func (ref *Funnel) GetName() string {
	return ref.Name
}

func (ref *LessonLocation) GetID() int64 {
	return ref.LocationID
}

func (ref *LessonLocation) GetName() string {
	return ref.Name
}

func (ref *LessonSubject) GetID() int64 {
	return ref.SubjectID
}

func (ref *LessonSubject) GetName() string {
	return ref.Name
}

func (ref *PaymentMethod) GetID() int64 {
	return ref.PaymentMethodID
}

func (ref *PaymentMethod) GetName() string {
	return ref.Name
}
