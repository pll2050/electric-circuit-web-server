package repositories

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
func (s *FirestoreService) CreateDocument(collection, docID string, data map[string]interface{}) error {
	// Add timestamp
	data["createdAt"] = time.Now()
	data["updatedAt"] = time.Now()

	_, err := s.client.Collection(collection).Doc(docID).Set(s.ctx, data)
	if err != nil {
		return fmt.Errorf("error creating document: %v", err)
	}
	return nil
}

// CreateDocumentWithAutoID creates a new document with auto-generated ID
func (s *FirestoreService) CreateDocumentWithAutoID(collection string, data map[string]interface{}) (string, error) {
	// Add timestamp
	data["createdAt"] = time.Now()
	data["updatedAt"] = time.Now()

	docRef, _, err := s.client.Collection(collection).Add(s.ctx, data)
	if err != nil {
		return "", fmt.Errorf("error creating document: %v", err)
	}
	return docRef.ID, nil
}

// GetDocument retrieves a document by ID
func (s *FirestoreService) GetDocument(collection, docID string) (map[string]interface{}, error) {
	doc, err := s.client.Collection(collection).Doc(docID).Get(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting document: %v", err)
	}

	if !doc.Exists() {
		return nil, fmt.Errorf("document not found")
	}

	return doc.Data(), nil
}

// UpdateDocument updates a document
func (s *FirestoreService) UpdateDocument(collection, docID string, data map[string]interface{}) error {
	// Add timestamp
	data["updatedAt"] = time.Now()

	updates := make([]firestore.Update, 0, len(data))
	for key, value := range data {
		updates = append(updates, firestore.Update{Path: key, Value: value})
	}

	_, err := s.client.Collection(collection).Doc(docID).Update(s.ctx, updates)
	if err != nil {
		return fmt.Errorf("error updating document: %v", err)
	}

	return nil
}

// DeleteDocument deletes a document
func (s *FirestoreService) DeleteDocument(collection, docID string) error {
	_, err := s.client.Collection(collection).Doc(docID).Delete(s.ctx)
	if err != nil {
		return fmt.Errorf("error deleting document: %v", err)
	}
	return nil
}

// GetCollection retrieves all documents in a collection
func (s *FirestoreService) GetCollection(collection string) ([]map[string]interface{}, error) {
	iter := s.client.Collection(collection).Documents(s.ctx)
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
func (s *FirestoreService) QueryCollection(collection, field, operator string, value interface{}) ([]map[string]interface{}, error) {
	var query firestore.Query

	switch operator {
	case "==":
		query = s.client.Collection(collection).Where(field, "==", value)
	case "!=":
		query = s.client.Collection(collection).Where(field, "!=", value)
	case "<":
		query = s.client.Collection(collection).Where(field, "<", value)
	case "<=":
		query = s.client.Collection(collection).Where(field, "<=", value)
	case ">":
		query = s.client.Collection(collection).Where(field, ">", value)
	case ">=":
		query = s.client.Collection(collection).Where(field, ">=", value)
	case "in":
		query = s.client.Collection(collection).Where(field, "in", value)
	case "array-contains":
		query = s.client.Collection(collection).Where(field, "array-contains", value)
	default:
		return nil, fmt.Errorf("unsupported operator: %s", operator)
	}

	iter := query.Documents(s.ctx)
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
func (s *FirestoreService) BatchWrite(operations []BatchOperation) error {
	batch := s.client.Batch()

	for _, op := range operations {
		docRef := s.client.Collection(op.Collection).Doc(op.DocID)

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

	_, err := batch.Commit(s.ctx)
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
