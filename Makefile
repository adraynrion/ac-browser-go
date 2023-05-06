migrate.up:
	./migrate -path migrations -database "mysql://ac_browser:gDTm8gSqh2iGZ#*QLCx6X4RB@tcp(localhost:3306)/ac_browser?collation=utf8mb4_general_ci" -verbose up

migrate.down:
	./migrate -path migrations -database "mysql://ac_browser:gDTm8gSqh2iGZ#*QLCx6X4RB@tcp(localhost:3306)/ac_browser?collation=utf8mb4_general_ci" -verbose down

migrate.force:
	./migrate -path migrations -database "mysql://ac_browser:gDTm8gSqh2iGZ#*QLCx6X4RB@tcp(localhost:3306)/ac_browser?collation=utf8mb4_general_ci" force 1
