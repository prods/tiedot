package modules


import (
	"github.com/HouzuoGuo/tiedot/db"

	"net/http"
	"github.com/julienschmidt/httprouter"
)

type QueryAPIModule struct {
	*APIModuleBase
}

func NewQueryAPIModule(db *db.DB) *QueryAPIModule {
	newInstance := &QueryAPIModule{&APIModuleBase{}}
	newInstance.db = db
	newInstance.routes = []APIRoute {
		POST("/collection/:collection_name/query", newInstance.Query, true),
		POST("/collection/:collection_name/query/count", newInstance.Count, true),
	}
	return newInstance
}

func (module QueryAPIModule) GetRoutes() []APIRoute {
	return module.routes
}

func (module QueryAPIModule) GetName() string {
	return "Query"
}

func (module QueryAPIModule) GetDocumentation() APIModuleDocumentation {
	return module.documentation
}

func (module QueryAPIModule) Query(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Query Implementation
}

func (module QueryAPIModule) Count(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Query Count Implementation
}