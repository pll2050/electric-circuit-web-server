package firebase

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreService handles Firestore database operations
type FirestoreService struct {
	client *firestore.Client
	ctx    context.Context
}

// NewFirestoreService creates a new Firestore service
func NewFirestoreService(client *firestore.Client, ctx context.Context) *FirestoreService {
	return &FirestoreService{
		client: client,
		ctx:    ctx,
	}
}

// CreateDocument creates a new document in a collection
func (fs *FirestoreService) CreateDocument(collection, docID string, data map[string]interface{}) error {
	// Add timestamp
	data["createdAt"] = time.Now()
	data["updatedAt"] = time.Now()

	_, err := fs.client.Collection(collection).Doc(docID).Set(fs.ctx, data)
	if err != nil {
		return fmt.Errorf("error creating document: %v", err)
	}
	return nil
}

// CreateDocumentWithAutoID creates a new document with auto-generated ID
func (fs *FirestoreService) CreateDocumentWithAutoID(collection string, data map[string]interface{}) (string, error) {
	// Add timestamp
	data["createdAt"] = time.Now()
	data["updatedAt"] = time.Now()

	docRef, _, err := fs.client.Collection(collection).Add(fs.ctx, data)
	if err != nil {
		return "", fmt.Errorf("error creating document: %v", err)
	}
	return docRef.ID, nil
}

// GetDocument retrieves a document by ID
func (fs *FirestoreService) GetDocument(collection, docID string) (map[string]interface{}, error) {
	doc, err := fs.client.Collection(collection).Doc(docID).Get(fs.ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting document: %v", err)
	}

	if !doc.Exists() {
		return nil, fmt.Errorf("document not found")
	}

	return doc.Data(), nil
}

// UpdateDocument updates a document
func (fs *FirestoreService) UpdateDocument(collection, docID string, data map[string]interface{}) error {
	// Add timestamp
	data["updatedAt"] = time.Now()

	_, err := fs.client.Collection(collection).Doc(docID).Update(fs.ctx,
		[]firestore.Update{
			{Path: "updatedAt", Value: time.Now()},
		})
	if err != nil {
		return fmt.Errorf("error updating document: %v", err)
	}

	// Update other fields
	for key, value := range data {
		if key != "updatedAt" {
			_, err = fs.client.Collection(collection).Doc(docID).Update(fs.ctx,
				[]firestore.Update{
					{Path: key, Value: value},
				})
			if err != nil {
				return fmt.Errorf("error updating field %s: %v", key, err)
			}
		}
	}

	return nil
}

// DeleteDocument deletes a document
func (fs *FirestoreService) DeleteDocument(collection, docID string) error {
	_, err := fs.client.Collection(collection).Doc(docID).Delete(fs.ctx)
	if err != nil {
		return fmt.Errorf("error deleting document: %v", err)
	}
	return nil
}

// GetCollection retrieves all documents in a collection
func (fs *FirestoreService) GetCollection(collection string) ([]map[string]interface{}, error) {
	iter := fs.client.Collection(collection).Documents(fs.ctx)
	defer iter.Stop()

	var documents []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating documents: %v", err)
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID // Add document ID to data
		documents = append(documents, data)
	}

	return documents, nil
}

// QueryCollection queries a collection with conditions
func (fs *FirestoreService) QueryCollection(collection, field, operator string, value interface{}) ([]map[string]interface{}, error) {
	var query firestore.Query

	switch operator {
	case "==":
		query = fs.client.Collection(collection).Where(field, "==", value)
	case "!=":
		query = fs.client.Collection(collection).Where(field, "!=", value)
	case "<":
		query = fs.client.Collection(collection).Where(field, "<", value)
	case "<=":
		query = fs.client.Collection(collection).Where(field, "<=", value)
	case ">":
		query = fs.client.Collection(collection).Where(field, ">", value)
	case ">=":
		query = fs.client.Collection(collection).Where(field, ">=", value)
	case "in":
		query = fs.client.Collection(collection).Where(field, "in", value)
	case "array-contains":
		query = fs.client.Collection(collection).Where(field, "array-contains", value)
	default:
		return nil, fmt.Errorf("unsupported operator: %s", operator)
	}

	iter := query.Documents(fs.ctx)
	defer iter.Stop()

	var documents []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating documents: %v", err)
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		documents = append(documents, data)
	}

	return documents, nil
}

// BatchWrite performs batch write operations
func (fs *FirestoreService) BatchWrite(operations []BatchOperation) error {
	batch := fs.client.Batch()

	for _, op := range operations {
		docRef := fs.client.Collection(op.Collection).Doc(op.DocID)

		switch op.Type {
		case "create":
			op.Data["createdAt"] = time.Now()
			op.Data["updatedAt"] = time.Now()
			batch.Set(docRef, op.Data)
		case "update":
			op.Data["updatedAt"] = time.Now()
			updates := make([]firestore.Update, 0, len(op.Data))
			for key, value := range op.Data {
				updates = append(updates, firestore.Update{Path: key, Value: value})
			}
			batch.Update(docRef, updates)
		case "delete":
			batch.Delete(docRef)
		}
	}

	_, err := batch.Commit(fs.ctx)
	if err != nil {
		return fmt.Errorf("error committing batch: %v", err)
	}
	return nil
}

// BatchOperation represents a batch operation
type BatchOperation struct {
	Type       string // "create", "update", "delete"
	Collection string
	DocID      string
	Data       map[string]interface{}
}
