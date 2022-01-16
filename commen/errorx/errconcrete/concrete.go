package errconcrete

/*
这个里面放的是错误的常量就定义在这个里面
*/

const (
	SqlNotFound   = "此用户还没有注册"
	InterErr      = "内部错误请等待"
	PasswordErr   = "密码输入错误"
	UserHasExit   = "账号重复"
	UserNotHasMsg = "此用户没有注册个人信息"
	UserUidValid  = "用户的uid无效"
)

const (
	RpcUidInvalid  = "传入的uid小于0"
	RpcUidNotFound = "没有这个人的信息"
	RpcInterErr    = "内部错误"
)

const (
	BookNotFount        = "没有这本书的信息"
	BookNameValid       = "书名无效传输"
	BookNotFound        = "图书馆内没有这本书"
	BookNotHasInventory = "这本书已经被借完了"
	BookHasRendByYou    = "这本书已经被你给借过了"
	BookNotRendByYou    = "你还没有借过这本书"
	BookNotAnyRend      = "你没有借过任何一本书"
	BookNotRend         = "这本书没有人来借"
)
