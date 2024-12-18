// db/metadata_search.go
package db

import (
    "encoding/json"
    "fmt"
)

// MetadataQuery represents a complex metadata query structure
type MetadataQuery struct {
    Must     []MetadataCondition            `json:"must"`      // AND conditions
    Should   []MetadataCondition            `json:"should"`    // OR conditions
    MustNot  []MetadataCondition            `json:"must_not"`  // NOT conditions
    Range    map[string]map[string]interface{} `json:"range"`    // Range queries
}

type MetadataCondition struct {
    Path  string      `json:"path"`   // JSON path
    Op    string      `json:"op"`     // Operator: eq, contains, exists, etc.
    Value interface{} `json:"value"`  // Value to compare against
}

// buildMetadataQuery constructs the PostgreSQL JSON query
func buildMetadataQuery(query MetadataQuery) (string, []interface{}) {
    var conditions []string
    var args []interface{}
    argCount := 1

    // Process MUST conditions (AND)
    for _, cond := range query.Must {
        condition, value := buildMetadataCondition(cond, argCount)
        conditions = append(conditions, condition)
        args = append(args, value)
        argCount++
    }

    // Process SHOULD conditions (OR)
    if len(query.Should) > 0 {
        var orConditions []string
        for _, cond := range query.Should {
            condition, value := buildMetadataCondition(cond, argCount)
            orConditions = append(orConditions, condition)
            args = append(args, value)
            argCount++
        }
        if len(orConditions) > 0 {
            conditions = append(conditions, "("+strings.Join(orConditions, " OR ")+")")
        }
    }

    // Process MUST NOT conditions (NOT)
    for _, cond := range query.MustNot {
        condition, value := buildMetadataCondition(cond, argCount)
        conditions = append(conditions, "NOT "+condition)
        args = append(args, value)
        argCount++
    }

    // Process Range queries
    for field, ranges := range query.Range {
        for op, value := range ranges {
            condition := fmt.Sprintf("(metadata->>'%s')::numeric", field)
            switch op {
            case "gt":
                condition += fmt.Sprintf(" > $%d", argCount)
            case "gte":
                condition += fmt.Sprintf(" >= $%d", argCount)
            case "lt":
                condition += fmt.Sprintf(" < $%d", argCount)
            case "lte":
                condition += fmt.Sprintf(" <= $%d", argCount)
            }
            conditions = append(conditions, condition)
            args = append(args, value)
            argCount++
        }
    }

    return strings.Join(conditions, " AND "), args
}

func buildMetadataCondition(cond MetadataCondition, argNum int) (string, interface{}) {
    switch cond.Op {
    case "eq":
        return fmt.Sprintf("metadata->>'%s' = $%d", cond.Path, argNum), cond.Value
    case "contains":
        return fmt.Sprintf("metadata->>'%s' LIKE '%%' || $%d || '%%'", cond.Path, argNum), cond.Value
    case "exists":
        return fmt.Sprintf("metadata ? '%s'", cond.Path), nil
    case "array_contains":
        return fmt.Sprintf("metadata->>'%s' @> $%d", cond.Path, argNum), cond.Value
    default:
        return fmt.Sprintf("metadata->>'%s' = $%d", cond.Path, argNum), cond.Value
    }
}

// Example usage with comments
func ExampleSearchQueries() {
    // Example 1: Basic metadata search
    basicSearch := SearchParams{
        MetadataQuery: MetadataQuery{
            Must: []MetadataCondition{
                {
                    Path:  "format",
                    Op:    "eq",
                    Value: "HD",
                },
            },
        },
    }

    // Example 2: Complex metadata search with multiple conditions
    complexSearch := SearchParams{
        Title: "News Broadcast",
        Type:  "video",
        MetadataQuery: MetadataQuery{
            Must: []MetadataCondition{
                {
                    Path:  "resolution",
                    Op:    "eq",
                    Value: "1920x1080",
                },
            },
            Should: []MetadataCondition{
                {
                    Path:  "tags",
                    Op:    "array_contains",
                    Value: []string{"news", "broadcast"},
                },
            },
            MustNot: []MetadataCondition{
                {
                    Path:  "status",
                    Op:    "eq",
                    Value: "archived",
                },
            },
            Range: map[string]map[string]interface{}{
                "duration": {
                    "gte": 300,  // Duration >= 5 minutes
                    "lte": 1800, // Duration <= 30 minutes
                },
            },
        },
    }

    // Example 3: Search for assets with specific technical metadata
    technicalSearch := SearchParams{
        MetadataQuery: MetadataQuery{
            Must: []MetadataCondition{
                {
                    Path:  "codec",
                    Op:    "eq",
                    Value: "H.264",
                },
                {
                    Path:  "bitrate",
                    Op:    "exists",
                },
            },
            Range: map[string]map[string]interface{}{
                "frame_rate": {
                    "gte": 24,
                },
            },
        },
    }

    // Example 4: Search for assets with specific production metadata
    productionSearch := SearchParams{
        MetadataQuery: MetadataQuery{
            Must: []MetadataCondition{
                {
                    Path:  "production_name",
                    Op:    "contains",
                    Value: "World Cup",
                },
            },
            Should: []MetadataCondition{
                {
                    Path:  "location",
                    Op:    "eq",
                    Value: "Stadium A",
                },
                {
                    Path:  "location",
                    Op:    "eq",
                    Value: "Stadium B",
                },
            },
        },
        CreatedAfter: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
        SortBy:       "created_at",
        SortOrder:    "DESC",
    }

    // Example 5: Complex rights management search
    rightsSearch := SearchParams{
        MetadataQuery: MetadataQuery{
            Must: []MetadataCondition{
                {
                    Path:  "rights_status",
                    Op:    "eq",
                    Value: "cleared",
                },
            },
            MustNot: []MetadataCondition{
                {
                    Path:  "rights_expiration",
                    Op:    "exists",
                },
            },
        },
        Facets: []string{"rights_territory", "rights_platform"},
    }

    // Example usage in code
    result, err := repo.SearchAssets(context.Background(), complexSearch)
    if err != nil {
        log.Fatal(err)
    }
}

// Example API usage showing how to construct and execute a search
func ExampleAPIUsage() {
    // Create a new search request
    searchRequest := `{
        "query": "football match",
        "type": "video",
        "metadata": {
            "must": [
                {
                    "path": "sport",
                    "op": "eq",
                    "value": "football"
                }
            ],
            "should": [
                {
                    "path": "competition",
                    "op": "eq",
                    "value": "Premier League"
                },
                {
                    "path": "competition",
                    "op": "eq",
                    "value": "Champions League"
                }
            ],
            "range": {
                "duration": {
                    "gte": 5400,
                    "lte": 7200
                }
            }
        },
        "created_after": "2024-01-01T00:00:00Z",
        "sort_by": "created_at",
        "sort_order": "DESC",
        "limit": 20,
        "offset": 0
    }`

    // Parse and execute the search
    var params SearchParams
    if err := json.Unmarshal([]byte(searchRequest), &params); err != nil {
        log.Fatal(err)
    }

    result, err := repo.SearchAssets(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }
}
