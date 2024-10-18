package command

import (
	"com.goldstar/goflyway/goflyway/consts"
	"com.goldstar/goflyway/goflyway/database"
	"com.goldstar/goflyway/goflyway/history"
	"com.goldstar/goflyway/goflyway/location"
	"com.goldstar/goflyway/goflyway/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

func init() {
	Registry(consts.CMD_NAME_MIGRATE, &Migrate{})
}

type Migrate struct {
}

func (m Migrate) Execute(ctx *Context) error {
	latestVersion, err := beforeExecute(ctx)
	if err != nil {
		return err
	}
	for _, l := range ctx.Options.Locations {
		for _, sql := range l.Sqls {
			startT := time.Now()
			err = m.invokeSql(ctx, sql, latestVersion)
			if err != nil {
				return errors.New(fmt.Sprintf("Failed to execute the SQL file:%s\nerror:%s", sql.Path, err.Error()))
			}
			tc := time.Since(startT)
			ctx.Logger.Info(ctx.Context, "invoke sql[%s] success,execution: %d ms", sql.Name, tc.Milliseconds())
		}
	}
	return nil
}

func (m Migrate) invokeSql(ctx *Context, sql location.SqlFile, latestVersion string) error {
	db := ctx.Database
	schemaHistory := ctx.SchemaHistory
	sd, err := schemaHistory.SelectVersion(sql.Version)
	if err != nil {
		return err
	}
	checksum, err := sql.CheckSum()
	var rank int64
	if sd != nil {
		rank = sd.InstalledRank
		if checksum != sd.Checksum {
			return errors.New(fmt.Sprintf("Flyway checksum mismatch error\n database: %d,local:%d", sd.Checksum, checksum))
		}
		if !sd.Success {
			content, err2 := m.readSqlContent(sql, ctx)
			if err2 != nil {
				return err2
			}
			d, err2 := m.invokeSqlContent(db, content)
			if err2 != nil {
				return err2
			} else {
				err = schemaHistory.UpdateSuccessAndTime(rank, true, d.Microseconds())
				if err != nil {
					return err
				}
			}
		}
	} else {
		// 传入的latestVersion不为空时，需要校验，新添加的sql版本是否高于latestVersion
		if latestVersion != "" {
			compare, err2 := utils.VersionCompare(sql.Version, latestVersion)
			if err2 != nil {
				return err2
			}
			if compare < 0 {
				return errors.New(fmt.Sprintf("The current version is %s. cannot execute %s", latestVersion, sql.Version))
			}
		}

		content, err2 := m.readSqlContent(sql, ctx)
		if err2 != nil {
			return err2
		}
		sd = &history.SchemaData{
			Version:       sql.Version,
			Description:   sql.Description,
			Type:          consts.SQL_TYPE,
			Script:        content,
			ExecutionTime: 0,
			Success:       false,
			Checksum:      checksum,
		}
		rank, err = schemaHistory.InsertData(*sd)
		if err != nil {
			return err
		}
		dur, err2 := m.invokeSqlContent(db, content)
		if err2 != nil {
			return err2
		} else {
			err = schemaHistory.UpdateSuccessAndTime(rank, true, dur.Microseconds())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m Migrate) invokeSqlContent(database database.Database, content string) (time.Duration, error) {
	start := time.Now()
	//err := database.Session().Exec(content)
	sqlArr := strings.Split(content, ";")
	var err error
	for _, sql := range sqlArr {
		if strings.TrimSpace(sql) != "" {
			fmt.Println("待执行目标sql：", sql)
			r, err := database.Session().ExecNew(sql)
			if err != nil {
				return time.Since(start), err
			}
			i, _ := r.LastInsertId()
			i1, _ := r.RowsAffected()
			fmt.Println("lastinsertId：", i, ",RowsAffected:", i1)
		}
	}
	since := time.Since(start)
	return since, err
}

func (m Migrate) readSqlContent(sql location.SqlFile, ctx *Context) (string, error) {
	content, err := sql.Content()
	if err != nil {
		return "", err
	}
	if ctx.Options.EnablePlaceholder {
		// 执行替换
		env, err := GenSqlPlaceholderEnv(ctx, sql)
		if err != nil {
			return "", err
		}
		t, err := utils.FormatTemplate(content, env)
		if err != nil {
			return "", err
		}
		content = t
	}
	return content, nil
}
