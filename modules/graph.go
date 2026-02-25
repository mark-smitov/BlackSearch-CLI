package modules

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type GraphNode struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
}

type GraphEdge struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	Target       string `json:"target"`
	Relationship string `json:"relationship"`
}

type GraphDB struct {
	Nodes map[string]GraphNode `json:"nodes"`
	Edges map[string]GraphEdge `json:"edges"`
}

const dbFile = "graph_data/graph.json"

var nodeColors = map[string]string{
	"person": "#ff79c6", "organization": "#8be9fd", "domain": "#50fa7b",
	"ip": "#ffb86c", "email": "#f1fa8c", "phone": "#bd93f9",
	"social": "#6be5fd", "crypto": "#ffd580", "default": "#6272a4",
}

var nodeIcons = map[string]string{
	"person": "👤", "organization": "🏢", "domain": "🌐",
	"ip": "💻", "email": "📧", "phone": "📱",
	"social": "🔗", "crypto": "🪙", "default": "📍",
}

func newUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func loadDB() GraphDB {
	db := GraphDB{
		Nodes: make(map[string]GraphNode),
		Edges: make(map[string]GraphEdge),
	}
	data, err := os.ReadFile(dbFile)
	if err != nil {
		return db
	}
	json.Unmarshal(data, &db)
	if db.Nodes == nil {
		db.Nodes = make(map[string]GraphNode)
	}
	if db.Edges == nil {
		db.Edges = make(map[string]GraphEdge)
	}
	return db
}

func saveDB(db GraphDB) {
	os.MkdirAll(filepath.Dir(dbFile), 0755)
	data, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile(dbFile, data, 0644)
}

func graphHTML() (string, error) {
	data, err := os.ReadFile("html/graph.html")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func RunGraphServer() {
	PrintInfo("Запуск веб-интерфейса на http://127.0.0.1:8765")
	PrintInfo("Данные: " + dbFile)
	PrintWarn("Ctrl+C для остановки")
	fmt.Println()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	r.GET("/", func(c *gin.Context) {
		html, err := graphHTML()
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка: файл html/graph.html не найден")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	r.GET("/api/graph", func(c *gin.Context) {
		db := loadDB()

		type VisNode struct {
			ID         string         `json:"id"`
			Label      string         `json:"label"`
			CleanLabel string         `json:"cleanLabel"`
			NodeType   string         `json:"nodeType"`
			NodeValue  string         `json:"nodeValue"`
			CreatedAt  string         `json:"createdAt"`
			Color      map[string]any `json:"color"`
			Shape      string         `json:"shape"`
			Title      string         `json:"title"`
		}
		type VisEdge struct {
			ID     string `json:"id"`
			From   string `json:"from"`
			To     string `json:"to"`
			Label  string `json:"label"`
			Source string `json:"source"`
			Target string `json:"target"`
		}

		var visNodes []VisNode
		var visEdges []VisEdge

		for _, n := range db.Nodes {
			clr := nodeColors["default"]
			if v, ok := nodeColors[n.Type]; ok {
				clr = v
			}
			icon := nodeIcons["default"]
			if ic, ok := nodeIcons[n.Type]; ok {
				icon = ic
			}
			tooltip := fmt.Sprintf("<b>%s</b><br>Тип: %s", n.Label, n.Type)
			if n.Value != "" {
				tooltip += "<br>Значение: " + n.Value
			}
			visNodes = append(visNodes, VisNode{
				ID:         n.ID,
				Label:      icon + " " + n.Label,
				CleanLabel: n.Label,
				NodeType:   n.Type,
				NodeValue:  n.Value,
				CreatedAt:  n.CreatedAt,
				Color: map[string]any{
					"background": "#141414",
					"border":     clr,
					"highlight":  map[string]string{"background": "#1e1e1e", "border": clr},
				},
				Shape: "dot",
				Title: tooltip,
			})
		}

		for _, e := range db.Edges {
			visEdges = append(visEdges, VisEdge{
				ID:     e.ID,
				From:   e.Source,
				To:     e.Target,
				Label:  e.Relationship,
				Source: e.Source,
				Target: e.Target,
			})
		}

		if visNodes == nil {
			visNodes = []VisNode{}
		}
		if visEdges == nil {
			visEdges = []VisEdge{}
		}

		c.JSON(http.StatusOK, gin.H{"nodes": visNodes, "edges": visEdges})
	})

	r.POST("/api/nodes", func(c *gin.Context) {
		var node GraphNode
		if err := c.ShouldBindJSON(&node); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
			return
		}
		if node.Label == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "label обязателен"})
			return
		}
		node.ID = newUUID()
		node.CreatedAt = time.Now().Format(time.RFC3339)

		db := loadDB()
		db.Nodes[node.ID] = node
		saveDB(db)

		c.JSON(http.StatusOK, node)
	})

	r.DELETE("/api/nodes/:id", func(c *gin.Context) {
		id := c.Param("id")
		db := loadDB()
		if _, ok := db.Nodes[id]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "узел не найден"})
			return
		}
		delete(db.Nodes, id)
		for eid, edge := range db.Edges {
			if edge.Source == id || edge.Target == id {
				delete(db.Edges, eid)
			}
		}
		saveDB(db)
		c.JSON(http.StatusOK, gin.H{"status": "deleted", "id": id})
	})

	r.POST("/api/edges", func(c *gin.Context) {
		var edge GraphEdge
		if err := c.ShouldBindJSON(&edge); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
			return
		}
		db := loadDB()
		if _, ok := db.Nodes[edge.Source]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "узел-источник не найден"})
			return
		}
		if _, ok := db.Nodes[edge.Target]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "узел-цель не найден"})
			return
		}
		edge.ID = newUUID()
		db.Edges[edge.ID] = edge
		saveDB(db)
		c.JSON(http.StatusOK, edge)
	})

	r.DELETE("/api/edges/:id", func(c *gin.Context) {
		id := c.Param("id")
		db := loadDB()
		if _, ok := db.Edges[id]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "связь не найдена"})
			return
		}
		delete(db.Edges, id)
		saveDB(db)
		c.JSON(http.StatusOK, gin.H{"status": "deleted", "id": id})
	})

	r.DELETE("/api/clear", func(c *gin.Context) {
		db := GraphDB{
			Nodes: make(map[string]GraphNode),
			Edges: make(map[string]GraphEdge),
		}
		saveDB(db)
		c.JSON(http.StatusOK, gin.H{"status": "cleared"})
	})

	if err := r.Run("127.0.0.1:8765"); err != nil {
		PrintError("Ошибка сервера: " + err.Error())
	}
}
