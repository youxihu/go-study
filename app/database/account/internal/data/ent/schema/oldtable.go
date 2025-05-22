package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() []ent.Field {

	return []ent.Field{

		field.Int32("id").SchemaType(map[string]string{
			dialect.MySQL: "int  ", // Override MySQL.
		}).Comment("用户id").Unique(),

		field.Int32("tenant_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Default(1),

		field.Time("updated_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("更新信息时间"),

		field.Time("created_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("注册时间"),

		field.Int32("account_integral").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Optional().Default(0).Comment("当前积分"),

		field.Float("account_balance").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2)", // Override MySQL.
		}).Optional().Comment("当前余额"),

		field.Int8("user_identity").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Comment("用户身份 1普通用户 2镖师 3业务员  4客服人员"),

		field.Float("freeze_balance").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2) ", // Override MySQL.
		}).Optional().Comment("冻结资金"),

		field.Int32("freeze_integral").SchemaType(map[string]string{
			dialect.MySQL: "int ", // Override MySQL.
		}).Default(0).Comment("冻结积分"),
	}

}
func (Account) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "account"},
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	return nil
}
