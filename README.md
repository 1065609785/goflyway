# GO flyway

![GitHub License](https://img.shields.io/github/license/goflyway/goflyway)
[![Static Badge](https://img.shields.io/badge/go.dev-reference-blue?style=flat)](https://pkg.go.dev/com.goldstar/goflyway/goflyway)

## 安装

```shell
go get -u com.goldstar/goflyway/goflyway
```

## 快速开始

```go
package main

import (
	"database/sql"
	_ "com.goldstar/mattn/go-sqlite3"
	"com.goldstar/goflyway/goflyway"
	"com.goldstar/goflyway/goflyway/database"
)

func main() {
	db, _ := sql.Open("sqlite3", "./flyway_test.db")
	// use database.T_SQLITE 、 database.T_MYSQL or "sqlite","mysql"
	f, _ := flyway.Open(database.T_SQLITE, db, &flyway.Config{Locations: []string{"db_migration"}})
	f.Migrate()
}
```

## 支持的数据库

- sqlite
- mysql

## 使用配置

示例:

```go
&flyway.Config{...}

```

配置项：

 名称                    | 类型       | 默认值              | 说明                                                                                          
-----------------------|----------|------------------|---------------------------------------------------------------------------------------------
 Locations             | []string | ["db_migration"] | 要递归扫描迁移的位置                                                                                  
 BaselineOnMigrate     | bool     | false            | 是否在对没有模式历史表的非空模式执行迁移时自动调用基线。在执行迁移之前，将使用baselineVersion对该模式进行基线化。只有baselinversion之上的迁移才会被应用。 
 BaselineVersion       | string   | 1                | 基线版本号，用于创建基线版本                                                                              
 Schemas               | []string | []               | 数据库连接的模式列表                                                                                  
 CreateSchemas         | bool     | false            | 是否创建 Schemas 指定的模式                                                                          
 DefaultSchema         | string   |                  | 默认的模式，为空时，默认为数据库连接的默认模式，如果指定了 Schemas 则取第一个为默认模式                                            
 CleanDisabled         | bool     | false            | 为ture时，会清空 Schemas 下所有表。注意：生产模式不要设置为true                                                    
 OutOfOrder            | bool     | false            | 是否允许版本乱序运行，为ture时，如果已经应用了1.0和3.0版本，现在发现了2.0版本，那么它也将被应用，而不是被忽略。                              
 EnablePlaceholder     | bool     | false            | 是否开启占位符替换，开启后将对sql文件进行处理                                                                    
 DisableCallbacks      | bool     | false            | 是否禁用callbacks方法执行，开启后，callbacks将不生效                                                         
 SqlMigrationSeparator | string   | __               | 脚本文件名中版本号和描述之间的分隔符                                                                          
 SqlMigrationPrefix    | string   | V                | 脚本文件名的前缀，用于标识脚本的版本号                                                                         



demo:
```go
func main() {
	db, _ := sql.Open(database.T_MYSQL, "root:xxxxx@tcp(192.168.xxx.xx:3306)/goflyway?charset=utf8")
	// use database.T_SQLITE 、 database.T_MYSQL or "sqlite","mysql"
	f, err := flyway.Open(database.T_MYSQL, db, &flyway.Config{Locations: []string{"db_migration"}, OutOfOrder: true, BaselineVersion: "1", BaselineOnMigrate: true, SqlMigrationPrefix: "V", SqlMigrationSeparator: "__"})
	if err != nil {
		panic(err)
	}
	err = f.Migrate()
	if err != nil {
		panic(err)
	}	
	fmt.Println("success")
}
```
