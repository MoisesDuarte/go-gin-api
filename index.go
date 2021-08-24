package main 

import (
	"errors"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const dbPath = "database.db"

	db := InitDB(dbPath)
	defer db.Close()

	CreateTable(db)

	GenerateProductsSheet(db)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is online!",
		})
	})

	// Create product
	router.POST("/products", func(c *gin.Context) {
		id := c.PostForm("id")
		name := c.PostForm("name")
		description := c.PostForm("description")
		price := c.PostForm("price")

		product := []ProductItem {
			ProductItem { id, name, description, price },
		}
		
		statementResult := StoreProduct(db, product)

		if statementResult != nil {
			c.JSON(200, gin.H{
				"status": "created",
				"name": name,
				"description": description,
				"price": price,
			})
		} else {
			c.JSON(500, gin.H{
				"status": "Error running create action on product",
			})
		}
	})

	// List all products
	router.GET("/products", func(c *gin.Context) {
		products, err := ReadProducts(db)

		if err != nil {
			c.JSON(500, gin.H{
				"message": "Error running statement",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Sucesfully fetched products",
				"products": products,
			})
		}
		
	})

	router.Run(":8080")
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func CreateTable(db *sql.DB) {
	sql_create_table := `
		CREATE TABLE IF NOT EXISTS products(
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT,
			description TEXT,
			price REAL,
			insertedAt DATETIME
		);
	`

	_, err := db.Exec(sql_create_table)
	if err != nil { panic(err) }
}

func StoreProduct(db *sql.DB, products []ProductItem) error {
	sql_add_product := `
		INSERT OR REPLACE INTO products(id, name, description, price, insertedAt)
		values (?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	statement, err := db.Prepare(sql_add_product)

	if err != nil {
		return errors.New("Error running create statement")
	}

	defer statement.Close()

	for _, product := range products {
		_, err2 := statement.Exec(product.Id, product.Name, product.Description, product.Price)
		if err2 != nil { panic(err2) }
	}

	return nil;
}

func ReadProducts(db *sql.DB) ([]ProductItem, error) {
	sql_read_products := `
		SELECT id, name, description, price FROM products
		ORDER BY datetime(insertedAt) DESC
	`

	rows, err := db.Query(sql_read_products)
	
	if err != nil { 
		return nil, errors.New("Error on query execution")
	}

	defer rows.Close()

	var result []ProductItem
	for rows.Next() {
		product := ProductItem{}

		err2 := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
		if err2 != nil {
			return nil, errors.New("Error on row scan")
		}

		result = append(result, product)
	}

	return result, nil
}

func GenerateProductsSheet(db *sql.DB) {
	result, err1 := ReadProducts(db)

	file := excelize.NewFile()

	file.SetCellValue("Sheet1", "A1", "Id")
	file.SetCellValue("Sheet1", "A2", "Name")
	file.SetCellValue("Sheet1", "A3", "Description")
	file.SetCellValue("Sheet1", "A4", "Price")

	

	if err := file.SaveAs("Book1.xlsx"); err != nil {}
        //fmt.Println(err)
    }
}

type ProductItem struct {
	Id string
	Name string
	Description string
	Price string
}
