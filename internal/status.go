package internal

var (
	StatusSuccess    int32 = 1000 // success
	StatusErr        int32 = 1001 // error
	StatusParamErr   int32 = 1002 // param error
	StatusAuthErr    int32 = 1003 // auth error
	StatusRpcErr     int32 = 1004 // rpc error
	StatusGatewayErr int32 = 1005 // gateway error

	// User
	StatusUserNotExist int32 = 2001 // username does not exist
	StatusPasswordErr  int32 = 2002 //  password is not correct

	StatusUserNameExist     int32 = 2003
	StatusEmailExist        int32 = 2004 // the email has been in the db
	StatusValidationCodeErr int32 = 2005

	// Webset

)
