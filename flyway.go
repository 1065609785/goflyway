package flyway

import (
	"com.goldstar/goflyway/goflyway/command"
	"com.goldstar/goflyway/goflyway/consts"
	"com.goldstar/goflyway/goflyway/database"
	_ "com.goldstar/goflyway/goflyway/init"
	"com.goldstar/goflyway/goflyway/logger"
	"database/sql"
)

type flyway struct {
	databaseType database.Type
	config       Config
	db           *sql.DB
}

func (f *flyway) Migrate() error {
	ctx, err := buildCommandCtx(consts.CMD_NAME_MIGRATE, f)
	if err != nil {
		return err
	}
	return command.Execute(ctx)
}

func (f *flyway) Validate() error {
	ctx, err := buildCommandCtx(consts.CMD_NAME_VALIDATE, f)
	if err != nil {
		return err
	}
	return command.Execute(ctx)
}

func (f *flyway) Callbacks() *command.CallbackDispatch {
	return command.Callbacks()
}

func Open(databaseType string, db *sql.DB, config *Config) (*flyway, error) {
	dbType, err := database.TypeValueOf(databaseType)
	if err != nil {
		return nil, err
	}
	err = configBuild(config)
	if err != nil {
		return nil, err
	}
	f := &flyway{
		databaseType: dbType,
		config:       *config,
		db:           db,
	}
	return f, nil
}

type Config struct {
	Logger                logger.Interface
	Locations             []string
	Table                 string
	BaselineOnMigrate     bool     // 是否使用基线迁移
	BaselineVersion       string   // 基线版本号，用于创建基线版本
	Schemas               []string // 连接的模式列表
	CreateSchemas         bool     // 是否创建 Schemas 指定的模式
	DefaultSchema         string   // 默认的模式，为空时，默认为数据库连接的默认模式，如果指定了 Schemas 则取第一个为默认模式
	CleanDisabled         bool     // 为ture时，会清空 Schemas 下所有表
	OutOfOrder            bool     // 是否允许版本乱序运行，为ture时，如果已经应用了1.0和3.0版本，现在发现了2.0版本，那么它也将被应用，而不是被忽略。
	EnablePlaceholder     bool     // 是否开启占位符替换
	DisableCallbacks      bool     // 是否禁用callback
	SqlMigrationSeparator string   // 脚本文件名中版本号和描述之间的分隔符
	SqlMigrationPrefix    string   // 脚本文件名的前缀，用于标识脚本的版本号
}
