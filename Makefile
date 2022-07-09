mysql:
	./cloud_sql_proxy -dir=./cloudsql --instances=i-matrix-349904:asia-east2:mac-sql=tcp:3306 --credential_file=serviceAcc.json
build:
	go build
dev:
	go run main.go