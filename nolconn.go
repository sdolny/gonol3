package gonol3

import (
	"encoding/xml"
	"fmt"
	"net"
)

const (
	initialBufferSize int = 65535
)

type NolConn struct {
	host      string
	syncPort  int
	asyncPort int

	syncConn       net.Conn
	syncConnBuffer []byte
}

func NolConnect(host string, syncPort int, asyncPort int) (*NolConn, error) {
	syncAddress := fmt.Sprintf("%s:%d", host, syncPort)

	syncConn, err := net.Dial("tcp", syncAddress)
	if err != nil {
		return nil, err
	}

	conn := NolConn{
		host:      host,
		syncPort:  syncPort,
		asyncPort: asyncPort,

		syncConn:       syncConn,
		syncConnBuffer: make([]byte, initialBufferSize),
	}

	return &conn, nil
}

func (this *NolConn) Close() {
	this.syncConn.Close()
}

func (this *NolConn) Login(username, password string) error {
	request := wrapFixmlRequest(fixmlUserRequest{
		RequestId: 1,
		Username:  username,
		Password:  password,
		Type:      userReqTypeLogin,
	})

	response := fixmlResponse{}
	err := exchangeFixmlMessages(this, request, &response)
	if err != nil {
		return err
	}

	if response.RejectMessage != nil {
		rejectReasonDesc := response.RejectMessage.RejectReason.desc()
		return fmt.Errorf("fixml request rejected with an error: %s", rejectReasonDesc)
	}

	if response.UserResponse.UserStatus != userStatusLoggedIn {
		userStatusDesc := response.UserResponse.UserStatus.desc()
		return fmt.Errorf("User cannot be logged in. Error code: %s", userStatusDesc)
	}

	return nil
}

func (this *NolConn) Logout(username, password string) error {
	request := wrapFixmlRequest(fixmlUserRequest{
		RequestId: 1,
		Username:  username,
		Password:  password,
		Type:      userReqTypeLogout,
	})

	response := fixmlResponse{}
	err := exchangeFixmlMessages(this, request, &response)
	if err != nil {
		return err
	}

	if response.RejectMessage != nil {
		rejectReasonDesc := response.RejectMessage.RejectReason.desc()
		return fmt.Errorf("FIXML request rejected with an error: %s", rejectReasonDesc)
	}

	if response.UserResponse.UserStatus != userStatusLoggedOut {
		userStatusDesc := response.UserResponse.UserStatus.desc()
		return fmt.Errorf("User cannot be logged in. Error code: %s", userStatusDesc)
	}

	return nil
}

func exchangeFixmlMessages[T any, R any](conn *NolConn, request T, response *R) error {
	// Send request
	xmlRequest, err := xml.Marshal(request)
	if err != nil {
		return err
	}

	err = connWriteFixmlMessage(conn.syncConn, xmlRequest)
	if err != nil {
		return err
	}

	// And now wait for response
	xmlResponse, err := connReadFixmlMessage(conn.syncConn)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(xmlResponse, response)
	if err != nil {
		return err
	}

	return nil
}

func connReadFixmlMessage(conn net.Conn) ([]byte, error) {
	expectedMessageLength, err := connReadInt(conn)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, expectedMessageLength)
	err = connReadAll(conn, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func connReadInt(conn net.Conn) (int, error) {
	buf := []byte{0, 0, 0, 0}
	err := connReadAll(conn, buf)
	if err != nil {
		return 0, err
	}

	return int(buf[0]) | int(buf[1])<<8 | int(buf[2])<<16 | int(buf[3])<<24, nil
}

func connReadAll(conn net.Conn, dst []byte) error {
	totalRead := 0
	for totalRead < len(dst) {
		n, err := conn.Read(dst[totalRead:])
		if err != nil {
			return err
		}

		totalRead += n
	}

	return nil
}

func connWriteFixmlMessage(conn net.Conn, data []byte) error {
	// Send the total data length
	err := connWriteInt(conn, len(data))
	if err != nil {
		return err
	}

	// Send the data
	err = connWriteAll(conn, data)
	if err != nil {
		return err
	}

	return nil
}

func connWriteInt(conn net.Conn, val int) error {
	buf := []byte{
		byte(val),
		byte(val >> 8),
		byte(val >> 16),
		byte(val >> 24),
	}

	return connWriteAll(conn, buf)
}

// Single call can be not enough - it's a helper method to write
// data in multiple calls
func connWriteAll(conn net.Conn, data []byte) error {
	totalWritten := 0
	for totalWritten < len(data) {
		n, err := conn.Write(data[totalWritten:])
		if err != nil {
			return err
		}

		totalWritten += n
	}

	return nil
}
