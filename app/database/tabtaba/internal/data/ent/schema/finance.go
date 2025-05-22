package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Finance holds the schema definition for the Finance entity.
type Finance struct {
	ent.Schema
}

// Fields of the Finance.
func (Finance) Fields() []ent.Field {

	return []ent.Field{

		field.Int32("id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("财务id").Unique(),

		field.Int32("tenant_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Default(1).Comment("租户id"),

		field.String("transact_no").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("交易号"),

		field.String("order_no").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("订单号"),

		field.Int32("project_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Optional().Comment("项目id"),

		field.Int32("account_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("账户id"),

		field.Int32("their_account_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("对方账户id"),

		field.Time("transact_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("交易时间"),

		field.Time("write_off_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("核销时间"),

		field.String("summary").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Optional().Comment("摘要说明"),

		field.Int8("pay_mode_type").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Optional().Comment("支付方式"),

		field.Text("pay_info").SchemaType(map[string]string{
			dialect.MySQL: "text", // Override MySQL.
		}).Optional().Comment("支付信息"),

		field.String("remark").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Optional().Comment("备注"),

		field.Int8("pay_status").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Optional().Comment("支付状态"),

		field.Float("amount").SchemaType(map[string]string{
			dialect.MySQL: "decimal(32,2)", // Override MySQL.
		}).Optional().Comment("金额"),

		field.Int8("change_type").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Optional().Comment("变动类型"),

		field.Time("created_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("创建时间"),

		field.Time("updated_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("更新时间"),
	}

}

// Edges of the Finance.
func (Finance) Edges() []ent.Edge {
	return nil
}
func (Finance) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "finance"},
	}
}
