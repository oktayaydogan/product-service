package api

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"product-service/database"
)

// ProductType, GraphQL şemasındaki ürün nesnesini temsil eder.
type ProductType struct {
	ID     int     `json:"id" example:"1"`
	Name   string  `json:"name" example:"Product 1"`
	Price  float32 `json:"price" example:"19.99"`
	GameID int     `json:"gameId" example:"1"`
}

// RootQuery, GraphQL şemasındaki kök sorgu nesnesini temsil eder.
type RootQuery struct {
	Products []ProductType `json:"products"`
	Product  ProductType   `json:"product"`
}

// IntrospectionType, GraphQL introspeksiyon özelliği için kullanılan nesneyi temsil eder.
type IntrospectionType struct {
	Typename string `json:"__typename"`
}

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"gameId": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		// "products" endpoint'ı, tüm ürünleri alır.
		"products": &graphql.Field{
			Type: graphql.NewList(productType),
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				gameID, ok := params.Args["gameId"].(int)
				if ok {
					// Eğer gameId belirtilmişse, ürünleri oyun ID'sine göre filtrele.
					return database.GetProductsByGameID(gameID)
				}
				// gameId belirtilmemişse, tüm ürünleri getir.
				return database.GetProducts()
			},
		},
		// "product" endpoint'ı, belirli bir ürünü alır.
		"product": &graphql.Field{
			Type: productType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, ok := params.Args["id"].(int)
				if !ok {
					// ID belirtilmemişse, null değeri döndür.
					return nil, nil
				}
				// Belirtilen ID'ye sahip ürünü getir.
				return database.GetProductByID(id)
			},
		},
	},
})

var introspectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "__Introspection",
	Fields: graphql.Fields{
		"__typename": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:        rootQuery,
	Mutation:     nil,
	Subscription: nil,
	Types: []graphql.Type{
		introspectionType,
	},
})

// @Summary Execute GraphQL queries.
// @Description This endpoint allows you to execute GraphQL queries and retrieve information.
// @Accept json
// @Produce json
// @Param input body string true "GraphQL query"
// @Success 200 {object} ProductType
// @Router /graphql [post]

// GraphQLHandler, GraphQL sorgularını işleyen bir HTTP handler'ıdır.
func GraphQLHandler(c *gin.Context) {
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	h.ServeHTTP(c.Writer, c.Request)
}

// @Summary Check the health status of the API.
// @Description This endpoint checks the health status of the API.
// @Produce json
// @Success 200 {object} gin.H
// @Router /health [get]

// HealthHandler, sağlık kontrol endpoint'ını işler.
func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "OK",
	})
}
