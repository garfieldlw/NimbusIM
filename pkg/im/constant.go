package im

const EtcdImCometServerPath string = "/grpc/im/comet/"
const EtcdImCometUserPath string = "/im/comet/user/"

func GetEtcdImCometServerPath() string {
	return EtcdImCometServerPath
}

func GetEtcdImCometUserPath() string {
	return EtcdImCometUserPath
}
