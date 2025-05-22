package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// AccountWallet holds the schema definition for the AccountWallet entity.
type AccountWallet struct {
	ent.Schema
}

// Fields of the AccountWallet.
func (AccountWallet) Fields() []ent.Field {

	return []ent.Field{

		field.Int32("id").SchemaType(map[string]string{
			dialect.MySQL: "int unsigned", // 修复为合法的类型
		}).Comment("用户id").Unique(),

		field.Int32("tenant_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}),

		field.Int32("account_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Optional().Comment("用户ID"),

		field.Int32("account_integral").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Optional().Default(0).Comment("当前积分"),

		field.Float("account_balance").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2)", // Override MySQL.
		}).Optional().Comment("当前余额"),

		field.Float("freeze_balance").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2)unsigned", // Override MySQL.
		}).Comment("冻结资金"),

		field.Int32("freeze_integral").SchemaType(map[string]string{
			dialect.MySQL: "int unsigned", // Override MySQL.
		}).Optional().Default(0).Comment("冻结积分"),

		field.Int8("open_installment").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Default(2).Comment("分期付款"),

		field.Int8("status").SchemaType(map[string]string{
			dialect.MySQL: "tinyint(1)", // Override MySQL.
		}).Optional(),

		field.Int8("state").SchemaType(map[string]string{
			dialect.MySQL: "tinyint(1)", // Override MySQL.
		}).Optional(),

		field.Time("created_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("注册时间"),

		field.Time("updated_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("更新信息时间"),
	}

}
func (AccountWallet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "account_wallet"},
	}
}

// Edges of the AccountWallet.
func (AccountWallet) Edges() []ent.Edge {
	return nil
}
