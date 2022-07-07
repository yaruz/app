generate:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative account/*.proto
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative session/*.proto
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative tg_account/*.proto
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative token/*.proto
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative user/*.proto
