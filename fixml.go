package gonol3

import "fmt"

type rejectReason string

const (
	rejectReasonOther              rejectReason = "0"
	rejectReasonUnknownId          rejectReason = "1"
	rejectReasonUnknownInstrument  rejectReason = "2"
	rejectReasonUnknownMessageType rejectReason = "3"
	rejectReasonNoAccess           rejectReason = "4"
	rejectReasonXmlError           rejectReason = "5"
	rejectReasonUnauthorized       rejectReason = "6"
	rejectReasonNotConnected       rejectReason = "7"
)

func (reason rejectReason) desc() string {
	switch reason {
	case rejectReasonOther:
		return fmt.Sprintf("Other error")
	case rejectReasonUnknownId:
		return "Unknown ID"
	case rejectReasonUnknownInstrument:
		return "Unknown instrument"
	case rejectReasonUnknownMessageType:
		return "Unknown message type"
	case rejectReasonNoAccess:
		return "No access"
	case rejectReasonXmlError:
		return "XML syntaxt error"
	case rejectReasonUnauthorized:
		return "Unauthorized"
	case rejectReasonNotConnected:
		return "Not connected"
	}

	return "Unknown error"
}

// Message reject structure
type fixmlBusinessMessageReject struct {
	XMLName struct{} `xml:"BizMsgRej"`

	RejectReason rejectReason `xml:"BizRejRsn,attr"`
	Text         string       `xml:"Txt,attr"`
}

// Main FIXML message structure
const (
	fixmlDefaultVersion = "5.0"
	fixmlRevisionDate   = "20080317"
	fixmlSchemaDate     = "20080314"
)

type fixmlRequest[T any] struct {
	XMLName struct{} `xml:"FIXML"`

	Version      string `xml:"v,attr"`
	RevisionDate string `xml:"r,attr"`
	SchemaDate   string `xml:"s,attr"`

	Message       *T                          `xml:",omitempty"`
	RejectMessage *fixmlBusinessMessageReject `xml:",omitemtpy"`
}

func wrapFixmlRequest[T any](v T) fixmlRequest[T] {
	return fixmlRequest[T]{
		Version:      fixmlDefaultVersion,
		RevisionDate: fixmlRevisionDate,
		SchemaDate:   fixmlSchemaDate,

		Message: &v,
	}
}

type fixmlResponse struct {
	XMLName struct{} `xml:"FIXML"`

	Version      string `xml:"v,attr"`
	RevisionDate string `xml:"r,attr"`
	SchemaDate   string `xml:"s,attr"`

	RejectMessage *fixmlBusinessMessageReject `xml:",omitemtpy"`
	UserResponse  *fixmlUserResponse
}

// User request / response
type userReqType int

const (
	userReqTypeLogin  userReqType = 1
	userReqTypeLogout userReqType = 2
	userReqTypeStatus userReqType = 4
)

type fixmlUserRequest struct {
	XMLName struct{} `xml:"UserReq"`

	RequestId int         `xml:"UserReqID,attr"`
	Type      userReqType `xml:"UserReqTyp,attr"`
	Username  string      `xml:"Username,attr"`
	Password  string      `xml:"Password,attr"`
}

type userStatus int

const (
	userStatusLoggedIn    userStatus = 1
	userStatusLoggedOut   userStatus = 2
	userStatusNotExists   userStatus = 3
	userStatusBadPassword userStatus = 4
	userStatusOffline     userStatus = 5
	userStatusOther       userStatus = 6
	userStatusNolOfflice  userStatus = 7
)

func (status userStatus) desc() string {
	switch status {
	case userStatusLoggedIn:
		return "User logged in"
	case userStatusLoggedOut:
		return "User loggged out"
	case userStatusNotExists:
		return "User does not exist"
	case userStatusBadPassword:
		return "Bad password"
	case userStatusOffline:
		return "User offline"
	case userStatusOther:
		return "Other"
	case userStatusNolOfflice:
		return "NOL Offline"
	}

	return "Unknown user status"
}

type fixmlUserResponse struct {
	XMLName struct{} `xml:"UserRsp"`

	RequestId   int        `xml:"UserReqID,attr"`
	Username    string     `xml:"Username,attr"`
	MarketDepth int        `xml:"MktDepth,attr"`
	UserStatus  userStatus `xml:"UserStat,attr"`
}

//
