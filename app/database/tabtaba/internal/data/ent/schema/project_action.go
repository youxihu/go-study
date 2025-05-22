package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ProjectActionRecord holds the schema definition for the ProjectActionRecord entity.
type ProjectActionRecord struct {
	ent.Schema
}

// Fields of the ProjectActionRecord.
func (ProjectActionRecord) Fields() []ent.Field {

	return []ent.Field{

		field.Int32("id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("动作ID").Unique(),

		field.Int32("tenant_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Default(1).Comment("租户id"),

		field.Int32("project_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("项目ID"),

		field.Int8("action_plat").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Default(1).Comment("动作平台 1 前台 2后台"),

		field.String("action").SchemaType(map[string]string{
			dialect.MySQL: "varchar(200)", // Override MySQL.
		}).Default("").Comment("动作内容"),

		field.Int32("account_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Optional().Default(0).Comment("操作人id"),

		field.String("account_name").SchemaType(map[string]string{
			dialect.MySQL: "varchar(50)", // Override MySQL.
		}).Optional().Default("").Comment("操作人"),

		field.Int32("project_status").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("订单状态"),

		field.String("project_status_text").SchemaType(map[string]string{
			dialect.MySQL: "varchar(200)", // Override MySQL.
		}).Default("").Comment("订单状态名称"),

		field.Int8("status").SchemaType(map[string]string{
			dialect.MySQL: "tinyint", // Override MySQL.
		}).Optional().Default(1).Comment("状态 1启用 2禁用 3删除"),

		field.Time("created_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Comment("创建时间"),

		field.Time("updated_at").SchemaType(map[string]string{
			dialect.MySQL: "timestamp", // Override MySQL.
		}).Comment("更新时间"),
	}

}

// Edges of the ProjectActionRecord.
func (ProjectActionRecord) Edges() []ent.Edge {
	return nil
}
func (ProjectActionRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "project_action_record"},
	}
}
