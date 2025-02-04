package nodes

import (
	"CivicCoinMain/pkg/models"
	"fmt"
	"log"
	"time"

	"CivicCoinMain/pkg/db/redis"
)

// SaveNode guarda la información de un nodo en Redis.
func SaveNode(node models.Node) error {

	// Guardar el nodo en un hash
	key := fmt.Sprintf("node:%s", node.Id)
	err := redis.Client.HSet(redis.Ctx, key, map[string]interface{}{
		"id":           node.Id,
		"ip_address":   node.Addr,
		"status":       node.Status,
		"last_updated": node.LastUpdated.Format(time.RFC3339),
	}).Err()
	if err != nil {
		return fmt.Errorf("error guardando nodo en Redis: %v", err)
	}

	log.Println("saved node id:", node.Id)

	return nil
}

// GetNode obtiene la información de un nodo por su ID.
func GetNode(id string) (*models.Node, error) {
	key := fmt.Sprintf("node:%s", id)

	// Obtener los campos del hash
	result, err := redis.Client.HGetAll(redis.Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo nodo de Redis: %v", err)
	}

	// Verificar si el nodo existe
	if len(result) == 0 {
		return nil, fmt.Errorf("nodo no encontrado")
	}

	// Convertir el resultado a un struct Node
	lastUpdated, err := time.Parse(time.RFC3339, result["last_updated"])
	if err != nil {
		return nil, fmt.Errorf("error parseando last_updated: %v", err)
	}

	node := &models.Node{
		Id:          result["id"],
		Addr:        result["ip_address"],
		Status:      result["status"],
		LastUpdated: lastUpdated,
	}

	return node, nil
}

// GetAllNodes obtiene todos los nodos almacenados en Redis.
func GetAllNodes() ([]models.Node, error) {

	log.Println("getting all nodes")

	// Obtener todas las claves que coinciden con el patrón "node:*"
	keys, err := redis.Client.Keys(redis.Ctx, "node:*").Result()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo claves de nodos: %v", err)
	}

	var nodes []models.Node
	for _, key := range keys {
		node, err := GetNode(key[len("node:"):]) // Extraer el ID del nodo
		if err != nil {
			return nil, fmt.Errorf("error obteniendo nodo %s: %v", key, err)
		}
		nodes = append(nodes, *node)
	}

	return nodes, nil
}
