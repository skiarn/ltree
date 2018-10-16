package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "testuser"
	password = "testuser"
	dbname   = "postgres"
)

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	defer db.Close()
	h := handler{db}
	http.HandleFunc("/hierarchy", h.hierarchyHandler)
	log.Fatal(http.ListenAndServe("localhost:9000", nil))
}

type handler struct {
	db *sql.DB
}

func (h handler) hierarchyHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		fmt.Printf("Request hierarchy %s took: %v \n", r.Method, time.Since(start))
	}()
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "PasreForm() err:%v", err)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		fmt.Fprintf(w, "missing id!")
		return
	}

	switch r.Method {
	case "GET":
		rows, err := h.db.Query("select id, nodeid, path from hierarchy where (select subpath(path, 0, 1) from hierarchy where nodeid = $1) @> path", id)
		if err != nil {
			fmt.Fprintf(w, "unable to query(%s)", id)
			return
		}

		lookupNode := make(map[int64]string)
		var rootNode *Node
		var curNode *Node
		for rows.Next() {
			var pk int64
			var nodeid string
			var path string

			err = rows.Scan(&pk, &nodeid, &path)
			if err != nil {
				fmt.Fprintf(w, "error scanning rows got err: %v", err)
			}

			lookupNode[pk] = nodeid
			if curNode == nil {
				rootNode = &Node{ID: nodeid}
				curNode = rootNode
				continue
			}
			for _, pathv := range strings.Split(path, ".") {
				pathValue, err := strconv.ParseInt(pathv, 10, 64)
				if err != nil {
					fmt.Fprintf(w, "Internal error parshing path %s", path)
					return
				}
				for _, n := range curNode.Children {
					if n.ID == lookupNode[pathValue] {
						curNode = n
						break
					}
				}
			}
			curNode.Children = append(curNode.Children, &Node{ID: nodeid})
			curNode = rootNode
		}
		err = json.NewEncoder(w).Encode(rootNode)
		if err != nil {
			fmt.Fprintf(w, "Error while encoding to json, got err: %v", err)
		}
		return
	case "POST":
		var parent string
		parent = r.FormValue("parent")
		if parent == "" {
			//insert into hierarchy(nodeid, path) values ('123e4567-e89b-12d3-a456-426655440001', CONCAT('1.', currval('hierarchy_id_seq')::text)::ltree);
			sqlStatement := "insert into hierarchy (nodeid, path) values ($1, currval('hierarchy_id_seq')::text::ltree)"
			_, err := h.db.Exec(sqlStatement, id)
			if err != nil {
				fmt.Fprintf(w, "unable to insert node %s got err: %v", id, err)
			}
		} else {
			sqlStatement := "insert into hierarchy (nodeid, path) values ($1, CONCAT(CONCAT((select path from hierarchy where nodeid = $2), '.'), currval('hierarchy_id_seq')::text)::ltree)"
			_, err := h.db.Exec(sqlStatement, id, parent)
			if err != nil {
				fmt.Fprintf(w, "unable to insert node %s got err: %v", id, err)
			}
		}

		break
	default:
		fmt.Fprintf(w, "unable to complete request unsupported method: %s", r.Method)
	}
	fmt.Fprintf(w, id)
	return
}

type Row struct {
	ID     int64
	Nodeid string
	Path   string
}

type Node struct {
	ID       string
	Children []*Node
}
