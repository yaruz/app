generate:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative account/*.proto
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative session/*.proto
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative sn_account/*.proto
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative token/*.proto
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative user/*.proto
