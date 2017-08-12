package modules


import (
	"github.com/HouzuoGuo/tiedot/db"

	"net/http"
	"github.com/julienschmidt/httprouter"
	"errors"
	"strconv"
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
	var requestBody interface{}

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get all documents", errors.New("No Collection was provided"), ""))
		return
	}

	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "query", nil); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}

	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("query", errors.New("Collection does not exist"), collectionName))
		return
	}

	// Evaluate the query
	queryResult := make(map[int]struct{})
	if err := db.EvalQuery(requestBody, dbcol, &queryResult); err != nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("query", err, collectionName))
		return
	}
	// Construct array of result
	resultDocs := make(map[string]interface{}, len(queryResult))
	counter := 0
	for docID := range queryResult {
		doc, _ := dbcol.Read(docID)
		if doc != nil {
			resultDocs[strconv.Itoa(docID)] = doc
			counter++
		}
	}

	RespondOk(w, resultDocs)
	return
}

func (module QueryAPIModule) Count(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestBody interface{}

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get all documents", errors.New("No Collection was provided"), ""))
		return
	}

	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "query", nil); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}

	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("query", errors.New("Collection does not exist"), collectionName))
		return
	}

	// Evaluate the query
	queryResult := make(map[int]struct{})
	if err := db.EvalQuery(requestBody, dbcol, &queryResult); err != nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("query", err, collectionName))
		return
	}

	RespondOk(w, map[string]interface{} {
		"collection": collectionName,
		"count": len(queryResult),
	})
	return
}