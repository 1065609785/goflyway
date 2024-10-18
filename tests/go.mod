module com.goldstar/goflyway/goflyway/tests

go 1.20

require (
	com.goldstar/goflyway/goflyway v0.0.0-20231227095744-97601a534699
	com.goldstar/mattn/go-sqlite3 v1.14.19
)

require com.goldstar/go-sql-driver/mysql v1.7.1

replace com.goldstar/goflyway/goflyway => ../
