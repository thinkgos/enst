# ormat

database/sql to golang struct, and generate enum annotation from comment to proto file.

## Features

## Usage

### Installation

```bash
    go install github.com/thinkgos/enst/cmd/ormat@latest
```

Example.

NOTE:

- database filed comment `[@jsontag: realjsontag]` will overwrite the filed json tags.
- database filed comment `[@affix]` will append `,string` to the filed json tags.

```go
// SysUser 用户表
type SysUser struct {
	Id        int64     `gorm:"column:id;type:bigint;autoIncrement;not null;primaryKey,priority:1" json:"id,omitempty"`
	Username  string    `gorm:"column:username;type:varchar(64);not null;primaryKey,priority:2;uniqueIndex:uk_username" json:"username,omitempty"`
	Password  string    `gorm:"column:password;type:varchar(255);not null" json:"password,omitempty"`
	Nickname  string    `gorm:"column:nickname;type:varchar(64);not null" json:"nickname,omitempty"`
	Phone     string    `gorm:"column:phone;type:varchar(16);not null" json:"phone,omitempty"`
	Avatar    string    `gorm:"column:avatar;type:varchar(255);not null" json:"avatar,omitempty"`
	Sex       int8      `gorm:"column:sex;type:tinyint;not null;default:3" json:"sex,omitempty"`
	Email     string    `gorm:"column:email;type:varchar(32);not null" json:"email,omitempty"`
	Status    string    `gorm:"column:status;type:varchar(1);not null;default:1" json:"status,omitempty"`
	Remark    string    `gorm:"column:remark;type:varchar(255);not null" json:"remark,omitempty"`
	Creator   string    `gorm:"column:creator;type:varchar(32);not null" json:"creator,omitempty"`
	Updator   string    `gorm:"column:updator;type:varchar(32);not null" json:"updator,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null" json:"updated_at,omitempty"`
}

// TableName implement schema.Tabler interface
func (*SysUser) TableName() string {
	return "sys_user"
}
```

### Help

```shell
$ ./ormat --help

database/sql to golang struct

Usage:
  ormat [command]

Available Commands:
  build       Generate model from sql
  completion  Generate the autocompletion script for the specified shell
  gen         Generate model from database
  help        Help about any command
  sql         Generate sql file
  upgrade     Upgrade ormat

Flags:
  -h, --help           help for ormat
  -l, --level string   log level(debug,info,warn,error) (default "info")
  -v, --version        version for ormat

Use "ormat [command] --help" for more information about a command.
```

### Build

```bash
goreleaser release --snapshot  --rm-dist
```

## References

## License

This project is under MIT License. See the [LICENSE](../../LICENSE) file for the full license text.
