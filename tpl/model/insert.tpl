func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	if session == nil{
		session=m.conn
	}
	{{if .withCache}}{{.keys}}
    ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
    ret,err:=session.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return ret,err
}
