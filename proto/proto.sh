find ./ | grep '\.proto' | xargs -n1 protoc  --go-grpc_out=require_unimplemented_servers=false:. --go_out=.

find ./ | grep '\.pb.go' | xargs -n1 sed -i '' -e 's/common \"common\/\"/common \"github.com\/garfieldlw\/NimbusIM\/proto\/common\"/g'

find ./ | grep '\.pb.go' | xargs -n1 sed -i '' -e 's/ws \"ws\/\"/ws \"github.com\/garfieldlw\/NimbusIM\/proto\/ws\"/g'
