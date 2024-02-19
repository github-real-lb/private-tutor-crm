package db

// ReferenceStruct provides all functions to get information of Reference Structs:
// College, Funnel, LessonLocation, LessonSubject or PaymentMethod
type ReferenceStruct interface {
	GetID() int64
	GetName() string
	GetReferenceStructName() ReferenceStructName
	SetID(id int64)
	SetName(name string)
}

// ReferenceStructName is a string type used to define the names of Reference Structs:
// College, Funnel, LessonLocation, LessonSubject or PaymentMethod
type ReferenceStructName string

// Names of Reference Structs: College, Funnel, LessonLocation, LessonSubject or PaymentMethod
const (
	ReferenceCollege        ReferenceStructName = "College"
	ReferenceFunnel         ReferenceStructName = "Funnel"
	ReferenceLessonLocation ReferenceStructName = "LessonLocation"
	ReferenceLessonSubject  ReferenceStructName = "LessonSubject"
	ReferencePaymentMethod  ReferenceStructName = "PaymentMethod"
)

// ReferencesStructMap maps the Reference Structs's names and a pointer to the struct itself
var ReferenceStructMap = map[ReferenceStructName]ReferenceStruct{
	ReferenceCollege:        new(College),
	ReferenceFunnel:         new(Funnel),
	ReferenceLessonLocation: new(LessonLocation),
	ReferenceLessonSubject:  new(LessonSubject),
	ReferencePaymentMethod:  new(PaymentMethod),
}

func (ref *College) GetID() int64 {
	return ref.CollegeID
}

func (ref *College) GetName() string {
	return ref.Name
}

func (ref *College) GetReferenceStructName() ReferenceStructName {
	return ReferenceCollege
}

func (ref *College) SetID(id int64) {
	ref.CollegeID = id
}

func (ref *College) SetName(name string) {
	ref.Name = name
}

func (ref *Funnel) GetID() int64 {
	return ref.FunnelID
}

func (ref *Funnel) GetName() string {
	return ref.Name
}

func (ref *Funnel) GetReferenceStructName() ReferenceStructName {
	return ReferenceFunnel
}

func (ref *Funnel) SetID(id int64) {
	ref.FunnelID = id
}

func (ref *Funnel) SetName(name string) {
	ref.Name = name
}

func (ref *LessonLocation) GetID() int64 {
	return ref.LocationID
}

func (ref *LessonLocation) GetName() string {
	return ref.Name
}

func (ref *LessonLocation) GetReferenceStructName() ReferenceStructName {
	return ReferenceLessonLocation
}

func (ref *LessonLocation) SetID(id int64) {
	ref.LocationID = id
}

func (ref *LessonLocation) SetName(name string) {
	ref.Name = name
}

func (ref *LessonSubject) GetID() int64 {
	return ref.SubjectID
}

func (ref *LessonSubject) GetName() string {
	return ref.Name
}

func (ref *LessonSubject) GetReferenceStructName() ReferenceStructName {
	return ReferenceLessonSubject
}

func (ref *LessonSubject) SetID(id int64) {
	ref.SubjectID = id
}

func (ref *LessonSubject) SetName(name string) {
	ref.Name = name
}

func (ref *PaymentMethod) GetID() int64 {
	return ref.PaymentMethodID
}

func (ref *PaymentMethod) GetName() string {
	return ref.Name
}

func (ref *PaymentMethod) GetReferenceStructName() ReferenceStructName {
	return ReferencePaymentMethod
}

func (ref *PaymentMethod) SetID(id int64) {
	ref.PaymentMethodID = id
}

func (ref *PaymentMethod) SetName(name string) {
	ref.Name = name
}
