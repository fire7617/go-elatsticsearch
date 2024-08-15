
# elasticsearch import file

curl -s -H "Content-Type: application/json" -XPOST  'http://localhost:9200/_bulk' --data-binary @2024-08-15.json
