package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Project holds the schema definition for the Project entity.
type Project struct {
	ent.Schema
}

// Fields of the Project.
func (Project) Fields() []ent.Field {

	return []ent.Field{

		field.Int32("id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("主键ID").Unique(),

		field.Int32("tenant_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Default(1).Comment("租户id"),

		field.Int32("account_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("发布账户ID"),

		field.String("account_phone").SchemaType(map[string]string{
			dialect.MySQL: "varchar(15)", // Override MySQL.
		}).Optional().Comment("发布账户手机号"),

		field.Int32("worker_account_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Optional().Comment("接单镖师ID"),

		field.Int8("project_type").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Comment("发布项目类型"),

		field.String("title").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Comment("项目名称"),

		field.Time("expect_deliver_time").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Comment("预期交付时间"),

		field.String("tender_files").SchemaType(map[string]string{
			dialect.MySQL: "json", // Override MySQL.
		}).Optional().Comment("招标资料"),

		field.String("deliver_bid_files").SchemaType(map[string]string{
			dialect.MySQL: "json", // Override MySQL.
		}).Optional().Comment("最终交付的标书"),

		field.Time("deliver_bid_time").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("标书最终交付时间"),

		field.String("remark").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Optional().Comment("备注(后台)"),

		field.Int32("project_status").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("项目状态"),

		field.Float("amount").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2)unsigned", // Override MySQL.
		}).Default(0.00).Comment("订单金额(选择镖师后报价反推金额)"),

		field.Float("paid_amount").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2)unsigned", // Override MySQL.
		}).Default(0.00).Comment("已支付金额"),

		field.Float("unpaid_amount").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2)unsigned", // Override MySQL.
		}).Default(0.00).Comment("未支付金额"),

		field.Float("commission").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2)", // Override MySQL.
		}).Optional().Comment("佣金"),

		field.Float("worker_commission").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2)", // Override MySQL.
		}).Optional().Comment("镖师实际佣金"),

		field.Int8("finance_invoice_status").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Default(1).Comment("财务开票状态 1未开票 2已开票 3开票失败 4开票中"),

		field.String("e_sign_id").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Optional().Comment("e签宝 签署流程id"),

		field.Int8("e_sign_status").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Optional().Default(1).Comment("e签宝 状态"),

		field.String("e_sign_notify").SchemaType(map[string]string{
			dialect.MySQL: "json", // Override MySQL.
		}).Optional().Comment("e签宝 流程结束回调"),

		field.Float("score").SchemaType(map[string]string{
			dialect.MySQL: "decimal(2,1)unsigned", // Override MySQL.
		}).Default(5.0).Comment("评分"),

		field.String("evaluate").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Optional().Comment("评价"),

		field.String("cancel_reason").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Optional().Comment("取消原因"),

		field.Time("close_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("超时关闭时间"),

		field.String("channel_code").SchemaType(map[string]string{
			dialect.MySQL: "varchar(255)", // Override MySQL.
		}).Optional().Comment("渠道code"),

		field.Time("created_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("创建时间"),
		field.Float("pay_money").SchemaType(map[string]string{
			dialect.MySQL: "decimal(10,2)", // Override MySQL.
		}).Default(0.00).Comment("订单实付金额"),
		field.Time("updated_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Optional().Comment("更新时间"),
	}

}

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return nil
}
func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "project"},
	}
}
