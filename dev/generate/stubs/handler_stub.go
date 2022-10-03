package stubs

var HandlerStruct = `
package {{PackageName}}

import (
	"git.tech.kora.id/go/rest"
)

type Handler struct{}

func (h *Handler) URLMapping(r *echo.Group) {
	{{ModulEndpoint}}
}
`

var HandlerGet = `
func (h *Handler) get(c *rest.Context) (e error) {
	var gr getRequest
	c.ResponseBody.Data, c.ResponseBody.Total, e = gr.get(c)

	return c.Serve(e)
}
`

var HandlerShow = `
func (h *Handler) show(c *rest.Context) (e error) {
	var sr showRequest

	sr.ID = c.Param("id")
	c.ResponseBody.Data, e = sr.get(id)
	return c.Serve(e)
}
`
var HandlerPost = `
func (h *Handler) create(c *rest.Context) (e error) {
	var cr createRequest
	if cr.Session, e = auth.RequestSession(c); e == nil {
		if e = c.Bind(&cr); e == nil {
			c.ResponseBody.Data, e = cr.Save()
		}
	}
	return c.Serve(e)
}
`

var HandlerPut = `
func (h *Handler) update(c *rest.Context) (e error) {
	var ur updateRequest
	ur.ID = c.Param("id")

	if ur.Session, e = auth.RequestSession(c); e == nil {
		if e = c.Bind(&ur); e == nil {
			c.ResponseBody.Data, e = ur.Save()
		}
	}
	return c.Serve(e)
}
`

var HandlerDelete = `
func (h *Handler) delete(c *rest.Context) (e error) {
	var dr deleteRequest
	dr.ID = c.Param("id")

	if dr.Session, e = auth.RequestSession(c); e == nil {
		if e = c.Bind(&dr); e == nil {
			e = dr.Save()
		}
	}
	return c.Serve(e)
}`
