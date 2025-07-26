package repository

import (
	"context"
	"orderfoodonline/internal/metrics"
	"orderfoodonline/internal/repository/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// productRepository provides MongoDB-backed access to product data.
type productRepository struct {
	collection           *mongo.Collection
	migrationsCollection *mongo.Collection
}

// NewProductRepository creates a new ProductRepository using the given Repository.
func NewProductRepository(repo *Repository) (ProductRepository, error) {
	collection := repo.db.Collection("products")
	migrationsCollection := repo.db.Collection("migrations")
	return &productRepository{
		collection:           collection,
		migrationsCollection: migrationsCollection,
	}, nil
}

// ListProducts returns all products from the database.
func (r *productRepository) ListProducts(ctx context.Context) ([]models.Product, error) {
	start := time.Now()

	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		metrics.RecordDatabaseQuery("find", "products", "error", time.Since(start).Seconds())
		return nil, err
	}
	defer cur.Close(ctx)

	var products []models.Product
	for cur.Next(ctx) {
		var p models.Product
		if err := cur.Decode(&p); err != nil {
			metrics.RecordDatabaseQuery("find", "products", "error", time.Since(start).Seconds())
			return nil, err
		}
		products = append(products, p)
	}
	if err := cur.Err(); err != nil {
		metrics.RecordDatabaseQuery("find", "products", "error", time.Since(start).Seconds())
		return nil, err
	}

	metrics.RecordDatabaseQuery("find", "products", "success", time.Since(start).Seconds())
	return products, nil
}

// FindProductByID returns a product by its ID, or nil if not found.
func (r *productRepository) FindProductByID(ctx context.Context, id string) (*models.Product, error) {
	start := time.Now()

	var p models.Product
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		metrics.RecordDatabaseQuery("find_one", "products", "not_found", time.Since(start).Seconds())
		return nil, nil
	}
	if err != nil {
		metrics.RecordDatabaseQuery("find_one", "products", "error", time.Since(start).Seconds())
		return nil, err
	}

	metrics.RecordDatabaseQuery("find_one", "products", "success", time.Since(start).Seconds())
	return &p, nil
}

// BulkInsertProducts inserts multiple products into the database.
func (r *productRepository) BulkInsertProducts(ctx context.Context, products []models.Product) error {
	if len(products) == 0 {
		return nil
	}

	start := time.Now()
	var documents []interface{}
	for _, product := range products {
		documents = append(documents, product)
	}
	_, err := r.collection.InsertMany(ctx, documents)

	if err != nil {
		metrics.RecordDatabaseQuery("insert_many", "products", "error", time.Since(start).Seconds())
		return err
	}

	metrics.RecordDatabaseQuery("insert_many", "products", "success", time.Since(start).Seconds())
	return nil
}

// GetAppliedMigrations returns all applied migrations from the database.
func (r *productRepository) GetAppliedMigrations(ctx context.Context) ([]models.Migration, error) {
	start := time.Now()

	cur, err := r.migrationsCollection.Find(ctx, bson.M{})
	if err != nil {
		metrics.RecordDatabaseQuery("find", "migrations", "error", time.Since(start).Seconds())
		return nil, err
	}
	defer cur.Close(ctx)

	var migrations []models.Migration
	for cur.Next(ctx) {
		var migration models.Migration
		if err := cur.Decode(&migration); err != nil {
			metrics.RecordDatabaseQuery("find", "migrations", "error", time.Since(start).Seconds())
			return nil, err
		}
		migrations = append(migrations, migration)
	}
	if err := cur.Err(); err != nil {
		metrics.RecordDatabaseQuery("find", "migrations", "error", time.Since(start).Seconds())
		return nil, err
	}

	metrics.RecordDatabaseQuery("find", "migrations", "success", time.Since(start).Seconds())
	return migrations, nil
}

// InsertMigration inserts a new migration record into the database.
func (r *productRepository) InsertMigration(ctx context.Context, migration *models.Migration) error {
	start := time.Now()

	_, err := r.migrationsCollection.InsertOne(ctx, migration)

	if err != nil {
		metrics.RecordDatabaseQuery("insert_one", "migrations", "error", time.Since(start).Seconds())
		return err
	}

	metrics.RecordDatabaseQuery("insert_one", "migrations", "success", time.Since(start).Seconds())
	return nil
}

// UpdateMigration updates an existing migration record in the database.
func (r *productRepository) UpdateMigration(ctx context.Context, migration *models.Migration) error {
	start := time.Now()

	filter := bson.M{"id": migration.ID}
	update := bson.M{"$set": bson.M{"status": migration.Status}}
	_, err := r.migrationsCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		metrics.RecordDatabaseQuery("update_one", "migrations", "error", time.Since(start).Seconds())
		return err
	}

	metrics.RecordDatabaseQuery("update_one", "migrations", "success", time.Since(start).Seconds())
	return nil
}
