package gonol3

import (
	"encoding/xml"
	"testing"
)

func TestFixmlMarshalling(t *testing.T) {
	userReq := fixmlUserRequest{
		RequestId: 1,
		Type:      userReqTypeLogin,
		Username:  "BOS",
		Password:  "BOS",
	}

	msg := wrapFixmlMessage(userReq)

	expectedContent := `<FIXML v="5.0" r="20080317" s="20080314"><UserReq UserReqID="1" UserReqTyp="1" Username="BOS" Password="BOS"></UserReq></FIXML>`
	content, err := xml.Marshal(msg)
	if err != nil {
		t.Fatalf(`Marshalling error: %v`, err)
	}

	contentStr := string(content)
	if contentStr != expectedContent {
		t.Fatalf("Marshalled into bad content (got: %s)", contentStr)
	}
}

func TestFixmlUnmarshalling(t *testing.T) {
	content := `<FIXML v="5.0" r="20080317" s="20080314"><UserReq UserReqID="1" UserReqTyp="1" Username="BOS" Password="BOS"></UserReq></FIXML>`
	dst := fixmlMessage[fixmlUserRequest]{}

	err := xml.Unmarshal([]byte(content), &dst)
	if err != nil {
		t.Fatalf("Error while unmarshalling: %v", err)
	}

	if dst.Message.Type != userReqTypeLogin {
		t.Fatalf("Unexpected request type: %d", dst.Message.Type)
	}
}

func TestFixmlUnmarshallingToRejectMessage(t *testing.T) {
	content := `<FIXML v="5.0" r="20080317" s="20080314"><BizMsgRej RefMsgTyp="BE" BizRejRsn="5"/></FIXML>`
	dst := fixmlMessage[fixmlUserRequest]{}

	err := xml.Unmarshal([]byte(content), &dst)
	if err != nil {
		t.Fatalf("Error while unmarshalling: %v", err)
	}

	if dst.RejectMessage == nil {
		t.Fatalf("MessageReject field must be set")
	}

	if dst.RejectMessage.RejectReason != rejectReasonXmlError {
		t.Fatalf("Reject reason is expected to be XML error but it's different (got: %s)", dst.RejectMessage.RejectReason)
	}
}
